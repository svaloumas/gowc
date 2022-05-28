package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/svaloumas/gowc/internal/gowc"
)

func main() {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var maxLine int
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

	for _, filepath := range opts.Filepaths {
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

				for _, chunk := range chunks {
					chunkCounter := <-chunk.CounterChan
					c.Aggregate(chunkCounter)
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
	wg.Wait()
}
