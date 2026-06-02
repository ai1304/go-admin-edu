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
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _2026052200030RestoreEduSettingMenus)
}

func _2026052200030RestoreEduSettingMenus(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Session(&gorm.Session{NewDB: true}).
			Model(&migrationModels.SysMenu{}).
			Where("menu_id in ?", []int{9007, 9008, 9012}).
			Update("visible", "0").Error; err != nil {
			return err
		}
		return tx.Session(&gorm.Session{NewDB: true}).Create(&common.Migration{Version: version}).Error
	})
}
