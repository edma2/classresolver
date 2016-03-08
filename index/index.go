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
	Name     string
	Children []string
	Path     string
}

func NewIndex() *Index {
	return &Index{
		tree: new(Node),
	}
}

func (idx *Index) Walk(name string, visit func(string)) {
	idx.Lock()
	root := idx.tree.Lookup(name)
	if root != nil {
		root.Walk(visit)
	}
	idx.Unlock()
}

func (idx *Index) Get(name string) *GetResult {
	idx.Lock()
	n := idx.tree.Lookup(name)
	defer idx.Unlock()
	if n == nil {
		return nil
	}
	get := new(GetResult)
	get.Name = name
	get.Path = n.path
	if len(n.kids) == 0 {
		return get
	}
	get.Children = make([]string, len(n.kids))
	i := 0
	for el, _ := range n.kids {
		get.Children[i] = name + "." + el
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
