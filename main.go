package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/edma2/pantsindex/analysis"
	"github.com/edma2/pantsindex/watch"
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

	stop := make(chan bool)
	pathChanges := watch.PathChanges(analysis.PantsRoot+"/.pants.d/compile/zinc/", stop)
	analysisFileChanges := watch.AnalysisFileChanges(pathChanges)
	analysisChanges := watch.AnalysisChanges(analysisFileChanges)

	go func() {
		for change := range analysisChanges {
			fmt.Printf("%s -> %s\n", change.Class, change.Path)
		}
	}()
	in := bufio.NewReader(os.Stdin)
	in.ReadString('\n')
	stop <- true
}
