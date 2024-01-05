package redis

import (
	"context"
	"go-redis-template/redis/option"
	"os"
	"time"
)

const redisKeyDefault = "test-key"

var redisTemplate Template

type testStruct struct {
	Name      string
	BirthDate time.Time
	Emails    []string
	Balance   float64
}

func initTemplate() {
	redisTemplate = NewTemplate(option.Client{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
}

func initTestStruct() testStruct {
	return testStruct{
		Name:      "Foo Bar",
		BirthDate: time.Now(),
		Emails:    []string{"foobar@gmail.com", "foobar2@hotmail.com"},
		Balance:   231.123,
	}
}

func initSet() {
	initTemplate()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	_ = redisTemplate.Set(ctx, redisKeyDefault, initTestStruct(), option.NewSet().SetExpiration(5*time.Minute))
}
