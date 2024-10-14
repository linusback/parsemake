package parsemake

import (
	"github.com/francoispqt/gojay"
)

const maxInt = int(^uint(0) >> 1)

var (
	spaceAsBytes   = []byte{' '}
	stringPrefix   = []byte{'['}
	stringSuffix   = []byte{']'}
	goStringPrefix = []byte(`[]StringNoAlloc{"`)
	goStringSep    = []byte(`", "`)
	goStringSuffix = []byte(`"}`)
)

type StringNoAlloc gojay.EmbeddedJSON

func (s StringNoAlloc) String() string {
	return string(s)
}

func (s StringNoAlloc) GoString() string {
	return "\"" + string(s) + "\""
}

type ArrStringNoAlloc []StringNoAlloc

func (s ArrStringNoAlloc) String() string {
	return string(join(s, spaceAsBytes, stringPrefix, stringSuffix))
}

func (s ArrStringNoAlloc) GoString() string {
	return string(join(s, goStringSep, goStringPrefix, goStringSuffix))
}

func StringToNoAlloc(s string) StringNoAlloc {
	return []byte(s)
}

func ArrStringToNoAlloc(sArr []string) ArrStringNoAlloc {
	res := make(ArrStringNoAlloc, len(sArr))
	for i, s := range sArr {
		res[i] = StringToNoAlloc(s)
	}
	return res
}

// Join concatenates the elements of s to create a new gojay.EmbeddedJSON var. The separator
// sep is placed between elements in the resulting slice code modified from bytes join.
func join(s ArrStringNoAlloc, sep, prefix, suffix []byte) gojay.EmbeddedJSON {
	if len(s) == 0 {
		return []byte{}
	}
	if len(s) == 1 {
		// Just return a copy.
		return append([]byte(nil), s[0]...)
	}

	var n int
	if len(sep) > 0 {
		if len(sep) >= maxInt/(len(s)-1) {
			panic("bytes: Join output length overflow")
		}
		n += len(sep) * (len(s) - 1)
	}

	if len(prefix) > maxInt-n {
		panic("bytes: Join output length overflow")
	}
	n += len(prefix)

	if len(suffix) > maxInt-n {
		panic("bytes: Join output length overflow")
	}
	n += len(suffix)

	for _, v := range s {
		if len(v) > maxInt-n {
			panic("bytes: Join output length overflow")
		}
		n += len(v)
	}

	var bp int
	b := make([]byte, n)
	bp += copy(b[bp:], prefix)
	bp += copy(b[bp:], s[0])
	for _, v := range s[1:] {
		bp += copy(b[bp:], sep)
		bp += copy(b[bp:], v)
	}
	bp += copy(b[bp:], suffix)
	return b
}
