package pluma

import (
	"fmt"
	"strconv"
)

type Option struct {
	Name  string
	Value interface{}
	Help  string
}

type AnyProvider interface {
	Getter
	Setter
}

type AllProvider interface {
	StringProvider
	NumberProvider
	BoolProvider
}

type BoolProvider interface {
	BoolGetter
	BoolSetter
}

type StringProvider interface {
	StringGetter
	StringSetter
}

type NumberProvider interface {
	IntProvider
	Float32Provider
	Float64Provider
}

type IntProvider interface {
	IntGetter
	IntSetter
}
type Float32Provider interface {
	Float32Getter
	Float32Setter
}

type Float64Provider interface {
	Float64Getter
	Float64Setter
}

type Setter interface {
	Set(k string, v interface{})
}

type Getter interface {
	Get(k string) interface{}
}

type StringGetter interface {
	GetString(k string) string
}

type StringSetter interface {
	SetString(k string, v string)
}

type BoolGetter interface {
	GetBool(k string) bool
}

type BoolSetter interface {
	SetBool(k string, v bool)
}

type IntGetter interface {
	GetInt(k string) int
}

type IntSetter interface {
	SetInt(k string, v int)
}

type Float32Getter interface {
	GetFloat32(k string) float32
}

type Float32Setter interface {
	SetFloat32(k string, v float32)
}

type Float64Getter interface {
	GetFloat64(k string) float64
}

type Float64Setter interface {
	SetFloat64(k string, v float64)
}

// Creates a new instance of a configuration provider
func DefaultProvider() (c *Config) {
	return &Config{make(map[string]interface{})}
}

// Creates a provider with given options by default
func WithOptions(options map[string]interface{}) (c *Config) {
	return &Config{options}
}

type Config struct {
	opts map[string]interface{}
}

// Returns the value of config with key k. The return type is undefined.
// To get the value in a specific type, use interface{} of the type specific methods:
// GetString,GetInt,GetFloat64,GetFloat32 and GetBool
func (c *Config) Get(k string) interface{} {
	return c.opts[k]
}
func (c *Config) Set(k string, v interface{}) {
	c.opts[k] = v
}

func (c *Config) Load() {}

// Returns the string version of the config with key k
// if type conversion is not possible or value is not set empty string "" is returned
func (c *Config) GetString(k string) (v string) {
	val := c.Get(k)
	if val == nil {
		return
	}
	return fmt.Sprint(val)

}

func (c *Config) SetString(k string, v string) {
	c.Set(k, v)
}

// Returns the int version of the config with key k
// if type conversion is not possible or value is not set, 0 (zero) is returned
func (c *Config) GetInt(k string) (v int) {
	pval := c.Get(k)
	return toInt(pval)
}

func (c *Config) SetInt(k string, v int) {
	c.Set(k, v)
}

// Returns the bool version of the config with key k
// if type conversion is not possible or value is not set, false is returned
func (c *Config) GetBool(k string) (v bool) {
	pval := c.Get(k)
	return toBool(pval)
}

func (c *Config) SetBool(k string, v bool) {
	c.Set(k, v)
}

// Returns the float64 version of the config with key k
// if type conversion is not possible or value is not set, 0.0 is returned
func (c *Config) GetFloat64(k string) (v float64) {
	pval := c.Get(k)
	return toFloat(pval, 64)
}

func (c *Config) SetFloat64(k string, v float64) {
	c.Set(k, v)
}

// Returns the float32 config with key k
// if type conversion is not possible or value is not set, 0.0 is returned
func (c *Config) GetFloat32(k string) (v float32) {
	return float32(toFloat(k, 32))
}

func (c *Config) SetFloat32(k string, v float32) {
	c.Set(k, v)
}

// Insert option `optn` in provider's available options
// If a value for key `name` exists, it will be replaced
func (c *Config) Insert(name string, optn *Option) {
	c.opts[name] = optn
}

// Removes option with key `name` from available options
func (c *Config) Remove(name string, optn *Option) {
	delete(c.opts, name)
}

func toInt(v interface{}) (i int) {
	if v == nil {
		return
	}
	if val, err := strconv.Atoi(fmt.Sprint(v)); err == nil {
		return val
	}
	return
}

func toFloat(v interface{}, bitSize int) (f float64) {
	if v == nil {
		return
	}
	if val, err := strconv.ParseFloat(fmt.Sprint(v), bitSize); err == nil {
		return val
	}
	return
}

func toBool(v interface{}) (b bool) {
	if v == nil {
		return
	}
	if val, err := strconv.ParseBool(fmt.Sprint(v)); err == nil {
		return val
	}
	return
}
