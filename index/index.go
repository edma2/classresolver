package index

import (
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/edma2/classresolver/analysis"
	"github.com/edma2/classresolver/watch"
)

type Index struct {
	tree  *Node
	stops []chan bool
	sync.Mutex
}

type GetResult struct {
	Children []string
	Path     string
}

func NewIndex() *Index {
	return &Index{
		stops: nil,
		tree:  new(Node),
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

func (idx *Index) Watch(path string) error {
	if err := readAnalysisFiles(idx, path); err != nil {
		return err
	}
	stop := make(chan bool)
	idx.stops = append(idx.stops, stop)
	pathChanges := watch.PathChanges(path, stop)
	analysisFileChanges := watch.AnalysisFileChanges(pathChanges)
	analysisChanges := watch.AnalysisChanges(analysisFileChanges)
	go func() {
		for change := range analysisChanges {
			idx.Lock()
			idx.tree.Insert(change.Class, change.Path)
			idx.Unlock()
		}
	}()
	return nil
}

func readAnalysisFiles(idx *Index, path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if analysis.IsAnalysisFile(path) {
			return analysis.ReadAnalysisFile(path, func(class, path string) {
				idx.tree.Insert(class, path)
			})
		}
		return nil
	})
}

func (idx *Index) Stop() {
	for _, stop := range idx.stops {
		stop <- true
	}
}
