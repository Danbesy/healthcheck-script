package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"
)

type Server struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

func loadServers(filename string) ([]Server, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return nil, err
	}

	var servers []Server
	err = json.Unmarshal(data, &servers)

	if err != nil {
		fmt.Println("Ошибка парсинга серверов:", err)
		return nil, err
	}

	return servers, nil
}

func checkServer(s Server) bool {
	address := net.JoinHostPort(s.IP, fmt.Sprintf("%d", s.Port))

	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err != nil {
		return false
	}

	conn.Close()
	return true
}

func main() {
	fmt.Println("Healthcheck script has started!")

	servers, err := loadServers("servers.json")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Список найденных серверов:")

	for _, server := range servers {
		fmt.Println("Имя сервера:", server.Name, "Айпи сервера:", server.IP, "Порт сервера:", server.Port)
	}

	fmt.Println("Проверка доступности серверов:")
	for _, server := range servers {
		if checkServer(server) {
			fmt.Println("OK", server.Name, "доступен")
		} else {
			fmt.Println("FAIL", server.Name, "недоступен")
		}
	}
}
