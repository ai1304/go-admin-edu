package version

import (
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
	"runtime"

	"gorm.io/gorm"
)

type sysUserEduScopeColumns struct {
	TenantId int    `json:"tenantId" gorm:"index;comment:租户ID"`
	RegionId int    `json:"regionId" gorm:"index;comment:区域ID"`
	SchoolId int    `json:"schoolId" gorm:"index;comment:学校ID"`
	UserType string `json:"userType" gorm:"size:32;index;comment:用户类型"`
}

func (sysUserEduScopeColumns) TableName() string {
	return "sys_user"
}

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _2026051600050SysUserEduScope)
}

func _2026051600050SysUserEduScope(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		model := new(sysUserEduScopeColumns)
		for _, column := range []string{"TenantId", "RegionId", "SchoolId", "UserType"} {
			if !tx.Migrator().HasColumn(model, column) {
				if err := tx.Migrator().AddColumn(model, column); err != nil {
					return err
				}
			}
		}
		return tx.Create(&common.Migration{Version: version}).Error
	})
}
