package util

import (
	"reflect"
	"sync"
)

var (
	reflectUtilInstance *Reflect
	reflectUtilOnce     sync.Once
)

func ReflectUtil() *Reflect {
	reflectUtilOnce.Do(func() {
		reflectUtilInstance = new(Reflect)
	})
	return reflectUtilInstance
}

// 反射工具
type Reflect struct {
}

// 引用对象方法
// @param object 结构体对象
// @param methodName 引用方法名
// @prams args 引用方法参数向量
func (u Reflect) InvokeMethod(object interface{}, methodName string, args ...interface{}) []reflect.Value {
	inputs := make([]reflect.Value, len(args))
	for i, _ := range args {
		inputs[i] = reflect.ValueOf(args[i])
	}
	//动态调用方法
	return reflect.ValueOf(object).MethodByName(methodName).Call(inputs)
}

// 引用对象属性
// @param object 结构体对象
// @prams elemName 引用属性名
func (u Reflect) InvokeElem(object interface{}, elemName string) reflect.Value {
	//动态访问属性
	return reflect.ValueOf(object).Elem().FieldByName(elemName)
}
