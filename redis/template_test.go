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
