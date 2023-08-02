package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
	"tarantula-v2/config"
	"tarantula-v2/models"
	"tarantula-v2/ossutil"
	"tarantula-v2/utils"
	"time"
)

const (
	Success         = "SUCCESS"
	PageError       = "PAGE_ERROR"
	PriceError      = "PRICE_ERROR"
	ScreenshotError = "SCREENSHOT_ERROR"
)

type AmazonCategory struct {
	// Category 请求参数
	Category models.CategoryInfoRequest

	// Errors returned error messages
	Errors []string
}

// NewCategoryService 创建service
func (amazon *AmazonCategory) NewCategoryService(cat models.CategoryInfoRequest) CategoryServiceImpl {
	return &AmazonCategory{
		Category: cat,
	}
}

// GetCategoryInfo Get the information about category
func (amazon *AmazonCategory) GetCategoryInfo() (info models.CategoryInfo, err error) {
	// 赋值CategoryInfo
	info = amazon.initCategoryInfo(amazon.Category)

	start := time.Now()
	// 创建一个chrome实例
	ctx, cancel, err := amazon.createContext()
	if err != nil {
		info.Status = PageError
		info.Errors = append(info.Errors, err.Error())
		return info, err
	}
	// cancel 不为空则需要在运行结束后关闭ctx
	if cancel != nil {
		defer cancel()
	}

	// 获取web link
	url, err := amazon.getWebLink()
	if err != nil {
		info.Status = PageError
		info.Errors = append(info.Errors, err.Error())
		return info, err
	}
	end := time.Now()
	log.Println("1. 创建Context,准备weblink耗时(ms)：", end.Sub(start))

	start = time.Now()
	err = utils.Navigate(ctx, url)
	if err != nil {
		info.Status = PageError
		info.Errors = append(info.Errors, "打开Amazon页面失败")
		return info, err
	}
	end = time.Now()
	log.Println("2. 打开weblink耗时(ms)：", end.Sub(start))

	start = time.Now()
	// 下载html
	html, err := amazon.downloadHtml(ctx)
	if err != nil {
		info.Status = PriceError
		info.Errors = append(info.Errors, "下载html失败")
		return info, err
	}

	// 解析html
	err = amazon.parseProductInfo(html, &info)
	if err != nil {
		info.Status = PriceError
		info.Errors = append(info.Errors, err.Error())
		return info, err
	}
	end = time.Now()
	log.Println("3. 下载html并解析Price耗时(ms)：", end.Sub(start))

	start = time.Now()
	// 保存截图
	filename, err := amazon.saveScreenshot(ctx)
	if err != nil {
		info.Status = ScreenshotError
		info.Errors = append(info.Errors, "保存截图失败")
		return info, err
	}
	info.Screenshot = filename
	end = time.Now()
	log.Println("4. 截图并保存总耗时(s)：", end.Sub(start).Seconds())

	if len(info.Errors) == 0 {
		info.Status = "SUCCESS"
	}

	return info, err
}

// 初始化返回结果
func (amazon *AmazonCategory) initCategoryInfo(category models.CategoryInfoRequest) models.CategoryInfo {
	return models.CategoryInfo{
		ProductNo:    category.ProductNo,
		Country:      category.Country,
		SalesChannel: category.SalesChannel,
		PriceNo:      category.PriceNo,
		Price:        category.Price,
	}
}

func (amazon *AmazonCategory) createContext() (ctx context.Context, cancel context.CancelFunc, err error) {
	if config.GlobalContext {
		return config.BrowserContext, nil, nil
	} else {
		// 创建一个chrome实例
		return utils.CreateBrowserContext(viper.GetString("chrome.url"))
	}
}

// getWebLink Get the web link of category
func (amazon *AmazonCategory) getWebLink() (string, error) {
	country := amazon.Category.Country
	productNo := amazon.Category.ProductNo

	url := viper.GetString("amazon.urls." + strings.ToLower(country))

	if url == "" {
		return "", errors.New("Can't get web link for country:" + country)
	}
	return strings.ReplaceAll(url, "ASIN", productNo), nil
}

// 下载html页面
func (amazon *AmazonCategory) downloadHtml(ctx context.Context) (html string, err error) {
	// 获取html
	html, err = utils.GetHtml(ctx)
	if err != nil {
		log.Println("get html error: ", err)
		return html, err
	}

	return html, err
}

// 截图并保存OSS
func (amazon *AmazonCategory) saveScreenshot(ctx context.Context) (filename string, err error) {
	height, err := utils.GetElementBottomRightHeight(ctx, viper.GetString("amazon.description-selectors"))
	if err != nil {
		log.Println("获取description元素高度失败，将截图全屏", err)
	}

	start := time.Now()
	// 截图
	bytes, err := utils.GetScreenshot(ctx, height)
	if err != nil {
		amazon.Errors = append(amazon.Errors, "截图失败")
		return
	}
	end := time.Now()
	log.Println("---- 4.1 截图耗时(s)：", end.Sub(start).Seconds())

	country := amazon.Category.Country
	productNo := amazon.Category.ProductNo

	// 保存到OSS
	filename = "AMAZON_O_" + country + "_" + productNo + "_" + time.Now().Format("060102150105") + ".png"

	if viper.GetBool("save-screenshot-on-disk") {
		err := os.WriteFile(filename, bytes, 0644)
		if err != nil {
			fmt.Println("开启调试模式，保存截图到磁盘失败，文件名：", filename, err)
		}
	}

	log.Println("开始上传截图到OSS...")
	start = time.Now()
	ali, err := ossutil.NewAliOss(viper.GetString("oss.endpoint"), viper.GetString("oss.access-key-id"), viper.GetString("oss.access-key-secret"))
	if err != nil {
		log.Println("创建OSS客户端失败", err)
		amazon.Errors = append(amazon.Errors, "创建OSS客户端失败")
		return
	}

	err = ali.UploadByte(viper.GetString("oss.bucket"), filename, bytes)
	if err != nil {
		log.Println("上传截图到OSS失败", err)
		amazon.Errors = append(amazon.Errors, "上传截图到OSS失败")
		return
	}
	end = time.Now()
	log.Println("---- 4.2 上传截图到OSS耗时(s)：", end.Sub(start).Seconds())

	return filename, err
}

// 解析商品信息
func (amazon *AmazonCategory) parseProductInfo(html string, info *models.CategoryInfo) (err error) {
	// goquery 解析html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("goquery解析html失败", err)
		return err
	}
	priceSelectors := viper.GetString("amazon.price-selectors")
	priceSelectorsArr := strings.Split(priceSelectors, ",")

	var text string
	for _, selector := range priceSelectorsArr {
		ele := doc.Find(selector)
		if ele != nil {
			text = ele.Text()
			break
		}
	}
	numbers := utils.GetFloat64sFromString(text)

	if len(numbers) > 0 {
		// float64转string
		info.NewPrice = strconv.FormatFloat(numbers[0], 'f', -1, 64)
	} else {
		log.Println("解析价格失败", err)
		return errors.New("解析价格失败, text: " + text)
	}
	return nil
}
