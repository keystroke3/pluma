package pluma

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Loader interface {
	Load()
}

// Sets config provider values for `keys` to those defined in the environment variables
// All environment keys are assumed to be defined in UPPERCASE and all passed keys passed
// to the function will be converted to uppercase before lookup.
// Optional string `prefix` can be passed to only fetch values that start with that prefix
// E.g. if prefix is set to `PROGRAM_` PROGRAM_RETRY will be set but RETRY will not
func FromEnv(keys []string, s Setter, prefix ...string) {
	var pfx string
	if len(prefix) == 0 {
		pfx = ""
	} else {
		pfx = prefix[0]
	}
	for _, k := range keys {
		val, set := os.LookupEnv(pfx + strings.ToUpper(k))
		if set {
			s.Set(k, val)
		}
	}
}

// Sets config provider values for the given `keys` from cli flags
func FromFlags(keys []string, p Setter) {
	for _, k := range keys {
		f := flag.Lookup(k)
		if f.DefValue != f.Value.String() {
			p.Set(k, f.Value.String())
		}
	}
}

// Loads listed options `opts` values from file f and sets them in provider p.
// Delimiter d is used to distinguish the end of key and start of value. For example "="
// An optional slice of comment prefixes can be passed to ignore lines prepended with those characters
// e.g, if comment prefixes include "#" or "//", all lines begging with those characters are ignored
func FromFile(f string, keys []string, p Setter, d string, c ...string) error {
	err := isFile(f)
	if err != nil {
		return err
	}
	_f, err := os.ReadFile(f)
	if err != nil {
		return err
	}
	lines := strings.Split(string(_f), "\n")
	cfgs := make(map[string]string)
	pattern := `^\s*(?:` + strings.Join(c, "|") + ").*"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("Error parsing comment prefixes, %v", err)
	}
	for _, l := range lines {
		if len(c) > 0 && regex.MatchString(l) {
			continue
		}
		_l := strings.SplitN(l, d, 2)
		key := strings.TrimSpace(_l[0])
		if len(_l) < 2 {
			continue
		}
		val := _l[1]
		if val != "" {
			cfgs[key] = val
		}
	}
	for _, k := range keys {
		v, set := cfgs[strings.ToUpper(k)]
		if set {
			p.Set(k, v)
		}
	}
	return nil
}

var Options = map[string]Option{
	"workDir": {
		Name:  "workDir",
		Value: ".",
		Help:  "The directory to place working files such as migrations and databases",
	},

	"allowedOrigins": {
		Name:  "allowedOrigins",
		Value: "http://localhost:*,https://sala.pm",
		Help:  "Comma separated allowed urls with one wildcard (*) per url for allowed. Use a single * to allow all",
	},
	"port": {
		Name:  "port",
		Value: "5000",
		Help:  "The port for the server to listen on",
	},
}
