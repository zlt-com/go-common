package common

import (
	"bufio"
	"bytes"
	"crypto/des"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

// GetRandomString 生成随机字符串
func GetRandomString(n int, letterBytes string) string {
	//str := "请填写接口配置信息此信息需要你有自己的服务器资源填写的URL需要正确响应微信发送的Token验证请阅读消息接口使用指南"
	if letterBytes == "" {
		letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}

	// bytes := []byte(str)
	// result := []byte{}
	// r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// for i := 0; i < leng; i++ {
	// 	result = append(result, bytes[r.Intn(len(bytes))])
	// }
	// return string(result)
	time.Sleep(time.Duration(1) * time.Nanosecond)
	var (
		src           = rand.NewSource(time.Now().UnixNano())
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(int(cache) & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// StringToUint 字符串转uint
func StringToUint(strNumber string, value interface{}) (err error) {
	var number interface{}
	number, err = strconv.ParseUint(strNumber, 10, 64)
	switch v := number.(type) {
	case uint64:
		switch d := value.(type) {
		case *uint64:
			*d = v
		case *uint:
			*d = uint(v)
		case *uint16:
			*d = uint16(v)
		case *uint32:
			*d = uint32(v)
		case *uint8:
			*d = uint8(v)
		}
	}
	return
}

// StrToUInt 字符串转uint
func StrToUInt(str string) uint {
	i, e := strconv.Atoi(str)
	if e != nil {
		return 0
	}
	return uint(i)
}

// MD5 加密
func MD5(source string) string {
	h := md5.New()
	h.Write([]byte(source))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

// Append 连接字符串，第一个是源字符串
func Append(args ...string) string {
	var buffer bytes.Buffer
	for _, arg := range args {
		buffer.WriteString(arg)
	}
	return buffer.String()
}

// Splits 分割字符串，返回string
func Splits(s string, sep string) []string {
	arr := strings.Split(s, sep)
	var args []string
	for _, v := range arr {
		args = append(args, v)
	}
	return args
}

// Split 分割字符串，返回interface
func Split(s string, sep string) []interface{} {
	arr := strings.Split(s, sep)
	var args []interface{}
	for _, v := range arr {
		args = append(args, v)
	}
	return args
}

// Seize 占位方法
//s 原始字符串
//se 占位符
//pos 位置，0：前 1：后
//length 共几位
func Seize(s interface{}, se string, pos, length int) string {
	var source string
	switch s.(type) {
	case string:
		source = s.(string)
	case int:
		source = strconv.Itoa(s.(int))
	}
	for i := 0; i < length-len(source); i++ {
		if pos == 0 {
			source = Append(se, source)
		} else {
			source = Append(source, se)
		}
	}
	return source
}

// Empty 判断字符串是否为空
func Empty(s string) bool {
	if strings.EqualFold(strings.Trim(s, " "), "") {
		return true
	}
	return false
}

// NotEmpty NotEmpty
func NotEmpty(s string) bool {
	if strings.EqualFold(strings.Trim(s, " "), "") {
		return false
	}
	return true
}

// ZeroPadding ZeroPadding
func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

// ZeroUnPadding ZeroUnPadding
func ZeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}

// Encrypt Encrypt加密
func Encrypt(text string, key []byte) (string, error) {
	src := []byte(text)
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	bs := block.BlockSize()
	src = ZeroPadding(src, bs)
	if len(src)%bs != 0 {
		return "", errors.New("Need a multiple of the blocksize")
	}
	out := make([]byte, len(src))
	dst := out
	for len(src) > 0 {
		block.Encrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	return hex.EncodeToString(out), nil
}

// Decrypt Decrypt解密
func Decrypt(decrypted string, key []byte) (string, error) {
	src, err := hex.DecodeString(decrypted)
	if err != nil {
		return "", err
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return "", err
	}
	out := make([]byte, len(src))
	dst := out
	bs := block.BlockSize()
	if len(src)%bs != 0 {
		return "", errors.New("crypto/cipher: input not full blocks")
	}
	for len(src) > 0 {
		block.Decrypt(dst, src[:bs])
		src = src[bs:]
		dst = dst[bs:]
	}
	out = ZeroUnPadding(out)
	return string(out), nil
}

// GetStructTagValue struct中tag的值
func GetStructTagValue(value interface{}, tagName string) (title []string) {
	e := reflect.ValueOf(value).Elem()
	for i := 0; i < e.NumField(); i++ {
		title = append(title, e.Type().Field(i).Tag.Get(tagName))
	}
	return title
}

// GetJSONString 获取json字符串
func GetJSONString(i interface{}) string {
	jsonBytes, err := json.Marshal(i)
	if err != nil {
		return err.Error()
	}
	return string(jsonBytes)
}

// ReadString ReadFile
func ReadString(str string) []string {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(str))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

// Int64TString Int64TString
func Int64TString(f int64) string {
	return strconv.FormatInt(f, 10)
}

// Float64ToString Float64ToString
func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', 10, 64)
}

// FilteredSQLInject 正则过滤sql注入的方法
// 参数 : 要匹配的语句
func FilteredSQLInject(source string) bool {
	//过滤 ‘
	//ORACLE 注解 --  /**/
	//关键字过滤 update ,delete
	// 正则的字符串, 不能用 " " 因为" "里面的内容会转义
	str := `(?:')|(?:--)|(/\\*(?:.|[\\n\\r])*?\\*/)|(\b(select|update|and|or|delete|insert|trancate|char|chr|into|substr|ascii|declare|exec|count|master|into|drop|execute)\b)`
	re, err := regexp.Compile(str)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return re.MatchString(source)
}
