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

type EbayCategory struct {
	// Category 请求参数
	Category models.CategoryInfoRequest

	// Errors returned error messages
	Errors []string
}

// NewCategoryService 创建service
func (ebay *EbayCategory) NewCategoryService(cat models.CategoryInfoRequest) CategoryServiceImpl {
	return &EbayCategory{
		Category: cat,
	}
}

// GetCategoryInfo Get the information about category
func (ebay *EbayCategory) GetCategoryInfo() (info models.CategoryInfo, err error) {
	start := time.Now()
	info = ebay.initCategoryInfo(ebay.Category)
	// 创建一个chrome实例
	ctx, cancel, err := ebay.createContext()
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
	url := viper.GetString("ebay.url")
	url = strings.ReplaceAll(url, "ASIN", ebay.Category.ProductNo)
	end := time.Now()
	log.Println("1. 获取weblink耗时：", end.Sub(start))

	start = time.Now()
	//err = utils.Navigate(ctx, url)
	//if err != nil {
	//	info.Status = PageError
	//	info.Errors = append(info.Errors, "1.打页面失败")
	//	return info, err
	//}
	//
	err = utils.NavigateAndWait(ctx, url, viper.GetString("ebay.content-selector"), 10*time.Second)
	if err != nil {
		log.Println("页面超时", err)
		info.Status = PageError
		info.Errors = append(info.Errors, "页面超时")
		return info, err
	}

	end = time.Now()
	log.Println("2. 打开weblink耗时：", end.Sub(start))

	// 下载html
	start = time.Now()
	html, err := ebay.downloadHtml(ctx)
	if err != nil {
		info.Status = PriceError
		info.Errors = append(info.Errors, "下载html失败")
		return info, err
	}
	end = time.Now()

	log.Println("3. 下载html耗时：", end.Sub(start))

	start = time.Now()
	// 解析html
	err = ebay.parseProductInfo(html, &info)
	if err != nil {
		info.Status = PriceError
		info.Errors = append(info.Errors, err.Error())
	}
	end = time.Now()
	log.Println("4. 解析Price耗时：", end.Sub(start))

	start = time.Now()
	// 保存截图
	filename, err := ebay.saveScreenshot(ctx)
	if err != nil {
		info.Status = ScreenshotError
		info.Errors = append(info.Errors, "保存截图失败")
		return info, err
	}
	info.Screenshot = filename
	end = time.Now()
	log.Println("5. 截图并保存总耗时：", end.Sub(start))

	if len(info.Errors) == 0 {
		info.Status = Success
	}

	return info, err
}

// 初始化返回结果
func (ebay *EbayCategory) initCategoryInfo(category models.CategoryInfoRequest) models.CategoryInfo {
	return models.CategoryInfo{
		ProductNo:    category.ProductNo,
		Country:      category.Country,
		SalesChannel: category.SalesChannel,
		PriceNo:      category.PriceNo,
		Price:        category.Price,
	}
}

// 创建一个chrome实例
func (ebay *EbayCategory) createContext() (ctx context.Context, cancel context.CancelFunc, err error) {
	if config.GlobalContext {
		return config.BrowserContext, nil, nil
	} else {
		// 创建一个chrome实例
		return utils.CreateBrowserContext(viper.GetString("chrome.url"))
	}
}

// 下载html页面
func (ebay *EbayCategory) downloadHtml(ctx context.Context) (html string, err error) {
	contentSel := viper.GetString("ebay.content-selector")

	// 获取html
	//html, err = utils.GetHtml(ctx)
	html, err = utils.GetHtmlBySelector(ctx, contentSel)
	if err != nil {
		log.Println("get html error: ", err)
		return html, err
	}

	return html, err
}

// 解析html页面，获取产品信息
func (ebay *EbayCategory) parseProductInfo(html string, info *models.CategoryInfo) error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		log.Println("goquery解析html失败", err)
		return err
	}
	priceSelectors := viper.GetString("ebay.price-selectors")
	priceSelectorsArr := strings.Split(priceSelectors, ",")

	var text string
	for _, selector := range priceSelectorsArr {
		log.Println("range selector: ", selector)
		ele := doc.Find(selector)
		if ele.Nodes != nil {
			log.Println("match selector: ", selector)
			text = ele.Text()
			break
		}
	}
	fmt.Println("text: ", text)
	text = strings.Trim(text, " \n\t")
	text = strings.ReplaceAll(text, ",", ".")

	//numbers := utils.GetFloat64sFromString(text)
	numbers := utils.GetPriceFromString(text, "EUR")

	if len(numbers) > 0 {
		// float64转string
		info.NewPrice = strconv.FormatFloat(numbers[0], 'f', -1, 64)
	} else {
		log.Println("解析价格失败", err)
		return errors.New("解析价格失败, text: " + text)
	}
	return nil
}

// 保存截图
func (ebay *EbayCategory) saveScreenshot(ctx context.Context) (filename string, err error) {
	height := viper.GetInt64("ebay.screenshot-height")
	if height == 0 {
		height = 960
	}

	start := time.Now()

	// 截图
	bytes, err := utils.GetScreenshot(ctx, height)
	if err != nil {
		ebay.Errors = append(ebay.Errors, "截图失败")
		return
	}

	end := time.Now()
	log.Println("---- 5.1 截图耗时：", end.Sub(start))

	country := ebay.Category.Country
	productNo := ebay.Category.ProductNo

	// 保存到OSS
	filename = "EBAY_O_" + country + "_" + productNo + "_" + time.Now().Format("060102150105") + ".png"

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
		ebay.Errors = append(ebay.Errors, "创建OSS客户端失败")
		return
	}

	err = ali.UploadByte(viper.GetString("oss.bucket"), filename, bytes)
	if err != nil {
		log.Println("上传截图到OSS失败", err)
		ebay.Errors = append(ebay.Errors, "上传截图到OSS失败")
		return
	}
	end = time.Now()
	log.Println("---- 5.2 上传截图到OSS耗时：", end.Sub(start))

	return filename, err
}
