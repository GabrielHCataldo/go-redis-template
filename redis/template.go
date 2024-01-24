package redis

import (
	"context"
	"fmt"
	"github.com/GabrielHCataldo/go-errors/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/go-redis-template/redis/option"
	"github.com/redis/go-redis/v9"
	"strings"
)

type MSetInput struct {
	// Key can be of any type, but cannot be null, and must be compatible with conversion to string (helper.ConvertToString).
	Key any
	// Value can be of any type, but cannot be null, and must be compatible with conversion to string (helper.ConvertToString).
	Value any
	// Opt to customize the operation (not required)
	Opt *option.Set
}

type MSetOutput struct {
	Key any
	Err error
}

type ScanOutput struct {
	Cursor uint64
	Page   []string
}

type Template struct {
	client *redis.Client
}

// NewTemplate create a new template instance
func NewTemplate(opts option.Client) *Template {
	client := redis.NewClient(opts.ParseToRedisOptions())
	return &Template{
		client: client,
	}
}

// Set supports all options that the SET command supports.
//
// The key and value parameters can be of any type, but cannot be nil, if an error occurs when converting the key
// or value, the error returned is ErrConvertKey or ErrConvertValue respectively.
//
// If the return is nil, the operation was carried out successfully, otherwise an error occurred in the operation.
//
// To customize the operation, use the opts parameter (option.Set).
func (t Template) Set(ctx context.Context, key, value any, opts ...*option.Set) error {
	result, err := t.set(ctx, key, value, false, opts...)
	if err == nil {
		err = errors.NewSkipCaller(2, result.Err())
	}
	return err
}

// MSet defines N values. (Multiple Set)
//
// Parameter values cannot be empty, and must follow the Set function documentation for each MSetInput.
//
// The return will have a list of MSetOutput with each key and the error // that occurred, if the MSetOutput.Err
// field is nil, it means that the operation for that key (MSetOutput.Key) was carried out successfully, otherwise
// it failed .
func (t Template) MSet(ctx context.Context, values ...MSetInput) []MSetOutput {
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

// SetGet supports all options that the SET command supports.
//
// The key and value parameters can be of any type, but cannot be null, in case an error occurs when converting the key
// or value, the error returned is ErrConvertKey or ErrConvertValue respectively.
//
// The dest parameter must be a pointer, not null, if we do not find a predecessor value to the set, dest will not
// have any modification
//
// If the return is null, the operation was performed successfully, otherwise an error occurred in the operation.
//
// To customize the operation, use the opts parameter (option.Set).
func (t Template) SetGet(ctx context.Context, key, value, dest any, opts ...*option.Set) error {
	result, err := t.set(ctx, key, value, true, opts...)
	if helper.IsNotNil(err) {
		return err
	} else if helper.IsNotNil(result.Err()) {
		return errors.NewSkipCaller(2, result.Err())
	}
	return errors.NewSkipCaller(2, helper.ConvertToDest(result.Val(), dest))
}

// Rename redis key.
//
// The key and newKey parameters can be of any type, but cannot be null, in case an error occurs when converting
// the key or newKey, the error returned is ErrConvertKey or ErrConvertNewKey respectively.
//
// If the return is null, the operation was performed successfully, otherwise an error occurred in the operation.
func (t Template) Rename(ctx context.Context, key, newKey any) error {
	sKey, err := helper.ConvertToString(key)
	if helper.IsNotNil(err) {
		return errConvertKey(2)
	}
	sNewKey, err := helper.ConvertToString(newKey)
	if helper.IsNotNil(err) {
		return errConvertNewKey(2)
	}
	return errors.NewSkipCaller(2, t.client.Rename(ctx, sKey, sNewKey).Err())
}

// Get redis `GET key` command.
//
// The key parameter can be of any type, but cannot be null, in case an error occurs when converting, the error
// returned is ErrConvertKey.
//
// The dest parameter must be a pointer, not null, if we do not find any value, dest will not have any modification.
//
// If the return is null, the operation was performed successfully, otherwise an error occurred in the operation.
func (t Template) Get(ctx context.Context, key, dest any) error {
	if !helper.IsPointer(dest) {
		return errDestIsNotPointer(2)
	}
	sKey, err := helper.ConvertToString(key)
	if helper.IsNotNil(err) {
		return errConvertKey(2)
	}
	result, err := t.client.Get(ctx, sKey).Result()
	if errors.Is(err, redis.Nil) {
		return nil
	}
	return helper.ConvertToDest(result, dest)
}

// GetDel get and delete value by key.
//
// The key parameter can be of any type, but cannot be null, if an error occurs during the conversion, the error
// returned is ErrConvertKey. If no registered key is found, the error ErrKeyNotFound is returned.
//
// The dest parameter must be a pointer, not null, if we do not find any value, dest will not have any modification.
//
// If the return is null, the operation was performed successfully, otherwise an error occurred in the operation.
func (t Template) GetDel(ctx context.Context, key, dest any) error {
	sKey, err := helper.ConvertToString(key)
	if helper.IsNotNil(err) {
		return errConvertKey(2)
	}
	result := t.client.GetDel(ctx, sKey)
	if helper.IsNotNil(result.Err()) {
		err = errors.NewSkipCaller(2, result.Err())
		if errors.Is(result.Err(), redis.Nil) {
			err = errKeyNotFound(2)
		}
		return err
	}
	return errors.NewSkipCaller(2, helper.ConvertToDest(result.Val(), dest))
}

// Exists redis values by key.
//
// The key parameter can be of any type, but cannot be null, if an error occurs during the conversion, the error
// returned is ErrConvertKey.
//
// The return if true means that the key exists, otherwise it returns false, and if an error occurs in the operation
// we return false with the second return parameter filled in
func (t Template) Exists(ctx context.Context, key any) (bool, error) {
	sKey, err := helper.ConvertToString(key)
	if helper.IsNotNil(err) {
		return false, errConvertKey(2)
	}
	result := t.client.Exists(ctx, sKey)
	return result.Val() > 0, errors.NewSkipCaller(2, result.Err())
}

// Keys return list of keys by pattern.
func (t Template) Keys(ctx context.Context, pattern string) ([]string, error) {
	return t.client.Keys(ctx, pattern).Result()
}

// Scan return list keys pageable by match
func (t Template) Scan(ctx context.Context, cursor uint64, match string, count int64) ScanOutput {
	result := t.client.Scan(ctx, cursor, match, count)
	keys, c := result.Val()
	return ScanOutput{
		Cursor: c,
		Page:   keys,
	}
}

// Del delete redis keys.
//
// The keys parameter can be of any type, but cannot be empty, if an error occurs during the conversion, the error
// returned is ErrConvertKeyIndex.
//
// If the return is null, the operation was performed successfully, otherwise an error occurred in the operation.
func (t Template) Del(ctx context.Context, keys ...any) error {
	var sKeys []string
	for i, key := range keys {
		sKey, err := helper.ConvertToString(key)
		if helper.IsNotNil(err) {
			return errConvertKeyIndex(2, i)
		}
		sKeys = append(sKeys, sKey)
	}
	return errors.NewSkipCaller(2, t.client.Del(ctx, sKeys...).Err())
}

// SprintKey format values as prefix in string for a future redis key, ex: "test", "test2" -> "test:test2"
func (t Template) SprintKey(vs ...any) string {
	var builder strings.Builder
	for _, v := range vs {
		s, err := helper.ConvertToString(v)
		if helper.IsNil(err) {
			if helper.IsNotEmpty(builder.Len()) {
				builder.WriteString(fmt.Sprint(":", s))
			} else {
				builder.WriteString(s)
			}
		}
	}
	return builder.String()
}

// Disconnect close connection to redis
func (t Template) Disconnect() error {
	return errors.NewSkipCaller(2, t.client.Close())
}

// SimpleDisconnect close connection to redis without error
func (t Template) SimpleDisconnect() {
	err := t.client.Close()
	if helper.IsNotNil(err) {
		logger.ErrorSkipCaller(2, "Error disconnect:", err)
		return
	}
	logger.InfoSkipCaller(2, "Connection to redis closed.")
}

func (t Template) set(ctx context.Context, key, value any, get bool, opts ...*option.Set) (*redis.StatusCmd, error) {
	opt := option.GetOptionSetByParams(opts)
	sKey, err := helper.ConvertToString(key)
	if helper.IsNotNil(err) {
		return nil, errConvertKey(3)
	}
	sValue, err := helper.ConvertToString(value)
	if helper.IsNotNil(err) {
		return nil, errConvertValue(3)
	}
	return t.client.SetArgs(ctx, sKey, sValue, redis.SetArgs{
		Mode:     opt.Mode.String(),
		TTL:      opt.TTL,
		ExpireAt: opt.ExpireAt,
		Get:      get,
		KeepTTL:  opt.KeepTTL,
	}), nil
}
