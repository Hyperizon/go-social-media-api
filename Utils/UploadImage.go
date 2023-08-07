package Utils

import (
	"encoding/base64"
	"io"
	"mime/multipart"
)

func UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()

	if err != nil {
		return "", err
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	imageBase64 := base64.StdEncoding.EncodeToString(imageBytes)
	return imageBase64, nil
}
