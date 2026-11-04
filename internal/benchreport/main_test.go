package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"slices"
	"strings"
	"testing"
	"time"
)

const input = `goos: linux
cpu: Test CPU
BenchmarkFindAuthorByID/*dbeval.Fast-4    100   1000 ns/op   100 B/op   2 allocs/op
BenchmarkFindAuthorByID/*dbeval.Slow-4    100   3000 ns/op   300 B/op   6 allocs/op
BenchmarkFindAuthorByID/*dbeval.Memory-4  100     10 ns/op     0 B/op   0 allocs/op
BenchmarkFindAuthorByID/*dbeval.Fast-4    100   1200 ns/op   100 B/op   2 allocs/op
2026/11/03 12:00:00 stray log line
BenchmarkRecentArticles/*dbeval.Fast-4    100   8000000 ns/op
BenchmarkRecentArticles/*dbeval.Slow-4    100   2000000 ns/op
`

func parsed(t *testing.T) *results {
	t.Helper()
	// A multi-megabyte log line, like an ORM logging a large INSERT, must not stop parsing.
	in := strings.Replace(input, "stray log line", "INSERT "+strings.Repeat("x", 8<<20), 1)
	res := &results{config: map[string]string{}, data: map[string]map[string]*samples{}}
	if err := parse(strings.NewReader(in), res); err != nil {
		t.Fatal(err)
	}
	return res
}

func TestRowsSortedFastestFirstReferenceLast(t *testing.T) {
	rs := parsed(t).rows("FindAuthorByID")
	var got []string
	for _, r := range rs {
		got = append(got, r.impl)
	}
	if want := "Fast Slow Memory"; strings.Join(got, " ") != want {
		t.Fatalf("order = %v, want %s", got, want)
	}
	// Fast's median is 1100ns, so Slow is 3000/1100 relative to it; the reference doesn't count as fastest.
	if got, want := rs[1].relative.mid, 3000.0/1100; math.Abs(got-want) > 1e-9 {
		t.Errorf("Slow relative = %v, want %v", got, want)
	}
	// Fast's runs were 1000 and 1200ns.
	if got, want := rs[0].relative, (estimate{1, 1000.0 / 1100, 1200.0 / 1100}); got != want {
		t.Errorf("Fast relative = %+v, want %+v", got, want)
	}
}

func TestSummaryAveragesAcrossBenchmarks(t *testing.T) {
	var buf bytes.Buffer
	write(&buf, parsed(t), "0123456789abcdef", "16", "", time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC))
	out := buf.String()

	// Fast: 1.00× and 4.00× (geomean 2.00×); Slow: 2.73× and 1.00× (geomean 1.65×), so Slow ranks first. Their
	// ranges (Slow ran once per benchmark; Fast's average spans 1.91×–2.09×) don't overlap, so there's no tie.
	for _, want := range []string{
		"| 1 | Slow | **1.65×** | 2.73× | 1.00× |",
		"| 2 | Fast | **2.00×** | 1.00× | 4.00× |",
		"| ref | Memory | 10.0ns | — | — | 0 B | 0 |",
		"from commit `0123456`",
		"CPU: Test CPU",
		"from 4 goroutines at once",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "| Memory | **") {
		t.Errorf("reference implementation ranked in summary:\n%s", out)
	}
}

func TestFormatDuration(t *testing.T) {
	for ns, want := range map[float64]string{
		38.1:    "38.1ns",
		8100:    "8.10µs",
		35000:   "35.0µs",
		275000:  "275µs",
		1700000: "1.70ms",
		2.5e9:   "2.50s",
	} {
		if got := formatDuration(ns); got != want {
			t.Errorf("formatDuration(%v) = %q, want %q", ns, got, want)
		}
	}
}

func TestTieRanks(t *testing.T) {
	es := []estimate{
		{1.00, 0.95, 1.05}, // group leader
		{1.03, 1.01, 1.06}, // overlaps the leader: tied
		{1.08, 1.06, 1.10}, // overlaps the second but not the leader: starts its own group, tied with the next
		{1.09, 1.07, 1.12},
		{1.50, 1.45, 1.55},
	}
	got := strings.Join(tieRanks(es), " ")
	if want := "=1 =1 =3 =3 5"; got != want {
		t.Errorf("tieRanks = %s, want %s", got, want)
	}
}

func TestOrdinal(t *testing.T) {
	for in, want := range map[string]string{"1": "1st", "=2": "=2nd", "3": "3rd", "4": "4th", "11": "11th", "12": "12th", "13": "13th", "21": "21st"} {
		if got := ordinal(in); got != want {
			t.Errorf("ordinal(%s) = %s, want %s", in, got, want)
		}
	}
}

func TestCategories(t *testing.T) {
	// A is fastest at every read, B at every write.
	var in strings.Builder
	for _, b := range benchmarkOrder {
		a, bt := 1000, 2000
		if b == "InsertAuthors" || b == "InsertArticles" {
			a, bt = 2000, 1000
		}
		fmt.Fprintf(&in, "Benchmark%s/*dbeval.A-4 100 %d ns/op\nBenchmark%s/*dbeval.B-4 100 %d ns/op\n", b, a, b, bt)
	}
	res := &results{config: map[string]string{}, data: map[string]map[string]*samples{}}
	if err := parse(strings.NewReader(in.String()), res); err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	write(&buf, res, "", "", "", time.Now())
	out := buf.String()
	// A: reads 1.00×, writes 2.00×, overall 2^(2/5) = 1.32×; B: reads 2.00×, writes 1.00×, overall 2^(3/5) = 1.52×.
	for _, want := range []string{
		"| # | Implementation | **Overall** | Reads | Writes |",
		"| 1 | A | **1.32×** | 1.00× (1st) | 2.00× (2nd) |",
		"| 2 | B | **1.52×** | 2.00× (2nd) | 1.00× (1st) |",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}

func TestSVG(t *testing.T) {
	res := parsed(t)
	res.data["FindAuthorByID"]["A<B&C"] = &samples{ns: []float64{5000}}
	var buf bytes.Buffer
	writeSVG(&buf, res, summarize(res, res.benchmarks()), "0123456789abcdef", time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC))
	out := buf.String()

	// Well-formed XML, so names from the benchmark output can't break the markup.
	d := xml.NewDecoder(strings.NewReader(out))
	var names []string
	inName := false
	for {
		tok, err := d.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("invalid SVG: %v\n%s", err, out)
		}
		switch tok := tok.(type) {
		case xml.StartElement:
			inName = tok.Name.Local == "text" && slices.Contains(tok.Attr, xml.Attr{Name: xml.Name{Local: "class"}, Value: "name"})
		case xml.CharData:
			if inName {
				names = append(names, string(tok))
			}
		case xml.EndElement:
			inName = false
		}
	}
	// Rows in Overall order (see TestSummaryAveragesAcrossBenchmarks), without the reference implementation, in the
	// leaderboard and then the heatmap.
	if got, want := strings.Join(names, " "), "Slow Fast A<B&C Slow Fast A<B&C"; got != want {
		t.Errorf("rows = %q, want %q", got, want)
	}
	if !strings.Contains(out, "commit 0123456") {
		t.Errorf("footer missing commit:\n%s", out)
	}
}

func TestSplitName(t *testing.T) {
	for in, want := range map[string]string{
		"InsertAuthors":    "Insert/Authors",
		"FindAuthorByName": "FindAuthor/ByName",
		"FindAuthorByID":   "FindAuthor/ByID",
		"Short":            "Short/",
	} {
		if a, b := splitName(in); a+"/"+b != want {
			t.Errorf("splitName(%s) = %s/%s, want %s", in, a, b, want)
		}
	}
}

func TestReportShowsChart(t *testing.T) {
	var buf bytes.Buffer
	write(&buf, parsed(t), "", "", "results.svg", time.Now())
	if !strings.Contains(buf.String(), "](results.svg)") {
		t.Errorf("report doesn't show the chart:\n%s", buf.String())
	}
}
