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
	root.Insert("com.twitter.util.Future", "/a/b/c")
	root.Insert("com.twitter.finagle.Addr", "/a/b/c")
	root.Insert("com.twitter.util", "/a/b/c")
	if root.String() != expected {
		t.Log(expected)
		t.Log(root.String())
		t.Fail()
	}
}

func TestLookup(t *testing.T) {
	root := new(Node)
	root.Insert("com.twitter.util.Future", "/a/b/c")
	root.Insert("com.twitter.finagle.Addr", "x")
	root.Insert("com.twitter.util", "/a/b/c")
	if root.Lookup("com.twitter.finagle.Addr").path != "x" {
		t.Error("Expected x")
	}
	if root.Lookup("com.twitter.bar.Addr") != nil {
		t.Error("Expected nil")
	}
}
