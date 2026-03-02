#!/bin/bash

RED="\033[31m"
GREEN="\033[32m"
BLUE="\033[34m"
CYAN="\033[36m"
RESET="\033[0m"

back_to_main_menu() {
    case "$action" in
        0)
            break
    esac
}

confirm_add_server() {
    while true; do
        read -p "Все данные указаны верно? Y/N " answer
        
        case "$answer" in
            [Yy])
                echo -e "${BLUE}Вы выбрали YES. Добавление сервера...${RESET}"
                go run main.go add "$name" "$ip" "$port"
                echo -e "${CYAN}===============================${RESET}"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            [Nn])
                echo -e "${GREEN}Отмена действия${RESET}"
                echo -e "${CYAN}===============================${RESET}"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            *)
                echo -e "${RED}Неверный ввод. Повторите попытку.${RESET}"
        esac
    done
}

confirm_delete_server() {
    while true; do
        read -p "Все данные указаны верно? Y/N " answer
        
        case "$answer" in
            [Yy])
                echo -e "${GREEN}Вы выбрали YES. Удаление сервера...${RESET}"
                go run main.go remove "$ip" "$port"
                echo -e "${CYAN}===============================${RESET}"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            [Nn])
                echo -e "${GREEN}Отмена действия${RESET}"
                echo -e "${CYAN}===============================${RESET}"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            *)
                echo -e "${RED}Неверный ввод. Повторите попытку.${RESET}"
        esac
    done
}

confirm_change_check_interval() {
    while true; do
        read -p "Вы действительно хотите изменить интервал проверки? Y/N " answer
        
        case "$answer" in
            [Yy])
                read -p "Вы выбрали YES. Введите новое значение в секундах: " interval
                sed -i "s|CHECK_INTERVAL=.*|CHECK_INTERVAL=$interval|" .env
                if [ $? -eq 0 ]; then
                    echo -e "${GREEN}Интервал проверки успешно изменен. Новое значение: $interval${RESET}"
                else
                    echo -e "${RED}Ошибка изменения.${RESET}"
                fi
                echo -e "${CYAN}===============================${RESET}"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            [Nn])
                echo -e "${GREEN}Отмена действия${RESET}"
                echo -e "${CYAN}===============================${RESET}"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            *)
                echo -e "${RED}Неверный ввод. Повторите попытку.${RESET}"
        esac
    done
}

confirm_change_bot_token() {
    while true; do
        read -p "Вы действительно хотите изменить токен бота? Y/N " answer
        
        case "$answer" in
            [Yy])
                read -p "Вы выбрали YES. Введите новое значение: " token
                sed -i "s|BOT_TOKEN=.*|BOT_TOKEN=$token|" .env
                if [ $? -eq 0 ]; then
                    echo -e "${GREEN}Токен бота успешно изменен. Новое значение: $token${RESET}"
                else
                    echo -e "${RED}Ошибка изменения.${RESET}"
                fi
                echo -e "${CYAN}===============================${RESET}"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            [Nn])
                echo -e "${GREEN}Отмена действия${RESET}"
                echo -e "${CYAN}===============================${RESET}"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            *)
                echo -e "${RED}Неверный ввод. Повторите попытку.${RESET}"
        esac
    done
}

confirm_change_chat_id() {
    while true; do
        read -p "Вы действительно хотите изменить чат айди? Y/N " answer
        
        case "$answer" in
            [Yy])
                read -p "Вы выбрали YES. Введите новое значение: " chatid
                sed -i "s|CHAT_ID=.*|CHAT_ID=$chatid|" .env
                if [ $? -eq 0 ]; then
                    echo -e "${GREEN}Айди чата успешно изменен. Новое значение: $chatid${RESET}"
                else
                    echo -e "${RED}Ошибка изменения.${RESET}"
                fi
                echo -e "${CYAN}===============================${RESET}"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            [Nn])
                echo -e "${GREEN}Отмена действия${RESET}"
                echo -e "${CYAN}===============================${RESET}"
                read -p "Нажмите Enter чтобы продолжить..."
                break
            ;;
            *)
                echo -e "${RED}Неверный ввод. Повторите попытку.${RESET}"
        esac
    done
}

while true; do
    clear
    echo -e "${CYAN}===============================${RESET}"
    echo "Меню управления Healthcheck скриптом"
    echo -e "${CYAN}===============================${RESET}"
    echo "0) Выйти"
    echo "1) Добавить сервер"
    echo "2) Удалить сервер"
    echo "3) Список серверов"
    echo "4) Проверить сервера"
    echo "5) Изменить интервал проверки"
    echo "6) Изменить BOT TOKEN"
    echo "7) Изменить CHAT ID"
    echo "8) Настройка отправки уведомлений"
    echo -e "${CYAN}===============================${RESET}"
    read -p "Выберите действие: " action
    echo -e "${CYAN}===============================${RESET}"

    if ! [[ "$action" =~ ^[0-8]+$ ]]; then
        echo "Ошибка: нужно ввести число от 0 до 8!"
        echo -e "${CYAN}===============================${RESET}"
        read -p "Нажмите Enter чтобы продолжить..."
        continue
    fi

    case "$action" in
        1)
            read -p "Введите название сервера (например: localhost)... " name
            echo -e "${CYAN}===============================${RESET}"
            read -p "Теперь введите айпи сервера который хотите добавить (например: 127.0.0.1)... " ip
            echo -e "${CYAN}===============================${RESET}"
            read -p "Укажите порт сервера (например: 80)... " port
            echo -e "${CYAN}===============================${RESET}"
            echo -e "${RED}Данные сервера который вы хотите добавить: $name - $ip:$port${RESET}"
            confirm_add_server
        ;;
        2)
            read -p "Введите айпи сервера который хотите удалить (например: 127.0.0.1)... " ip
            echo -e "${CYAN}===============================${RESET}"
            read -p "Теперь укажите порт сервера (например: 80)... " port
            echo -e "${CYAN}===============================${RESET}"
            echo -e "${RED}Данные сервера который вы хотите удалить: $ip:$port${RESET}"
            confirm_delete_server
        ;;
        3)
            go run main.go list
            echo -e "${CYAN}===============================${RESET}"
            read -p "Нажмите Enter чтобы продолжить..."
        ;;
        4)
            echo "0) Назад в главное меню:"
            echo "1) Одноразовая проверка"
            echo "2) Многоразовая проверка с указанным интервалом"
            echo -e "${CYAN}===============================${RESET}"
            read -p "Выберите действие: " sub_action

            case "$sub_action" in
                0)
                    back_to_main_menu
                ;;
                1)
                    echo -e "${CYAN}===============================${RESET}"
                    go run main.go check --once
                    echo -e "${CYAN}===============================${RESET}"
                    read -p "Нажмите Enter чтобы продолжить..."
                ;;
                2)
                    echo -e "${CYAN}===============================${RESET}"
                    echo -e "${RED}Для завершения проверки используйте сочетание клавиш: Ctrl + C${RESET}"
                    echo -e "${CYAN}===============================${RESET}"
                    go run main.go check
                ;;
                *)
                    echo -e "${CYAN}===============================${RESET}"
                    echo -e "${RED}Данное действие не существует. Повторите попытку.${RESET}"
                    echo -e "${CYAN}===============================${RESET}"
                    read -p "Нажмите Enter чтобы продолжить..."
                esac
        ;;
        5)
            source .env
            echo -e "${GREEN}Текущий интервал проверки: $CHECK_INTERVAL${RESET}"
            echo -e "${CYAN}===============================${RESET}"
            confirm_change_check_interval
        ;;
        6)
            source .env
            echo -e "${GREEN}Текущий BOT TOKEN: $BOT_TOKEN${RESET}"
            echo -e "${CYAN}===============================${RESET}"
            confirm_change_bot_token
        ;;
        7)
            source .env
            echo -e "${GREEN}Текущий CHAT ID: $CHAT_ID${RESET}"
            echo -e "${CYAN}===============================${RESET}"
            confirm_change_chat_id
        ;;
        8)
            source .env
            if [[ "$SEND_SUCCESS_NOTIFICATION" == "yes" ]]; then
                message_success="Отключить уведомления об успешной проверке сервера (OK)"
            else
                message_success="Включить уведомления об успешной проверке сервера (OK)"
            fi

            if [[ "$SEND_FAILURE_NOTIFICATION" == "yes" ]]; then
                message_failure="Отключить уведомления о недоступности сервера (FAIL)"
            else
                message_failure="Включить уведомления о недоступности сервера (FAIL)"
            fi

            echo "0) Назад в главное меню:"
            echo "1) $message_success"
            echo "2) $message_failure"
            echo -e "${CYAN}===============================${RESET}"
            read -p "Выберите действие: " sub_action

            case "$sub_action" in
                0)
                    back_to_main_menu
                ;;
                1)
                    if [[ "$SEND_SUCCESS_NOTIFICATION" == "yes" ]]; then
                        sed -i "s|SEND_SUCCESS_NOTIFICATION=.*|SEND_SUCCESS_NOTIFICATION=no|" .env
                    else
                        sed -i "s|SEND_SUCCESS_NOTIFICATION=.*|SEND_SUCCESS_NOTIFICATION=yes|" .env
                    fi
                ;;
                2)
                    if [[ "$SEND_FAILURE_NOTIFICATION" == "yes" ]]; then
                        sed -i "s|SEND_FAILURE_NOTIFICATION=.*|SEND_FAILURE_NOTIFICATION=no|" .env
                    else
                        sed -i "s|SEND_FAILURE_NOTIFICATION=.*|SEND_FAILURE_NOTIFICATION=yes|" .env
                    fi
                ;;
                *)
                    echo -e "${CYAN}===============================${RESET}"
                    echo -e "${RED}Данное действие не существует. Повторите попытку.${RESET}"
                    echo -e "${CYAN}===============================${RESET}"
                    read -p "Нажмите Enter чтобы продолжить..."
                esac
        ;;
        0)
            clear
            exit 0
        ;;
        *)
            echo -e "${RED}Данное действие не существует. Повторите попытку.${RESET}"
            echo -e "${CYAN}===============================${RESET}"
            read -p "Нажмите Enter чтобы продолжить..."
        esac
    done