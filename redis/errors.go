package redis

import (
	"github.com/GabrielHCataldo/go-errors/errors"
)

var MsgErrConvertKey = "redis: error convert key to string"
var MsgErrConvertKeyIndex = "redis: error convert key to string on index:"
var MsgErrConvertNewKey = "redis: error convert new key to string"
var MsgErrConvertValue = "redis: error convert value"
var MsgErrDestIsNotPointer = "redis: dest is not pointer"
var MsgErrKeyNotFound = "redis: key not exists"

var ErrConvertKey = errors.New(MsgErrConvertKey)
var ErrConvertKeyIndex = errors.New(MsgErrConvertKeyIndex)
var ErrConvertNewKey = errors.New(MsgErrConvertNewKey)
var ErrConvertValue = errors.New(MsgErrConvertValue)
var ErrDestIsNotPointer = errors.New(MsgErrDestIsNotPointer)
var ErrKeyNotFound = errors.New(MsgErrKeyNotFound)

func errConvertKey(skip int) error {
	ErrConvertKey = errors.NewSkipCaller(skip+1, MsgErrConvertKey)
	return ErrConvertKey
}

func errConvertKeyIndex(skip, index int) error {
	ErrConvertKeyIndex = errors.NewSkipCaller(skip+1, MsgErrConvertKeyIndex, index)
	return ErrConvertKeyIndex
}

func errConvertNewKey(skip int) error {
	ErrConvertNewKey = errors.NewSkipCaller(skip+1, MsgErrConvertNewKey)
	return ErrConvertNewKey
}

func errConvertValue(skip int) error {
	ErrConvertValue = errors.NewSkipCaller(skip+1, MsgErrConvertValue)
	return ErrConvertValue
}

func errDestIsNotPointer(skip int) error {
	ErrDestIsNotPointer = errors.NewSkipCaller(skip+1, MsgErrDestIsNotPointer)
	return ErrDestIsNotPointer
}

func errKeyNotFound(skip int) error {
	ErrKeyNotFound = errors.NewSkipCaller(skip+1, MsgErrKeyNotFound)
	return ErrKeyNotFound
}
