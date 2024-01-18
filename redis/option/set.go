package option

import (
	"github.com/GabrielHCataldo/go-helper/helper"
	"time"
)

// Set represents options that can be used to configure an 'Set' operation.
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

// NewSet creates a new Set instance.
func NewSet() *Set {
	return &Set{}
}

// SetMode sets value for the Mode field.
func (s *Set) SetMode(mode SetMode) *Set {
	s.Mode = mode
	return s
}

// SetTTL sets value for the TTL field.
func (s *Set) SetTTL(ttl time.Duration) *Set {
	s.TTL = ttl
	return s
}

// SetExpireAt sets value for the ExpireAt field.
func (s *Set) SetExpireAt(expAt time.Time) *Set {
	s.ExpireAt = expAt
	return s
}

// SetKeepTTL sets value for the KeepTTL field.
func (s *Set) SetKeepTTL(keepTTL bool) *Set {
	s.KeepTTL = keepTTL
	return s
}

// GetOptionSetByParams assembles the Set object from optional parameters.
func GetOptionSetByParams(opts []*Set) *Set {
	result := &Set{}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
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
