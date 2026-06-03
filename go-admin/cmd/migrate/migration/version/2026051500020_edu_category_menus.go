package version

import (
	"go-admin/cmd/migrate/migration"
	migrationModels "go-admin/cmd/migrate/migration/models"
	common "go-admin/common/models"
	"runtime"

	"gorm.io/gorm"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _2026051500020EduCategoryMenus)
}

func _2026051500020EduCategoryMenus(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		menus := []migrationModels.SysMenu{
			{
				MenuId:     9009,
				MenuName:   "EduCategory",
				Title:      "资源分类",
				Icon:       "IconApps",
				Path:       "/edu/category",
				Paths:      "/0/9000/9009",
				MenuType:   "C",
				Action:     "无",
				Permission: "edu:category:list",
				ParentId:   9000,
				Component:  "/edu/category/index",
				Sort:       11,
				Visible:    "0",
				IsFrame:    "1",
			},
			{
				MenuId:     9010,
				MenuName:   "EduTag",
				Title:      "资源标签",
				Icon:       "IconTags",
				Path:       "/edu/tag",
				Paths:      "/0/9000/9010",
				MenuType:   "C",
				Action:     "无",
				Permission: "edu:tag:list",
				ParentId:   9000,
				Component:  "/edu/tag/index",
				Sort:       12,
				Visible:    "0",
				IsFrame:    "1",
			},
		}

		for _, menu := range menus {
			var count int64
			if err := tx.Model(&migrationModels.SysMenu{}).Where("menu_id = ?", menu.MenuId).Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				continue
			}
			if err := tx.Create(&menu).Error; err != nil {
				return err
			}
		}

		return tx.Create(&common.Migration{Version: version}).Error
	})
}
