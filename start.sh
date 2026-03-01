#!/bin/bash

confirm_add_server() {
    while true; do
        read -p "Все данные указаны верно? Y/N " answer
        
        case "$answer" in
            [Yy])
                echo "Вы выбрали YES. Добавление сервера..."
                go run main.go add "$name" "$ip" "$port"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            [Nn])
                echo "Отмена действия"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            *)
                echo "Неверный ввод. Повторите попытку."
        esac
    done
}

confirm_delete_server() {
    while true; do
        read -p "Все данные указаны верно? Y/N " answer
        
        case "$answer" in
            [Yy])
                echo "Вы выбрали YES. Удаление сервера..."
                go run main.go remove "$ip" "$port"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            [Nn])
                echo "Отмена действия"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            *)
                echo "Неверный ввод. Повторите попытку."
        esac
    done
}

while true; do
    clear
    echo "Меню управления Healthcheck скриптом."
    echo "-------------------------"
    echo "Навигация:"
    echo "1. Добавить сервер"
    echo "2. Удалить сервер"
    echo "3. Список серверов"
    echo "4. Проверить сервера"
    echo "-------------------------"
    read -p "Выберите действие: " action

    if ! [[ "$action" =~ ^[0-9]+$ ]]; then
        echo "Ошибка: нужно ввести число от 1 до 2!"
        read -p "Нажмите Enter чтобы продолжить..."
        continue
    fi

    case "$action" in
        1)
            read -p "Введите название сервера (например: localhost)... " name
            echo "-------------------------"
            read -p "Теперь введите айпи сервера который хотите добавить (например: 127.0.0.1)... " ip
            echo "-------------------------"
            read -p "Укажите порт сервера (например: 80)... " port
            echo "-------------------------"
            echo "Данные сервера который вы хотите добавить: $name - $ip:$port"
            confirm_add_server
        ;;
        2)
            read -p "Введите айпи сервера который хотите удалить (например: 127.0.0.1)... " ip
            echo "-------------------------"
            read -p "Теперь укажите порт сервера (например: 80)... " port
            echo "-------------------------"
            echo "Данные сервера который вы хотите удалить: $ip:$port"
            confirm_delete_server
        ;;
        3)
            go run main.go list
            read -p "Нажмите Enter чтобы продолжить..."
        ;;
        4)
            go run main.go check
            read -p "Нажмите Enter чтобы продолжить..."
        ;;
        *)
            echo "Данное действие не существует. Повторите попытку."
            read -p "Нажмите Enter чтобы продолжить..."
        esac
    done