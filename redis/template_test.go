package redis

import (
	"context"
	"github.com/GabrielHCataldo/go-logger/logger"
	"testing"
	"time"
)

func TestTemplateSet(t *testing.T) {
	initTemplate()
	for _, tt := range initListTestSet() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			err := redisTemplate.Set(ctx, tt.key, tt.value, tt.opt)
			if (err != nil) != tt.wantErr {
				logger.Errorf("Set() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
			}
		})
	}
}

func TestTemplateMSet(t *testing.T) {
	initTemplate()
	for _, tt := range initListTestMSet() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			result := redisTemplate.MSet(ctx, tt.values...)
			logger.Infof("MSet() result = %v", result)
		})
	}
}

func TestTemplateSetGet(t *testing.T) {
	initTemplate()
	for _, tt := range initListTestSetGet() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			err := redisTemplate.SetGet(ctx, tt.key, tt.value, tt.dest, tt.opt)
			if (err != nil) != tt.wantErr {
				logger.Errorf("SetGet() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
				return
			}
			logger.Infof("SetGet() result = %v", tt.dest)
		})
	}
}

func TestTemplateRename(t *testing.T) {
	initSet()
	for _, tt := range initListTestRename() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			err := redisTemplate.Rename(ctx, tt.key, tt.newKey)
			if (err != nil) != tt.wantErr {
				logger.Errorf("Rename() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
			}
			logger.Infof("Rename() err = %v", err)

		})
	}
}

func TestTemplateGet(t *testing.T) {
	initSet()
	for _, tt := range initListTestGet() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			err := redisTemplate.Get(ctx, tt.key, tt.dest)
			if (err != nil) != tt.wantErr {
				logger.Errorf("Get() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
				return
			}
			logger.Infof("Get() result = %v err = %v", tt.dest, err)
		})
	}
}

func TestTemplateGetDel(t *testing.T) {
	initSet()
	for _, tt := range initListTestGetDel() {
		t.Run(tt.name, func(t *testing.T) {
			d := 5 * time.Second
			if tt.wantTimoutErr {
				d = 1 * time.Nanosecond
			}
			ctx, cancel := context.WithTimeout(context.TODO(), d)
			defer cancel()
			err := redisTemplate.GetDel(ctx, tt.key, tt.dest)
			if (err != nil) != tt.wantErr {
				logger.Errorf("GetDel() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
				return
			}
			logger.Infof("GetDel() result = %v err = %v", tt.dest, err)
		})
	}
}

func TestTemplateExists(t *testing.T) {
	initSet()
	for _, tt := range initListTestExists() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			result, err := redisTemplate.Exists(ctx, tt.key)
			if (err != nil) != tt.wantErr {
				logger.Errorf("Exists() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
				return
			}
			logger.Infof("Exists() result = %v err = %v", result, err)
		})
	}
}

func TestTemplateKeys(t *testing.T) {
	initSet()
	for _, tt := range initListTestKeys() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			result, err := redisTemplate.Keys(ctx, tt.patten)
			logger.Infof("Keys() result = %v err = %v", result, err)
		})
	}
}

func TestTemplateScan(t *testing.T) {
	initSet()
	for _, tt := range initListTestKeys() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			result := redisTemplate.Scan(ctx, 0, tt.patten, 10)
			logger.Infof("Scan() result = %v", result)
		})
	}
}

func TestTemplateDel(t *testing.T) {
	initSet()
	for _, tt := range initListTestDel() {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
			defer cancel()
			err := redisTemplate.Del(ctx, tt.keys...)
			if (err != nil) != tt.wantErr {
				logger.Errorf("Del() err = %v, wantErr = %v", err, tt.wantErr)
				t.Fail()
			}
		})
	}
}

func TestTemplateSprintKey(t *testing.T) {
	initTemplate()
	for _, tt := range initListTestSprintKey() {
		t.Run(tt.name, func(t *testing.T) {
			result := redisTemplate.SprintKey(tt.values...)
			logger.Info("result spring key:", result)
		})
	}
}

func TestTemplateDisconnect(t *testing.T) {
	initTemplate()
	err := redisTemplate.Disconnect()
	logger.Infof("Del() err = %v", err)
}

func TestTemplateSimpleDisconnect(t *testing.T) {
	initTemplate()
	redisTemplate.SimpleDisconnect()
	redisTemplate.SimpleDisconnect()
}
