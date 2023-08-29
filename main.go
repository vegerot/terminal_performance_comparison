package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	start := time.Now()
	chunkSize := 12
	files := 100
	if len(os.Args) == 1 {
		benchmark(1)
		return
	}

	args := os.Args[1:]

	switch arg := args[0]; arg {
	case "create":
		createFiles(files)
	case "nodefer":
		OpenAllFilesNoDefer()
	case "defer":
		OpenAllFilesDefer()
	case "nodefer-chunk":
		OpenAllFilesNoDeferChunked(chunkSize)
	case "defer-chunk":
		OpenAllFilesDeferChunked(chunkSize)
	case "bench":
		var runs int
		if len(args) < 2 {
			runs = 1
		} else {
			runss, err := strconv.Atoi(args[1])
			runs = runss
			if err != nil {
				panic(err)
			}
		}
		benchmark(runs)
	default:
		log.Fatalf("unknown command: %s", arg)
	}
	fmt.Printf("program time: %s\n", time.Since(start))
}

func benchmark(N int) {
	now := time.Now()
	for i := 0; i < N; i++ {
		OpenAllFilesNoDefer()
	}
	durationPerRun := time.Since(now) / time.Duration(N)
	fmt.Printf("OpenAllFilesNoDefer: %s\n", durationPerRun)

	now = time.Now()
	for i := 0; i < N; i++ {
		OpenAllFilesDefer()
	}
	durationPerRun = time.Since(now) / time.Duration(N)
	fmt.Printf("OpenAllFilesDefer: %s\n", durationPerRun)

	// NOTE: OpenAllFilesDefer and OpenAllFilesLeakFileHandles have the same
	// speed, but it depends on which order you run them in.  Uncomment just
	// one of these blocks and you'll see the same speed

	/*{
		now := time.Now()
		for i := 0; i < N; i++ {
			OpenAllFilesLeakFileHandles()
		}
		durationPerRun := time.Since(now) / time.Duration(N)
		fmt.Printf("OpenAllFilesLeakFileHandles: %s\n", durationPerRun)
	}*/

}

func createFiles(n int) {
	// create files directory
	err := os.Mkdir("files", 0700)
	if err != nil {
		panic(err)
	}
	for i := 0; i < n; i++ {
		// create file
		f, err := os.Create(fmt.Sprintf("files/file%d.txt", i))
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

func OpenAllFilesNoDefer() {
	files, err := os.ReadDir("files/")
	if err != nil {
		panic(err)
	}

	openTheseFilesNoDefer(files)
}

func OpenAllFilesDefer() {
	files, err := os.ReadDir("files/")
	if err != nil {
		panic(err)
	}

	openTheseFilesDefer(files)
}

func OpenAllFilesLeakFileHandles() {
	files, err := os.ReadDir("files/")
	if err != nil {
		panic(err)
	}

	openTheseFilesLeakFileHandles(files)
}

func OpenAllFilesNoDeferChunked(chunkSize int) {
	files, err := os.ReadDir("files/")
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
	files, err := os.ReadDir("files/")
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
	for _, file := range files {
		f, err := os.Open(fmt.Sprintf("files/%s", file.Name()))
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func openTheseFilesDefer(files []os.DirEntry) {
	for _, file := range files {
		f, err := os.Open(fmt.Sprintf("files/%s", file.Name()))
		defer f.Close()
		if err != nil {
			panic(err)
		}
	}
}

func openTheseFilesLeakFileHandles(files []os.DirEntry) {
	for _, file := range files {
		_, err := os.Open(fmt.Sprintf("files/%s", file.Name()))
		if err != nil {
			panic(err)
		}
	}
}
