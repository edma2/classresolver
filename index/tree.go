package index

import (
	"fmt"
	"sort"
)

type Node struct {
	kids map[string]*Node
	path string
}

func (n *Node) Insert(name []string, path string) {
	if len(name) == 0 {
		n.path = path
		return
	}
	if n.kids == nil {
		n.kids = make(map[string]*Node)
	}
	var k *Node
	var ok bool
	if k, ok = n.kids[name[0]]; !ok {
		k = new(Node)
		n.kids[name[0]] = k
	}
	k.Insert(name[1:], path)
}

func (n *Node) Lookup(name []string) string {
	if len(name) == 0 {
		return n.path
	}
	if k, ok := n.kids[name[0]]; ok {
		return k.Lookup(name[1:])
	}
	return ""
}

func (n *Node) String() string {
	return n.string(0)
}

func (n *Node) string(depth int) string {
	s := ""
	elems := make([]string, len(n.kids))
	i := 0
	for k, _ := range n.kids {
		elems[i] = k
		i++
	}
	sort.Strings(elems)
	for _, elem := range elems {
		k := n.kids[elem]
		for i := 0; i < depth; i++ {
			s += "\t"
		}
		s += elem
		if k.path != "" {
			s += fmt.Sprintf(" (%s)", k.path)
		}
		s += "\n"
		s += k.string(depth + 1)
	}
	return s
}
