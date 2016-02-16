package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"path"
	"sort"
	"strings"

	"9fans.net/go/plan9"
	"9fans.net/go/plumb"

	"github.com/edma2/pantsindex/index"
)

var (
	root = flag.String("root", "/Users/ema/src/source/", "pants root directory")
)

func leafOf(class string) string {
	if i := strings.LastIndexByte(class, '.'); i != -1 && i+1 <= len(class) {
		return class[i+1:]
	}
	return ""
}

func candidatesOf(class string) []string {
	candidates := []string{}
	elems := strings.Split(class, ".")
	for i, _ := range elems {
		candidates = append(candidates, strings.Join(elems[0:i+1], "."))
	}
	sort.Sort(sort.Reverse(sort.StringSlice(candidates)))
	return candidates
}

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
		path := ""
		for _, c := range candidatesOf(class) {
			if path = idx.Get(c); path != "" {
				break
			}
		}
		if path != "" {
			m.Src = "pantsindex"
			m.Dst = ""
			m.Data = []byte(path)
			if leafName := leafOf(class); leafName != "" {
				addr := fmt.Sprintf("/(trait|class|object|interface)[ 	]*%s/", leafName)
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
