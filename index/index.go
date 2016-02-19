package index

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/edma2/pantsindex/analysis"
	"github.com/edma2/pantsindex/watch"
)

type Index struct {
	tree *Node
	stop chan bool
	sync.Mutex
}

func NewIndex() *Index {
	return &Index{
		stop: make(chan bool),
		tree: new(Node),
	}
}

func (idx *Index) Get(class string) string {
	idx.Lock()
	path := idx.tree.Lookup(class)
	idx.Unlock()
	return path
}

func (idx *Index) Watch(path string) {
	if err := readAnalysisFiles(idx, path); err != nil {
		log.Fatal(err)
	}
	pathChanges := watch.PathChanges(path, idx.stop)
	analysisFileChanges := watch.AnalysisFileChanges(pathChanges)
	analysisChanges := watch.AnalysisChanges(analysisFileChanges)
	go func() {
		for change := range analysisChanges {
			idx.Lock()
			idx.tree.Insert(change.Class, change.Path)
			idx.Unlock()
		}
	}()
}

func readAnalysisFiles(idx *Index, path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if analysis.IsAnalysisFile(path) {
			return analysis.ReadAnalysisFile(path, func(class, path string) {
				idx.tree.Insert(class, path)
			})
		}
		return nil
	})
}

func (idx *Index) Stop() {
	idx.stop <- true
}
