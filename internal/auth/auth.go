package auth

import (
	"bufio"
	"context"
	"fmt"
	"log"
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
	// Убедимся, что телефон в международном формате
	phone := a.phone
	if !strings.HasPrefix(phone, "+") {
		// Для России добавляем +7
		if strings.HasPrefix(phone, "8") {
			phone = "+7" + phone[1:]
		} else if strings.HasPrefix(phone, "7") {
			phone = "+" + phone
		} else {
			phone = "+7" + phone
		}
	}

	log.Printf("📱 Используем номер телефона: %s", phone)
	return phone, nil
}

func (a *AuthHandler) Code(ctx context.Context, sentCode *tg.AuthSentCode) (string, error) {
	log.Printf("📱 Код отправлен через %T", sentCode.Type)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите код из Telegram: ")
	code, _ := reader.ReadString('\n')
	code = strings.TrimSpace(code)
	log.Printf("📱 Получен код: %s", code)
	return code, nil
}

func (a *AuthHandler) Password(ctx context.Context) (string, error) {
	log.Println("🔒 Запрошен пароль 2FA")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите пароль (если включена 2FA): ")
	password, _ := reader.ReadString('\n')
	return strings.TrimSpace(password), nil
}

func (a *AuthHandler) SignUp(ctx context.Context) (auth.UserInfo, error) {
	log.Println("⚠️ Запрошена регистрация нового пользователя")
	return auth.UserInfo{}, fmt.Errorf("регистрация нового пользователя не поддерживается")
}

func (a *AuthHandler) AcceptTermsOfService(ctx context.Context, tos tg.HelpTermsOfService) error {
	log.Println("⚠️ Запрошено принятие условий использования")
	fmt.Println("Необходимо принять условия использования:", tos.Text)
	return fmt.Errorf("пользователь должен вручную принять условия")
}

// Login выполняет авторизацию через номер телефона
func Login(ctx context.Context, authClient *auth.Client) error {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите ваш номер телефона: ")
	phone, _ := reader.ReadString('\n')
	phone = strings.TrimSpace(phone)

	log.Printf("📱 Начинаем авторизацию для номера: %s", phone)

	// Создаем обработчик авторизации
	handler := &AuthHandler{phone: phone}
	flow := auth.NewFlow(handler, auth.SendCodeOptions{})

	log.Println("📱 Запускаем процесс авторизации...")
	// Используем `IfNecessary`, который корректно выполняет вход
	if err := authClient.IfNecessary(ctx, flow); err != nil {
		log.Printf("⚠️ Ошибка авторизации: %v", err)
		return fmt.Errorf("Ошибка авторизации: %w", err)
	}

	log.Println("✅ Авторизация завершена успешно")
	return nil
}
