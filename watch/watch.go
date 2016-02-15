package watch

import (
	"time"

	"github.com/go-fsnotify/fsevents"
)

func Watch(path string, stop chan bool) chan string {
	es := &fsevents.EventStream{
		Paths:   []string{path},
		Latency: 500 * time.Millisecond,
		Flags:   fsevents.FileEvents | fsevents.WatchRoot}
	es.Start()
	paths := make(chan string)
	go func() {
		for {
			select {
			case <-stop:
				es.Stop()
				close(paths)
			case events := <-es.Events:
				for _, e := range events {
					paths <- e.Path
				}
			}
		}
	}()
	return paths
}
