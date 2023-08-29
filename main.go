package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	chunkSize := 100
	files := 100
	if len(os.Args) == 1 {
		OpenAllFilesNoDeferChunked(chunkSize)
		return
	}

	args := os.Args[1:]

	switch arg := args[0]; arg {
	case "create":
		createFiles(files)
	case "nodefer":
		OpenAllFilesNoDeferChunked(chunkSize)
	case "defer":
		OpenAllFilesDeferChunked(chunkSize)
	default:
		log.Fatalf("unknown arg: %s", arg)
	}
	fmt.Println(time.Since(start))
}

func createFiles(n int) {
	for i := 0; i < n; i++ {
		// create file
		f, err := os.Create(fmt.Sprintf("files/file%d.txt", i))
		if err != nil {
			panic(err)
		}

		// write to file
		longString := strings.Repeat("a", 2<<16)
		_, err = f.Write([]byte(longString))
		if err != nil {
			panic(err)
		}

		// close file
		err = f.Close()
		if err != nil {
			panic(err)
		}
	}
}
func OpenAllFilesNoDeferChunked(chunkSize int) {
	cwd, _ := os.Getwd()
	files, err := os.ReadDir(filepath.Join(cwd, "files"))
	if err != nil {
		panic(err)
	}

	fileBaskets := chunk(files, chunkSize)

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

func OpenAllFilesDeferChunked(chunkSize int) {
	cwd, _ := os.Getwd()
	files, err := os.ReadDir(filepath.Join(cwd, "files"))
	if err != nil {
		panic(err)
	}

	fileBaskets := chunk(files, chunkSize)

	var wg sync.WaitGroup
	for _, basket := range fileBaskets {
		wg.Add(1)
		go func(basket []os.DirEntry, wg *sync.WaitGroup) {
			defer wg.Done()
			openTheseFilesDefer(basket)
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

func openTheseFilesNoDefer(files []os.DirEntry) {
	cwd, _ := os.Getwd()
	for _, file := range files {
		f, err := os.Open(filepath.Join(cwd, "files", file.Name()))
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 1024)
		_, err = f.Read(buf)
		addOneToEverything(buf)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func openTheseFilesDefer(files []os.DirEntry) {
	cwd, _ := os.Getwd()
	for _, file := range files {
		f, err := os.Open(filepath.Join(cwd, "files", file.Name()))
		defer f.Close()
		if err != nil {
			panic(err)
		}
		buf := make([]byte, 1024)
		_, err = f.Read(buf)
		addOneToEverything(buf)
		if err != nil {
			panic(err)
		}
	}
}

func addOneToEverything(buf []byte) {
	for i := range buf {
		buf[i]++
	}
}
