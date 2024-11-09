package filea

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// 文件是否存在
func FileExists(file string) bool {
	_, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// 目录是否存在
func DirExists(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			return true
		}
	} else if !info.IsDir() {
		return false
	} else {
		return true
	}
}

// 创建目录
func Mkdir(dir string) error {
	if !DirExists(dir) {
		err := os.MkdirAll(dir, 0666)
		if err != nil {
			return err
		}
		err = os.Chmod(dir, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

// 创建文件
func CreateFile(filename string) error {
	if !FileExists(filename) {
		file, err := os.Create(filename)
		if err != nil {
			return err
		}
		file.Close()
		return os.Chmod(filename, 0666)
	} else {
		return os.Chmod(filename, 0666)
	}

}

// 读取文件内容
func ReadFile(filename string) (string, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// 读取文件内容 返回文本切片
func ReadFileSlice(filename string) ([]string, error) {
	var s []string
	file, err := os.Open(filename)
	if err != nil {
		return s, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s = append(s, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return s, err
	}
	return s, nil
}

// 写文件内容 覆盖
func WriteFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	// 确保在函数结束时关闭文件
	defer file.Close()
	_, err = io.WriteString(file, content)
	if err != nil {
		return err
	}
	return nil
}
func WriteFileByte(filename string, content []byte) error {
	// 打开文件，设置为只写、截断和创建
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	// 确保在函数结束时关闭文件
	defer file.Close()
	// 写入内容
	_, err = file.Write(content)
	if err != nil {
		return err
	}
	return nil
}

// 写文件内容 追加
func WriteFileAppEnd(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}
	return nil

}

// 解压 ZIP 文件到指定目录
func Unzip(src, dest string) error {
	err := Mkdir(dest)
	if err != nil {
		return err
	}
	// 打开 ZIP 文件
	zipReader, err := zip.OpenReader(src)
	if err != nil {
		return fmt.Errorf("无法打开 ZIP 文件: %v", err)
	}
	defer zipReader.Close()

	// 遍历 ZIP 文件中的每一个文件
	for _, file := range zipReader.File {
		// 创建目标文件的完整路径
		destPath := filepath.Join(dest, file.Name)

		// 判断文件是否是目录
		if file.FileInfo().IsDir() {
			// 如果是目录，创建目录
			if err := os.MkdirAll(destPath, file.Mode()); err != nil {
				return fmt.Errorf("无法创建目录 %v: %v", destPath, err)
			}
			continue
		}

		// 如果是文件，解压文件
		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("无法创建文件 %v: %v", destPath, err)
		}
		defer destFile.Close()

		// 打开 ZIP 文件中的文件
		zipFile, err := file.Open()
		if err != nil {
			return fmt.Errorf("无法打开 ZIP 中的文件: %v", err)
		}
		defer zipFile.Close()

		// 将文件内容拷贝到目标文件
		if _, err := io.Copy(destFile, zipFile); err != nil {
			return fmt.Errorf("解压文件时出错: %v", err)
		}
	}

	return nil
}
