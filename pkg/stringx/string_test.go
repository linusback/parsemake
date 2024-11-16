package stringx

import (
	"fmt"
	"strings"
	"testing"
)

var (
	testStrings = []string{"", "Hello", "my", "name", "is"}
	nilVal      StringNoAlloc
)

func Test_StringNoAlloc_String(t *testing.T) {
	for _, expected := range testStrings {
		testString(t, expected, StringToNoAlloc(expected))
	}
	testString(t, "", nilVal)
}

func Test_StringNoAlloc_GoString(t *testing.T) {
	for _, expected := range testStrings {
		testGoString(t, expected, StringToNoAlloc(expected))
	}
	testString(t, "", nilVal)
}

func Test_ArrStringNoAlloc_String(t *testing.T) {
	test := ArrStringToNoAlloc(testStrings)
	testString(t, testStrings, test)
	test = append(test, nilVal, nilVal, StringNoAlloc("hello"))
	result := append(testStrings, "", "", "hello")
	testString(t, result, test)
}

func Test_ArrStringNoAlloc_GoString(t *testing.T) {
	test := ArrStringToNoAlloc(testStrings)
	testGoString(t, testStrings, test)
	test = append(test, nilVal, nilVal, StringNoAlloc("hello"))
	result := append(testStrings, "", "", "hello")
	testGoString(t, result, test)
}

func testString(t *testing.T, expected any, test fmt.Stringer) {
	if fmt.Sprintf("%v", expected) != test.String() {
		t.Errorf("err String() expected: %v, got: %v", expected, test)
	}
}

func testGoString(t *testing.T, expected any, test fmt.GoStringer) {
	expString := fmt.Sprintf("%#v", expected)
	expTest := test.GoString()
	if strings.HasPrefix(expString, "[]string") {
		expString = strings.TrimPrefix(expString, "[]string")
		expTest = strings.TrimPrefix(expTest, "[]StringNoAlloc")
	}
	if expString != expTest {
		t.Errorf("err GoString() expected: %s, got: %s", expString, expTest)
	}
}
