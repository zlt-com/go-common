package common

import (
	"encoding/json"
	"go/ast"
	"reflect"
)

//浅克隆，可以克隆任意数据类型，对指针类型子元素无法克隆
//获取类型：如果类型是指针类型，需要使用Elem()获取对象实际类型
//获取实际值：如果值是指针类型，需要使用Elem()获取实际数据
//说白了，Elem()就是获取反射数据的实际类型和实际值
func Clone(src interface{}) interface{} {
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr { //如果是指针类型
		typ = typ.Elem()               //获取源实际类型(否则为指针类型)
		dst := reflect.New(typ).Elem() //创建对象
		data := reflect.ValueOf(src)   //源数据值
		data = data.Elem()             //源数据实际值（否则为指针）
		dst.Set(data)                  //设置数据
		dst = dst.Addr()               //创建对象的地址（否则返回值）
		return dst.Interface()         //返回地址
	} else {
		dst := reflect.New(typ).Elem() //创建对象
		data := reflect.ValueOf(src)   //源数据值
		dst.Set(data)                  //设置数据
		return dst.Interface()         //返回
	}
}

//深度克隆，可以克隆任意数据类型
func DeepClone(src interface{}) interface{} {
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr { //如果是指针类型
		typ = typ.Elem()                          //获取源实际类型(否则为指针类型)
		dst := reflect.New(typ).Elem()            //创建对象
		b, _ := json.Marshal(src)                 //导出json
		json.Unmarshal(b, dst.Addr().Interface()) //json序列化
		return dst.Addr().Interface()             //返回指针
	} else {
		dst := reflect.New(typ).Elem()            //创建对象
		b, _ := json.Marshal(src)                 //导出json
		json.Unmarshal(b, dst.Addr().Interface()) //json序列化
		return dst.Interface()                    //返回值
	}
}

//反射字段值
func ReflectFilde(i interface{}, name string) (field map[string]interface{}) {
	field = map[string]interface{}{}
	reflectType := reflect.ValueOf(i).Type()
	refValue := reflect.ValueOf(i)
	for reflectType.Kind() == reflect.Slice || reflectType.Kind() == reflect.Ptr {
		reflectType = reflectType.Elem()
	}

	for i := 0; i < reflectType.NumField(); i++ {
		if fieldStruct := reflectType.Field(i); ast.IsExported(fieldStruct.Name) && name == fieldStruct.Name {
			field[name] = refValue.Field(i).Interface()
		}
	}
	return
}

//反射字段值
func ReflectMethod(i interface{}, name string) []reflect.Value {
	refValue := reflect.ValueOf(i)
	method := refValue.MethodByName(name)
	return method.Call(nil)
}
