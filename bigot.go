package bigot

import (
	"strings"
	"time"

	"sync"

	"github.com/mitchellh/mapstructure"
	"github.com/rexlv/bigot/provider"
	"github.com/rexlv/bigot/provider/file"
	"github.com/rexlv/gutil/cast"
	"github.com/rexlv/gutil/gmap"
)

// Bigot bigot
type Bigot struct {
	provider provider.Provider

	mu sync.RWMutex

	data     map[string]interface{}
	keyDelim string
}

// New returns Bigot instance
func New(p provider.Provider) *Bigot {
	return &Bigot{
		provider: p,
		data:     make(map[string]interface{}),
		keyDelim: ".",
	}
}

// SetKeyDelim set delim for key
func (b *Bigot) SetKeyDelim(delim string) {
	b.keyDelim = delim
}

// InitFromFile returns
func InitFromFile(path string) *Bigot {
	p := file.NewProvider(path)
	return New(p)
}

// ReadInConfig load config from provider
func (b *Bigot) ReadInConfig() (err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	data, err := b.provider.Read()
	if err != nil {
		return err
	}

	b.data = data
	return nil
}

// Get returns the value associated with the key
func (b *Bigot) Get(key string) interface{} {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if val := b.find(key); val != nil {
		return val
	}

	return nil
}

// GetString returns the value associated with the key as a string.
func (b *Bigot) GetString(key string) string {
	return cast.ToString(b.Get(key))
}

// GetBool returns the value associated with the key as a boolean.
func (b *Bigot) GetBool(key string) bool {
	return cast.ToBool(b.Get(key))
}

// GetInt returns the value associated with the key as an integer.
func (b *Bigot) GetInt(key string) int {
	return cast.ToInt(b.Get(key))
}

// GetInt64 returns the value associated with the key as an integer.
func (b *Bigot) GetInt64(key string) int64 {
	return cast.ToInt64(b.Get(key))
}

// GetFloat64 returns the value associated with the key as a float64.
func (b *Bigot) GetFloat64(key string) float64 {
	return cast.ToFloat64(b.Get(key))
}

// GetTime returns the value associated with the key as time.
func (b *Bigot) GetTime(key string) time.Time {
	return cast.ToTime(b.Get(key))
}

// GetDuration returns the value associated with the key as a duration.
func (b *Bigot) GetDuration(key string) time.Duration {
	return cast.ToDuration(b.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (b *Bigot) GetStringSlice(key string) []string {
	return cast.ToStringSlice(b.Get(key))
}

// GetSlice returns the value associated with the key as a slice of strings.
func (b *Bigot) GetSlice(key string) []interface{} {
	return cast.ToSlice(b.Get(key))
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (b *Bigot) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(b.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (b *Bigot) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(b.Get(key))
}

func (b *Bigot) GetSliceStringMap(key string) []map[string]interface{} {
	return cast.ToSliceStringMap(b.Get(key))
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (b *Bigot) GetStringMapStringSlice(key string) map[string][]string {
	return cast.ToStringMapStringSlice(b.Get(key))
}

// UnmarshalKey takes a single key and unmarshals it into a Struct.
func (b *Bigot) UnmarshalKey(key string, rawVal interface{}) error {
	if key == "" {
		return mapstructure.Decode(b.data, rawVal)
	}
	return mapstructure.Decode(b.Get(key), rawVal)
}

func (b *Bigot) find(key string) interface{} {
	paths := strings.Split(key, b.keyDelim)
	m := gmap.DeepSearchInMap(b.data, paths[:len(paths)-1]...)
	return m[paths[len(paths)-1]]
}
