package service

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"tarantula-v2/models"
)

// NewCategoryService 创建service
func NewCategoryService(cat models.CategoryInfoRequest) CategoryServiceImpl {
	fmt.Println("debug:", viper.GetBool("debug"))
	channel := cat.SalesChannel
	switch channel {
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
