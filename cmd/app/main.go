package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"tgmessenger/internal/app"
	"tgmessenger/internal/bot"
	"tgmessenger/internal/config"
)

func main() {
	log.Println("🚀 Начинаем запуск бота...")

	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("❌ Ошибка загрузки конфигурации: %v", err)
	}

	// Выводим информацию о конфигурации
	log.Printf("📋 Загружена конфигурация:")
	log.Printf("   - Количество сообщений: %d", cfg.MessagesPerDialog)
	log.Printf("   - Интервал обновления: %d сек", cfg.UpdateInterval)
	log.Printf("   - Форматы вывода: %v", cfg.OutputFormats)
	if cfg.Debug {
		log.Printf("   - Режим отладки: включен")
	}

	// Получаем клиента
	client, ctx, cancel := app.InitBot()
	defer cancel()

	// Настраиваем перехват сигналов для корректного завершения
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	log.Printf("🚀 Бот запущен! Проверка сообщений каждые %d секунд.", cfg.UpdateInterval)

	ticker := time.NewTicker(time.Duration(cfg.UpdateInterval) * time.Second)
	defer ticker.Stop()

	// Мьютекс для защиты от одновременных выполнений bot.Run
	var mu sync.Mutex
	// Флаг, показывающий, что операция в процессе выполнения
	isRunning := false

	// Запустим первую проверку сразу
	log.Println("🔄 Запускаем первую проверку сообщений...")
	go func() {
		mu.Lock()
		isRunning = true
		log.Println("🔄 Запускаем первый bot.Run...")
		if err := bot.Run(ctx, client, cfg); err != nil {
			log.Printf("⚠️ Ошибка первого запуска бота: %v", err)
		}
		log.Println("✅ Первый bot.Run завершился")
		isRunning = false
		mu.Unlock()
	}()

	for {
		select {
		case <-ctx.Done():
			log.Println("⏹️ Завершаем работу бота...")
			return
		case sig := <-sigs:
			log.Printf("⏹️ Получен сигнал %v, завершаем работу...", sig)
			cancel()
			return
		case <-ticker.C:
			// Проверяем, не выполняется ли уже команда
			mu.Lock()
			if isRunning {
				log.Println("⏸️ Предыдущая команда ещё выполняется, пропускаем запуск")
				mu.Unlock()
				continue
			}
			isRunning = true
			mu.Unlock()

			// Создаем новый контекст для каждого запуска
			runCtx, runCancel := context.WithTimeout(ctx, time.Duration(cfg.UpdateInterval-5)*time.Second)

			log.Println("🔄 Запускаем bot.Run...")
			go func() {
				defer runCancel() // Отмена контекста при завершении
				log.Println("🔄 Запуск проверки сообщений...")
				if err := bot.Run(runCtx, client, cfg); err != nil {
					log.Printf("⚠️ Ошибка работы бота: %v", err)
				}
				log.Println("✅ bot.Run завершился")

				mu.Lock()
				isRunning = false
				mu.Unlock()
			}()
		}
	}
}
