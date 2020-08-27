package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	KB int = 1000
	MB int = 1000000
	GB int = 1000000000
)

func main() {
	flag.Parse()
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"./"}
	}
	// key: rootディレクトリ, value: 各ディレクトリのサイズのchannel
	sizeChMap := make(map[string]chan int)
	for _, root := range roots {
		sizeChMap[root] = make(chan int)
		go func(root string) {
			size := calcFileSize(root)
			sizeChMap[root] <- size
			close(sizeChMap[root])
		}(root)
	}

	for root, ch := range sizeChMap {
		size := <-ch
		print(root, size)
	}
}

func calcFileSize(dir string) int {
	var total int
	for _, entry := range entries(dir) {
		if entry.IsDir() {
			dir := filepath.Join(dir, entry.Name())
			total += calcFileSize(dir)
		} else {
			total += int(entry.Size())
		}
	}

	return total
}

func entries(dir string) []os.FileInfo {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		return nil
	}

	return entries
}

func print(root string, size int) {
	fmt.Printf("%s\t%s\n", root, humanizeSize(size))
}

func humanizeSize(size int) string {
	switch {
	case size >= GB:
		return fmt.Sprintf("%.1f G", float64(size)/float64(GB))
	case size >= MB && size < GB:
		return fmt.Sprintf("%.1f M", float64(size)/float64(MB))
	case size >= KB && size < MB:
		return fmt.Sprintf("%.1f K", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
