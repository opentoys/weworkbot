package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

var reg = regexp.MustCompile(`-{1,}(.*)=(.*)`)

func ParseArgs(lst []string) (data map[string]string) {
	var keys, vals []string
	data = make(map[string]string)
	for k := range lst {
		if reg.MatchString(lst[k]) {
			var parse = reg.FindStringSubmatch(lst[k])
			keys = append(keys, parse[1])
			vals = append(vals, parse[2])
			continue
		}

		if strings.HasPrefix(lst[k], "--") {
			keys = append(keys, lst[k][2:])
			vals = append(vals, "")
			continue
		}

		if strings.HasPrefix(lst[k], "-") {
			keys = append(keys, lst[k][1:])
			vals = append(vals, "")
			continue
		}

		if vals[len(keys)-1] == "" {
			vals[len(keys)-1] = lst[k]
		}
	}

	for k := range keys {
		data[keys[k]] = vals[k]
	}
	return
}

func JSONBytes(dst any) (buf []byte) {
	buf, _ = json.Marshal(dst)
	return
}

func Path2bytes(s string) (buf []byte, name string, e error) {
	if s == "" {
		e = errors.New("input file error")
		return
	}
	if strings.HasPrefix(s, "http") {
		var resp *http.Response
		resp, e = http.DefaultClient.Get(s)
		if e != nil {
			return nil, "", e
		}
		buf, e = io.ReadAll(resp.Body)
		name = path.Base(resp.Request.URL.Path)
		return
	}

	if strings.HasPrefix(s, "data:") {
		var lst = strings.Split(s, ",")
		if len(lst) == 0 {
			e = errors.New("base64 error")
			return
		}
		buf, e = base64.RawStdEncoding.DecodeString(lst[1])
		if e != nil {
			buf, e = base64.RawURLEncoding.DecodeString(lst[1])
		}
		return
	}

	buf, e = os.ReadFile(s)
	name = filepath.Base(s)
	return
}
