package utils

import (
	"os"
	"path/filepath"
)

func CopyFile(src, dst string) error {
	// 读取源文件内容
	content, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// 创建目标文件
	err = os.WriteFile(dst, content, 0644)
	if err != nil {
		return err
	}

	return nil
}

func CopyDir(srcDir, destDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			// 如果是目录，则递归拷贝
			err = os.MkdirAll(destPath, 0755)
			if err != nil {
				return err
			}
			err = CopyFile(srcPath, destPath)
			if err != nil {
				return err
			}
		} else {
			// 如果是文件，则直接拷贝
			err = CopyFile(srcPath, destPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
