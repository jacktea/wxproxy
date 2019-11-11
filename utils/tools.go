package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

//解析过期时间
func ParseExpire(expiresIn int64) time.Time {
	now := time.Now()
	d, _ := time.ParseDuration(fmt.Sprintf("%ds", expiresIn))
	return now.Add(d)
}

func Random(c int) (b []byte) {
	src := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789abcdefghijklmnopqrstuvwxyz")
	n := len(src)
	buf := &bytes.Buffer{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < c; i += 1 {
		index := r.Intn(n)
		buf.WriteByte(src[index])
	}
	b = buf.Bytes()
	return
}

func Scheme(req *http.Request) string {
	if req.TLS != nil {
		return "https"
	}
	if scheme := req.Header.Get("X-Forwarded-Proto"); scheme != "" {
		return scheme
	}
	if scheme := req.Header.Get("X-Forwarded-Protocol"); scheme != "" {
		return scheme
	}
	if ssl := req.Header.Get("X-Forwarded-Ssl"); ssl == "on" {
		return "https"
	}
	if scheme := req.Header.Get("X-Url-Scheme"); scheme != "" {
		return scheme
	}
	return "http"
}

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ValidatePath(dirs []string) (string, bool) {
	for _, v := range dirs {
		if len(v) > 0 {
			v, _ = filepath.Abs(v)
			if FileExists(v) {
				return v, true
			}
		}
	}
	return "", false
}

func Md5(in string) string {
	has := md5.Sum([]byte(in))
	return fmt.Sprintf("%x", has)
}
