package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/edma2/pantsindex/watch"
)

func main() {
	stop := make(chan bool)
	paths := watch.Watch("/Users/ema/src/source/.pants.d/compile/zinc/", stop)
	go func() {
		for p := range paths {
			fmt.Println(p)
		}
	}()
	in := bufio.NewReader(os.Stdin)
	in.ReadString('\n')
	stop <- true
}
