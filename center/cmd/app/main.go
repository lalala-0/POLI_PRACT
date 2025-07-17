package main

import (
	"center/internal/app"
	"center/internal/config"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

///home/polina/Pract/POLI_PRACT/center

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config/config.yml" // путь по умолчанию
	}
	// Загрузка конфигурации из переменных окружения или флагов
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// // Валидация конфигурации БД
	// if err := cfg.Postgres.Validate(); err != nil {
	//     log.Fatalf("Invalid DB config: %v", err)
	// }

	// Создание контекста с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация приложения
	application := app.NewApp(cfg)

	// Запуск горутин с использованием WaitGroup
	var wg sync.WaitGroup

	// Запуск приложения (внутри создаются две горутины)
	application.Run(ctx, &wg)

	// Обработка сигналов завершения
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Получен сигнал завершения, остановка серверв...")
	cancel() // Отмена контекста приведет к завершению горутин

	// Ожидание завершения всех горутин
	wg.Wait()
	log.Println("Сервер завершил работу")
}
