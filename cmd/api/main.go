package main

import (
	"flag"
	"log"

	"go-api/internal/server"
	"go-api/pkg/config"
	"go-api/pkg/db/postgres"
	"go-api/pkg/db/redis"
)

var (
	configName *string
	configPath *string
)

func init() {
	configName = flag.String("cfg-name", "local", "")
	configPath = flag.String("cfg-path", ".", "")
	flag.Parse()
}

func main() {
	cfgFile, err := config.LoadConfig(*configName, *configPath)
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	pqDB, err := postgres.NewDB(cfg)
	if err != nil {
		log.Fatalf("Postgresql init: %s", err)
	} else {
		log.Printf("Postgres connected, Status: %#v", pqDB.Stats())
	}
	defer pqDB.Close()

	redisClient := redis.NewClient(cfg)
	defer redisClient.Close()
	log.Printf("Redis connected")

	s := server.NewServer(cfg, pqDB, redisClient)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
