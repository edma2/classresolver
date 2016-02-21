package main

import (
	"flag"
	"log"

	"github.com/edma2/zincindexd/index"
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
	defer idx.Stop()
	for _, path := range paths {
		if err := idx.Watch(path); err != nil {
			return err
		}
	}
	return serve(idx)
}

func main() {
	if err := Main(); err != nil {
		log.Fatal(err)
	}
}
