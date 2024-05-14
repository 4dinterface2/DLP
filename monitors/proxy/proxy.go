package proxy

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func parse(input string) string {
	// Разбиваем строку на строки
	lines := strings.Split(input, "\r\n")

	// Ищем строку, содержащую слово "Host"
	for _, line := range lines {
		if strings.Contains(line, "Host") {
			// Удаляем лишние пробелы
			hostLine := strings.TrimSpace(line)
			// Возвращаем текст после строки с Host
			return strings.TrimSpace(strings.TrimPrefix(hostLine, "Host:"))
		}
	}
	return ""
}

func HandleClient(clientConn net.Conn) {
	// Read the URL and port from the client
	buffer := make([]byte, 1024)
	_, err := clientConn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading client request:", err)
		clientConn.Close()
		return
	}
	fmt.Println("proxy request")

	// Parse the URL from the client request
	clientURL := string(buffer)
	fmt.Println(clientURL)

	host := parse(clientURL)
	fmt.Println(host)

	//targetURL, err := url.Parse(clientURL)
	//
	//if err != nil {
	//	fmt.Println("Error parsing client URL:", err)
	//	clientConn.Close()
	//	return
	//}

	// Connect to the destination server
	destServer, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("Error connecting to destination server:", err)
		clientConn.Close()
		return
	}
	defer destServer.Close()

	// Copy data from client to destination and vice versa
	go io.Copy(destServer, clientConn)
	go io.Copy(clientConn, destServer)
}
