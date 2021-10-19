package file

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/ville-vv/vilgo/vfile"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

// AutoSplitFile 自动根据文件大小分割文件，在写文件前会判断一下当前文件大小，如果超过设定值会先分割一下
// 如果在写入之前没有超过设定的大小，会正常写入，在下一次写入的时候会判断文件大小，所有可能出现文件分割不是精准的
type AutoSplitFile struct {
	*os.File
	filePath           string
	splitIdx           int
	singleFileMaxSize  int64 // 单个文件大小
	dataFileName       string
	fileType           string
	FileHeaderWriteFun func(w io.Writer) error
}

// NewAutoSplitFile
// fileName 文件名称，可以包含文件路径
// maxSize 文档大小最小单位是 byte 字节  *1024 为kb  *1024*1024 为mb
func NewAutoSplitFile(fileName string, maxSize int64) (*AutoSplitFile, error) {
	dirPath, fileName := path.Split(fileName)
	dirPath = strings.TrimSpace(dirPath)

	if dirPath != "" {
		if !vfile.PathExists(dirPath) {
			err := os.Mkdir(dirPath, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
	}
	fileType := path.Ext(fileName)
	if len(fileType) > 0 {
		fileName = fileName[:len(fileName)-len(fileType)]
	}

	cntPath := dirPath
	if cntPath == "" {
		cntPath = "./"
	}

	fileInfoList, err := ioutil.ReadDir(cntPath)
	if err != nil {
		return nil, err
	}

	af := &AutoSplitFile{
		dataFileName:      fileName,
		filePath:          dirPath,
		singleFileMaxSize: maxSize,
		splitIdx:          len(fileInfoList),
		fileType:          fileType,
	}

	return af, nil
}

func (sel *AutoSplitFile) Init() error {
	return sel.createDataFile()
}

// SetSingleFileMaxSize
func (sel *AutoSplitFile) SetSingleFileMaxSize(singleFileMaxSize int64) {
	sel.singleFileMaxSize = singleFileMaxSize
}

func (sel *AutoSplitFile) SetMaxSizeMb(size int64) {
	sel.singleFileMaxSize = size * 1024 * 1024
}

func (sel *AutoSplitFile) SetMaxSizeKb(size int64) {
	sel.singleFileMaxSize = size * 1024
}

// Write
func (sel *AutoSplitFile) Write(b []byte) (n int, err error) {
	if err := sel.split(); err != nil {
		return 0, err
	}
	return sel.File.Write(b)
}

func (sel *AutoSplitFile) WriteString(s string) (n int, err error) {
	if err := sel.split(); err != nil {
		return 0, err
	}
	return sel.File.WriteString(s)
}

func (sel *AutoSplitFile) split() error {
	info, err := sel.File.Stat()
	if err != nil {
		return err
	}
	size := info.Size()
	if size > sel.singleFileMaxSize && sel.singleFileMaxSize > 0 {
		_ = sel.File.Close()
		oldName := path.Join(sel.filePath, info.Name())
		newName := path.Join(sel.filePath, fmt.Sprintf("%s_%s%s", sel.dataFileName, time.Now().Format("060102150405"), sel.fileType))
		err := os.Rename(oldName, newName)
		if err != nil {
			return errors.Wrap(err, "rename data file")
		}
		sel.splitIdx++
		return sel.createDataFile()
	}
	return nil
}

func (sel *AutoSplitFile) splitNewFileName() string {
	return path.Join(sel.filePath, fmt.Sprintf("%s_%s%s", sel.dataFileName, time.Now().Format("060102150405"), sel.fileType))
}

func (sel *AutoSplitFile) createDataFile() error {
	fileName := path.Join(sel.filePath, sel.dataFileName+sel.fileType)
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	sel.File = f
	stat, err := sel.File.Stat()
	if err != nil {
		return err
	}

	//fmt.Println("file size", stat.Size())
	if stat.Size() == 0 {
		if sel.FileHeaderWriteFun != nil {
			return sel.FileHeaderWriteFun(f)
		}
	}
	return nil
}
