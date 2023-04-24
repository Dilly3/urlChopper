package repository

import (
	"log"
	"os"
	"strconv"

	"github.com/dilly3/urlshortner/internal"
	"github.com/dilly3/urlshortner/repository/mongo"
	"github.com/dilly3/urlshortner/repository/postgres"
	"github.com/dilly3/urlshortner/repository/redis"
)

func ChooseRepo() internal.RedirectRepositoryPort {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := redis.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mongo.NewMongoRepository(mongoURL, mongoTimeout, mongodb)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "postgres":
		port := os.Getenv("POSTGRES_PORT")
		pgDbName := os.Getenv("POSTGRES_DBNAME")
		password := os.Getenv("POSTGRES_PASSWORD")
		user := os.Getenv("POSTGRES_USER")

		repo, err := postgres.NewPostgreRepository(user, pgDbName, password, port)
		if err != nil {
			log.Fatal(err)
		}
		return repo

	}
	return nil
}
