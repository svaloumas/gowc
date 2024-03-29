package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	gowc "github.com/svaloumas/gowc/internal"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var filepaths []string

	opts := gowc.ParseOptions()
	if opts.Version {
		fmt.Printf("gowc %s", gowc.Version)
		return
	}

	if opts.FilesFrom != "" {
		filefrom, err := os.Open(opts.FilesFrom)
		if err != nil {
			log.Fatalf("could not open files-from file: %s", err)
		}
		scanner := bufio.NewScanner(filefrom)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			filepaths = append(filepaths, scanner.Text())
		}
		filefrom.Close()
	}
	opts.Filepaths = append(opts.Filepaths, filepaths...)

	if len(opts.Filepaths) == 0 ||
		(len(opts.Filepaths) == 1 &&
			(opts.Filepaths[0] == "/dev/stdin" || opts.Filepaths[0] == "/dev/stdin/")) {

		var maxLine int

		fp := os.Stdin
		defer fp.Close()
		fi, err := fp.Stat()
		if err != nil {
			log.Fatal(err)
		}
		fileSize := fi.Size()

		file := fp.Name()
		wg.Add(1)
		go func(string) {
			defer wg.Done()
			c := &gowc.Counter{}
			c.Bytes = int(fileSize)
			if !opts.BytesOnly {
				chunks := gowc.ReadFileInChunks(fp, int(fileSize), opts)

				for chunk := range chunks {
					c.Aggregate(chunk.Counter)
				}
				if opts.MaxLine {
					maxLine = c.MaxLineLength()
				}
			}
			c.PrintCounter(maxLine, file, opts)
		}(file)
	} else {
		for _, filepath := range opts.Filepaths {
			var maxLine int

			fp, err := os.Open(filepath)
			if err != nil {
				log.Fatal(err)
			}
			defer fp.Close()
			fi, err := fp.Stat()
			if err != nil {
				log.Fatal(err)
			}
			fileSize := fi.Size()

			file := filepath
			wg.Add(1)
			go func(string) {
				defer wg.Done()
				c := &gowc.Counter{}
				c.Bytes = int(fileSize)
				if !opts.BytesOnly {
					chunks := gowc.ReadFileInChunks(fp, int(fileSize), opts)

					for chunk := range chunks {
						c.Aggregate(chunk.Counter)
					}
					if opts.MaxLine {
						maxLine = c.MaxLineLength()
					}
				}
				mu.Lock()
				c.PrintCounter(maxLine, file, opts)
				mu.Unlock()
			}(file)
		}
	}
	wg.Wait()
}
