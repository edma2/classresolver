package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/edma2/pantsindex/watch"
)

func main() {
	stop := make(chan bool)
	pathChanges := watch.PathChanges("/Users/ema/src/source/.pants.d/compile/zinc/", stop)
	analysisFileChanges := watch.AnalysisFileChanges(pathChanges)
	go func() {
		for path := range analysisFileChanges {
			fmt.Println(path)
		}
	}()
	in := bufio.NewReader(os.Stdin)
	in.ReadString('\n')
	stop <- true
}
