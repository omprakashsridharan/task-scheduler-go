package main

import (
	"context"
	"log"
	"task-scheduler/config"
	"task-scheduler/repository"
)

func main() {
	log.Println("Task scheduler")
	initConfig, err := config.InitConfig(".env", "initConfig/data.json")
	if err != nil {
		log.Fatalf("Error while loading initConfig %s", err)
	}
	log.Println("Config loaded")
	ctx := context.Background()
	_ = repository.NewRedisStorage(initConfig.Redis, ctx)

}
