package util

import (
	"sync/atomic"
)

func NewRoundRobinUtil(number uint) RoundRobinUtil {
	return RoundRobinUtil{
		number: number,
		count: 0,
	}
}

// 轮询工具
type RoundRobinUtil struct {
	// 轮询数量
	number uint
	// 轮询计数
	count uint64
}

// 获取下一个轮询索引
func (u *RoundRobinUtil) Next() uint64 {
	if u.number <= 0 {
		return 0
	}

	// 轮询计数递增,并取模
	return (atomic.AddUint64(&u.count, 1) - 1) % uint64(u.number)
}