package config

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
)

const (
	ChromeRemoteDebuggingUrl = "http://localhost:9222"
)

var (
	// GlobalContext 全局上下文
	GlobalContext bool

	// BrowserContext 浏览器上下文
	BrowserContext context.Context
	// BrowserCancel 浏览器上下文取消函数
	BrowserCancel context.CancelFunc
)

// InitBrowserContext 初始化浏览器上下文
func InitBrowserContext(debugUrl string) error {
	if debugUrl == "" {
		debugUrl = ChromeRemoteDebuggingUrl
	}
	fmt.Println("初始化浏览器上下文...")
	// 创建一个chrome实例
	BrowserContext, BrowserCancel = chromedp.NewRemoteAllocator(context.Background(), debugUrl)

	// create a new chrome instance
	BrowserContext, BrowserCancel = chromedp.NewContext(BrowserContext)

	return BrowserContext.Err()
}
