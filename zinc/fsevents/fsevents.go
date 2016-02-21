package fsevents

import (
	"time"

	"github.com/go-fsnotify/fsevents"
)

func Watch(root string) chan string {
	es := &fsevents.EventStream{
		Paths:   []string{root},
		Latency: 500 * time.Millisecond,
		Flags:   fsevents.FileEvents | fsevents.WatchRoot}
	es.Start()
	paths := make(chan string)
	go func() {
		for events := range es.Events {
			for _, event := range events {
				paths <- event.Path
			}
		}
	}()
	return paths
}
