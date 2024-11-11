package parsemake

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

func Benchmark_Read_Scanner(b *testing.B) {
	toRun := createParserBenchmarks(scannerBenchmark)
	b.ResetTimer()
	for _, bench := range toRun {
		b.Run(bench.name, bench.f)
	}
}

func Benchmark_Read_Simple(b *testing.B) {
	toRun := createParserBenchmarks(simpleBenchmark)
	b.ResetTimer()
	for _, bench := range toRun {
		b.Run(bench.name, bench.f)
	}
}

func Benchmark_Read_Custom(b *testing.B) {
	toRun := createParserBenchmarks(customBenchmark)
	b.ResetTimer()
	for _, bench := range toRun {
		b.Run(bench.name, bench.f)
	}
}

func Benchmark_Read_Include_Scanner(b *testing.B) {
	toRun := createParserBenchmarks(scannerBenchmarkInclude)
	b.ResetTimer()
	for _, bench := range toRun {
		b.Run(bench.name, bench.f)
	}
}

func Benchmark_Read_Include_Simple(b *testing.B) {
	toRun := createParserBenchmarks(simpleBenchmarkInclude)
	b.ResetTimer()
	for _, bench := range toRun {
		b.Run(bench.name, bench.f)
	}
}

func Benchmark_Read_Include_Custom(b *testing.B) {
	toRun := createParserBenchmarks(customBenchmarkInclude)
	b.ResetTimer()
	for _, bench := range toRun {
		b.Run(bench.name, bench.f)
	}
}

func createParserBenchmarks(getBench func(string, int64) func(*testing.B)) (b []benchmark) {
	b = make([]benchmark, len(benchmarks))
	for i, bench := range benchmarks {
		b[i] = bench // copies values
		b[i].f = getBench(bench.fileName, bench.byteSize)
	}
	return b
}

func scannerBenchmark(fileName string, byteSize int64) func(*testing.B) {
	var (
		r       = new(bytes.Reader)
		scanner = new(bufio.Scanner)
		err2    error
		buff    []byte
	)
	buff, err2 = readFile(fileName)
	if err2 != nil {
		panic(err2)
	}
	r.Reset(buff)

	return func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(byteSize)
		for i := 0; i < b.N; i++ {
			r.Reset(buff)
			scanner = bufio.NewScanner(r)
			readWithScanner(scanner)
		}
	}
}

func simpleBenchmark(fileName string, byteSize int64) func(*testing.B) {
	var (
		err2 error
		buff []byte
	)
	buff, err2 = readFile(fileName)
	if err2 != nil {
		panic(err2)
	}

	return func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(byteSize)
		for i := 0; i < b.N; i++ {
			readSimple(buff)
		}
	}
}

func scannerBenchmarkInclude(fileName string, byteSize int64) func(*testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(byteSize)
		for i := 0; i < b.N; i++ {
			file, err := os.Open(fileName)
			if err != nil {
				b.Error(err)
			}
			scanner := bufio.NewScanner(file)
			readWithScanner(scanner)
			file.Close()
		}
	}
}

func simpleBenchmarkInclude(fileName string, byteSize int64) func(*testing.B) {
	return func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(byteSize)
		for i := 0; i < b.N; i++ {
			buff, err := os.ReadFile(fileName)
			if err != nil {
				b.Error(err)
			}
			readSimple(buff)
		}
	}
}

func customBenchmark(fileName string, byteSize int64) func(*testing.B) {
	var fileInfo os.FileInfo

	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	fileInfo, err = file.Stat()
	if err != nil {
		panic(err)
	}
	buff := make([]byte, fileInfo.Size())
	_, err = io.ReadFull(file, buff)
	if err != nil {
		file.Close()
		panic(err)
	}
	file.Close()
	return func(b *testing.B) {
		b.ReportAllocs()
		b.SetBytes(byteSize)
		for i := 0; i < b.N; i++ {
			readSimple(buff)
		}
	}
}

func customBenchmarkInclude(fileName string, byteSize int64) func(*testing.B) {
	return func(b *testing.B) {
		var (
			n        int
			err      error
			file     *os.File
			buff     []byte
			fileInfo os.FileInfo
		)
		b.ReportAllocs()
		b.SetBytes(byteSize)
		for i := 0; i < b.N; i++ {
			file, err = os.Open(fileName)
			if err != nil {
				b.Error(err)
			}
			fileInfo, err = file.Stat()
			if err != nil {
				b.Error(err)
			}
			buff = make([]byte, fileInfo.Size())
			n, err = io.ReadFull(file, buff)
			if err != nil {
				file.Close()
				b.Error(err)
			}
			if n != len(buff) {
				b.Errorf("values of n and len(buff) are different %d != %d", n, buff)
			}
			file.Close()
			readSimple(buff)
		}
	}
}

func Test_Parser_readFile(t *testing.T) {
	p := NewParser()
	r := strings.NewReader("hellohellohello")
	n, err := p.readFile(r, len("hellohellohello"))
	if err != nil {
		t.Error(err)
	}
	t.Logf("total bytes read %d, size of data %d\n", n, len(p.data))
	t.Logf("val: %s", p.data)

	p = NewParser()
	r = strings.NewReader("hellohellohello")
	n, err = p.readFile(r, len("hellohellohello")*30)
	if err != nil {
		t.Error(err)
	}
	t.Logf("total bytes read %d, size of data %d\n", n, len(p.data))
	t.Logf("val: %s", p.data)
}

func Test_Parser_New(t *testing.T) {
	var (
		b        []byte
		err      error
		rows     int
		n        int
		fileInfo os.FileInfo
	)
	b, err = readFile(makefile)
	if err != nil {
		t.Error(err)
	}
	rows = readWithScanner(bufio.NewScanner(bytes.NewReader(b)))
	t.Logf("rows from scanner %d\n", rows)
	rows = readSimple(b)
	t.Logf("rows from simple %d\n", rows)
	file, err := os.Open(makefile)
	if err != nil {
		t.Error(err)
	}
	fileInfo, err = file.Stat()
	if err != nil {
		t.Error(err)
	}
	buff := make([]byte, 1024*16)
	n, err = io.ReadFull(file, buff)
	if err != io.ErrUnexpectedEOF && err != nil {
		file.Close()
		t.Error(err)
	}
	file.Close()
	buff = buff[:n]
	rows = readSimple(buff)
	t.Logf("rows read from cusom %d", rows)
	t.Logf("size of bytes %d, fileinfo says %d, error is %v", len(buff), fileInfo.Size(), err)

	b, err = readFile(checkMakefile)
	if err != nil {
		t.Error(err)
	}
	rows = readWithScanner(bufio.NewScanner(bytes.NewReader(b)))
	t.Logf("rows from scanner %d\n", rows)
	rows = readSimple(b)
	t.Logf("rows from simple %d\n", rows)
	file, err = os.Open(checkMakefile)
	if err != nil {
		t.Error(err)
	}
	fileInfo, err = file.Stat()
	if err != nil {
		t.Error(err)
	}
	buff = make([]byte, fileInfo.Size())
	n, err = io.ReadFull(file, buff)
	if err != nil {
		file.Close()
		t.Error(err)
	}
	file.Close()
	buff = buff[:n]

	rows = readSimple(buff)
	t.Logf("rows read from cusom %d", rows)
	t.Logf("size of bytes %d, fileinfo says %d, error is %v", len(buff), fileInfo.Size(), err)
}

func readFile(fileName string) (b []byte, err error) {
	var file *os.File
	file, err = os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	b, err = io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return b, nil

}

func readWithScanner(s *bufio.Scanner) (rows int) {
	for more := s.Scan(); more; more = s.Scan() {
		s.Bytes()
		rows++
	}
	return rows
}

func readSimple(buff []byte) (rows int) {
	const stop = '\n'
	var start int
	for i := 0; i < len(buff); i++ {
		if buff[i] == stop {
			rows++
			//fmt.Println(string(buff[start:i]))
			start = i
		}
	}
	if start < len(buff)-1 {
		//fmt.Println(string(buff[start:]))
		rows++
	}
	return rows
}
