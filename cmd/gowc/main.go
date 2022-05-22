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
	opts := gowc.ParseOptions()
	if opts.Version {
		fmt.Printf("gowc %s", gowc.Version)
		return
	}

	var filepaths []string
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

		file := filepath
		wg.Add(1)
		go func(string) {
			c := &gowc.Counter{}
			chunks := gowc.ReadFileInChunks(fp, int(fi.Size()), opts.BufferSize)
			chunkCounterChan := gowc.ProcessChunks(chunks, opts)

			for chunkCounter := range chunkCounterChan {
				gowc.Aggregate(c, chunkCounter)
			}
			if opts.MaxLine {
				maxLine = gowc.MaxLineLength(c)
			}
			mu.Lock()
			gowc.PrintCounter(c, maxLine, file, opts)
			mu.Unlock()
			wg.Done()
		}(file)
	}
	wg.Wait()
}
