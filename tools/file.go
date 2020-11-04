package tools

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func PathCreate(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

// PathExist 判断目录是否存在
func PathExist(addr string) bool {
	s, err := os.Stat(addr)
	if err != nil {
		log.Println(err)
		return false
	}
	return s.IsDir()
}

func FileCreate(content bytes.Buffer, name string) {
	file, err := os.Create(name)
	if err != nil {
		log.Println(err)
	}
	_, err = file.WriteString(content.String())
	if err != nil {
		log.Println(err)
	}
	file.Close()
}

type ReplaceHelper struct {
	Root    string //路径
	OldText string //需要替换的文本
	NewText string //新的文本
}

func (h *ReplaceHelper) DoWrok() error {

	return filepath.Walk(h.Root, h.walkCallback)

}

func (h ReplaceHelper) walkCallback(path string, f os.FileInfo, err error) error {

	if err != nil {
		return err
	}
	if f == nil {
		return nil
	}
	if f.IsDir() {
		log.Println("DIR:", path)
		return nil
	}

	//文件类型需要进行过滤

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		//err
		return err
	}
	content := string(buf)
	log.Printf("h.OldText: %s \n", h.OldText)
	log.Printf("h.NewText: %s \n", h.NewText)

	//替换
	newContent := strings.Replace(content, h.OldText, h.NewText, -1)

	//重新写入
	ioutil.WriteFile(path, []byte(newContent), 0)

	return err
}

func FileMonitoringById(ctx context.Context, filePth string, id string, group string, hookfn func(context.Context, string, string, []byte)) {
	f, err := os.Open(filePth)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	rd := bufio.NewReader(f)
	f.Seek(0, 2)
	for {
		if ctx.Err() != nil {
			break
		}
		line, err := rd.ReadBytes('\n')
		// 如果是文件末尾不返回
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
			continue
		} else if err != nil {
			log.Fatalln(err)
		}
		go hookfn(ctx, id, group, line)
	}
}

// 获取文件大小
func GetFileSize(filename string) int64 {
	var result int64
	filepath.Walk(filename, func(path string, f os.FileInfo, err error) error {
		result = f.Size()
		return nil
	})
	return result
}

//获取当前路径，比如：E:/abc/data/test
func GetCurrentPath() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
