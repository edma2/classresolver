package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"

	"github.com/edma2/classresolver/index"

	"9fans.net/go/plan9"
	"9fans.net/go/plumb"
)

func leafOf(name string) string {
	if i := strings.LastIndexByte(name, '.'); i != -1 && i+1 <= len(name) {
		return name[i+1:]
	}
	return ""
}

func candidatesOf(name string) []string {
	candidates := []string{}
	elems := strings.Split(name, ".")
	for i, _ := range elems {
		candidates = append(candidates, strings.Join(elems[0:i+1], "."))
	}
	sort.Sort(sort.Reverse(sort.StringSlice(candidates)))
	return candidates
}

func plumbFile(m *plumb.Message, w io.Writer, name, path string) error {
	m.Src = "classresolver"
	m.Dst = ""
	m.Data = []byte(path)
	var attr *plumb.Attribute
	for attr = m.Attr; attr != nil; attr = attr.Next {
		if attr.Name == "addr" {
			break
		}
	}
	if attr == nil {
		if leafName := leafOf(name); leafName != "" {
			addr := fmt.Sprintf("/(trait|class|object|interface)[ 	]*%s/", leafName)
			m.Attr = &plumb.Attribute{Name: "addr", Value: addr, Next: m.Attr}
		}
	}
	return m.Send(w)
}

func serve(idx *index.Index) error {
	fid, err := plumb.Open("classedit", plan9.OREAD)
	if err != nil {
		return err
	}
	defer fid.Close()
	r := bufio.NewReader(fid)
	w, err := plumb.Open("send", plan9.OWRITE)
	if err != nil {
		return err
	}
	defer w.Close()
	for {
		m := plumb.Message{}
		err := m.Recv(r)
		if err != nil {
			log.Println(err)
		}
		name := string(m.Data)
		var get *index.GetResult
		for _, c := range candidatesOf(name) {
			if get = idx.Get(c); get != nil {
				break
			}
		}
		if get == nil {
			continue
		}
		if get.Path != "" {
			if err := plumbFile(&m, w, name, get.Path); err != nil {
				log.Println(err)
			}
		}
		if get.Children != nil {
			if err := openWin(name, get.Children); err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}
