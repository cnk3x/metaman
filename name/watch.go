package name

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

type WatchEvent func(fsw *fsnotify.Watcher, evts []fsnotify.Event)

func Watch(ctx context.Context, srcRoot string, onEvent WatchEvent, delay time.Duration) (err error) {
	if delay <= 0 {
		delay = time.Second
	}

	if onEvent == nil {
		onEvent = func(fsw *fsnotify.Watcher, evts []fsnotify.Event) {}
	}

	var fsw *fsnotify.Watcher
	if fsw, err = fsnotify.NewWatcher(); err != nil {
		return
	}

	fswAdd := func(path string) {
		if err := fsw.Add(path); err != nil {
			log.Printf("[fsw] add %s: %v", path, err)
		}
	}

	fswDel := func(path string) {
		if err := fsw.Remove(path); err != nil {
			if s := err.Error(); !strings.Contains(s, "can't remove non-existent watcher") {
				log.Printf("[fsw] del %s: %s", path, s)
			}
		}
	}

	defer fsw.Close()

	err = filepath.WalkDir(srcRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if err = ctx.Err(); err != nil {
			return fs.SkipAll
		}

		if !d.IsDir() {
			return nil
		}

		if !allowPath(path) {
			return nil
		}

		return fsw.Add(path)
	})

	if err != nil {
		return
	}

	backlog := make([]fsnotify.Event, 0, 10)
	l := &sync.Mutex{}
	idle := true

	readBacklog := func() (out []fsnotify.Event) {
		l.Lock()
		out = make([]fsnotify.Event, len(backlog))
		copy(out, backlog)
		backlog = backlog[:0]
		l.Unlock()
		return
	}

	t := time.AfterFunc(time.Minute, func() {
		onEvent(fsw, readBacklog())
		l.Lock()
		idle = true
		l.Unlock()
	})
	defer t.Stop()
	t.Stop()

	addBacklog := func(evt fsnotify.Event) {
		l.Lock()
		backlog = append(backlog, evt)
		if idle {
			idle = false
			t.Reset(delay)
		}
		l.Unlock()
	}

ON_EVENT:
	for {
		select {
		case evt := <-fsw.Events:
			if !allowPath(evt.Name) {
				continue
			}

			switch {
			case evt.Has(fsnotify.Create):
				if info, _ := os.Stat(evt.Name); info != nil && info.IsDir() {
					fswAdd(evt.Name)
				} else if allowFile(evt.Name, info, 200) {
					log.Printf("[fsw] %s", evt.String())
					addBacklog(evt)
				}
			case evt.Has(fsnotify.Remove), evt.Has(fsnotify.Rename):
				fswDel(evt.Name)
			}
		case err = <-fsw.Errors:
			break ON_EVENT
		case <-ctx.Done():
			break ON_EVENT
		}
	}

	return err
}

func allowPath(p string) bool {
	if p == "" {
		return false
	}
	if v := filepath.VolumeName(p); v != "" {
		p = strings.TrimPrefix(p, v)
	}
	ps := strings.Split(filepath.ToSlash(p), "/")
	ss := make([]rune, 0, len(ps))
	for _, s := range ps {
		if s != "" {
			ss = append(ss, []rune(s)[0])
		}
	}
	return !strings.ContainsAny(string(ss), "._-$#~")
}
