package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os/exec"
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
			plumbEdit(path)
		} else {
			log.Println("couldn't find " + class)
		}
	}
}

// TODO: don't spawn external process
func plumbEdit(path string) {
	out, err := exec.Command("plumb", "-d", "edit", path).CombinedOutput()
	if err != nil {
		log.Fatalf("plumb: %v\n%s", err, out)
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
