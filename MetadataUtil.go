package util

import (
	"context"
	"fmt"
	"google.golang.org/grpc/metadata"
	"strconv"
)

func NewMetadataUtil() MetadataUtil {
	return MetadataUtil{}
}

type MetadataUtil struct {
	md    metadata.MD
	ctxMd metadata.MD
}

// 设置 Metadata
func (u *MetadataUtil) SetMetadata(kv map[string]interface{}) {
	var kvs []string
	for k, v := range kv {
		switch v := v.(type) {
		case string:
			kvs = append(kvs, k, v)
		case int:
			kvs = append(kvs, k, strconv.Itoa(v))
		case int32:
			kvs = append(kvs, k, strconv.FormatInt(int64(v), 10))
		case int64:
			kvs = append(kvs, k, strconv.FormatInt(v, 10))
		case uint:
			kvs = append(kvs, k, strconv.Itoa(int(v)))
		case uint32:
			kvs = append(kvs, k, strconv.FormatInt(int64(v), 10))
		case uint64:
			kvs = append(kvs, k, strconv.FormatInt(int64(v), 10))
		case float32:
			kvs = append(kvs, k, strconv.FormatFloat(float64(v), 'E', -1, 32))
		case float64:
			kvs = append(kvs, k, strconv.FormatFloat(v, 'E', -1, 64))
		case []string:
			for _, item := range v {
				kvs = append(kvs, k, item)
			}
		case []int:
			for _, item := range v {
				kvs = append(kvs, k, strconv.Itoa(item))
			}
		case []int32:
			for _, item := range v {
				kvs = append(kvs, k, strconv.FormatInt(int64(item), 10))
			}
		case []int64:
			for _, item := range v {
				kvs = append(kvs, k, strconv.FormatInt(item, 10))
			}
		case []uint:
			for _, item := range v {
				kvs = append(kvs, k, strconv.Itoa(int(item)))
			}
		case []uint32:
			for _, item := range v {
				kvs = append(kvs, k, strconv.FormatInt(int64(item), 10))
			}
		case []uint64:
			for _, item := range v {
				kvs = append(kvs, k, strconv.FormatInt(int64(item), 10))
			}
		case []float32:
			for _, item := range v {
				kvs = append(kvs, k, strconv.FormatFloat(float64(item), 'E', -1, 32))
			}
		case []float64:
			for _, item := range v {
				kvs = append(kvs, k, strconv.FormatFloat(item, 'E', -1, 64))
			}
		default:
			kvs = append(kvs, k, fmt.Sprint(v))
		}
	}

	u.md = metadata.Pairs(kvs...)
}

// 获取 Metadata
func (u MetadataUtil) GetMetadata() metadata.MD {
	return u.md
}

// 解析 Metadata
func (u *MetadataUtil) ParseMetadata(ctx context.Context) (metadata.MD, bool) {
	var ok bool
	u.ctxMd, ok = metadata.FromIncomingContext(ctx)
	return u.ctxMd, ok
}

func (u MetadataUtil) GetStringValueSlice(key string) []string {
	return u.ctxMd.Get(key)
}

func (u MetadataUtil) GetStringValue(key string) string {
	if values := u.GetStringValueSlice(key); values != nil {
		return values[0]
	}
	return ""
}

func (u MetadataUtil) GetIntValueSlice(key string) []int {
	var values []int
	for _, v := range u.ctxMd.Get(key) {
		if intV, err := strconv.Atoi(v); err == nil {
			values = append(values, intV)
		}
	}
	return values
}

func (u MetadataUtil) GetIntValue(key string) int {
	if values := u.GetStringValueSlice(key); values != nil {
		if val, err := strconv.Atoi(values[0]); err == nil {
			return val
		}
	}
	return 0
}

func (u MetadataUtil) GetInt32ValueSlice(key string) []int32 {
	var values []int32
	for _, v := range u.ctxMd.Get(key) {
		if intV, err := strconv.ParseInt(v, 10, 32); err == nil {
			values = append(values, int32(intV))
		}
	}
	return values
}

func (u MetadataUtil) GetInt32Value(key string) int32 {
	if values := u.GetStringValueSlice(key); values != nil {
		if val, err := strconv.ParseInt(values[0], 10, 32); err == nil {
			return int32(val)
		}
	}
	return 0
}

func (u MetadataUtil) GetInt64ValueSlice(key string) []int64 {
	var values []int64
	for _, v := range u.ctxMd.Get(key) {
		if intV, err := strconv.ParseInt(v, 10, 64); err == nil {
			values = append(values, intV)
		}
	}
	return values
}

func (u MetadataUtil) GetInt64Value(key string) int64 {
	if values := u.GetStringValueSlice(key); values != nil {
		if val, err := strconv.ParseInt(values[0], 10, 64); err == nil {
			return val
		}
	}
	return 0
}

func (u MetadataUtil) GetUintValueSlice(key string) []uint {
	var values []uint
	for _, v := range u.ctxMd.Get(key) {
		if intV, err := strconv.Atoi(v); err == nil {
			values = append(values, uint(intV))
		}
	}
	return values
}

func (u MetadataUtil) GetUintValue(key string) uint {
	if values := u.GetStringValueSlice(key); values != nil {
		if val, err := strconv.Atoi(values[0]); err == nil {
			return uint(val)
		}
	}
	return 0
}

func (u MetadataUtil) GetUint32ValueSlice(key string) []uint32 {
	var values []uint32
	for _, v := range u.ctxMd.Get(key) {
		if intV, err := strconv.ParseInt(v, 10, 32); err == nil {
			values = append(values, uint32(intV))
		}
	}
	return values
}

func (u MetadataUtil) GetUint32Value(key string) uint32 {
	if values := u.GetStringValueSlice(key); values != nil {
		if val, err := strconv.ParseInt(values[0], 10, 32); err == nil {
			return uint32(val)
		}
	}
	return 0
}

func (u MetadataUtil) GetUint64ValueSlice(key string) []uint64 {
	var values []uint64
	for _, v := range u.ctxMd.Get(key) {
		if intV, err := strconv.ParseInt(v, 10, 64); err == nil {
			values = append(values, uint64(intV))
		}
	}
	return values
}

func (u MetadataUtil) GetUint64Value(key string) uint64 {
	if values := u.GetStringValueSlice(key); values != nil {
		if val, err := strconv.ParseInt(values[0], 10, 64); err == nil {
			return uint64(val)
		}
	}
	return 0
}

func (u MetadataUtil) GetFloat32ValueSlice(key string) []float32 {
	var values []float32
	for _, v := range u.ctxMd.Get(key) {
		if intV, err := strconv.ParseFloat(v, 32); err == nil {
			values = append(values, float32(intV))
		}
	}
	return values
}

func (u MetadataUtil) GetFloat32Value(key string) float32 {
	if values := u.GetStringValueSlice(key); values != nil {
		if val, err := strconv.ParseFloat(values[0], 32); err == nil {
			return float32(val)
		}
	}
	return 0
}

func (u MetadataUtil) GetFloat64ValueSlice(key string) []float64 {
	var values []float64
	for _, v := range u.ctxMd.Get(key) {
		if intV, err := strconv.ParseFloat(v, 64); err == nil {
			values = append(values, intV)
		}
	}
	return values
}

func (u MetadataUtil) GetFloat64Value(key string) float64 {
	if values := u.GetStringValueSlice(key); values != nil {
		if val, err := strconv.ParseFloat(values[0], 64); err == nil {
			return val
		}
	}
	return 0
}
