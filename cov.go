package common

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// JSON2Object 类型转换
func JSON2Object(js string, v interface{}) error {
	return json.Unmarshal([]byte(js), v)
}

// Byte2Object Byte2Object
func Byte2Object(by []byte, v interface{}) interface{} {
	json.Unmarshal(by, v)
	return v
}

// Object2Byte Object2Byte
func Object2Byte(v interface{}) []byte {
	value, err := json.Marshal(v)
	if err != nil {
		fmt.Println("utils.Object2Byte", err)
	}
	return value
}

// JSON2Map JSON2Map
func JSON2Map(jsonStr string) map[string]float64 {
	m := make(map[string]float64)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		fmt.Println(err)
	}
	return m
}

// Uint8FromInt int转unit
func Uint8FromInt(n int) (uint8, error) {
	if 0 <= n && n <= math.MaxUint8 { // conversion is safe
		return uint8(n), nil
	}
	return 0, fmt.Errorf("%d is out of the uint8 range", n)
}

// IntFromFloat64 float转int
func IntFromFloat64(x float64) int {
	if math.MinInt32 <= x && x <= math.MaxInt32 { // x lies in the integer range
		whole, fraction := math.Modf(x)
		if fraction >= 0.5 {
			whole++
		}
		return int(whole)
	}
	return 0
}

// IntFromInterface IntFromInterface
func IntFromInterface(x interface{}) int {
	switch x.(type) {
	case string:
		v, _ := strconv.Atoi(x.(string))
		return v
	case float64:
		f := x.(float64)
		if math.MinInt32 <= f && f <= math.MaxInt32 { // x lies in the integer range
			whole, fraction := math.Modf(f)
			if fraction >= 0.5 {
				whole++
			}
			return int(whole)
		}
		break
	default:
		return 0
	}
	return 0
}

// UnicodeEmojiDecode 表情解码
func UnicodeEmojiDecode(s string) string {
	//emoji表情的数据表达式
	re := regexp.MustCompile("\\[[\\\\u0-9a-zA-Z]+\\]")
	//提取emoji数据表达式
	reg := regexp.MustCompile("\\[\\\\u|]")
	src := re.FindAllString(s, -1)
	for i := 0; i < len(src); i++ {
		e := reg.ReplaceAllString(src[i], "")
		p, err := strconv.ParseInt(e, 16, 32)
		if err == nil {
			s = strings.Replace(s, src[i], string(rune(p)), -1)
		}
	}
	return s
}

// UnicodeEmojiCode 表情转换
func UnicodeEmojiCode(s string) string {
	ret := ""
	rs := []rune(s)
	for i := 0; i < len(rs); i++ {
		fmt.Println(string(rs[i]), len(string(rs[i])))
		if len(string(rs[i])) == 4 {
			u := `[\u` + strconv.FormatInt(int64(rs[i]), 16) + `]`
			ret += u

		} else {
			ret += string(rs[i])
		}
	}
	return ret
}

//CoventArray 转换数组类型
func CoventArray(value interface{}) (list []interface{}) {
	e := reflect.ValueOf(value).Elem()
	for i := 0; i < e.NumField(); i++ {
		varType := e.Type().Field(i).Type
		// title = append(title, e.Type().Field(i).Tag.Get("note"))
		// fmt.Println(varType.Name(), e.Type().Field(i).Tag.Get("note"), e.Field(i).Interface())
		switch varType.Name() {
		case "string":
			varName := e.Type().Field(i).Name
			if Contains(varName, []string{"Fulltextstr"}) {
				list = append(list, "")
				continue
			}
			varValue := e.Field(i).Interface()
			list = append(list, varValue.(string))
		case "int":
			varValue := e.Field(i).Interface()
			list = append(list, strconv.Itoa(varValue.(int)))
		case "uint":
			varValue := e.Field(i).Interface()
			list = append(list, strconv.Itoa(int(varValue.(uint))))
		case "Time":
			varValue := e.Field(i).Interface()
			list = append(list, (varValue.(time.Time)).Format(Datetime))
		case "bool":
			varValue := e.Field(i).Interface()
			b := "否"
			if varValue.(bool) {
				b = "是"
			}
			list = append(list, b)
		}
	}
	return list
}

//Covent2Map 转换数组类型
func Covent2Map(value interface{}) map[string]interface{} {
	list := make(map[string]interface{}, 0)
	e := reflect.ValueOf(value).Elem()
	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		varValue := e.Field(i).Interface()
		list[varName] = varValue
	}
	return list
}

// Struct2Arrary 反射模型所有属性
func Struct2Arrary(i interface{}) (kv []interface{}) {
	e := reflect.ValueOf(i).Elem()
	for i := 0; i < e.NumField(); i++ {
		varName := e.Type().Field(i).Name
		varType := e.Type().Field(i).Type
		varValue := e.Field(i).Interface()
		fmt.Printf("%v %v %v\n", varName, varType, varValue)
	}
	return kv
}

// Object2JSON Object2JSON
func Object2JSON(obj interface{}) string {
	j, err := json.Marshal(obj)
	if err != nil {
		return err.Error()
	}
	return string(j)
}
