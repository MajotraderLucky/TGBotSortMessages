package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/session"
	"github.com/gotd/td/tg"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Преобразуем API_ID в int
	apiID, err := strconv.Atoi(os.Getenv("API_ID"))
	if err != nil {
		log.Fatal("Ошибка конвертации API_ID в int:", err)
	}
	apiHash := os.Getenv("API_HASH")
	if apiHash == "" {
		log.Fatal("API_HASH не указан в .env")
	}

	// Создаем хранилище сессии
	sessStorage := &session.FileStorage{Path: "session.json"}
	client := telegram.NewClient(apiID, apiHash, telegram.Options{
		SessionStorage: sessStorage,
	})

	ctx := context.Background()

	// Запускаем клиент
	err = client.Run(ctx, func(ctx context.Context) error {
		log.Println("Бот запущен и подключен к Telegram")
		api := tg.NewClient(client)

		// Проверяем авторизацию
		authClient := client.Auth()
		status, err := authClient.Status(ctx)
		if err != nil {
			return fmt.Errorf("Ошибка проверки статуса авторизации: %w", err)
		}
		if !status.Authorized {
			if err := login(ctx, authClient); err != nil {
				return fmt.Errorf("Ошибка авторизации: %w", err)
			}
		}

		// Получаем информацию о текущем пользователе
		users, err := api.UsersGetUsers(ctx, []tg.InputUserClass{&tg.InputUserSelf{}})
		if err != nil {
			return fmt.Errorf("Ошибка получения информации о себе: %w", err)
		}
		user, ok := users[0].(*tg.User)
		if !ok {
			return fmt.Errorf("Ошибка приведения типа к *tg.User")
		}
		log.Printf("✅ Успешно авторизован как %s", user.Username)

		return nil
	})

	if err != nil {
		log.Fatalf("Ошибка работы клиента: %v", err)
	}
}

// Реализация UserAuthenticator
type authHandler struct {
	phone string
}

func (a *authHandler) Phone(ctx context.Context) (string, error) {
	return a.phone, nil
}

func (a *authHandler) Code(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите код из Telegram: ")
	code, _ := reader.ReadString('\n')
	return strings.TrimSpace(code), nil
}

func (a *authHandler) Password(ctx context.Context) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите пароль (если включена 2FA): ")
	password, _ := reader.ReadString('\n')
	return strings.TrimSpace(password), nil
}

func (a *authHandler) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, fmt.Errorf("регистрация нового пользователя не поддерживается")
}

func (a *authHandler) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	fmt.Println("Необходимо принять условия использования:", tos.Text)
	return fmt.Errorf("пользователь должен вручную принять условия")
}

// Функция авторизации через номер телефона
func login(ctx context.Context, authClient *auth.Client) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите ваш номер телефона: ")
	phone, _ := reader.ReadString('\n')
	phone = strings.TrimSpace(phone)

	flow := auth.NewFlow(&authHandler{phone: phone}, auth.SendCodeOptions{})

	// Используем `IfNecessary`, который корректно выполняет вход
	if err := authClient.IfNecessary(ctx, flow); err != nil {
		return fmt.Errorf("Ошибка авторизации: %w", err)
	}

	return nil
}

