package fs

import (
	"github.com/rjeczalik/notify"
	"log"
)

func FSWatch() {
	// Создаем новый канал уведомлений
	events := make(chan notify.EventInfo, 1)

	// Путь к корневой папке, которую мы хотим отслеживать
	rootPath := "/Users/evgenijlesovyh"

	// Начинаем отслеживание событий для корневой папки
	if err := notify.Watch(rootPath, events, notify.All); err != nil {
		log.Fatal(err)
	}
	defer notify.Stop(events)

	// что-бы работали подпапки надо дать доступ к ним

	// Рекурсивно добавляем все подпапки в наблюдение
	//err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
	//	if err != nil {
	//		log.Printf("Error accessing path %q: %v\n", path, err)
	//		return err
	//	}
	//
	//	// Если это директория, добавляем её в наблюдение
	//	if info.IsDir() {
	//		if err := notify.Watch(path, events, notify.All); err != nil {
	//			log.Printf("Error watching path %q: %v\n", path, err)
	//			return err
	//		}
	//	}
	//	return nil
	//})
	//if err != nil {
	//	log.Fatal(err)
	//}

	// Бесконечный цикл для обработки событий
	for {
		select {
		case event := <-events:
			// Обработка событий
			log.Println("Event:", event)

			// Вы можете проверять тип события и действовать соответственно
		}
	}
}
