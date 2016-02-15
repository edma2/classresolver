package main

import (
	"bufio"
	"flag"
	"os"
	"path"

	"github.com/edma2/pantsindex/index"
)

var (
	root = flag.String("root", ".", "pants root directory")
)

func main() {
	flag.Parse()

	idx := index.NewIndex()
	idx.Watch(path.Join(*root, ".pants.d", "compile", "zinc"))

	in := bufio.NewReader(os.Stdin)
	in.ReadString('\n')

	idx.Stop()
}
