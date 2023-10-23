package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "sync"
)

var (
    clients    = make(map[net.Conn]bool)
    clientsMtx sync.Mutex
)

func main() {
    server, err := net.Listen("tcp", ":8080")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer server.Close()

    for {
        conn, err := server.Accept()
        if err != nil {
            fmt.Println(err)
            continue
        }
        clientsMtx.Lock()
        clients[conn] = true
        clientsMtx.Unlock()

        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()
		ip := conn.RemoteAddr().String()
		conn.Write([]byte("Welcome " + ip + "!\n"))
    reader := bufio.NewReader(conn)
    for {
        message, err := reader.ReadString('\n')
        if err != nil {
            clientsMtx.Lock()
            delete(clients, conn)
            clientsMtx.Unlock()
            break
        }

        clientsMtx.Lock()
        for client := range clients {
            if client != conn {
                client.Write([]byte(message))
            }
        }
        clientsMtx.Unlock()
    }
}
