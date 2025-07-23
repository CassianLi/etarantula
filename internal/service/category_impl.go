package service

import "etarantula/internal/models"

type CategoryServiceImpl interface {

	// GetCategoryInfo Get the information about category
	GetCategoryInfo() (info models.CategoryInfo, err error)
}
