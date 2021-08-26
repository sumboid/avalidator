package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	_ "github.com/joho/godotenv/autoload"
	"github.com/machinebox/graphql"
	log "github.com/sirupsen/logrus"
)

var Config *ConfigModel

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	var err error
	Config, err = CreateConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     Config.Redis.Host + ":" + Config.Redis.Port,
		Password: Config.Redis.Password,
		DB:       Config.Redis.DB,
	})
	defer redisClient.Close()

	_, err := redisClient.Ping().Result()
	if err != nil {
		panic(err)
	}

	gqlClient := graphql.NewClient(Config.GraphQL.URL)

	jwtService := NewJWTService(redisClient)

	router := gin.Default()

	ConfigureAuthEndpoint(router.Group("/"), jwtService, gqlClient)

	addr := "0.0.0.0:" + Config.Port

	log.WithField("addr", addr).Info("Server has been started")

	router.Use(gin.Recovery())
	router.Run(addr)
}
