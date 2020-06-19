package tools

import (
	"archive/zip"
	"bytes"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func PathCreate(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func FileCreate(content bytes.Buffer, name string) {
	file, err := os.Create(name)
	if err != nil {
		log.Println(err)
	}
	file.WriteString(content.String())
	//for i := 0; i < len(content); i++ {
	//	//写入byte的slice数据
	//	file.Write(content)
	//	//写入字符串
	//	//
	//}
	file.Close()
}

func PathRemove(name string) {
	err := os.RemoveAll(name)
	if err != nil {
		log.Println(err)
	}
}

func FileRemove(name string) {
	err := os.Remove(name)
	if err != nil {
		log.Println(err)
	}
}

func FileZip(dst, src string,notContPath string) (err error) {
	//创建准备写入的文件
	fw, err := os.Create(dst)
	defer fw.Close()
	if err != nil {
		return err
	}

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			log.Fatalln(err)
		}
	}()
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return
		}

		fh.Name = strings.TrimPrefix(path, string(filepath.Separator))

		if fi.IsDir() {
			fh.Name += "/"
		}
		fh.Name = strings.Replace(fh.Name, notContPath, "", -1)

		w, err := zw.CreateHeader(fh)
		if err != nil {
			return
		}

		if !fh.Mode().IsRegular() {
			return nil
		}

		fr, err := os.Open(path)
		defer fr.Close()
		if err != nil {
			return
		}

		n, err := io.Copy(w, fr)
		if err != nil {
			return
		}

		log.Printf("成功压缩文件： %s, 共写入了 %d 个字符的数据\n", path, n)

		return nil
	})
}
