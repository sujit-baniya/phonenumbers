package main

import (
	"github.com/sujit-baniya/phonenumbers"
	"testing"
)

func BenchmarkCalculate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		phonenumbers.Parse("+9779856034616", "NP")
	}
}
