package main

import (
	"agent/internal/app"
	"agent/internal/config"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// Загрузка конфигурации из файла
	configPath := flag.String("config", "./config/config.yml", "path to config file")
	flag.Parse()
	cfg, err := config.LoadAgentConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

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

	log.Println("Получен сигнал завершения, остановка агента...")
	cancel() // Отмена контекста приведет к завершению горутин

	// Ожидание завершения всех горутин
	wg.Wait()
	log.Println("Агент завершил работу")
}
