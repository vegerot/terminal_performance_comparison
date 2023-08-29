package main_test

import (
	"testing"

	main "github.com/vegerot/terminal_performance_comparison"
)

func BenchmarkOpenAllFilesNoDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main.OpenAllFilesNoDefer(100)
	}
}

func BenchmarkOpenAllFilesDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		main.OpenAllFilesDefer(100)
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
