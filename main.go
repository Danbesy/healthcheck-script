package main

import (
	"encoding/json"
	"fmt"
	"os"
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
}
