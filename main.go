package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog/log"
)

func main() {
	fileToWatch := flag.String("fileToWatch", "", "the name of the file to watch")
	watchParent := flag.Bool("watchParent", false, "watch the parent dir of fileToWatch")
	flag.Parse()

	if *fileToWatch == "" {
		log.Error().Msg("cannot watch file: \"\"")
		os.Exit(1)
	}

	log.Info().Msgf("file: %s", *fileToWatch)
	contents, err := readFile(*fileToWatch)
	if err != nil {
		log.Error().Msgf("failed to read file: %s, got err: %w", *fileToWatch, err)
		os.Exit(1)
	}
	log.Info().Msgf("file contents: %s", contents)

	// Create a watcher
	// If your are developing on a Mac, you may want to export GOOS=linx or
	// your LSP may take you down the BSD kqueue road which is
	// confusing|fustrating.
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panic().Msgf("got err: %w, creating watcher", err)
	}

	if *watchParent {
		parentDir := filepath.Dir(*fileToWatch)
		log.Info().Msgf("watching parent dir: %s", parentDir)
		w.Add(parentDir)
	}

	log.Info().Msgf("watching: %s", *fileToWatch)
	w.Add(*fileToWatch)

	// Log events
	for {
		e := <-w.Events
		log.Info().Msgf("got event: %s for file: %s", e.Op.String(), e.Name)
		contents, err := readFile(*fileToWatch)
		if err != nil {
			log.Error().Msgf("failed to read file: %s, got err: %w", *fileToWatch, err)
			os.Exit(1)
		}
		log.Info().Msgf("file contents: %s", contents)
	}
}

func readFile(fileToWatch string) (string, error) {
	b, err := os.ReadFile(fileToWatch)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
