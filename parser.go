package parsemake

import (
	"io"
	"os"
)

const (
	minBuff     = 512 - 1
	defaultBuff = 1 << 14
)

type Parser struct {
	data    []byte
	size, i int64
}

func NewParser() *Parser {
	return NewParseWithBuffSize(defaultBuff)
}

func NewParseWithBuffSize(buff int64) (p *Parser) {
	if buff < minBuff {
		buff = minBuff
	}
	p = new(Parser)
	p.size = buff
	return p
}

func (p *Parser) populateData(name string) (err error) {
	var (
		file  *os.File
		fInfo os.FileInfo
		n     int
	)
	if len(name) == 0 {
		file = os.Stdin
	} else {
		file, err = os.Open(name)
		if err != nil {
			return err
		}
	}
	defer file.Close()
	fInfo, err = file.Stat()
	if err != nil {
		return err
	}
	n = int(fInfo.Size())
	if n < minBuff {
		n = minBuff
	}
	n, err = p.readFile(file, n)
	if err != nil {
		return err
	}
	return nil
}

// func (p *Parser) Parse(fileName string) (m *Makefile, err error) {
func (p *Parser) Parse(fileName string) (m []byte, err error) {
	err = p.populateData(fileName)
	if err != nil {
		return nil, err
	}
	return p.data, nil
}

func (p *Parser) readFile(r io.Reader, min int) (n int, err error) {
	p.data = make([]byte, min)
	var nn int
	for err == nil {
		if n >= len(p.data) {
			d := append(p.data[:cap(p.data)], 0)
			p.data = d
		}
		nn, err = r.Read(p.data[n:])
		n += nn
	}
	if err == io.EOF {
		err = nil
	}
	p.data = p.data[:n]
	return
}
