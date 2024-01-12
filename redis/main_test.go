package redis

import (
	"context"
	"github.com/GabrielHCataldo/go-redis-template/redis/option"
	"os"
	"time"
)

const redisKeyDefault = "test-key"
const redisDurationDefault = 5 * time.Minute

var redisTemplate Template

type testStruct struct {
	Name      string
	BirthDate time.Time
	Emails    []string
	Balance   float64
}

type testSet struct {
	name    string
	key     any
	value   any
	opt     option.Set
	wantErr bool
}

type testMSet struct {
	name   string
	values []MSetInput
}

type testSetGet struct {
	name    string
	key     any
	value   any
	dest    any
	opt     option.Set
	wantErr bool
}

type testRename struct {
	name    string
	key     any
	newKey  any
	wantErr bool
}

type testGet struct {
	name          string
	key           any
	dest          any
	wantErr       bool
	wantTimoutErr bool
}

type testDel struct {
	name    string
	keys    []any
	wantErr bool
}

type testKeys struct {
	name   string
	patten string
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
	_ = redisTemplate.Set(ctx, redisKeyDefault, initTestStruct(), option.NewSet().SetTTL(redisDurationDefault))
}

func initListTestSet() []testSet {
	return []testSet{
		{
			name:  "success",
			key:   redisKeyDefault,
			value: initTestStruct(),
			opt:   initOptionSet(),
		},
		{
			name:  "failed opts",
			key:   redisKeyDefault,
			value: initTestStruct(),
			opt: initOptionSet().
				SetMode(option.SetModeNx).
				SetMode(option.SetModeXx).
				SetKeepTTL(true).
				SetExpireAt(time.Now()),
			wantErr: true,
		},
		{
			name:    "failed key",
			key:     nil,
			value:   initTestStruct(),
			wantErr: true,
		},
		{
			name:    "failed value",
			key:     redisKeyDefault,
			value:   nil,
			wantErr: true,
		},
	}
}

func initListTestMSet() []testMSet {
	return []testMSet{
		{
			name:   "success",
			values: initMSetInputs(),
		},
		{
			name:   "failed empty",
			values: []MSetInput{},
		},
	}
}

func initListTestSetGet() []testSetGet {
	return []testSetGet{
		{
			name:  "success",
			key:   redisKeyDefault,
			value: initTestStruct(),
			dest:  &testStruct{},
			opt:   initOptionSet(),
		},
		{
			name:    "failed dest not pointer",
			key:     redisKeyDefault,
			value:   initTestStruct(),
			dest:    testStruct{},
			opt:     initOptionSet().SetKeepTTL(true),
			wantErr: true,
		},
		{
			name:    "failed key",
			key:     nil,
			value:   initTestStruct(),
			wantErr: true,
		},
		{
			name:    "failed value",
			key:     redisKeyDefault,
			value:   nil,
			wantErr: true,
		},
	}
}

func initListTestRename() []testRename {
	return []testRename{
		{
			name:   "success",
			key:    redisKeyDefault,
			newKey: "test-rename",
		},
		{
			name:    "failed not exists",
			key:     "",
			newKey:  "test-rename",
			wantErr: true,
		},
		{
			name:    "failed key",
			key:     nil,
			newKey:  "test-rename",
			wantErr: true,
		},
		{
			name:    "failed new key",
			key:     "test",
			newKey:  nil,
			wantErr: true,
		},
	}
}

func initListTestGet() []testGet {
	return []testGet{
		{
			name: "success",
			key:  redisKeyDefault,
			dest: &testStruct{},
		},
		{
			name: "success empty",
			key:  "",
			dest: &testStruct{},
		},
		{
			name:    "failed key",
			key:     nil,
			dest:    &testStruct{},
			wantErr: true,
		},
		{
			name:    "failed key not pointer",
			key:     "test",
			dest:    testStruct{},
			wantErr: true,
		},
	}
}

func initListTestGetDel() []testGet {
	return []testGet{
		{
			name: "success",
			key:  redisKeyDefault,
			dest: &testStruct{},
		},
		{
			name:          "failed timout",
			key:           redisKeyDefault,
			dest:          &testStruct{},
			wantTimoutErr: true,
			wantErr:       true,
		},
		{
			name:    "failed not exists",
			key:     "",
			dest:    &testStruct{},
			wantErr: true,
		},
		{
			name:    "failed key",
			key:     nil,
			dest:    &testStruct{},
			wantErr: true,
		},
		{
			name:    "failed key not pointer",
			key:     "test",
			dest:    testStruct{},
			wantErr: true,
		},
	}
}

func initListTestDel() []testDel {
	return []testDel{
		{
			name: "success",
			keys: []any{redisKeyDefault},
		},
		{
			name:    "failed keys",
			keys:    nil,
			wantErr: true,
		},
		{
			name:    "failed keys empty",
			keys:    []any{},
			wantErr: true,
		},
		{
			name:    "failed key",
			keys:    []any{nil},
			wantErr: true,
		},
	}
}

func initListTestExists() []testGet {
	return []testGet{
		{
			name: "success",
			key:  redisKeyDefault,
			dest: &testStruct{},
		},
		{
			name:    "failed key",
			key:     nil,
			dest:    &testStruct{},
			wantErr: true,
		},
	}
}

func initListTestKeys() []testKeys {
	return []testKeys{
		{
			name:   "success",
			patten: redisKeyDefault,
		},
		{
			name:   "success empty",
			patten: "",
		},
	}
}

func initMSetInputs() []MSetInput {
	return []MSetInput{
		{
			Key:   redisKeyDefault,
			Value: initTestStruct(),
			Opt:   initOptionSet(),
		},
		{
			Key:   "test-1",
			Value: initTestStruct(),
			Opt:   initOptionSet(),
		},
		{
			Key:   "test-2",
			Value: initTestStruct(),
			Opt:   initOptionSet().SetKeepTTL(true),
		},
		{
			Value: initTestStruct(),
			Opt:   initOptionSet().SetKeepTTL(true),
		},
		{
			Opt: initOptionSet().SetKeepTTL(true),
		},
	}
}

func initOptionSet() option.Set {
	return option.NewSet().
		SetMode(option.SetModeDefault).
		SetTTL(redisDurationDefault).
		SetKeepTTL(false)
}
