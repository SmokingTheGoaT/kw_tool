package common

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

func ConcatenateMapIntoString(m map[string]string) (str string) {
	str = "?"
	c := true
	for i, v := range m {
		if c {
			str = str + i + "=" + v
			c = false
		}
		str = str + fmt.Sprintf("&%s=%s", i, v)
	}
	return
}

func CreateFilePath(s string) (path string, err error) {
	if s == "" {
		err = errors.New("path value is empty please provide a value")
	} else {
		if path, err = os.Getwd(); err == nil {
			path = path + s
		}
	}
	return
}

func NormalizeQuery(s string) string {
	return strings.ReplaceAll(s, " ", "+")
}

func GetNextTaskExecutionTime(min int64, max int64) (t time.Duration) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int63n(max-min+1) + min
	fmt.Println(n)
	t = time.Duration(n) * time.Millisecond
	return
}
