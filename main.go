package main

import (
	"PriemBot/internal/di"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Инициализация контейнера
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}
	defer container.Close()

	// Запуск бота
	container.Bot.Start()

	// Ожидание сигнала завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}
