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
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _2026052000020EduNewsAIMenus)
}

func _2026052000020EduNewsAIMenus(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		menus := []migrationModels.SysMenu{
			{
				MenuId:     9011,
				MenuName:   "EduNews",
				Title:      "资讯管理",
				Icon:       "IconFile",
				Path:       "/edu/news",
				Paths:      "/0/9000/9011",
				MenuType:   "C",
				Action:     "无",
				Permission: "edu:news:list",
				ParentId:   9000,
				Component:  "/edu/news/index",
				Sort:       80,
				Visible:    "0",
				IsFrame:    "1",
			},
			{
				MenuId:     9012,
				MenuName:   "EduAIMonitor",
				Title:      "AI 监控",
				Icon:       "IconRobot",
				Path:       "/edu/ai",
				Paths:      "/0/9000/9012",
				MenuType:   "C",
				Action:     "无",
				Permission: "edu:ai:list",
				ParentId:   9000,
				Component:  "/edu/ai/index",
				Sort:       90,
				Visible:    "0",
				IsFrame:    "1",
			},
		}
		for _, menu := range menus {
			var count int64
			if err := tx.Model(&migrationModels.SysMenu{}).Where("menu_id = ?", menu.MenuId).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				if err := tx.Create(&menu).Error; err != nil {
					return err
				}
			}
		}

		buttons := []migrationModels.SysMenu{
			eduLocalButtonMenu(9181, 9011, 1, "EduNewsQuery", "查询资讯", "edu:news:query"),
			eduLocalButtonMenu(9182, 9011, 2, "EduNewsAdd", "新增资讯", "edu:news:add"),
			eduLocalButtonMenu(9183, 9011, 3, "EduNewsEdit", "编辑资讯", "edu:news:edit"),
			eduLocalButtonMenu(9184, 9011, 4, "EduNewsRemove", "删除资讯", "edu:news:remove"),
			eduLocalButtonMenu(9185, 9012, 1, "EduAIQuery", "查询 AI 会话", "edu:ai:query"),
			eduLocalButtonMenu(9186, 9012, 2, "EduAIRemove", "删除 AI 会话", "edu:ai:remove"),
		}
		for _, button := range buttons {
			var count int64
			if err := tx.Model(&migrationModels.SysMenu{}).Where("menu_id = ?", button.MenuId).Count(&count).Error; err != nil {
				return err
			}
			if count == 0 {
				if err := tx.Create(&button).Error; err != nil {
					return err
				}
			}
		}
		return tx.Create(&common.Migration{Version: version}).Error
	})
}

func eduLocalButtonMenu(menuId int, parentId int, sort int, name string, title string, permission string) migrationModels.SysMenu {
	return migrationModels.SysMenu{
		MenuId:     menuId,
		MenuName:   name,
		Title:      title,
		MenuType:   "F",
		Action:     "无",
		Permission: permission,
		ParentId:   parentId,
		Sort:       sort,
		Visible:    "0",
		IsFrame:    "1",
	}
}
