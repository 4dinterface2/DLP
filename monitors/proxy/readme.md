вроде как можно извлечь инфу из тела запроса
```
import (
    "bufio"
    "fmt"
    "net"
    "strings"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()

    // Создаем Reader для чтения данных из соединения
    reader := bufio.NewReader(conn)

    // Читаем первую строку, которая обычно содержит заголовок запроса HTTP
    requestLine, err := reader.ReadString('\n')
    if err != nil {
        fmt.Println("Error reading request line:", err)
        return
    }

    // Печатаем заголовок запроса
    fmt.Println("Request Line:", requestLine)

    // Прочитаем остальные заголовки
    for {
        headerLine, err := reader.ReadString('\n')
        if err != nil || headerLine == "\r\n" {
            break
        }
        fmt.Println("Header:", headerLine)
    }

    // Прочитаем тело запроса (если есть)
    for {
        bodyLine, err := reader.ReadString('\n')
        if err != nil || bodyLine == "\r\n" {
            break
        }
        fmt.Println("Body Line:", bodyLine)
    }
}

func main() {
    // Пример использования listener.Accept()
    listener, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println("Error listening:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Listening on :8080...")

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }

        go handleConnection(conn)
    }
}

```