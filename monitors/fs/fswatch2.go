package fs

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func FSWatch2() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	rootPath := "/Users/evgenijlesovyh"

	// Рекурсивно добавляем все подпапки в наблюдение
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error accessing path %q: %v\n", path, err)
			return err
		}

		// Если это директория, добавляем её в наблюдение
		if info.IsDir() {
			if err := watcher.Add(path); err != nil {
				log.Printf("Error watching path %q: %v\n", path, err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Бесконечный цикл для обработки событий
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Println("Event:", event)

			// Вы можете проверять тип события и действовать соответственно

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Error:", err)
		}
	}
}
