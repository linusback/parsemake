package main

import (
	"fmt"
	"github.com/mrtazz/checkmake/parser"
	"io"
	"os"
)

func main() {
	//err := readFile("./Makefile")
	err := parseCheckMake("./Makefile")

	if err != nil {
		fmt.Println("err: ", err)
		os.Exit(1)
	}
	os.Exit(0)
	return
	//filename := ""
	//if len(os.Args) > 1 {
	//	filename = os.Args[1]
	//}
	//p := parsemake.NewParser()
	//b, err := p.Parse(filename)
	//if err != nil {
	//	os.Stderr.WriteString(fmt.Sprintf("err: %v", err))
	//	os.Stderr.Close()
	//	os.Exit(1)
	//}
	//_, err = os.Stdout.Write(b)
	//if err != nil {
	//	os.Stderr.WriteString(fmt.Sprintf("err writing string to std out: %v", err))
	//	os.Stderr.Close()
	//	os.Exit(1)
	//}
	//os.Exit(0)
}

func parseCheckMake(filename string) error {
	m, err := parser.Parse(filename)
	if err != nil {
		return err
	}
	fmt.Println(len(m.Rules) + len(m.Variables))
	return nil
}

func readFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}
	buff := make([]byte, fileInfo.Size())
	n, err := io.ReadFull(file, buff)
	if err != nil {
		return err
	}
	if n != len(buff) {
		return fmt.Errorf("values of n and len(buff) are different %d != %d", n, buff)
	}
	rows := readSimple(buff)
	fmt.Println(rows)
	return nil
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
