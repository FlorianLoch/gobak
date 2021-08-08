package internal

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"

	"github.com/florianloch/gobak/internal/config"
)

var opsMapping = map[string]watcher.Op{
	"create": watcher.Create,
	"write":  watcher.Write,
	"remove": watcher.Remove,
	"rename": watcher.Rename,
	"chmod":  watcher.Chmod,
	"move":   watcher.Move,
}

type NotifyFn func(item *WorksetItem)

type PollWatcher struct {
	interval time.Duration
	watcher  *watcher.Watcher
}

func NewPollWatcher(config *config.Config) (*PollWatcher, error) {
	re, err := regexp.Compile(config.ExlcudePattern)
	if err != nil {
		return nil, err
	}

	w := watcher.New()

	w.AddFilterHook(excludeFile(re))

	ops, err := parseOps(config.MonitoredOperations)
	if err != nil {
		return nil, err
	}

	w.FilterOps(ops...)

	return &PollWatcher{
		interval: config.Interval,
		watcher:  w,
	}, nil
}

func excludeFile(re *regexp.Regexp) watcher.FilterFileHookFunc {
	return func(info os.FileInfo, fullPath string) error {
		if re.MatchString(fullPath) {
			return watcher.ErrSkip
		}

		return nil
	}
}

func parseOps(opsAsStrings []string) ([]watcher.Op, error) {
	var ops []watcher.Op

	for _, opStr := range opsAsStrings {
		if val, found := opsMapping[strings.ToLower(opStr)]; found {
			ops = append(ops, val)

			continue
		}

		return nil, fmt.Errorf("unsupported operation: '%s'", opStr)
	}

	return ops, nil
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
