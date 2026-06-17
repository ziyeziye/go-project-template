package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// RandomString 随机数生成
func RandomString(n int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// RandomInt 随机数生成
func RandomInt(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	var letter = []rune("0123456789")
	b := make([]rune, length)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// Random 根据区间产生随机数
func Random(min int, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(max-min) + min
}

// Sha1 随机数生成
func Sha1(data string) string {
	s1 := sha1.New()
	s1.Write([]byte(data))
	return hex.EncodeToString(s1.Sum([]byte("")))
}

func StrArrRandom(arr []string) string {
	arrLen := len(arr)
	randN := Random(0, arrLen)
	n := randN % arrLen
	return arr[n]
}

func JsonEncode(v any) string {
	marshal, _ := json.Marshal(v)
	return string(marshal)
}

// ArrayUnique 数组去重
func ArrayUnique[T comparable](list []T) []T {
	// 创建一个临时map用来存储数组元素
	temp := make(map[T]struct{})
	index := 0
	// 将元素放入map中
	for _, v := range list {
		temp[v] = struct{}{}
	}
	tempList := make([]T, len(temp))
	for key := range temp {
		tempList[index] = key
		index++
	}
	return tempList
}

func IfCase[T any](condition bool, trueVal, falseVal T) T {
	if condition {
		return trueVal
	}

	return falseVal
}

// DoTimeout 限制程序运行超时

func DoTimeout(duration time.Duration, task func(), doTimeout ...func()) error {
	timer := time.NewTimer(duration) // 设置定时器的超时时间，主线程只等5秒
	final := make(chan bool)
	go func() {
		task()
		final <- true
	}()
	// 知识点：主协程等待子线程，并有超时机制
	select {
	case <-final:
		return nil
	case <-timer.C: // 定时器也是一个通道
		if len(doTimeout) > 0 {
			doTimeout[0]()
		}
		return errors.New("do timeout")
	}
}

// RecoverRun 以保护方式运行一个函数
func RecoverRun(name string, entry func(), params ...interface{}) (result bool) {
	result = true
	// 延迟处理的函数
	defer func() {
		// 发生宕机时，获取panic传递的上下文并打印
		if err := recover(); err != nil {
			log.Printf("【%s】Recover#, params: %v, error: %v \n", name, params, err)

			var skip int
			for {
				skip++
				_, file, line, ok := runtime.Caller(skip)
				if ok {
					log.Printf("%s:%d \n", file, line)
				} else {
					break
				}
			}

			result = false
		}
	}()

	entry()
	return
}

func FilePutContents(filename string, data []byte) error {
	if dir := filepath.Dir(filename); dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return os.WriteFile(filename, data, 0644)
}

func IsZeroRef[T any](v T) bool {
	return reflect.ValueOf(&v).Elem().IsZero()
}

func VarDump(v ...interface{}) {
	spew.Dump(v)
}

func VarDumpDie(v ...interface{}) {
	spew.Dump(v)
	os.Exit(1)
}

func Base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func Base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}
