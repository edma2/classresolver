package main

import (
	"flag"
	"log"

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
		if err := idx.Watch(zinc.Watch(fsevents.Watch(path))); err != nil {
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
