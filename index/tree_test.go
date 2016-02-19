package index

import "testing"

func TestInsert(t *testing.T) {
	expected := `com
	twitter
		finagle
			Addr (/a/b/c)
		util (/a/b/c)
			Future (/a/b/c)
`
	root := new(Node)
	root.Insert([]string{"com", "twitter", "util", "Future"}, "/a/b/c")
	root.Insert([]string{"com", "twitter", "finagle", "Addr"}, "/a/b/c")
	root.Insert([]string{"com", "twitter", "util"}, "/a/b/c")
	if root.String() != expected {
		t.Log(expected)
		t.Log(root.String())
		t.Fail()
	}
}
