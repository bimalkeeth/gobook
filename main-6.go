package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func walkDir(dir string, wg *sync.WaitGroup, fileSizes chan<- int64) {
	defer wg.Done()

	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			subdir := filepath.Join(dir, entry.Name())
			wg.Add(1)
			walkDir(subdir, wg, fileSizes)
		} else {
			fileSizes <- entry.Size()
		}
	}
}

func dirents(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "du1:%v\n", err)
		return nil
	}

	return entries
}

func main() {
	fileSizes := make(chan int64)
	roots := []string{"/Users/bimalkeeth/Personal", "/Users/bimalkeeth/go"}

	var n sync.WaitGroup

	for _, root := range roots {
		n.Add(1)
		go walkDir(root, &n, fileSizes)
	}

	go func() {
		n.Wait()
		close(fileSizes)
	}()

	var nfiles, nbytes int64

	var tick <-chan time.Time
	tick = time.Tick(500 * time.Millisecond)

	for {

		select {
		case size, ok := <-fileSizes:
			if !ok {
				break
			}
			nfiles++
			nbytes += size
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}

	//printDiskUsage(nfiles, nbytes)
}

func printDiskUsage(nfiles, nbytes int64) {
	fmt.Printf("%d files  %.1f GB\n", nfiles, float64(nbytes)/1e9)
}
