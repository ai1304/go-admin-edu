package search

import (
	"go-admin/app/edu/models"

	"gorm.io/gorm"
)

type ResourceQuery struct {
	Keyword          string
	Status           string
	Sort             string
	TagId            int
	SchoolId         int
	StageCategoryId  int
	DisabilityTypeId int
	ResourceTypeId   int
	AbilityDomainId  int
	TopicCategoryId  int
	PageIndex        int
	PageSize         int
}

type ResourceResult struct {
	List  []models.EduResource
	Count int64
}

type ResourceSearcher interface {
	SearchResources(query ResourceQuery) (ResourceResult, error)
}

type MySQLSearcher struct {
	DB *gorm.DB
}

func NewMySQLSearcher(db *gorm.DB) ResourceSearcher {
	return MySQLSearcher{DB: db}
}

func (s MySQLSearcher) SearchResources(query ResourceQuery) (ResourceResult, error) {
	if query.PageIndex <= 0 {
		query.PageIndex = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	list := make([]models.EduResource, 0)
	db := applyResourceFilters(s.DB.Model(&models.EduResource{}), query)
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return ResourceResult{}, err
	}
	if err := db.Order(resourceOrder(query.Sort)).
		Limit(query.PageSize).
		Offset((query.PageIndex - 1) * query.PageSize).
		Find(&list).Error; err != nil {
		return ResourceResult{}, err
	}
	return ResourceResult{List: list, Count: count}, nil
}

func applyResourceFilters(db *gorm.DB, query ResourceQuery) *gorm.DB {
	if query.Keyword != "" {
		like := "%" + query.Keyword + "%"
		db = db.Where("title like ? or summary like ? or keywords like ? or author_name like ?", like, like, like, like)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.SchoolId != 0 {
		db = db.Where("school_id = ?", query.SchoolId)
	}
	if query.StageCategoryId != 0 {
		db = db.Where("stage_category_id = ?", query.StageCategoryId)
	}
	if query.DisabilityTypeId != 0 {
		db = db.Where("disability_type_id = ?", query.DisabilityTypeId)
	}
	if query.ResourceTypeId != 0 {
		db = db.Where("resource_type_id = ?", query.ResourceTypeId)
	}
	if query.AbilityDomainId != 0 {
		db = db.Where("ability_domain_id = ?", query.AbilityDomainId)
	}
	if query.TopicCategoryId != 0 {
		db = db.Where("topic_category_id = ?", query.TopicCategoryId)
	}
	if query.TagId != 0 {
		db = db.Where("id in (?)", db.Session(&gorm.Session{}).
			Model(&models.EduResourceTagRelation{}).
			Select("resource_id").
			Where("tag_id = ?", query.TagId))
	}
	return db
}

func resourceOrder(sort string) string {
	switch sort {
	case "view":
		return "view_count desc,id desc"
	case "download":
		return "download_count desc,id desc"
	case "favorite":
		return "favorite_count desc,id desc"
	default:
		return "id desc"
	}
}
