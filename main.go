package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
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

func getEnvInt(key string, defaultVal int) int {
	_ = godotenv.Load()
	valStr := os.Getenv(key)
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}

func getEnv(key string) string {
	_ = godotenv.Load()

	value := os.Getenv(key)
	if value == "" {
		fmt.Printf("Переменная %s не задана", key)
	}
	return value
}

func sendMessage(botToken string, chatID int, text string) error {
	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    text,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Не удалось отправить сообщение: %s", string(body))
	}
	return nil
}

func hasFlag(flag string) bool {
	for _, arg := range os.Args[2:] {
		if arg == flag {
			return true
		}
	}
	return false
}

func main() {
	interval := getEnvInt("CHECK_INTERVAL", 0)
	botToken := getEnv("BOT_TOKEN")
	chatID := getEnvInt("CHAT_ID", 0)
	sendSuccessNotification := getEnv("SEND_SUCCESS_NOTIFICATION")
	sendFailureNotification := getEnv("SEND_FAILURE_NOTIFICATION")

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
				fmt.Printf("OK. %s - %s:%d доступен\n", server.Name, server.IP, server.Port)
				if sendSuccessNotification == "yes" {
					err := sendMessage(botToken, chatID, fmt.Sprintf("OK. %s - %s:%d доступен\n", server.Name, server.IP, server.Port))
					if err != nil {
						fmt.Println("Ошибка отправки сообщения:", err)
					}
				}
			} else {
				fmt.Printf("FAIL. %s - %s:%d недоступен\n", server.Name, server.IP, server.Port)
				if sendFailureNotification == "yes" {
					err := sendMessage(botToken, chatID, fmt.Sprintf("FAIL. %s - %s:%d недоступен\n", server.Name, server.IP, server.Port))
					if err != nil {
						fmt.Println("Ошибка отправки сообщения:", err)
					}
				}
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
		if len(os.Args) > 3 {
			fmt.Println("Использование: check (флаг --once для одноразового запуска)")
			return
		}

		if interval <= 0 || hasFlag("--once") {
			fmt.Println("Проверка доступности серверов:")
			for _, server := range servers {
				if checkServer(server) {
					fmt.Printf("OK. %s - %s:%d доступен\n", server.Name, server.IP, server.Port)
					if sendSuccessNotification == "yes" {
						err := sendMessage(botToken, chatID, fmt.Sprintf("OK. %s - %s:%d доступен\n", server.Name, server.IP, server.Port))
						if err != nil {
							fmt.Println("Ошибка отправки сообщения:", err)
						}
					}
				} else {
					fmt.Printf("FAIL. %s - %s:%d недоступен\n", server.Name, server.IP, server.Port)
					if sendFailureNotification == "yes" {
						err := sendMessage(botToken, chatID, fmt.Sprintf("FAIL. %s - %s:%d недоступен\n", server.Name, server.IP, server.Port))
						if err != nil {
							fmt.Println("Ошибка отправки сообщения:", err)
						}
					}
				}
			}
			return
		} else {
			for {
				fmt.Println("Проверка доступности серверов:")
				for _, server := range servers {
					if checkServer(server) {
						fmt.Printf("OK. %s - %s:%d доступен\n", server.Name, server.IP, server.Port)
						if sendSuccessNotification == "yes" {
							err := sendMessage(botToken, chatID, fmt.Sprintf("OK. %s - %s:%d доступен\n", server.Name, server.IP, server.Port))
							if err != nil {
								fmt.Println("Ошибка отправки сообщения:", err)
							}
						}
					} else {
						fmt.Printf("FAIL. %s - %s:%d недоступен\n", server.Name, server.IP, server.Port)
						if sendFailureNotification == "yes" {
							err := sendMessage(botToken, chatID, fmt.Sprintf("FAIL. %s - %s:%d недоступен\n", server.Name, server.IP, server.Port))
							if err != nil {
								fmt.Println("Ошибка отправки сообщения:", err)
							}
						}
					}
				}
				time.Sleep(time.Duration(interval) * time.Second)
			}
		}

	default:
		fmt.Println("Неизвестная команда")
	}
}
