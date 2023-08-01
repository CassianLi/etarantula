package utils

import (
	"encoding/base64"
	"os"
)

// SaveBase64ToImage saves the base64 string to an image
func SaveBase64ToImage(baseEncodeStr string, savePath string) error {
	data, err := base64.StdEncoding.DecodeString(baseEncodeStr)
	if err != nil {
		return err
	}
	err = os.WriteFile(savePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// MergeBase64Image 将两张base64图片拼接为一张图片
func MergeBase64Image(base64Img1, base64Img2 string) (string, error) {
	img1, err := base64.StdEncoding.DecodeString(base64Img1)
	if err != nil {
		return "", err
	}
	img2, err := base64.StdEncoding.DecodeString(base64Img2)
	if err != nil {
		return "", err
	}
	img := append(img1, img2...)
	return base64.StdEncoding.EncodeToString(img), nil
}

// MergeBase64ImageBytes 将两张base64图片拼接为一张图片
func MergeBase64ImageBytes(base64Img1, base64Img2 []byte) []byte {
	return append(base64Img1, base64Img2...)
}
