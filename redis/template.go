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

type MSetInput struct {
	Key   any
	Value any
	Opt   option.Set
}

type MSetOutput struct {
	Key any
	Err error
}

type ScanOutput struct {
	Cursor uint64
	Page   []string
}

type template struct {
	client *redis.Client
}

type Template interface {
	// Set redis value, custom options command set using option.Set
	Set(ctx context.Context, key, value any, opts ...option.Set) error
	// MSet multiple set redis value, custom options command set MSetInput.Opt using option.Set
	MSet(ctx context.Context, values ...MSetInput) []MSetOutput
	// SetGet set redis value and fill dest with old value
	SetGet(ctx context.Context, key, value, dest any, opts ...option.Set) error
	// Rename redis key
	Rename(ctx context.Context, key, newKey any) error
	// Get redis value by key, dest need a pointer
	Get(ctx context.Context, key, dest any) error
	// GetDel delete redis keys and fill dest with values
	GetDel(ctx context.Context, dest any, keys any) error
	// Count redis values by key
	Count(ctx context.Context, key any) (int64, error)
	// Exists redis values by key
	Exists(ctx context.Context, key any) (bool, error)
	// Keys return list of keys by pattern
	Keys(ctx context.Context, pattern string) ([]string, error)
	// Del delete redis keys
	Del(ctx context.Context, keys ...any) error
	// Scan return list keys pageable by match
	Scan(ctx context.Context, cursor uint64, match string, count int64) ScanOutput
	// Disconnect close connection to redis
	Disconnect() error
	// SimpleDisconnect close connection to redis without error
	SimpleDisconnect()
}

func NewTemplate(opts option.Client) Template {
	client := redis.NewClient(opts.ParseToRedisOptions())
	return template{
		client: client,
	}
}

func (t template) Set(ctx context.Context, key, value any, opts ...option.Set) error {
	result, err := t.set(ctx, key, value, false, opts...)
	if err == nil {
		err = result.Err()
	}
	return err
}

func (t template) MSet(ctx context.Context, values ...MSetInput) []MSetOutput {
	var output []MSetOutput
	for _, v := range values {
		err := t.Set(ctx, v.Key, v.Value, v.Opt)
		output = append(output, MSetOutput{
			Key: v.Key,
			Err: err,
		})
	}
	return output
}

func (t template) SetGet(ctx context.Context, key, value, dest any, opts ...option.Set) error {
	result, err := t.set(ctx, key, value, true, opts...)
	if err != nil {
		return err
	} else if result.Err() != nil {
		return result.Err()
	}
	return helper.ConvertToDest(result.Val(), dest)
}

func (t template) Rename(ctx context.Context, key, newKey any) error {
	sKey, err := helper.ConvertToString(key)
	if err != nil {
		return ErrConvertKey
	}
	sNewKey, err := helper.ConvertToString(newKey)
	if err != nil {
		return ErrConvertNewKey
	}
	return t.client.Rename(ctx, sKey, sNewKey).Err()
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
	}
	return helper.ConvertToDest(result, dest)
}

func (t template) GetDel(ctx context.Context, dest any, key any) error {
	sKey, err := helper.ConvertToString(key)
	if err != nil {
		return ErrConvertKey
	}
	result := t.client.GetDel(ctx, sKey)
	if result.Err() != nil {
		if errors.Is(result.Err(), redis.Nil) {
			return ErrKeyNotFound
		}
		return result.Err()
	}
	return helper.ConvertToDest(result.Val(), dest)
}

func (t template) Count(ctx context.Context, key any) (int64, error) {
	sKey, err := helper.ConvertToString(key)
	if err != nil {
		return 0, ErrConvertKey
	}
	result := t.client.Exists(ctx, sKey)
	if result.Err() != nil && !errors.Is(result.Err(), redis.Nil) {
		return 0, result.Err()
	}
	return result.Val(), nil
}

func (t template) Exists(ctx context.Context, key any) (bool, error) {
	count, err := t.Count(ctx, key)
	return count > 0, err
}

func (t template) Keys(ctx context.Context, pattern string) ([]string, error) {
	result := t.client.Keys(ctx, pattern)
	return result.Val(), result.Err()
}

func (t template) Scan(ctx context.Context, cursor uint64, match string, count int64) ScanOutput {
	result := t.client.Scan(ctx, cursor, match, count)
	keys, c := result.Val()
	return ScanOutput{
		Cursor: c,
		Page:   keys,
	}
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

func (t template) set(ctx context.Context, key, value any, get bool, opts ...option.Set) (*redis.StatusCmd, error) {
	opt := option.GetOptionSetByParams(opts)
	sKey, err := helper.ConvertToString(key)
	if err != nil {
		return nil, ErrConvertKey
	}
	sValue, err := helper.ConvertToString(value)
	if err != nil {
		return nil, ErrConvertValue
	}
	return t.client.SetArgs(ctx, sKey, sValue, redis.SetArgs{
		Mode:     opt.Mode.String(),
		TTL:      opt.TTL,
		ExpireAt: opt.ExpireAt,
		Get:      get,
		KeepTTL:  opt.KeepTTL,
	}), nil
}
