package redis

import (
	"context"
	"github.com/GabrielHCataldo/go-logger/logger"
	"go-redis-template/redis/option"
	"testing"
	"time"
)

func TestTemplateSet(t *testing.T) {
	initTemplate()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	err := redisTemplate.Set(ctx, redisKeyDefault, initTestStruct(), option.NewSet().SetExpiration(5*time.Minute))
	if err != nil {
		logger.Error("error test set:", err)
		t.Fail()
	} else {
		logger.Info("test completed successfully!")
	}
}

func TestTemplateGet(t *testing.T) {
	initSet()
	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()
	var dest testStruct
	err := redisTemplate.Get(ctx, redisKeyDefault, &dest)
	if err != nil {
		logger.Error("error test get:", err)
		t.Fail()
	} else {
		logger.Info("test completed successfully! ", dest)
	}
}
