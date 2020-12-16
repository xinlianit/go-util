package util

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
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

// 创建目录
// @param dirPath 目录路径
// @param recursive 递归创建目录
func (u File) CreateDir(dirPath string, recursive bool) error {
	if !u.FileExist(dirPath) {
		if recursive {
			return os.MkdirAll(dirPath, os.ModePerm)
		} else {
			return os.Mkdir(dirPath, os.ModePerm)
		}
	}
	return nil
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

// 列出目录文件
// @param dirPath 目录路径
// @param readChildDir 读取子目录
func (u File) LsDir(dirPath string, readChildDir bool) ([]string, error) {

	var fileList []string

	rd, err := ioutil.ReadDir(dirPath)

	if err != nil {
		return fileList, err
	}

	for _, fi := range rd {
		if fi.IsDir() {
			// 递归读取子目录
			if readChildDir {
				childFileList, err := u.LsDir(strings.TrimRight(dirPath, "/")+"/"+fi.Name(), readChildDir)
				if err != nil {
					return fileList, err
				}

				for _, childFile := range childFileList {
					fileList = append(fileList, fi.Name()+"/"+childFile)
				}
			}
		} else {
			fileList = append(fileList, fi.Name())
		}
	}

	return fileList, nil
}

// 递归拷贝目录
// @param src 源目录
// @param dst 目标目录
func (u File) CopyDir(src string, dst string) error {
	var err error
	var fds []os.FileInfo
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}
	for _, fd := range fds {
		srcFp := path.Join(src, fd.Name())
		dstFp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = u.CopyDir(srcFp, dstFp); err != nil {
				return err
			}
		} else {
			if err = u.CopyFile(srcFp, dstFp); err != nil {
				return err
			}
		}
	}
	return nil
}

// 递归拷贝目录
// @param src 源目录
// @param dst 目标目录
// @param filePrefix 文件前缀
func (u File) CopyDirWithFix(src string, dst string, filePrefix string) error {
	var err error
	var fds []os.FileInfo
	var srcInfo os.FileInfo

	if srcInfo, err = os.Stat(src); err != nil {
		return err
	}

	if err = os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	if fds, err = ioutil.ReadDir(src); err != nil {
		return err
	}

	for _, fd := range fds {
		srcFp := path.Join(src, fd.Name())
		dstFp := path.Join(dst, fd.Name())

		if fd.IsDir() {
			if err = u.CopyDirWithFix(srcFp, dstFp, filePrefix); err != nil {
				return err
			}
		} else {
			if err = u.CopyFileWithFix(srcFp, dstFp, filePrefix); err != nil {
				return err
			}
		}
	}
	return nil
}
