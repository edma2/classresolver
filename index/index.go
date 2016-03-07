package index

import (
	"sort"
	"sync"
)

type Update struct {
	Class string
	Path  string
}

type Index struct {
	tree *Node
	sync.Mutex
}

type GetResult struct {
	Children []string
	Path     string
}

func NewIndex() *Index {
	return &Index{
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
		get.Children[i] = elem
		i++
	}
	sort.Strings(get.Children)
	return get
}

func (idx *Index) Watch(updates chan *Update) {
	go func() {
		for update := range updates {
			idx.Lock()
			idx.tree.Insert(update.Class, update.Path)
			idx.Unlock()
		}
	}()
}
