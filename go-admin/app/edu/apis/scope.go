package apis

import (
	adminModels "go-admin/app/admin/models"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"gorm.io/gorm"
)

func applyEduUserScope(c *gin.Context, db *gorm.DB) *gorm.DB {
	if isEduScopeBypass(c) {
		return db
	}

	var current adminModels.SysUser
	if err := db.Session(&gorm.Session{NewDB: true}).First(&current, user.GetUserId(c)).Error; err != nil {
		return db
	}
	if current.UserType == "super_admin" {
		return db
	}
	if current.SchoolId != 0 {
		return db.Where("school_id = ?", current.SchoolId)
	}
	if current.RegionId != 0 {
		return db.Where("region_id = ?", current.RegionId)
	}
	return db
}

func isEduScopeBypass(c *gin.Context) bool {
	roleName := user.GetRoleName(c)
	return roleName == "admin" || roleName == "系统管理员" || user.GetRoleId(c) == 1
}
