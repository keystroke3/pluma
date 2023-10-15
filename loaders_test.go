package pluma

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

var optionsString = `
{
	"port": "4928",
	"database": "sqlite",
	"email":"test@example.com",
	"pi": "3.142",
	"verbose": "false"
}
`

type testConfig struct{ opts map[string]string }

func (c *testConfig) Set(k string, v interface{}) {
	c.opts[k] = v.(string)
}

func TestFromEnv(t *testing.T) {
	prefix := "TEST_"
	options := make(map[string]string)
	upperOptions := make(map[string]string)
	err := json.Unmarshal([]byte(optionsString), &options)
	if err != nil {
		t.Error("Failed to load json", err)
	}
	type Test struct {
		keys   []string
		prefix string
		expect map[string]string
		name   string
	}
	var keys []string
	var upper_keys []string
	for k, v := range options {
		keys = append(keys, k)
		upper_k := strings.ToUpper(k)
		upper_keys = append(upper_keys, upper_k)
		upperOptions[upper_k] = v
		os.Setenv(upper_k, v)
		os.Setenv(fmt.Sprintf("%v_%v", prefix, upper_k), v)
	}
	unsetMap := map[string]string{
		"LOREM": "",
		"IPSUM": "",
	}

	tests := []Test{
		{
			keys:   keys,
			expect: options,
			name:   "Lowercase keys, expect full config",
		},
		{
			keys:   keys,
			expect: options,
			prefix: prefix,
			name:   "Lowercase keys with prefix, expect full config",
		},
		{
			keys:   []string{},
			expect: make(map[string]string),
			name:   "Empty string slice, expect empty config",
		},
		{
			keys:   upper_keys,
			expect: upperOptions,
			name:   "Uppercase keys, expect full config",
		},
		{
			keys:   upper_keys,
			prefix: prefix,
			expect: upperOptions,
			name:   "Uppercase with prefix keys, expect full config",
		},
		{
			keys:   []string{"LOREM", "IPSUM"},
			expect: unsetMap,
			name:   "Unset values, expect map with empty values",
		},
	}
	for _, test := range tests {
		cfg := testConfig{opts: make(map[string]string)}
		t.Run(test.name, func(t *testing.T) {
			var bf bytes.Buffer
			log.SetOutput(&bf)
			t.Cleanup(func() {
				log.SetOutput(os.Stdout)
			})
			FromEnv(test.keys, &cfg, test.prefix)
			if !reflect.DeepEqual(cfg.opts, test.expect) {
				t.Errorf("Test Failed, expect %v, have %v", test.expect, cfg.opts)
			}
			t.Log(bf.String())
		})
	}

}
