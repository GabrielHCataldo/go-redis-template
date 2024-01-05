package option

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"time"
)

type Set struct {
	// Mode can be SetModeNx, SetModeXx or SetModeDefault.
	Mode SetMode
	// Zero `TTL` or `Expiration` means that the key has no expiration time.
	TTL      time.Duration
	ExpireAt time.Time
	// KeepTTL is a Redis KEEPTTL option to keep existing TTL, it requires your redis-server version >= 6.0,
	// otherwise you will receive an error: (error) ERR syntax error.
	KeepTTL bool
}

func NewSet() Set {
	return Set{}
}

func (s Set) SetMode(mode SetMode) Set {
	s.Mode = mode
	return s
}

func (s Set) SetTTL(ttl time.Duration) Set {
	s.TTL = ttl
	return s
}

func (s Set) SetExpireAt(expAt time.Time) Set {
	s.ExpireAt = expAt
	return s
}

func (s Set) SetKeepTTL(keepTTL bool) Set {
	s.KeepTTL = keepTTL
	return s
}

func GetOptionSetByParams(opts []Set) Set {
	result := Set{}
	for _, opt := range opts {
		if helper.IsNotEmpty(opt.Mode) {
			result.Mode = opt.Mode
		}
		if helper.IsNotEmpty(opt.TTL) {
			result.TTL = opt.TTL
		}
		if helper.IsNotEmpty(opt.ExpireAt) {
			result.ExpireAt = opt.ExpireAt
		}
		if opt.KeepTTL {
			result.KeepTTL = opt.KeepTTL
		}
	}
	return result
}
