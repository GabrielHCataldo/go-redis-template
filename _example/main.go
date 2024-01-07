package main

import (
	"context"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"go-redis-template/redis"
	"go-redis-template/redis/option"
	"os"
	"time"
)

type testStruct struct {
	Name      string
	BirthDate time.Time
}

func main() {
	set()
	get()
	exists()
	del()
}

func set() {
	key := "example-struct"
	value := testStruct{
		Name:      "Foo bar",
		BirthDate: time.Now(),
	}
	redisTemplate := redis.NewTemplate(option.Client{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	defer redisTemplate.SimpleDisconnect()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	opt := option.NewSet()
	opt.SetMode(option.SetModeDefault)
	opt.SetTTL(5 * time.Minute)
	opt.SetExpireAt(time.Time{})
	opt.SetKeepTTL(false)
	err := redisTemplate.Set(ctx, key, value, opt)
	if helper.IsNotNil(err) {
		logger.Error("error set redis value:", err)
		return
	}
	logger.Info("set", key, "value redis completed successfully!")
}

func get() {
	key := "example-struct"
	var dest testStruct
	redisTemplate := redis.NewTemplate(option.Client{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	defer redisTemplate.SimpleDisconnect()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	err := redisTemplate.Get(ctx, key, &dest)
	if helper.IsNotNil(err) {
		logger.Error("error get redis value:", err)
		return
	}
	logger.Info("get", key, "value redis completed successfully!", dest)
}

func exists() {
	key := "example-struct"
	redisTemplate := redis.NewTemplate(option.Client{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	defer redisTemplate.SimpleDisconnect()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	alreadyExists, err := redisTemplate.Exists(ctx, key)
	if helper.IsNotNil(err) {
		logger.Error("error exists redis value:", err)
		return
	}
	logger.Info("already exists", key, "?", alreadyExists)
}

func del() {
	key1 := "example-struct"
	key2 := "example-struct-2"
	redisTemplate := redis.NewTemplate(option.Client{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	err := redisTemplate.Del(ctx, key1, key2)
	if helper.IsNotNil(err) {
		logger.Error("error delete redis keys:", err)
		return
	}
	logger.Info("delete redis keys (", key1, "-", key2, ") completed successfully!")
}
