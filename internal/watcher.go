package internal

import (
	"log"
	"os"
	"time"

	"github.com/radovskyb/watcher"
)

type NotifyFn func(item *WorksetItem)

type PollWatcher struct {
	interval time.Duration
	watcher  *watcher.Watcher
}

func NewPollWatcher(interval time.Duration) *PollWatcher {
	// TODO: set appropriate filters

	return &PollWatcher{
		interval: interval,
		watcher:  watcher.New(),
	}
}

func (p *PollWatcher) StartWatching(notifyFn NotifyFn) error {
	go func() {
		for {
			select {
			case event := <-p.watcher.Event:
				notifyFn(&WorksetItem{
					path: event.Path,
					info: event.FileInfo,
				})
			case err := <-p.watcher.Error:
				// TODO: add proper logging
				log.Fatalln(err)
			case <-p.watcher.Closed:
				// TODO: add logging
				return
			}
		}
	}()

	return p.watcher.Start(p.interval)
}

func (p *PollWatcher) AddRecursively(path string) error {
	return p.watcher.AddRecursive(path)
}

func (p *PollWatcher) WatchedFiles() map[string]os.FileInfo {
	return p.watcher.WatchedFiles()
}
