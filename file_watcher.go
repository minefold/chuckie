package main

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type Watcher struct {
	C chan string

	cancelled   bool
	root        string
	observation time.Duration
}

// pushes new files onto c when they've stopped getting bigger after observation time
func NewWatcher(path string, observation time.Duration) *Watcher {
	return &Watcher{
		C:           make(chan string, 0),
		root:        path,
		observation: observation,
	}
}

func (w *Watcher) Watch() {
	var wg sync.WaitGroup

	w.loop(&wg)
	wg.Wait()
	close(w.C)
}

func (w *Watcher) Cancel() {
	w.cancelled = true
}

func (w *Watcher) loop(wg *sync.WaitGroup) {
	files := make(map[string]bool)
	for !w.cancelled {
		w.walk(wg, files)
		time.Sleep(100 * time.Millisecond)
	}
  w.walk(wg, files)
}

func (w *Watcher) walk(wg *sync.WaitGroup, files map[string]bool) {
	filepath.Walk(w.root,
		func(path string, info os.FileInfo, err error) error {
			lower := strings.ToLower(path) // multiple filewalks return files with different case...
			if !info.IsDir() {
				if _, exists := files[lower]; !exists {
					files[lower] = true
					wg.Add(1)
					go pushWhenWriteStops(w.C, wg, path, w.observation)
				}
			}
			return nil
		})

}

func pushWhenWriteStops(c chan string, wg *sync.WaitGroup, path string, observation time.Duration) {
	fi, err := os.Stat(path)
	if err != nil {
		return
	}

	size := fi.Size()
	mtime := time.Now()

	for {
		time.Sleep(100 * time.Millisecond)

		fi, err := os.Stat(path)
		if err != nil {
			return
		}

		age := time.Now().Sub(mtime)
		if age > observation && size == fi.Size() {
			c <- path
			wg.Done()
			return
		}

		if fi.Size() != size {
			mtime = time.Now()
			size = fi.Size()
		}
	}

}
