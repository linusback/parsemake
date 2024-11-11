package main

import (
	"fmt"
	"github.com/linusback/parsemake"
	"io"
	"os"
)

func main() {
	filename := ""
	if len(os.Args) > 1 {
		filename = os.Args[1]
	}
	p := parsemake.NewParser()
	b, err := p.Parse(filename)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("err: %v", err))
		os.Stderr.Close()
		os.Exit(1)
	}
	_, err = os.Stdout.Write(b)
	if err != nil {
		os.Stderr.WriteString(fmt.Sprintf("err writing string to std out: %v", err))
		os.Stderr.Close()
		os.Exit(1)
	}
	os.Exit(0)
}

func test() error {
	fmt.Println("hello world")
	for i, s := range os.Args {
		fmt.Printf("arg %d: %s\n", i, s)
	}

	file := os.Stdin
	fStat, err := file.Stat()
	if err != nil {
		return err
	}
	fmt.Println("fstat: ", fStat.Size())
	b, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error while reading all: %w", err)

	}
	fmt.Println(len(b))
	return nil
	//entries, err := os.ReadDir("./")
	//if err != nil {
	//	return err
	//}
	//
	////files, err := filepath.Glob("[Mm]akefile")
	//
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//for _, e := range files {
	//	fmt.Println(e)
	//}
	//b, ok := getMakeFile()
	//if ok {
	//	fmt.Println(string(b))
	//}
}
