package internal

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/rs/zerolog/log"

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
	re, err := regexp.Compile(config.ExcludePattern)
	if err != nil {
		return nil, err
	}

	w := watcher.New()

	w.AddFilterHook(excludeFile(re))

	ops, err := parseOps(config.MonitoredOperations)
	if err != nil {
		return nil, err
	}

	log.Info().Strs("operations", config.MonitoredOperations).Msg("Going to monitor the following operations:")

	w.FilterOps(ops...)

	return &PollWatcher{
		interval: config.Interval,
		watcher:  w,
	}, nil
}

func excludeFile(re *regexp.Regexp) watcher.FilterFileHookFunc {
	// Empty regex allows all files, so we return a dummy FilterFileHookFunc
	if re.String() == "" {
		return func(info os.FileInfo, fullPath string) error {
			return nil
		}
	}

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
				// We do not want to be notified in case of directory changes as we do only care for files
				if event.IsDir() {
					continue
				}

				notifyFn(&WorksetItem{
					path:      event.Path,
					oldPath:   event.OldPath,
					operation: Operation(event.Op),
					info:      event.FileInfo,
				})
			case err := <-p.watcher.Error:
				// TODO: add proper logging
				panic(err)
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
