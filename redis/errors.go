package redis

import "errors"

var ErrConvertKey = errors.New("redis: error convert key to string")
var ErrConvertValue = errors.New("redis: error convert value")
var ErrDestIsNotPointer = errors.New("redis: dest is not pointer")
