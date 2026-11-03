package main

import (
	"bytes"
	"math"
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
	if got, want := rs[1].relative, 3000.0/1100; math.Abs(got-want) > 1e-9 {
		t.Errorf("Slow relative = %v, want %v", got, want)
	}
	if rs[0].relative != 1 {
		t.Errorf("Fast relative = %v, want 1", rs[0].relative)
	}
}

func TestSummaryAveragesAcrossBenchmarks(t *testing.T) {
	var buf bytes.Buffer
	write(&buf, parsed(t), "0123456789abcdef", "16", time.Date(2026, 11, 3, 0, 0, 0, 0, time.UTC))
	out := buf.String()

	// Fast: 1.00× and 4.00× (geomean 2.00×); Slow: 2.73× and 1.00× (geomean 1.65×), so Slow ranks first.
	for _, want := range []string{
		"| 1 | Slow | **1.65×** | 2.73× | 1.00× |",
		"| 2 | Fast | **2.00×** | 1.00× | 4.00× |",
		"| ref | Memory | 10.0ns | — | — | 0 B | 0 |",
		"from commit `0123456`",
		"CPU: Test CPU",
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
