package util

import (
	"errors"
	"strconv"
	"sync"
)

var (
	sliceUtilInstance *Slice
	sliceOnce sync.Once
)

// SliceUtil 切片工具
func SliceUtil() *Slice {
	sliceOnce.Do(func() {
		sliceUtilInstance = new(Slice)
	})
	return sliceUtilInstance
}

type Slice struct {

}

// StringSliceToMap 切片转Map
// @param stringSlice 字符串切片
// @return map[string]string 字符串map
func (u Slice) StringSliceToMap(stringSlice []string) map[string]string {
	stringMap := make(map[string]string)
	for i := range stringSlice {
		stringMap[strconv.Itoa(i)] = stringSlice[i]
	}
	return stringMap
}

// StringSliceToMapWithKey 切片转Map
// @param stringSlice 字符串切片
// @param keys keys切片
// @return stringMap 字符串map
// @return error
func (u Slice) StringSliceToMapWithKey(stringSlice []string, keys []string) (stringMap map[string]string, err error) {
	if len(stringSlice) != len(keys) {
		return nil, errors.New("keys slice length discord")
	}

	stringMap = make(map[string]string)
	for i := range stringSlice {
		stringMap[keys[i]] = stringSlice[i]
	}

	return stringMap, nil
}
