package test

import (
	"github.com/xinlianit/go-util"
	"testing"
)

func TestRoundRobin(t *testing.T) {
	list := []string{"192.168.1.1","192.168.1.2","192.168.1.3","192.168.1.4"}
	roundRobinUtil := util.NewRoundRobinUtil(uint(len(list)))

	for i:=0; i < 10; i++ {
		n := roundRobinUtil.Next()
		t.Log(list[n])
	}
}
