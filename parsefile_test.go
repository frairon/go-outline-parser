package main

import "testing"

func TestParseFile(t *testing.T) {
	parseFile("test/test0.go")
}

func TestParseFile2(t *testing.T) {
	parseFile("test/test1.go")
}
