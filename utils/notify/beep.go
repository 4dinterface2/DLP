package notify

import "github.com/gen2brain/beeep"

func InitBeep() func() {
	err2 := beeep.Beep(beeep.DefaultFreq, beeep.DefaultDuration*3)
	if err2 != nil {
		panic(err2)
	}
	return func() {
		err := beeep.Notify("Доступ к whatsup", "отсканируйте qr ", "assets/information.png")
		if err != nil {
			panic(err)
		}
	}
}
