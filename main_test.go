package main_test

import (
	"testing"
)

func BenchmarkOpenAllFilesNoDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OpenAllFilesNoDeferChunked()
	}
}

func BenchmarkOpenAllFilesDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		OpenAllFilesDeferChunked()
	}
}
