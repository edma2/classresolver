package main

import (
	"bufio"
	"flag"
	"os"
	"path"

	"github.com/edma2/pantsindex/analysis"
	"github.com/edma2/pantsindex/index"
)

var (
	pantsRootFlag = flag.String("root", "", "pants root directory")
)

func main() {
	flag.Parse()

	if *pantsRootFlag == "" {
		if analysis.PantsRoot = os.Getenv("PANTSROOT"); analysis.PantsRoot == "" {
			cwd, _ := os.Getwd()
			analysis.PantsRoot = cwd
		}
	} else {
		analysis.PantsRoot = *pantsRootFlag
	}

	idx := index.NewIndex()
	idx.Watch(path.Join(analysis.PantsRoot, ".pants.d", "compile", "zinc"))

	in := bufio.NewReader(os.Stdin)
	in.ReadString('\n')

	idx.Stop()
}
