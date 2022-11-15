package main

import (
	"context"
	"log"
	"task-scheduler/config"
	"task-scheduler/repository"
)

func main() {
	log.Println("Task scheduler")
	config, err := config.InitConfig(".env", "config/data.json")
	if err != nil {
		log.Fatalf("Error while loading config %s", err)
	}
	log.Println("Config loaded")
	ctx := context.Background()
	_ = repository.NewRedisStorage(config.Redis, ctx)

}
