package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

type outlineHelper FileOutline

func (oh outlineHelper) firstEntryByName(name string) *Entry {
	for _, e := range oh.Entries {
		if e.Name == name {
			return e
		}
	}
	return nil
}

type testType struct {
	value       int
	structValue int
}

func (t *testType) foo() {
	fmt.Println("foo of test type in test package")
}

func TestParseFile(t *testing.T) {
	parseFile("subjects/test0.go")
}

func TestParseFile2(t *testing.T) {
	parseFile("subjects/test1.go")
}

func TestParseInvalidFile(t *testing.T) {
	parsed, err := parseFile("nonexistent-file")

	assertEquals(t, parsed, "")
	assertNotNil(t, err)
}

// parses a file with multiple symbols defined with the same name (different receivers)
// and makes sure all appear in the parsed output.
func TestDoubleSymbol(t *testing.T) {
	*verbose = true
	parsed, err := parseFile("subjects/doublesymbol.go")
	var outline FileOutline
	err = json.Unmarshal([]byte(parsed), &outline)
	assertErrNil(t, err)
	assertEquals(t, len(outline.Entries), 5)
}

func TestSliceType(t *testing.T) {
	parsed, err := parseFile("subjects/arraysslices.go")
	assertErrNil(t, err)

	var outline FileOutline
	err = json.Unmarshal([]byte(parsed), &outline)
	foo := outlineHelper(outline).firstEntryByName("foo")
	assertNotNil(t, foo)

	assertEquals(t, foo.Receiver, "sliceType")
	assertEquals(t, len(outline.Entries), 3)

	st := outlineHelper(outline).firstEntryByName("sliceType")
	assertNotNil(t, st)
	assertEquals(t, st.Receiver, "")
}

func TestInterfaceType(t *testing.T) {
	parsed, err := parseFile("subjects/interfacetest.go")
	assertErrNil(t, err)
	var outline FileOutline
	err = json.Unmarshal([]byte(parsed), &outline)
	assertErrNil(t, err)

	iface := outlineHelper(outline).firstEntryByName("iface")
	assertNotNil(t, iface)
	assertEquals(t, iface.Receiver, "")

	ifacechild := outlineHelper(outline).firstEntryByName("foo")
	assertNotNil(t, ifacechild)
	assertEquals(t, ifacechild.Receiver, "iface")
}
func assertErrNil(t *testing.T, err error) {
	if err != nil {
		t.Errorf("Err was not nil, but %v", err)
	}
}

func assertNotNil(t *testing.T, v interface{}) {
	if v == nil {
		t.Error("Value was nil, but shouldnt be.")
	}
}
func assertEquals(t *testing.T, a, b interface{}) {
	if a != b {
		t.Errorf("%v does not equal %v", a, b)
	}
}

type SomeInterface interface {
	InterfaceFuncA(int)
	InterfaceFuncB(int64)
}
