package redis

import (
	"context"
	"github.com/GabrielHCataldo/go-error-detail/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/redis/go-redis/v9"
	"go-redis-template/redis/option"
	"strconv"
)

type template struct {
	client *redis.Client
}

type Template interface {
	Set(ctx context.Context, key, value any, opts ...option.Set) error
	Get(ctx context.Context, key, dest any) error
	Del(ctx context.Context, keys ...any) error
	Disconnect() error
	SimpleDisconnect()
}

func NewTemplate(opts option.Client) Template {
	client := redis.NewClient(opts.ParseToRedisOptions())
	return template{
		client: client,
	}
}

func (t template) Set(ctx context.Context, key, value any, opts ...option.Set) error {
	opt := option.GetOptionSetByParams(opts)
	sKey, err := helper.ConvertToString(key)
	if err != nil {
		return ErrConvertKey
	}
	sValue, err := helper.ConvertToString(value)
	if err != nil {
		return ErrConvertValue
	}
	return t.client.Set(ctx, sKey, sValue, opt.Expiration).Err()
}

func (t template) Get(ctx context.Context, key, dest any) error {
	if !helper.IsPointer(dest) {
		return ErrDestIsNotPointer
	}
	sKey, err := helper.ConvertToString(key)
	if err != nil {
		return ErrConvertKey
	}
	result, err := t.client.Get(ctx, sKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	} else if err != nil {
		return err
	}
	return helper.ConvertToDest(result, dest)
}

func (t template) Del(ctx context.Context, keys ...any) error {
	var sKeys []string
	for i, key := range keys {
		sKey, err := helper.ConvertToString(key)
		if err != nil {
			return errors.New(ErrConvertKey, "index:"+strconv.Itoa(i))
		}
		sKeys = append(sKeys, sKey)
	}
	return t.client.Del(ctx, sKeys...).Err()
}

func (t template) Disconnect() error {
	return t.client.Close()
}

func (t template) SimpleDisconnect() {
	err := t.client.Close()
	if err != nil {
		logger.Error("error disconnect redis:", err)
		return
	}
	logger.Info("connection to redis closed.")
}
