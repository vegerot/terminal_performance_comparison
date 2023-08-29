package main_test

import (
	"fmt"
	"os"
	"sync"
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

func OpenAllFilesNoDeferChunked() {
	files, err := os.ReadDir("files/")
	if err != nil {
		panic(err)
	}

	fileBaskets := chunk(files, 5)

	var wg sync.WaitGroup
	for _, basket := range fileBaskets {
		wg.Add(1)
		go func(basket []os.DirEntry, wg *sync.WaitGroup) {
			openTheseFilesNoDefer(basket)
			wg.Done()
		}(basket, &wg)
	}
	wg.Wait()
}

func OpenAllFilesDeferChunked() {
	files, err := os.ReadDir("files/")
	if err != nil {
		panic(err)
	}

	fileBaskets := chunk(files, 5)

	var wg sync.WaitGroup
	for _, basket := range fileBaskets {
		wg.Add(1)
		go func(basket []os.DirEntry, wg *sync.WaitGroup) {
			defer wg.Done()
			openTheseFilesNoDefer(basket)
		}(basket, &wg)
	}
	wg.Wait()
}

func chunk[T interface{}](list []T, size int) [][]T {
	var chunks [][]T
	for size < len(list) {
		list, chunks = list[size:], append(chunks, list[0:size:size])
	}
	chunks = append(chunks, list)
	return chunks
}

func openTheseFilesNoDefer(basket []os.DirEntry) {
	for _, file := range basket {
		// open file
		f, err := os.Open(fmt.Sprintf("files/%s", file.Name()))
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func openTheseFilesDefer(basket []os.DirEntry) {
	for _, file := range basket {
		// open file
		f, err := os.Open(fmt.Sprintf("files/%s", file.Name()))
		defer f.Close()
		if err != nil {
			panic(err)
		}
	}
}
