package main

import (
	"bufio"
	"flag"
	"fmt"
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

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "stop" {
			break
		}
		fmt.Println(idx.Get(text))
	}
	idx.Stop()
}
