package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/svaloumas/gowc/internal/gowc"
)

func main() {
	var wg sync.WaitGroup
	var maxLine int
	opts := gowc.ParseOptions()
	if opts.Version {
		fmt.Printf("gowc %s", gowc.Version)
		return
	}

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
			gowc.PrintCounter(c, maxLine, file, opts)
			wg.Done()
		}(file)
	}
	wg.Wait()
}
