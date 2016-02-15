package index

import (
	"log"
	"os"
	"path/filepath"

	"github.com/edma2/pantsindex/analysis"
	"github.com/edma2/pantsindex/watch"
)

type Index struct {
	m    map[string]string
	stop chan bool
}

func NewIndex() *Index {
	return &Index{
		stop: make(chan bool),
		m:    make(map[string]string),
	}
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
			// TODO: mutex
			idx.m[change.Class] = change.Path
		}
	}()
}

func readAnalysisFiles(idx *Index, path string) error {
	return filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if analysis.IsAnalysisFile(path) {
			return analysis.ReadAnalysisFile(path, func(class, path string) {
				idx.m[class] = path
			})
		}
		return nil
	})
}

func (idx *Index) Stop() {
	idx.stop <- true
}
