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
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _2026051500010EduMenus)
}

func _2026051500010EduMenus(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		menus := []migrationModels.SysMenu{
			{
				MenuId:    9000,
				MenuName:  "Edu",
				Title:     "特殊教育",
				Icon:      "IconBook",
				Path:      "/edu",
				Paths:     "/0/9000",
				MenuType:  "M",
				Action:    "无",
				ParentId:  0,
				Component: "Layout",
				Sort:      15,
				Visible:   "0",
				IsFrame:   "1",
			},
			{
				MenuId:     9001,
				MenuName:   "EduStats",
				Title:      "数据概览",
				Icon:       "IconDashboard",
				Path:       "/edu/stats",
				Paths:      "/0/9000/9001",
				MenuType:   "C",
				Action:     "无",
				Permission: "edu:stats:list",
				ParentId:   9000,
				Component:  "/edu/stats/index",
				Sort:       1,
				Visible:    "0",
				IsFrame:    "1",
			},
			{
				MenuId:     9002,
				MenuName:   "EduResource",
				Title:      "资源管理",
				Icon:       "IconFile",
				Path:       "/edu/resource",
				Paths:      "/0/9000/9002",
				MenuType:   "C",
				Action:     "无",
				Permission: "edu:resource:list",
				ParentId:   9000,
				Component:  "/edu/resource/index",
				Sort:       10,
				Visible:    "0",
				IsFrame:    "1",
			},
			{
				MenuId:     9003,
				MenuName:   "EduCourse",
				Title:      "课程管理",
				Icon:       "IconBook",
				Path:       "/edu/course",
				Paths:      "/0/9000/9003",
				MenuType:   "C",
				Action:     "无",
				Permission: "edu:course:list",
				ParentId:   9000,
				Component:  "/edu/course/index",
				Sort:       20,
				Visible:    "0",
				IsFrame:    "1",
			},
			{
				MenuId:     9004,
				MenuName:   "EduActivity",
				Title:      "活动管理",
				Icon:       "IconCalendar",
				Path:       "/edu/activity",
				Paths:      "/0/9000/9004",
				MenuType:   "C",
				Action:     "无",
				Permission: "edu:activity:list",
				ParentId:   9000,
				Component:  "/edu/activity/index",
				Sort:       30,
				Visible:    "0",
				IsFrame:    "1",
			},
			{
				MenuId:     9005,
				MenuName:   "EduCase",
				Title:      "个案管理",
				Icon:       "IconUserGroup",
				Path:       "/edu/case",
				Paths:      "/0/9000/9005",
				MenuType:   "C",
				Action:     "无",
				Permission: "edu:case:list",
				ParentId:   9000,
				Component:  "/edu/case/index",
				Sort:       40,
				Visible:    "0",
				IsFrame:    "1",
			},
			{
				MenuId:     9006,
				MenuName:   "EduExpert",
				Title:      "专家管理",
				Icon:       "IconUser",
				Path:       "/edu/expert",
				Paths:      "/0/9000/9006",
				MenuType:   "C",
				Action:     "无",
				Permission: "edu:expert:list",
				ParentId:   9000,
				Component:  "/edu/expert/index",
				Sort:       50,
				Visible:    "0",
				IsFrame:    "1",
			},
			{
				MenuId:     9007,
				MenuName:   "EduSchool",
				Title:      "学校管理",
				Icon:       "IconHome",
				Path:       "/edu/school",
				Paths:      "/0/9000/9007",
				MenuType:   "C",
				Action:     "无",
				Permission: "edu:school:list",
				ParentId:   9000,
				Component:  "/edu/school/index",
				Sort:       60,
				Visible:    "0",
				IsFrame:    "1",
			},
			{
				MenuId:     9008,
				MenuName:   "EduRegion",
				Title:      "区域管理",
				Icon:       "IconLocation",
				Path:       "/edu/region",
				Paths:      "/0/9000/9008",
				MenuType:   "C",
				Action:     "无",
				Permission: "edu:region:list",
				ParentId:   9000,
				Component:  "/edu/region/index",
				Sort:       70,
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
