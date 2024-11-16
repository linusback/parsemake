package parsemake

import (
	"fmt"
	"github.com/mrtazz/checkmake/parser"
	"log/slog"
	"os"
	"testing"
)

const (
	checkMakefile = "./test/CheckmakeMakefile"
	makefile      = "./Makefile"
)

var testFiles = [...]string{
	checkMakefile,
	makefile,
}

var benchmarks = []benchmark{
	newBenchmark("standard", makefile),
	newBenchmark("large", checkMakefile),
}

func Benchmark_Parse(b *testing.B) {
	toRun := getParseBenchmarks(Parse)
	for _, bench := range toRun {
		b.Run(bench.name, bench.f)
	}
}

func Benchmark_Checkmake_Parse(b *testing.B) {
	toRun := getParseBenchmarks(parser.Parse)
	b.ResetTimer()
	for _, bench := range toRun {
		b.Run(bench.name, bench.f)
	}
}

type benchmark struct {
	name     string
	fileName string
	byteSize int64
	f        func(*testing.B)
}

func newBenchmark(name, filename string) benchmark {
	return benchmark{name, filename, getByteSize(filename), nil}
}

func getParseBenchmarks[T any](parse func(string) (T, error)) (b []benchmark) {
	b = make([]benchmark, len(benchmarks))
	for i, bench := range benchmarks {
		b[i] = bench // copies values
		b[i].f = getBenchmarks(bench.fileName, bench.byteSize, parse)
	}
	return b
}

func getBenchmarks[T any](fileName string, byteSize int64, parse func(string) (T, error)) func(*testing.B) {
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
	t.Run("makefile", testCheckmakeParse(makefile))
	t.Run("checkmake", testCheckmakeParse(checkMakefile))
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
	m, err := getMake(checkMakefile)
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

func getByteSize(fileName string) int64 {
	var (
		file     *os.File
		fileInfo os.FileInfo
		err      error
	)
	file, err = os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileInfo, err = file.Stat()
	if err != nil {
		panic(err)
	}
	return fileInfo.Size()
}
