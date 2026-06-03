package version

import (
	"go-admin/app/edu/models"
	"go-admin/cmd/migrate/migration"
	migrationModels "go-admin/cmd/migrate/migration/models"
	common "go-admin/common/models"
	"runtime"

	"gorm.io/gorm"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _2026060200010EduRecruit)
}

func _2026060200010EduRecruit(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.AutoMigrate(
			new(models.EduRecruitCompany),
			new(models.EduRecruitJob),
			new(models.EduRecruitReview),
		); err != nil {
			return err
		}
		menus := []migrationModels.SysMenu{
			recruitMenu(9400, 0, 16, "Recruit", "招聘管理", "IconBriefcase", "/recruit", "/0/9400", "M", "", "Layout", "recruit:manage"),
			recruitMenu(9401, 9400, 1, "RecruitReview", "信息审核", "IconCheckCircle", "/recruit/review", "/0/9400/9401", "C", "/recruit/review/index", "/recruit/review/index", "recruit:review:list"),
			recruitMenu(9402, 9400, 2, "RecruitCompany", "企业管理", "IconHome", "/recruit/company", "/0/9400/9402", "C", "/recruit/company/index", "/recruit/company/index", "recruit:company:list"),
			recruitMenu(9403, 9400, 3, "RecruitJob", "岗位管理", "IconFile", "/recruit/job", "/0/9400/9403", "C", "/recruit/job/index", "/recruit/job/index", "recruit:job:list"),
		}
		for i := range menus {
			if err := tx.Where("menu_id = ?", menus[i].MenuId).Assign(menus[i]).FirstOrCreate(&menus[i]).Error; err != nil {
				return err
			}
		}
		return tx.Create(&common.Migration{Version: version}).Error
	})
}

func recruitMenu(id, parentId, sort int, name, title, icon, path, paths, menuType, component, menuComponent, permission string) migrationModels.SysMenu {
	if component == "" && menuType == "M" {
		component = "Layout"
	}
	return migrationModels.SysMenu{
		MenuId:     id,
		MenuName:   name,
		Title:      title,
		Icon:       icon,
		Path:       path,
		Paths:      paths,
		MenuType:   menuType,
		Action:     "无",
		Permission: permission,
		ParentId:   parentId,
		Component:  component,
		Sort:       sort,
		Visible:    "0",
		IsFrame:    "1",
	}
}
