package service

import (
	"etarantula/internal/models"
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// NewCategoryService 创建service
func NewCategoryService(cat models.CategoryInfoRequest) CategoryServiceImpl {
	fmt.Println("debug:", viper.GetBool("debug"))
	channel := cat.SalesChannel
	switch channel {
	// tarantula2 不再提供Amazon的获取方式，采用tarantula3 通过亚马逊计算器页面获取商品费率信息
	case "amazon":
		return &AmazonCategory{
			Category: cat,
		}
	case "ebay":
		return &EbayCategory{
			Category: cat,
		}
	default:
		log.Println("Can't find the sales channel: ", channel)
		return nil
	}
}
