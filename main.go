package main

import (
	"fmt"
	"os"
)

func main() {
	create100Files()
}

func create100Files() {
	for i := 0; i < 100; i++ {
		// create file
		f, err := os.Create(fmt.Sprintf("files/file%d.txt", i))
		if err != nil {
			panic(err)
		}

		// write to file
		_, err = f.Write([]byte(fmt.Sprintf("This is file %d\n", i)))
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
	for _, file := range files {
		// open file
		f, err := os.Open(fmt.Sprintf("files/%s", file.Name()))
		if err != nil {
			panic(err)
		}
		f.Close()

	}
}

func OpenAllFilesDefer() {
	files, err := os.ReadDir("files/")
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		// open file
		f, err := os.Open(fmt.Sprintf("files/%s", file.Name()))
		defer f.Close()
		if err != nil {
			panic(err)
		}
	}
}
