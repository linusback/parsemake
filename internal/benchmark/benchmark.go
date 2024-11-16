package benchmark

import (
	"fmt"
	"os"
	"path"
	"testing"
)

const (
	checkMakefile   = "./test/CheckmakeMakefile"
	projectMakefile = "./Makefile"
)

var (
	CheckMakefile   = getRelativeFile(checkMakefile)
	ProjectMakefile = getRelativeFile(projectMakefile)
)

var TestFiles = [...]string{
	CheckMakefile,
	ProjectMakefile,
}

var Benchmarks = []Benchmark{
	NewBenchmark("standard", ProjectMakefile),
	NewBenchmark("large", CheckMakefile),
}

type Benchmark struct {
	Name     string
	Filename string
	ByteSize int64
	F        func(*testing.B)
}

func getRelativeFile(filename string) string {
	return path.Join(getBinPath(), filename)
}

func getBinPath() string {
	e, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dirPath := path.Dir(e)
	return dirPath
}

func NewBenchmark(name, filename string) Benchmark {
	return Benchmark{
		Name:     name,
		Filename: filename,
		ByteSize: GetByteSize(filename)}
}

func GetByteSize(fileName string) int64 {
	var (
		file     *os.File
		fileInfo os.FileInfo
		err      error
	)
	fmt.Println(fileName)
	file, err = os.Open(fileName)
	if err != nil {
		os.di
		panic(err)
	}
	defer file.Close()
	fileInfo, err = file.Stat()
	if err != nil {
		panic(err)
	}
	return fileInfo.Size()

}
