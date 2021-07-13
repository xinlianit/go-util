package util

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var (
	fileUtilInstance *File
	fileUtilOnce     sync.Once
)

func FileUtil() *File {
	fileUtilOnce.Do(func() {
		fileUtilInstance = new(File)
	})
	return fileUtilInstance
}

// 文件工具
type File struct {
}

// 文件是否存在
// @param file 文件
func (u File) FileExist(file string) bool {
	_, err := os.Stat(file)
	// 判断报错信息是否为文件不存在
	return err == nil || os.IsExist(err)

}

// 读取文件内容
// @param fileName 文件名
func (u File) Read(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	fd, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(fd), nil
}

// 文件写入
// @param string fileName 文件名
func (u File) Write(fileName string, context string) (int, error) {
	// 打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// 写入文件
	len, err := f.WriteString(context)
	if err != nil {
		return 0, err
	}

	return len, nil
}

// 文件追加写入
// @param string fileName 文件名
func (u File) AppendWrite(fileName string, context string) (int, error) {
	// 打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	// 写入文件
	len, err := f.WriteString(context)
	if err != nil {
		return 0, err
	}

	return len, nil
}

// 拷贝文件
// @param src 源文件
// @param dst 目标文件
func (u File) CopyFile(src string, dst string) error {
	var err error
	var srcFd *os.File
	var dstFd *os.File
	var srcInfo os.FileInfo

	if srcFd, err = os.Open(src); err != nil {
		return err
	}
	defer srcFd.Close()

	if dstFd, err = os.Create(dst); err != nil {
		return err
	}
	defer dstFd.Close()

	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

// 拷贝文件
// @param src 源文件
// @param dst 目标文件
// @param prefix 文件名前缀
func (u File) CopyFileWithFix(src string, dst string, prefix string) error {
	var err error
	var srcFd *os.File
	var dstFd *os.File
	var srcInfo os.FileInfo

	if srcFd, err = os.Open(src); err != nil {
		return err
	}
	defer srcFd.Close()

	newDst := filepath.Join(filepath.Dir(dst), prefix + filepath.Base(dst))
	if dstFd, err = os.Create(newDst); err != nil {
		return err
	}
	defer dstFd.Close()

	if _, err = io.Copy(dstFd, srcFd); err != nil {
		return err
	}
	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}
	return os.Chmod(newDst, srcInfo.Mode())
}
