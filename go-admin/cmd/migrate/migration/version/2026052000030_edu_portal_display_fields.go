package version

import (
	"go-admin/app/edu/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
	"runtime"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	"gorm.io/gorm"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _2026052000030EduPortalDisplayFields)
}

func _2026052000030EduPortalDisplayFields(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if config.DatabaseConfig.Driver == "mysql" {
			tx = tx.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
		}
		if err := tx.AutoMigrate(
			new(models.EduCourse),
			new(models.EduActivity),
			new(models.EduCase),
			new(models.EduExpert),
		); err != nil {
			return err
		}
		return tx.Create(&common.Migration{Version: version}).Error
	})
}
