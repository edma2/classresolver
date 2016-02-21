package main

import (
	"log"
	"strings"

	"9fans.net/go/acme"
)

func newWin(title string) (*acme.Win, error) {
	win, err := acme.New()
	if err != nil {
		return nil, err
	}
	win.Name(title)
	return win, nil
}

func openWin(name string, childNames []string) {
	w, err := newWin("/zinc/" + name)
	if err != nil {
		log.Printf("acme win: %s\n", err)
	}
	for _, name := range childNames {
		if !strings.ContainsRune(name, '$') {
			w.Fprintf("body", "%s\n", name)
		}
	}
	w.Ctl("clean")
	w.Addr("#0")
	w.Ctl("show")
}
