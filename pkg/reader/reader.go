package reader

import (
	"io"
	"os"
)

const (
	minBuff = 512
)

type Reader struct {
	Size   int
	n      int
	isDone bool
	Data   []byte
	file   *os.File
	err    error
}

func New(filename string) (p *Reader, err error) {
	p = new(Reader)
	if len(filename) == 0 {
		p.file = os.Stdin
	} else {
		p.file, err = os.Open(filename)
		if err != nil {
			return
		}
	}
	err = p.readInitial()
	return
}

func (r *Reader) Reset(filename string) (err error) {
	r.err = nil
	r.isDone = false
	r.n = 0
	if len(filename) == 0 {
		r.file = os.Stdin
	} else {
		r.file, err = os.Open(filename)
		if err != nil {
			return
		}
	}
	err = r.readInitial()
	return
}

func (r *Reader) readInitial() (err error) {
	var fInfo os.FileInfo
	fInfo, err = r.file.Stat()
	if err != nil {
		return
	}
	r.Size = int(fInfo.Size())
	r.Size++ // one byte for final read at EOF

	// If a file claims a small size, read at least 512 bytes.
	// In particular, files in Linux's /proc claim size 0 but
	// then do not work right if read in small pieces,
	// so an initial read of 1 byte would not work correctly.
	if r.Size < minBuff {
		r.Size = minBuff
	}
	if cap(r.Data) < r.Size {
		r.Data = make([]byte, r.Size)
	} else {
		r.Data = r.Data[:r.Size]
	}

	//_, err = r.ReadMore()
	//if err != nil {
	//	return err
	//}
	//_, err = r.ReadMore()
	return
}

func (r *Reader) ReadMore() (bool, error) {
	if r.isDone {
		return true, nil
	}
	var nn int

	if r.err == nil {
		nn, r.err = r.file.Read(r.Data[r.n:])
		r.n += nn
		r.Data = r.Data[:r.n]
	}
	if r.err == io.EOF {
		r.isDone = true
		return true, r.file.Close()
	}
	if r.err != nil {
		return true, r.err
	}
	if r.n >= len(r.Data) {
		d := append(r.Data[:cap(r.Data)], 0)
		r.Data = d
	}
	return false, nil
}

func (r *Reader) ReadAll() error {
	for !r.isDone {
		_, err := r.ReadMore()
		if err != nil {
			return err
		}
	}
	return nil
}
