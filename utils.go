package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
)

// ExistPath 判断路径是否存在
func ExistPath(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) { //如果返回的错误类型使用os.isNotExist()判断为true，说明文件或者文件夹不存在
		return false, nil
	}
	return false, err //如果有错误了，但是不是不存在的错误，所以把这个错误原封不动的返回
}

// DeletePath 删除路径
func DeletePath(path string) error {
	return os.RemoveAll(path)
}

// NewPath 新建路径
func NewPath(path string) error {
	return os.MkdirAll(path, 0755)
}

func RunCmd(cmd string) error {
	var outerr bytes.Buffer
	c := exec.Command("/bin/sh", "-c", cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = &outerr

	if c.Run() != nil {
		return fmt.Errorf(outerr.String())
	}
	return nil
}

func getVersion(link string) (string, error) {
	// https://minecraft.azureedge.net/bin-linux/bedrock-server-1.19.52.01.zip
	re, err := regexp.Compile(`server-(.+?).zip`)
	if err != nil {
		return "", err
	}
	result := re.FindAllStringSubmatch(link, 1)

	if len(result) == 0 {
		return "", fmt.Errorf("not found version")
	}

	return result[0][1], nil
}

func NewDir(dir string) func() {
	exist, err := ExistPath(dir)
	if err != nil {
		log.Fatal(err)
	}
	if exist { // 这个目录可以安全删除
		if err := DeletePath(dir); err != nil {
			log.Fatal(err)
		}
	}
	if err := NewPath(dir); err != nil {
		log.Fatal(err)
	}

	return func() {
		_ = DeletePath(dir)
	}
}
