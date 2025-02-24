package auth

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
)

// AuthHandler реализует auth.UserAuthenticator
type AuthHandler struct {
	phone string
}

func (a *AuthHandler) Phone(ctx context.Context) (string, error) {
	return a.phone, nil
}

func (a *AuthHandler) Code(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите код из Telegram: ")
	code, _ := reader.ReadString('\n')
	return strings.TrimSpace(code), nil
}

func (a *AuthHandler) Password(ctx context.Context) (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите пароль (если включена 2FA): ")
	password, _ := reader.ReadString('\n')
	return strings.TrimSpace(password), nil
}

func (a *AuthHandler) SignUp(ctx context.Context) (auth.UserInfo, error) {
	return auth.UserInfo{}, fmt.Errorf("регистрация нового пользователя не поддерживается")
}

func (a *AuthHandler) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	fmt.Println("Необходимо принять условия использования:", tos.Text)
	return fmt.Errorf("пользователь должен вручную принять условия")
}

// Login выполняет авторизацию через номер телефона
func Login(ctx context.Context, authClient *auth.Client) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите ваш номер телефона: ")
	phone, _ := reader.ReadString('\n')
	phone = strings.TrimSpace(phone)

	flow := auth.NewFlow(&AuthHandler{phone: phone}, auth.SendCodeOptions{})

	// Используем `IfNecessary`, который корректно выполняет вход
	if err := authClient.IfNecessary(ctx, flow); err != nil {
		return fmt.Errorf("Ошибка авторизации: %w", err)
	}

	return nil
}

