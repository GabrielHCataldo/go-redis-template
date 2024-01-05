package option

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"time"
)

type Set struct {
	Expiration time.Duration
}

func NewSet() Set {
	return Set{}
}

func (s Set) SetExpiration(exp time.Duration) Set {
	s.Expiration = exp
	return s
}

func GetOptionSetByParams(opts []Set) Set {
	result := Set{}
	for _, opt := range opts {
		if helper.IsNotEmpty(opt.Expiration) {
			result.Expiration = opt.Expiration
		}
	}
	return result
}
