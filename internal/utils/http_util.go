package utils

import (
	"io"
	"net/http"
)

// DownloadHtml Download html from url
func DownloadHtml(url string) (body string, err error) {
	// 构造HTTP GET请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept-Language", "en-US")
	// 设置请求头中的User-Agent，模拟真实的浏览器访问
	req.Header.Set("User-Agent", GetUserAgent())

	// 发送HTTP请求，并获取响应
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 从响应中获取HTML源代码
	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(htmlBytes), nil
}
