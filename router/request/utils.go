package request

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/rand"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
)

type Params struct {
	Keys  []string
	Match *regexp.Regexp
}

func randFileName(filename string) string {
	//randBytes := make([]byte, 16)
	//rand.Read(randBytes)
	var buff = strconv.FormatInt(time.Now().UnixNano(), 10)

	var hash = sha1.New()
	hash.Write([]byte(buff))
	return fmt.Sprintf("%s%s", hex.EncodeToString(hash.Sum(nil)), filepath.Ext(filename))

	return fmt.Sprintf(
		"%d%s",
		rand.Int63n(time.Now().UnixNano()),
		filepath.Ext(filename),
	)
}

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

func InArray[T Ordered](needles T, array []T) bool {
	for _, t := range array {
		if t == needles {
			return true
		}
	}
	return false
}
