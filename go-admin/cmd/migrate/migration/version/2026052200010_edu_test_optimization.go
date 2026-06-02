package version

import (
	"go-admin/app/edu/models"
	"go-admin/cmd/migrate/migration"
	migrationModels "go-admin/cmd/migrate/migration/models"
	common "go-admin/common/models"
	"runtime"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	"gorm.io/gorm"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _2026052200010EduTestOptimization)
}

func _2026052200010EduTestOptimization(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		schemaTx := tx.Session(&gorm.Session{NewDB: true})
		if config.DatabaseConfig.Driver == "mysql" {
			schemaTx = schemaTx.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
		}
		if err := schemaTx.AutoMigrate(new(models.EduExpertResource)); err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{NewDB: true}).Model(&migrationModels.SysMenu{}).Where("menu_id = ?", 9000).Update("title", "系统设置").Error; err != nil {
			return err
		}
		if err := tx.Session(&gorm.Session{NewDB: true}).Model(&migrationModels.SysMenu{}).Where("menu_id in ?", []int{9007, 9008, 9012}).Update("visible", "1").Error; err != nil {
			return err
		}
		return tx.Session(&gorm.Session{NewDB: true}).Create(&common.Migration{Version: version}).Error
	})
}
