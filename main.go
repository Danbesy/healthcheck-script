package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
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

func saveServer(filename string, servers []Server) error {
	data, err := json.MarshalIndent(servers, "", " ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func removeServer(servers []Server, ip string, port int) ([]Server, error) {
	var updated []Server
	found := false

	for _, s := range servers {
		if s.IP == ip && s.Port == port {
			found = true
			continue
		}
		updated = append(updated, s)
	}

	if !found {
		return servers, fmt.Errorf("Сервер %s:%d не найден", ip, port)
	}

	return updated, nil
}

func main() {
	servers, err := loadServers("servers.json")
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Healthcheck скрипт успешно запущен!")
		fmt.Println("-------------------------")
		fmt.Println("Использование: add, remove, list, check")
		fmt.Println("-------------------------")
		fmt.Println("Список найденных серверов:")

		for _, server := range servers {
			fmt.Println("Имя сервера:", server.Name, "Айпи сервера:", server.IP, "Порт сервера:", server.Port)
		}

		fmt.Println("-------------------------")

		fmt.Println("Проверка доступности серверов:")
		for _, server := range servers {
			if checkServer(server) {
				fmt.Println("OK", server.Name, "доступен")
			} else {
				fmt.Println("FAIL", server.Name, "недоступен")
			}
		}
		return
	}

	command := os.Args[1]

	switch command {

	case "add":
		if len(os.Args) < 5 {
			fmt.Println("Использование: add <name> <ip> <port>")
			return
		}

		name := os.Args[2]
		ip := os.Args[3]

		port, err := strconv.Atoi(os.Args[4])
		if err != nil {
			fmt.Println("Порт должен быть числом")
			return
		}

		newServer := Server{
			Name: name,
			IP:   ip,
			Port: port,
		}

		servers = append(servers, newServer)

		err = saveServer("servers.json", servers)
		if err != nil {
			fmt.Println("Ошибка сохранения сервера:", err)
			return
		}

		fmt.Println("Сервер добавлен:", name)

	case "remove":
		if len(os.Args) < 4 {
			fmt.Println("Использование: remove <ip> <port>")
			return
		}

		ip := os.Args[2]

		port, err := strconv.Atoi(os.Args[3])
		if err != nil {
			fmt.Println("Порт должен быть числом")
			return
		}

		servers, err = removeServer(servers, ip, port)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = saveServer("servers.json", servers)
		if err != nil {
			fmt.Println("Ошибка сохранения в файл:", err)
			return
		}

		fmt.Println("Сервер удален:", ip)

	case "list":
		if len(os.Args) > 2 {
			fmt.Println("Использование: list (без дополнительных аргументов)")
			return
		}

		fmt.Println("Список найденных серверов:")

		for _, server := range servers {
			fmt.Println("Имя сервера:", server.Name, "Айпи сервера:", server.IP, "Порт сервера:", server.Port)
		}
		return

	case "check":
		if len(os.Args) > 2 {
			fmt.Println("Использование: check (без дополнительных аргументов)")
			return
		}

		fmt.Println("Проверка доступности серверов:")
		for _, server := range servers {
			if checkServer(server) {
				fmt.Printf("OK. %s - %s:%d доступен\n", server.Name, server.IP, server.Port)
			} else {
				fmt.Printf("FAIL. %s - %s:%d недоступен\n", server.Name, server.IP, server.Port)

			}
		}
		return

	default:
		fmt.Println("Неизвестная команда")
	}
}
