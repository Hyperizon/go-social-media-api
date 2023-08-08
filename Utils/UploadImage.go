package Utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	fileExtension := filepath.Ext(fileHeader.Filename)
	newFileName := fmt.Sprintf("%d%s", time.Now().Unix(), fileExtension)

	uploadDir := "./Public/Images"
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	filePath := filepath.Join(uploadDir, newFileName)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
