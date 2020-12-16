package test

import (
	"github.com/xinlianit/go-util/util"
	"testing"
)

func TestPrintf(t *testing.T) {
	len, err := util.FileUtil().Write("./test.log", "test")
	if err != nil {
		t.Errorf("file write error: %v", err)
	}
	t.Log("file write success, len: ", len)
}
