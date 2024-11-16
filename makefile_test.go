package parsemake

import (
	"fmt"
	"github.com/linusback/parsemake/internal/benchmark"
	"github.com/mrtazz/checkmake/parser"
	"log/slog"
	"testing"
)

func Benchmark_Parse(b *testing.B) {
	toRun := getParseBenchmarks(Parse)
	for _, bench := range toRun {
		b.Run(bench.Name, bench.F)
	}
}

func Benchmark_Checkmake_Parse(b *testing.B) {
	toRun := getParseBenchmarks(parser.Parse)
	b.ResetTimer()
	for _, bench := range toRun {
		b.Run(bench.Name, bench.F)
	}
}

func getParseBenchmarks[T any](parse func(string) (T, error)) (b []benchmark.Benchmark) {
	b = make([]benchmark.Benchmark, len(benchmark.Benchmarks))
	for i, bench := range benchmark.Benchmarks {
		b[i] = bench // copies values
		b[i].F = getBenchmarks(bench.Filename, bench.ByteSize, parse)
	}
	return b
}

func getBenchmarks[T any](fileName string, byteSize int64, parse func(string) (T, error)) func(*testing.B) {
	fileName = fileName[4:]
	return func(b *testing.B) {
		var err error
		b.ReportAllocs()
		b.SetBytes(byteSize)
		for i := 0; i < b.N; i++ {
			_, err = parse(fileName)
			if err != nil {
				b.Error(err)
			}
		}
	}
}

// testing

// TODO write equality comparisons for makefile and underlying types
func Test_Checkmake_Parse(t *testing.T) {
	t.Logf(benchmark.ProjectMakefile[4:])
	//t.Run("makefile", testCheckmakeParse(benchmark.ProjectMakefile[4:]))
	//t.Run("checkmake", testCheckmakeParse(benchmark.CheckMakefile[4:]))
}

func testCheckmakeParse(filename string) func(t *testing.T) {
	return func(t *testing.T) {
		m, err := getCheckMake(filename)
		if err != nil {
			t.Error(err)
		}

		fmt.Println("")
		for _, v := range m.Variables {
			fmt.Printf("variable: %#v\n", v)
		}

		fmt.Println("rules: ")
		for _, r := range m.Rules {
			fmt.Printf("rules: %#v\n", r)
		}
	}
}

// TODO write equality comparisons for makefile and underlying types
func Test_Parse(t *testing.T) {
	m, err := getMake(benchmark.CheckMakefile)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("")
	for _, v := range m.Variables {
		fmt.Printf("variable: %#v\n", v)
	}

	fmt.Println("rules: ")
	for _, r := range m.Rules {
		fmt.Printf("rules: %#v\n", r)
	}
}

func getCheckMake(filename string) (f parser.Makefile, err error) {
	f, err = parser.Parse(filename)
	if err != nil {
		return parser.Makefile{}, err
	}
	return f, nil
}

func getMake(filename string) (f *Makefile, err error) {
	f, err = ParseLog(filename, slog.LevelInfo)
	if err != nil {
		return nil, err
	}
	return f, nil
}
