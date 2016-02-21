package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/edma2/classresolver/index"
	"github.com/edma2/classresolver/zinc"
	"github.com/edma2/classresolver/zinc/fsevents"
)

func Main() error {
	flag.Parse()
	paths := flag.Args()
	if len(paths) == 0 {
		return nil
	}
	for _, path := range paths {
		log.Println("Watching " + path)
	}
	idx := index.NewIndex()
	for _, path := range paths {
		if err := idx.Watch(zinc.Watch(watch(path))); err != nil {
			return err
		}
	}
	return serve(idx)
}

func watch(root string) chan string {
	paths := make(chan string)
	go func() {
		// TODO: handle walk errors
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			paths <- path
			return nil
		})
		for path := range fsevents.Watch(root) {
			paths <- path
		}
	}()
	return paths
}

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}
