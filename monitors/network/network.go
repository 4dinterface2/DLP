package network

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"log"
)

func networkMonitor() {
	// Указываем интерфейс для прослушивания трафика
	iface := "en0"

	// Открываем сессию захвата на указанном интерфейсе
	handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Настраиваем фильтр для захвата трафика
	err = handle.SetBPFFilter("tcp")
	if err != nil {
		log.Fatal(err)
	}

	// Создаем пакетный декодер
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Начинаем бесконечный цикл обработки пакетов
	for packet := range packetSource.Packets() {
		// Выводим информацию о пакете
		fmt.Println(packet)
		// Для примера, остановим обработку после первого пакета
		// break
	}

	// Подождем некоторое время перед завершением программы
	// time.Sleep(5 * time.Second)
}
