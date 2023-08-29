package main_test

import (
	"testing"

	main "github.com/vegerot/terminal_performance_comparison"
)

/// MUST run with `go test -bench=. -benchtime 1x` to see difference

func BenchmarkOpenAllFilesNoDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main.OpenAllFilesNoDefer()
	}
}

func BenchmarkOpenAllFilesDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main.OpenAllFilesDefer()
	}
}

func BenchmarkOpenAllFilesNoDeferChunked(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main.OpenAllFilesNoDeferChunked(100)
	}
}

func BenchmarkOpenAllFilesDeferChunked(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main.OpenAllFilesDeferChunked(100)
	}
}
