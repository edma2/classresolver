package index

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/edma2/zincindexd/analysis"
	"github.com/edma2/zincindexd/watch"
)

type Index struct {
	tree *Node
	stop chan bool
	sync.Mutex
}

type GetResult struct {
	Children []string
	Path     string
}

func NewIndex() *Index {
	return &Index{
		stop: make(chan bool),
		tree: new(Node),
	}
}

func (idx *Index) Get(name string) *GetResult {
	idx.Lock()
	n := idx.tree.Lookup(name)
	defer idx.Unlock()
	if n == nil {
		return nil
	}
	get := new(GetResult)
	get.Path = n.path
	if len(n.kids) == 0 {
		return get
	}
	get.Children = make([]string, len(n.kids))
	i := 0
	for elem, _ := range n.kids {
		get.Children[i] = name + "." + elem
		i++
	}
	sort.Strings(get.Children)
	return get
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
