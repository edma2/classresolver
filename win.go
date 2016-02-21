package main

import (
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

func openWin(name string, names []string) error {
	w, err := newWin("/class/" + name)
	if err != nil {
		return err
	}
	for _, name := range names {
		if !strings.ContainsRune(name, '$') {
			w.Fprintf("body", "%s\n", name)
		}
	}
	w.Ctl("clean")
	w.Addr("#0")
	w.Ctl("show")
	return nil
}
