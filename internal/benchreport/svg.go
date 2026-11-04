package main

import (
	"fmt"
	"html"
	"io"
	"math"
	"regexp"
	"slices"
	"strings"
	"time"
)

// The chart is laid out in pixels. Text widths are estimated from character counts, since an SVG can't measure its
// text ahead of time; the estimate is generous for the system sans-serif fonts the chart asks for.
const (
	chartPad   = 28.0
	charWidth  = 7.6 // at the 13px name size
	rowHeight  = 28.0
	plotWidth  = 340.0
	catWidth   = 110.0
	cellHeight = 28.0
	cellGap    = 2.0
)

// ramp is the sequential blue scale for the heatmap, lightest first. In the light scheme a cell's step is its index;
// in the dark scheme the ramp runs the other way, so fast cells recede into the dark surface as they do into the
// light one. Text on a step is dark up to rampDarkText and white past it, which keeps it above 4.5:1.
var ramp = []string{"#cde2fb", "#b7d3f6", "#9ec5f4", "#86b6ef", "#6da7ec", "#5598e7", "#3987e5", "#2a78d6", "#256abf", "#1c5cab", "#184f95", "#104281", "#0d366b"}

const rampDarkText = 6

func chartStyle() string {
	var b strings.Builder
	b.WriteString(`
svg { --surface: #fcfcfb; --rule: #d6d5cf; --rule-strong: #8a8984; --grid: #e6e5e0; --one: #b9b8b2;
      --text-primary: #0b0b0b; --text-secondary: #52514e; --text-muted: #6f6e69; --series: #2a78d6; }
@media (prefers-color-scheme: dark) {
  svg { --surface: #1a1a19; --rule: #383835; --rule-strong: #77766f; --grid: #2a2a28; --one: #4a4a47;
        --text-primary: #f0efec; --text-secondary: #c3c2b7; --text-muted: #9a9991; --series: #3987e5; }
}
text { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif; fill: var(--text-primary); }
.bg { fill: var(--surface); }
.rule { stroke: var(--rule); stroke-width: 1; }
.rule-strong { stroke: var(--rule-strong); stroke-width: 1; }
.grid { stroke: var(--grid); stroke-width: 1; }
.one { stroke: var(--one); stroke-width: 1; }
.whisker { stroke: var(--series); stroke-width: 2; stroke-linecap: round; }
.dot { fill: var(--series); stroke: var(--surface); stroke-width: 2; }
.title { font-size: 20px; font-weight: 600; }
.fig { font-size: 16px; font-weight: 600; }
.sub { font-size: 12.5px; fill: var(--text-secondary); }
.col { font-size: 10.5px; fill: var(--text-muted); letter-spacing: 0.08em; }
.tick { font-size: 11px; fill: var(--text-muted); }
.rank { font-size: 12.5px; fill: var(--text-muted); font-variant-numeric: tabular-nums; }
.name { font-size: 13px; font-weight: 600; }
.overall { font-size: 13px; font-weight: 600; font-variant-numeric: tabular-nums; }
.value { font-size: 12.5px; font-variant-numeric: tabular-nums; }
.ord { font-size: 11px; fill: var(--text-muted); }
.bench { font-size: 12px; font-weight: 600; }
.fastest { font-size: 11px; fill: var(--text-muted); font-variant-numeric: tabular-nums; }
.cell { font-size: 12px; font-variant-numeric: tabular-nums; }
.foot { font-size: 11px; fill: var(--text-muted); }
`)
	textFor := func(step int) string {
		if step <= rampDarkText {
			return "#0b0b0b"
		}
		return "#ffffff"
	}
	var light, dark strings.Builder
	for i := range ramp {
		fmt.Fprintf(&light, ".h%d rect { fill: %s; } .h%d text { fill: %s; }\n", i, ramp[i], i, textFor(i))
		j := len(ramp) - 1 - i
		fmt.Fprintf(&dark, "  .h%d rect { fill: %s; } .h%d text { fill: %s; }\n", i, ramp[j], i, textFor(j))
	}
	b.WriteString(light.String())
	b.WriteString("@media (prefers-color-scheme: dark) {\n")
	b.WriteString(dark.String())
	b.WriteString("}\n")
	return b.String()
}

// logScale maps relative times onto [x0, x0+width] logarithmically: scores compare by how many times slower, so
// 1×→2× spans the same distance as 2×→4×.
type logScale struct{ lo, hi, x0, width float64 }

func (s logScale) x(v float64) float64 {
	return s.x0 + s.width*(math.Log(v)-math.Log(s.lo))/(math.Log(s.hi)-math.Log(s.lo))
}

// ticks returns round values within the scale, coarsest first, skipping any that would crowd one already placed.
func (s logScale) ticks() []float64 {
	var out []float64
	for _, level := range [][]float64{{1, 2, 4, 8, 16, 32, 64}, {1.5, 3, 6, 12, 24, 48}, {1.25, 2.5, 5, 10, 20, 40}, {0.5, 0.75, 0.9}} {
		for _, v := range level {
			if v < s.lo || v > s.hi {
				continue
			}
			if !slices.ContainsFunc(out, func(o float64) bool { return math.Abs(s.x(o)-s.x(v)) < 44 }) {
				out = append(out, v)
			}
		}
	}
	slices.Sort(out)
	return out
}

func formatTick(v float64) string {
	return strings.TrimSuffix(strings.TrimRight(fmt.Sprintf("%.2f", v), "0"), ".") + "×"
}

var wordRE = regexp.MustCompile(`[A-Z][a-z]+|[A-Z]+|[^A-Z]+`)

// splitName breaks a CamelCase benchmark name into two lines of about equal length, e.g. FindAuthor / ByName. A run
// of capitals (ID) is one word, and on a tie the first line takes more.
func splitName(name string) (string, string) {
	words := wordRE.FindAllString(name, -1)
	best, bestLen := 0, len(name)
	for i := 1; i < len(words); i++ {
		a := len(strings.Join(words[:i], ""))
		if l := max(a, len(name)-a); l <= bestLen {
			best, bestLen = i, l
		}
	}
	if best == 0 {
		return name, ""
	}
	return strings.Join(words[:best], ""), strings.Join(words[best:], "")
}

// writeSVG draws the summary in two figures. The leaderboard has a row per implementation in Overall order: its
// Overall score as a dot, with a whisker from the average of each benchmark's fastest run to the average of its
// slowest, and the category averages with their ranks; rules separate tie groups. The heatmap shows each
// implementation's score on each benchmark.
func writeSVG(w io.Writer, res *results, sum summary, commit string, now time.Time) {
	benches := res.benchmarks()
	esc := html.EscapeString

	nameWidth := 0.0
	for _, s := range sum.impls {
		nameWidth = max(nameWidth, float64(len([]rune(s.impl)))*charWidth)
	}
	rankX := chartPad
	nameX := rankX + 40
	overallRight := nameX + nameWidth + 66
	plotX := overallRight + 28
	catX := plotX + plotWidth + 20
	width := catX + float64(len(sum.cats))*catWidth + chartPad
	right := width - chartPad

	var body strings.Builder
	y := chartPad

	// Header.
	fmt.Fprintf(&body, `<text class="title" x="%.0f" y="%.0f">Go libraries for PostgreSQL</text>`+"\n", chartPad, y+18)
	y += 44
	for _, line := range []string{
		fmt.Sprintf("%d libraries run the same %d operations against a local PostgreSQL server, %d runs each.", len(sum.impls), len(benches), res.runs()),
		"A score is time relative to the fastest library on that benchmark: 1.00× is fastest, 2.00× takes twice as long.",
	} {
		fmt.Fprintf(&body, `<text class="sub" x="%.0f" y="%.0f">%s</text>`+"\n", chartPad, y, esc(line))
		y += 18
	}
	y += 22

	// Figure 1: leaderboard.
	fmt.Fprintf(&body, `<text class="fig" x="%.0f" y="%.0f">Leaderboard</text>`+"\n", chartPad, y)
	y += 22
	for _, line := range []string{
		"Overall is the geometric mean of the per-benchmark scores (dot). The line spans the range of runs, averaged",
		"the same way. Rules separate ties: rows in a group have runs that overlap the group's first row. Log scale.",
	} {
		fmt.Fprintf(&body, `<text class="sub" x="%.0f" y="%.0f">%s</text>`+"\n", chartPad, y, esc(line))
		y += 17
	}
	y += 18

	lo, hi := 1.0, 1.0
	for _, s := range sum.impls {
		lo, hi = min(lo, s.overall.lo), max(hi, s.overall.hi)
	}
	sc := logScale{lo / 1.03, hi * 1.03, plotX, plotWidth}

	fmt.Fprintf(&body, `<text class="col" x="%.0f" y="%.0f">RANK</text>`+"\n", rankX, y)
	fmt.Fprintf(&body, `<text class="col" x="%.0f" y="%.0f">LIBRARY</text>`+"\n", nameX, y)
	fmt.Fprintf(&body, `<text class="col" x="%.0f" y="%.0f" text-anchor="end">OVERALL</text>`+"\n", overallRight, y)
	for i, c := range sum.cats {
		fmt.Fprintf(&body, `<text class="col" x="%.0f" y="%.0f" text-anchor="end">%s</text>`+"\n", catX+float64(i)*catWidth+48, y, esc(strings.ToUpper(c)))
	}
	ticks := sc.ticks()
	for _, t := range ticks {
		fmt.Fprintf(&body, `<text class="tick" x="%.1f" y="%.0f" text-anchor="middle">%s</text>`+"\n", sc.x(t), y, formatTick(t))
	}
	y += 10
	top := y
	bottom := top + float64(len(sum.impls))*rowHeight
	for _, t := range ticks {
		class := "grid"
		if t == 1 {
			class = "one"
		}
		fmt.Fprintf(&body, `<line class="%s" x1="%.1f" x2="%.1f" y1="%.0f" y2="%.0f"/>`+"\n", class, sc.x(t), sc.x(t), top, bottom)
	}
	fmt.Fprintf(&body, `<line class="rule-strong" x1="%.0f" x2="%.0f" y1="%.0f" y2="%.0f"/>`+"\n", chartPad, right, top, top)
	for i, s := range sum.impls {
		ry := top + float64(i)*rowHeight
		if i > 0 && sum.ranks[i] != sum.ranks[i-1] {
			fmt.Fprintf(&body, `<line class="rule" x1="%.0f" x2="%.0f" y1="%.0f" y2="%.0f"/>`+"\n", chartPad, right, ry, ry)
		}
		cy := ry + rowHeight/2
		e := s.overall
		fmt.Fprintf(&body, "<g><title>%s: Overall %s, range %s–%s, rank %s</title>\n", esc(s.impl), formatRelative(e.mid), formatRelative(e.lo), formatRelative(e.hi), esc(sum.ranks[i]))
		fmt.Fprintf(&body, `<text class="rank" x="%.0f" y="%.1f" dominant-baseline="central">%s</text>`+"\n", rankX, cy, esc(sum.ranks[i]))
		fmt.Fprintf(&body, `<text class="name" x="%.0f" y="%.1f" dominant-baseline="central">%s</text>`+"\n", nameX, cy, esc(s.impl))
		fmt.Fprintf(&body, `<text class="overall" x="%.0f" y="%.1f" dominant-baseline="central" text-anchor="end">%s</text>`+"\n", overallRight, cy, formatRelative(e.mid))
		fmt.Fprintf(&body, `<line class="whisker" x1="%.1f" x2="%.1f" y1="%.1f" y2="%.1f"/>`+"\n", sc.x(e.lo), sc.x(e.hi), cy, cy)
		fmt.Fprintf(&body, `<circle class="dot" cx="%.1f" cy="%.1f" r="5"/>`+"\n", sc.x(e.mid), cy)
		for ci, c := range sum.cats {
			x := catX + float64(ci)*catWidth + 48
			ce, ok := s.category[c]
			if !ok {
				fmt.Fprintf(&body, `<text class="value" x="%.0f" y="%.1f" dominant-baseline="central" text-anchor="end">—</text>`+"\n", x, cy)
				continue
			}
			fmt.Fprintf(&body, `<text class="value" x="%.0f" y="%.1f" dominant-baseline="central" text-anchor="end">%s</text>`+"\n", x, cy, formatRelative(ce.mid))
			fmt.Fprintf(&body, `<text class="ord" x="%.0f" y="%.1f" dominant-baseline="central">%s</text>`+"\n", x+6, cy, esc(ordinal(s.catRank[c])))
		}
		fmt.Fprintln(&body, "</g>")
	}
	fmt.Fprintf(&body, `<line class="rule-strong" x1="%.0f" x2="%.0f" y1="%.0f" y2="%.0f"/>`+"\n", chartPad, right, bottom, bottom)
	y = bottom + 52

	// Figure 2: heatmap of scores per benchmark.
	fmt.Fprintf(&body, `<text class="fig" x="%.0f" y="%.0f">Where the time goes</text>`+"\n", chartPad, y)
	y += 22
	for _, line := range []string{
		"Each library's score on each benchmark: its median time relative to the fastest library's. The more a cell",
		"stands out from the background, the slower (key below, log scale). Under each benchmark is its fastest median.",
	} {
		fmt.Fprintf(&body, `<text class="sub" x="%.0f" y="%.0f">%s</text>`+"\n", chartPad, y, esc(line))
		y += 17
	}
	y += 18

	cellX := nameX + nameWidth + 20
	cellW := (right - cellX) / float64(len(benches))
	category := map[string]string{}
	for _, c := range categories {
		for _, b := range c.benches {
			category[b] = c.name
		}
	}
	maxScore := 1.0
	fastest := map[string]float64{}
	for _, b := range benches {
		for _, r := range res.rows(b) {
			if !r.isReference {
				fastest[b] = r.ns
				break
			}
		}
	}
	for _, s := range sum.impls {
		for _, e := range s.relative {
			maxScore = max(maxScore, e.mid)
		}
	}
	step := func(v float64) int {
		if maxScore <= 1 {
			return 0
		}
		return min(len(ramp)-1, max(0, int(math.Round(float64(len(ramp)-1)*math.Log(v)/math.Log(maxScore)))))
	}
	for i, b := range benches {
		cx := cellX + (float64(i)+0.5)*cellW
		if c := category[b]; c != "" {
			fmt.Fprintf(&body, `<text class="col" x="%.1f" y="%.0f" text-anchor="middle">%s</text>`+"\n", cx, y, esc(strings.ToUpper(c)))
		}
		l1, l2 := splitName(b)
		fmt.Fprintf(&body, `<text class="bench" x="%.1f" y="%.0f" text-anchor="middle">%s</text>`+"\n", cx, y+17, esc(l1))
		if l2 != "" {
			fmt.Fprintf(&body, `<text class="bench" x="%.1f" y="%.0f" text-anchor="middle">%s</text>`+"\n", cx, y+32, esc(l2))
		}
		fmt.Fprintf(&body, `<text class="fastest" x="%.1f" y="%.0f" text-anchor="middle">%s</text>`+"\n", cx, y+48, formatDuration(fastest[b]))
	}
	y += 60
	for i, s := range sum.impls {
		ry := y + float64(i)*cellHeight
		fmt.Fprintf(&body, `<text class="name" x="%.0f" y="%.1f" dominant-baseline="central">%s</text>`+"\n", nameX, ry+cellHeight/2, esc(s.impl))
		for j, b := range benches {
			x := cellX + float64(j)*cellW
			e, ok := s.relative[b]
			if !ok {
				fmt.Fprintf(&body, `<text class="tick" x="%.1f" y="%.1f" dominant-baseline="central" text-anchor="middle">—</text>`+"\n", x+cellW/2, ry+cellHeight/2)
				continue
			}
			fmt.Fprintf(&body, `<g class="h%d"><title>%s · %s: %s</title>`, step(e.mid), esc(s.impl), esc(b), formatRelative(e.mid))
			fmt.Fprintf(&body, `<rect x="%.1f" y="%.1f" width="%.1f" height="%.1f" rx="3"/>`, x+cellGap/2, ry+cellGap/2, cellW-cellGap, cellHeight-cellGap)
			fmt.Fprintf(&body, `<text class="cell" x="%.1f" y="%.1f" dominant-baseline="central" text-anchor="middle">%s</text></g>`+"\n", x+cellW/2, ry+cellHeight/2, formatRelative(e.mid))
		}
	}
	y += float64(len(sum.impls))*cellHeight + 22

	// Legend: the ramp's steps, from fastest to the slowest score shown.
	fmt.Fprintf(&body, `<text class="tick" x="%.1f" y="%.1f" dominant-baseline="central" text-anchor="end">1.00×</text>`+"\n", cellX-6, y)
	for i := range ramp {
		fmt.Fprintf(&body, `<g class="h%d"><rect x="%.1f" y="%.1f" width="16" height="10" rx="2"/></g>`+"\n", i, cellX+float64(i)*18, y-5)
	}
	fmt.Fprintf(&body, `<text class="tick" x="%.1f" y="%.1f" dominant-baseline="central">%s (slowest shown)</text>`+"\n", cellX+float64(len(ramp))*18+4, y, formatRelative(maxScore))
	y += 34

	var foot []string
	if commit != "" {
		foot = append(foot, fmt.Sprintf("commit %.7s", commit))
	}
	foot = append(foot, now.UTC().Format("2006-01-02"))
	if cpu := res.config["cpu"]; cpu != "" {
		foot = append(foot, cpu)
	}
	if p := res.config["procs"]; p != "" {
		foot = append(foot, p+" goroutines per benchmark")
	}
	foot = append(foot, "full tables in results.md")
	fmt.Fprintf(&body, `<line class="rule" x1="%.0f" x2="%.0f" y1="%.0f" y2="%.0f"/>`+"\n", chartPad, right, y-14, y-14)
	fmt.Fprintf(&body, `<text class="foot" x="%.0f" y="%.0f">%s</text>`+"\n", chartPad, y+2, esc(strings.Join(foot, " · ")))
	height := y + chartPad - 6

	var desc []string
	for i, s := range sum.impls {
		desc = append(desc, fmt.Sprintf("%s %s %s", sum.ranks[i], s.impl, formatRelative(s.overall.mid)))
	}
	fmt.Fprintf(w, `<svg xmlns="http://www.w3.org/2000/svg" width="%.0f" height="%.0f" viewBox="0 0 %.0f %.0f" role="img" aria-labelledby="t d">`+"\n", width, height, width, height)
	fmt.Fprintln(w, `<title id="t">Benchmark leaderboard and per-benchmark scores</title>`)
	fmt.Fprintf(w, "<desc id=\"d\">Overall ranking, fastest first: %s.</desc>\n", esc(strings.Join(desc, "; ")))
	fmt.Fprintf(w, "<style>%s</style>\n", chartStyle())
	fmt.Fprintf(w, `<rect class="bg" width="%.0f" height="%.0f" rx="8"/>`+"\n", width, height)
	io.WriteString(w, body.String())
	fmt.Fprintln(w, "</svg>")
}
