package redis

import (
	"context"
	"github.com/GabrielHCataldo/go-error-detail/errors"
	"github.com/GabrielHCataldo/go-helper/helper"
	"github.com/GabrielHCataldo/go-logger/logger"
	"github.com/GabrielHCataldo/go-redis-template/redis/option"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type MSetInput struct {
	// Key can be of any type, but cannot be null, and must be compatible with conversion to string (helper.ConvertToString).
	Key any
	// Value can be of any type, but cannot be null, and must be compatible with conversion to string (helper.ConvertToString).
	Value any
	// Opt to customize the operation (not required)
	Opt option.Set
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
	// Set supports all options that the SET command supports.
	//
	// The key and value parameters can be of any type, but cannot be nil, if an error occurs when converting the key
	// or value, the error returned is ErrConvertKey or ErrConvertValue respectively.
	//
	// If the return is nil, the operation was carried out successfully, otherwise an error occurred in the operation.
	//
	// To customize the operation, use the opts parameter (option.Set).
	Set(ctx context.Context, key, value any, opts ...option.Set) error
	// MSet defines N values. (Multiple Set)
	//
	// Parameter values cannot be empty, and must follow the Set function documentation for each MSetInput.
	//
	// The return will have a list of MSetOutput with each key and the error // that occurred, if the MSetOutput.Err
	// field is nil, it means that the operation for that key (MSetOutput.Key) was carried out successfully, otherwise
	// it failed .
	MSet(ctx context.Context, values ...MSetInput) []MSetOutput
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
	SetGet(ctx context.Context, key, value, dest any, opts ...option.Set) error
	// Rename redis key.
	//
	// The key and newKey parameters can be of any type, but cannot be null, in case an error occurs when converting
	// the key or newKey, the error returned is ErrConvertKey or ErrConvertNewKey respectively.
	//
	// If the return is null, the operation was performed successfully, otherwise an error occurred in the operation.
	Rename(ctx context.Context, key, newKey any) error
	// Get redis `GET key` command.
	//
	// The key parameter can be of any type, but cannot be null, in case an error occurs when converting, the error
	// returned is ErrConvertKey.
	//
	// The dest parameter must be a pointer, not null, if we do not find any value, dest will not have any modification.
	//
	// If the return is null, the operation was performed successfully, otherwise an error occurred in the operation.
	Get(ctx context.Context, key, dest any) error
	// GetDel get and delete value by key.
	//
	// The key parameter can be of any type, but cannot be null, if an error occurs during the conversion, the error
	// returned is ErrConvertKey. If no registered key is found, the error ErrKeyNotFound is returned.
	//
	// The dest parameter must be a pointer, not null, if we do not find any value, dest will not have any modification.
	//
	// If the return is null, the operation was performed successfully, otherwise an error occurred in the operation.
	GetDel(ctx context.Context, key, dest any) error
	// Exists redis values by key.
	//
	// The key parameter can be of any type, but cannot be null, if an error occurs during the conversion, the error
	// returned is ErrConvertKey.
	//
	// The return if true means that the key exists, otherwise it returns false, and if an error occurs in the operation
	// we return false with the second return parameter filled in
	Exists(ctx context.Context, key any) (bool, error)
	// Keys return list of keys by pattern.
	Keys(ctx context.Context, pattern string) ([]string, error)
	// Scan return list keys pageable by match
	Scan(ctx context.Context, cursor uint64, match string, count int64) ScanOutput
	// Del delete redis keys.
	//
	// The keys parameter can be of any type, but cannot be empty, if an error occurs during the conversion, the error
	// returned is ErrConvertKey.
	//
	// If the return is null, the operation was performed successfully, otherwise an error occurred in the operation.
	Del(ctx context.Context, keys ...any) error
	// Disconnect close connection to redis
	Disconnect() error
	// SimpleDisconnect close connection to redis without error
	SimpleDisconnect()
}

// NewTemplate create a new template instance
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

func (t template) GetDel(ctx context.Context, key, dest any) error {
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

func (t template) Exists(ctx context.Context, key any) (bool, error) {
	sKey, err := helper.ConvertToString(key)
	if err != nil {
		return false, ErrConvertKey
	}
	result := t.client.Exists(ctx, sKey)
	return result.Val() > 0, result.Err()
}

func (t template) Keys(ctx context.Context, pattern string) ([]string, error) {
	return t.client.Keys(ctx, pattern).Result()
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
