package internal

const WorksetBufferSize = 1E5

type Gobak struct {
	watcher *PollWatcher
	copier *Copier
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

	// TODO: log amount of files

	for k,v := range sourceFilesMap {
		if err := g.copier.Copy(k, v); err != nil {
			// TODO: log error
		}
	}

	// Second, watch files for changes
	workset := NewWorkset(WorksetBufferSize)

	go func() {
		for item := range workset.C() {
			if err := g.copier.Copy(item.path, item.info); err != nil {
				// TODO: log error

				continue
			}

			// TODO: log file successfully copied, also log size and time required
		}

		// TODO: log that copy loop has been shut down
	}()

	g.watcher.StartWatching(workset.Put)
}
