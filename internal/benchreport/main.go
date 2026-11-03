// Command benchreport turns this repository's `go test -bench` output into a Markdown report: a summary ranking
// each implementation by its time relative to the fastest implementation, averaged across all benchmarks, followed
// by a table per benchmark sorted fastest first.
//
// Usage:
//
//	go test -run '^$' -bench . -benchmem -count 6 . | tee bench.txt
//	go run ./internal/benchreport bench.txt > results.md
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"regexp"
	"runtime"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
)

// reference implementations are shown for scale but excluded from "fastest" and from the ranking.
var reference = map[string]bool{"Memory": true}

// benchmarkOrder and descriptions document the workloads in dbeval_test.go.
var benchmarkOrder = []string{"InsertAuthors", "InsertArticles", "FindAuthorByID", "FindAuthorByName", "RecentArticles"}

var descriptions = map[string]string{
	"InsertAuthors":    "insert a batch of 10 authors",
	"InsertArticles":   "insert a batch of 10 articles with multi-paragraph bodies",
	"FindAuthorByID":   "look up one author by primary key, out of 10,000",
	"FindAuthorByName": "find the ~1,000 authors sharing a name, out of 10,000",
	"RecentArticles":   "read the 2,000 most recent articles, out of 10,000",
}

var (
	resultRE = regexp.MustCompile(`^Benchmark([^/\s]+)/(\S+?)(?:-\d+)?\s+\d+\s+([\d.]+) ns/op(?:\s+([\d.]+) B/op)?(?:\s+([\d.]+) allocs/op)?`)
	configRE = regexp.MustCompile(`^(goos|goarch|cpu|pkg):\s*(.*?)\s*$`)
)

type samples struct{ ns, bytes, allocs []float64 }

type results struct {
	config map[string]string
	// bench -> implementation -> samples
	data map[string]map[string]*samples
}

func parse(r io.Reader, res *results) error {
	// bufio.Reader rather than bufio.Scanner: benchmark output can contain arbitrarily long log lines (e.g. an ORM
	// logging a 10,000-row INSERT), which would exceed a Scanner's token limit.
	br := bufio.NewReader(r)
	for {
		line, err := br.ReadString('\n')
		if line != "" {
			res.addLine(strings.TrimRight(line, "\r\n"))
		}
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
	}
}

func (res *results) addLine(line string) {
	if m := configRE.FindStringSubmatch(line); m != nil {
		res.config[m[1]] = m[2]
		return
	}
	m := resultRE.FindStringSubmatch(line)
	if m == nil {
		return
	}
	bench := m[1]
	impl := strings.TrimPrefix(strings.TrimPrefix(m[2], "*"), "dbeval.")
	if res.data[bench] == nil {
		res.data[bench] = map[string]*samples{}
	}
	s := res.data[bench][impl]
	if s == nil {
		s = &samples{}
		res.data[bench][impl] = s
	}
	s.ns = append(s.ns, mustFloat(m[3]))
	if m[4] != "" {
		s.bytes = append(s.bytes, mustFloat(m[4]))
	}
	if m[5] != "" {
		s.allocs = append(s.allocs, mustFloat(m[5]))
	}
}

func mustFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func median(xs []float64) float64 {
	if len(xs) == 0 {
		return math.NaN()
	}
	s := slices.Sorted(slices.Values(xs))
	if n := len(s); n%2 == 1 {
		return s[n/2]
	} else {
		return (s[n/2-1] + s[n/2]) / 2
	}
}

// spread is the largest distance from the median to the slowest or fastest run, as a percentage of the median.
func spread(xs []float64) float64 {
	if len(xs) < 2 {
		return math.NaN()
	}
	m := median(xs)
	return 100 * max(slices.Max(xs)-m, m-slices.Min(xs)) / m
}

func geomean(xs []float64) float64 {
	sum := 0.0
	for _, x := range xs {
		sum += math.Log(x)
	}
	return math.Exp(sum / float64(len(xs)))
}

func (res *results) benchmarks() []string {
	var known, other []string
	for _, b := range benchmarkOrder {
		if res.data[b] != nil {
			known = append(known, b)
		}
	}
	for b := range res.data {
		if !slices.Contains(benchmarkOrder, b) {
			other = append(other, b)
		}
	}
	sort.Strings(other)
	return append(known, other...)
}

// row is one implementation's medians in one benchmark.
type row struct {
	impl                string
	ns, bytes, allocs   float64
	spread, relative    float64
	isReference, hasMem bool
}

// rows returns bench's implementations sorted fastest first (reference implementations last), with times relative
// to the fastest non-reference implementation.
func (res *results) rows(bench string) []row {
	var rs []row
	best := math.Inf(1)
	for impl, s := range res.data[bench] {
		r := row{
			impl:        impl,
			ns:          median(s.ns),
			spread:      spread(s.ns),
			isReference: reference[impl],
			hasMem:      len(s.allocs) > 0,
		}
		if r.hasMem {
			r.bytes, r.allocs = median(s.bytes), median(s.allocs)
		}
		if !r.isReference {
			best = min(best, r.ns)
		}
		rs = append(rs, r)
	}
	for i := range rs {
		rs[i].relative = rs[i].ns / best
	}
	sort.Slice(rs, func(i, j int) bool {
		if rs[i].isReference != rs[j].isReference {
			return !rs[i].isReference
		}
		return rs[i].ns < rs[j].ns
	})
	return rs
}

// formatDuration formats ns with three significant figures, e.g. 8.10µs, 35.0µs, 275µs.
func formatDuration(ns float64) string {
	units := []struct {
		name  string
		scale float64
	}{{"ns", 1}, {"µs", 1e3}, {"ms", 1e6}, {"s", 1e9}}
	u := units[0]
	for _, c := range units[1:] {
		if ns >= c.scale {
			u = c
		}
	}
	v := ns / u.scale
	switch {
	case v >= 100:
		return fmt.Sprintf("%.0f%s", v, u.name)
	case v >= 10:
		return fmt.Sprintf("%.1f%s", v, u.name)
	default:
		return fmt.Sprintf("%.2f%s", v, u.name)
	}
}

func formatBytes(b float64) string {
	switch {
	case b < 1<<10:
		return fmt.Sprintf("%.0f B", b)
	case b < 1<<20:
		return fmt.Sprintf("%.1f KiB", b/(1<<10))
	default:
		return fmt.Sprintf("%.1f MiB", b/(1<<20))
	}
}

func formatRelative(x float64) string {
	if x < 10 {
		return fmt.Sprintf("%.2f×", x)
	}
	return fmt.Sprintf("%.1f×", x)
}

func write(w io.Writer, res *results, commit, pg string, now time.Time) {
	benches := res.benchmarks()
	runs := 0
	for _, b := range benches {
		for _, s := range res.data[b] {
			runs = max(runs, len(s.ns))
		}
	}

	fmt.Fprintln(w, "# Benchmark results")
	fmt.Fprintln(w)
	generated := "Generated by [`.github/workflows/bench.yml`](.github/workflows/bench.yml)"
	if commit != "" {
		generated += fmt.Sprintf(" from commit `%.7s`", commit)
	}
	fmt.Fprintf(w, "%s on %s. Raw output is in the workflow run's log.\n\n", generated, now.UTC().Format("2006-01-02 15:04 MST"))
	env := []string{}
	if cpu := res.config["cpu"]; cpu != "" {
		env = append(env, "CPU: "+cpu)
	}
	env = append(env, "Go: "+runtime.Version())
	if pg != "" {
		env = append(env, "PostgreSQL: "+pg)
	}
	env = append(env, fmt.Sprintf("%d runs per benchmark", runs))
	fmt.Fprintf(w, "%s.\n\n", strings.Join(env, " · "))
	fmt.Fprintln(w, "Times are the median of the runs. **×** is time relative to the fastest implementation for that")
	fmt.Fprintln(w, "benchmark (1.00× is fastest; 2.00× takes twice as long). Memory is an in-process reference, not a")
	fmt.Fprintln(w, "database, so it's listed for scale but excluded from the rankings. Runs happen on a shared GitHub runner:")
	fmt.Fprintln(w, "compare implementations within this report rather than absolute numbers across reports.")
	fmt.Fprintln(w)

	// Summary: each implementation's relative time per benchmark, averaged (geometric mean, the right average for
	// ratios) across benchmarks.
	type summary struct {
		impl     string
		relative map[string]float64
		mean     float64
	}
	byImpl := map[string]*summary{}
	for _, b := range benches {
		for _, r := range res.rows(b) {
			if r.isReference {
				continue
			}
			s := byImpl[r.impl]
			if s == nil {
				s = &summary{impl: r.impl, relative: map[string]float64{}}
				byImpl[r.impl] = s
			}
			s.relative[b] = r.relative
		}
	}
	var sums []*summary
	for _, s := range byImpl {
		var xs []float64
		for _, x := range s.relative {
			xs = append(xs, x)
		}
		s.mean = geomean(xs)
		sums = append(sums, s)
	}
	sort.Slice(sums, func(i, j int) bool { return sums[i].mean < sums[j].mean })

	fmt.Fprintln(w, "## Summary")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Implementations ranked by their average relative time across all benchmarks (geometric mean of the")
	fmt.Fprintln(w, "per-benchmark **×** values, so each benchmark counts equally regardless of how long it takes).")
	fmt.Fprintln(w)
	fmt.Fprint(w, "| # | Implementation | **Average** |")
	for _, b := range benches {
		fmt.Fprintf(w, " %s |", b)
	}
	fmt.Fprint(w, "\n|--:|---|--:|")
	for range benches {
		fmt.Fprint(w, "--:|")
	}
	fmt.Fprintln(w)
	for i, s := range sums {
		fmt.Fprintf(w, "| %d | %s | **%s** |", i+1, s.impl, formatRelative(s.mean))
		for _, b := range benches {
			if x, ok := s.relative[b]; ok {
				fmt.Fprintf(w, " %s |", formatRelative(x))
			} else {
				fmt.Fprint(w, " — |")
			}
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w)

	for _, b := range benches {
		fmt.Fprintf(w, "## %s\n\n", b)
		if d := descriptions[b]; d != "" {
			fmt.Fprintf(w, "Each operation: %s.\n\n", d)
		}
		fmt.Fprintln(w, "| # | Implementation | Time/op | ± | × | Memory/op | Allocs/op |")
		fmt.Fprintln(w, "|--:|---|--:|--:|--:|--:|--:|")
		rank := 0
		for _, r := range res.rows(b) {
			pos := "ref"
			rel := "—"
			if !r.isReference {
				rank++
				pos = strconv.Itoa(rank)
				rel = formatRelative(r.relative)
			}
			sp := "—"
			if !math.IsNaN(r.spread) {
				sp = fmt.Sprintf("%.0f%%", r.spread)
			}
			mem, allocs := "—", "—"
			if r.hasMem {
				mem, allocs = formatBytes(r.bytes), strconv.FormatFloat(r.allocs, 'f', 0, 64)
			}
			fmt.Fprintf(w, "| %s | %s | %s | %s | %s | %s | %s |\n", pos, r.impl, formatDuration(r.ns), sp, rel, mem, allocs)
		}
		fmt.Fprintln(w)
	}
	fmt.Fprintln(w, "**±** is the largest distance from the median to the fastest or slowest run.")
}

func main() {
	commit := flag.String("commit", "", "commit the benchmarks ran on")
	pg := flag.String("pg", "", "PostgreSQL server version")
	flag.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), "usage: benchreport [-commit sha] [-pg version] [bench.txt ...]")
		flag.PrintDefaults()
	}
	flag.Parse()

	res := &results{config: map[string]string{}, data: map[string]map[string]*samples{}}
	if flag.NArg() == 0 {
		if err := parse(os.Stdin, res); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	for _, name := range flag.Args() {
		f, err := os.Open(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		err = parse(f, res)
		f.Close()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	if len(res.data) == 0 {
		fmt.Fprintln(os.Stderr, "benchreport: no benchmark results found")
		os.Exit(1)
	}
	write(os.Stdout, res, *commit, *pg, time.Now())
}
