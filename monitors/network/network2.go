package network

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
)

func NetworkMonitor2() {
	// Указываем интерфейс для прослушивания трафика
	iface := "en0"

	// Открываем сессию захвата на указанном интерфейсе
	handle, err := pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Создаем пакетный декодер
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	// Начинаем бесконечный цикл обработки пакетов
	for packet := range packetSource.Packets() {

		// получаем ip слой пакета
		if ipLayer := packet.Layer(layers.LayerTypeIPv4); ipLayer != nil {
			ip, _ := ipLayer.(*layers.IPv4)
			fmt.Println("Destination IP:", ip.DstIP)
		} else {
			log.Println("IPv4 layer not found")
		}

		// Получаем TCP слой пакета
		tcpLayer := packet.Layer(layers.LayerTypeTCP)
		if tcpLayer == nil {
			continue // Пропускаем пакеты без TCP слоя
		}
		tcp, _ := tcpLayer.(*layers.TCP)

		// Модифицируем данные TCP пакета, добавляя новый текст
		if tcp != nil {
			fmt.Println("Исходные данные TCP:", string(tcp.Payload))
		}

		// Выводим информацию о пакете
		fmt.Println(packet)
	}
	// Подождем некоторое время перед завершением программы
	// time.Sleep(5 * time.Second)
}
