package service

import "tarantula-v2/models"

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

	return info, err
}
