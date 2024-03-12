package redis

import (
	"errors"
)

var MsgErrConvertKey = "redis: error convert key to string"
var MsgErrConvertNewKey = "redis: error convert new key to string"
var MsgErrConvertValue = "redis: error convert value"
var MsgErrDestIsNotPointer = "redis: dest is not pointer"
var MsgErrKeyNotFound = "redis: key not found"

var ErrConvertKey = errors.New(MsgErrConvertKey)
var ErrConvertNewKey = errors.New(MsgErrConvertNewKey)
var ErrConvertValue = errors.New(MsgErrConvertValue)
var ErrDestIsNotPointer = errors.New(MsgErrDestIsNotPointer)
var ErrKeyNotFound = errors.New(MsgErrKeyNotFound)
