package main

import (
	"agent/internal/app"
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	// Создание контекста с возможностью отмены
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализация приложения
	application := app.NewApp()

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
