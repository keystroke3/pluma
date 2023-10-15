package pluma

import (
	"flag"
	"fmt"
	"io"
	"log"
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
// E.g. if prefix is set to `PROGRAM_` PROGRAM_RETRY will be used but RETRY will not
func FromEnv(keys []string, s Setter, prefix ...string) {
	var pfx string
	if len(prefix) == 0 {
		pfx = ""
	} else {
		pfx = prefix[0]
	}
	log.Println("Prefix:", pfx)
	for _, k := range keys {
		lookupStr := pfx + strings.ToUpper(k)
		val, _ := os.LookupEnv(lookupStr)
		s.Set(k, val)
	}
}

// Sets config provider values for the given `keys` from cli flags
// e.g. --retry=true  will load as retry="true"
func FromFlags(keys []string, p Setter) {
	for _, k := range keys {
		f := flag.Lookup(k)
		if f.DefValue != f.Value.String() {
			p.Set(k, f.Value.String())
		}
	}
}

// Wrapper around FromReader that takes a file name f as a string and creates a file reader
// It then calls FromReader with the os.File as the reader
func FromFile(f string, keys []string, p Setter, d string, c ...string) error {
	err := isFile(f)
	if err != nil {
		return err
	}
	_f, err := os.Open(f)
	if err != nil {
		return err
	}
	return FromReader(_f, keys, p, d, c...)
}

// Loads listed options `opts` values from reader r and sets them in provider p.
// Delimiter d is used to distinguish the end of key and start of value. For example "="
// An optional slice of comment prefixes can be passed to ignore lines prepended with those characters
// e.g, if comment prefixes include "#" or "//", all lines begging with those characters are ignored
func FromReader(r io.Reader, keys []string, p Setter, d string, c ...string) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	lines := strings.Split(string(b), "\n")
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
