package internal

import (
	"time"

	"github.com/rs/zerolog/log"
)

const WorksetBufferSize = 1e5

type Gobak struct {
	watcher *PollWatcher
	copier  *Copier
}

func NewGobak(watcher *PollWatcher, copier *Copier) *Gobak {
	return &Gobak{
		watcher: watcher,
		copier:  copier,
	}
}

func (g *Gobak) Run() {
	// First, make sure all currently existing files are back-upped
	sourceFilesMap := g.watcher.WatchedFiles()

	var copiedCounter, errCounter, fileCounter, dirCounter int

	for k, v := range sourceFilesMap {
		if v.IsDir() {
			dirCounter++

			// We do not copy directories, it is easier to just copy the files and ensure that the containing
			// folders (and their parents) exist. This would be necessary anyways because the map is not sorted in any
			// way, so we are not guaranteed to handle directories before files.
			continue
		}

		if copied, _, err := g.copier.Copy(k, v); err != nil {
			log.Error().Err(err).Str("path", k).Str("name", v.Name()).Msg("Failed to copy file.")
			errCounter++
		} else if copied {
			copiedCounter++
		}

		fileCounter++
	}

	log.Info().Msgf("Going to watch %d files in %d directories.", fileCounter, dirCounter)

	log.Info().
		Int("copied", copiedCounter).
		Int("errors", errCounter).
		Int("already_present", fileCounter-copiedCounter).
		Msg("Done with initial mirroring.")

	// Second, watch files for changes
	workset := NewWorkset(WorksetBufferSize)

	go func() {
		// TODO: Allow graceful shutdown
		for item := range workset.C() {

			log := log.With().
				Str("path", item.path).
				Str("name", item.info.Name()).
				Stringer("op", item.operation).
				Logger()

			start := time.Now()

			switch item.operation {
			case Remove:
				// TODO: Think about having an additional mechanism like suffixing deleted files or mainting a list of
				// deleted files in order to clean up at some point.
				log.Warn().Msg("Original image has been deleted.")

			case Rename:
				// TODO: Investigate on what's the difference between Rename and Move operation
				log.With().Str("oldPath", item.oldPath).Logger()

				if err := g.copier.Rename(item.oldPath, item.path); err != nil {
					log.Error().Err(err).Msg("Failed to rename file.")
				} else {
					log.Info().Msg("File renamed.")
				}

			case Chmod:
				// Noop, we do not consider this

			case Move:
				log.With().Str("oldPath", item.oldPath).Logger()

				if moved, bytesMoved, err := g.copier.Move(item.oldPath, item.path, item.info); err != nil {
					log.Error().Err(err).Msg("Failed to move file.")
				} else if !moved {
					log.Error().Msg("File already present, therefore neither copied nor removed. Though it had been scheduled by the watcher.")
				} else {
					log.Info().
						Dur("duration", time.Since(start)).
						Int("bytes_copied", bytesMoved).
						Msg("File moved.")
				}

			default: // Write, Create
				if copied, bytesCopied, err := g.copier.Copy(item.path, item.info); err != nil {
					log.Error().Err(err).Msg("Failed to copy file.")

					continue
				} else if !copied {
					log.Error().Msg("File already present, therefore not copied. Though it had been scheduled by the watcher.")
				} else {
					log.Info().
						Dur("duration", time.Since(start)).
						Int("bytes_copied", bytesCopied).
						Msg("File copied.")
				}
			}
		}
	}()

	g.watcher.StartWatching(workset.Put)
}
