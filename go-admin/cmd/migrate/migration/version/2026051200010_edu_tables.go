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
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _2026051200010EduTables)
}

func _2026051200010EduTables(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if config.DatabaseConfig.Driver == "mysql" {
			tx = tx.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4")
		}
		if err := tx.AutoMigrate(
			new(models.EduRegion),
			new(models.EduSchool),
			new(models.EduResourceCategory),
			new(models.EduResourceTag),
			new(models.EduResource),
			new(models.EduResourceFile),
			new(models.EduResourceReview),
		); err != nil {
			return err
		}
		return tx.Create(&common.Migration{Version: version}).Error
	})
}
