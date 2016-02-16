package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"path"
	"strings"

	"9fans.net/go/plan9"
	"9fans.net/go/plumb"

	"github.com/edma2/pantsindex/index"
)

var (
	root = flag.String("root", "/Users/ema/src/source/", "pants root directory")
)

func servePlumber(idx *index.Index, r io.ByteReader) {
	send, err := plumb.Open("send", plan9.OWRITE)
	if err != nil {
		log.Fatalf("error opening plumb/send: %s\n", err)
	}
	defer send.Close()
	for {
		m := plumb.Message{}
		err := m.Recv(r)
		if err != nil {
			log.Printf("recv error: %s\n", err)
		}
		class := string(m.Data)
		path := idx.Get(class)
		if path == "" {
			// a.b.c -> a.b.$c
			if i := strings.LastIndexByte(class, '.'); i != -1 {
				path = idx.Get(class[0:i] + "$" + class[i+1:])
			}
		}
		if path != "" {
			m.Src = "pantsindex"
			m.Dst = ""
			m.Data = []byte(path)
			if i := strings.LastIndexByte(class, '.'); i != -1 {
				leafName := class[i+1:]
				addr := fmt.Sprintf("/(trait|class|object|interface) %s/", leafName)
				m.Attr = &plumb.Attribute{Name: "addr", Value: addr}
			}
			if err := m.Send(send); err != nil {
				log.Printf("send error: %s\n", err)
			}
		} else {
			log.Println("couldn't find " + class)
		}
	}
}

func main() {
	flag.Parse()
	plumber, err := plumb.Open("pantsindex", plan9.OREAD)
	if err != nil {
		log.Fatalf("error opening plumb/pantsindex: %s\n", err)
	}
	defer plumber.Close()

	idx := index.NewIndex()
	idx.Watch(path.Join(*root, ".pants.d", "compile", "zinc"))
	servePlumber(idx, bufio.NewReader(plumber))
	idx.Stop() // not reached
}
