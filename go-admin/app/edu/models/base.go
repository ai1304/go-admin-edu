package models

import common "go-admin/common/models"

type TenantScope struct {
	TenantId int `json:"tenantId" gorm:"index;comment:租户ID"`
	RegionId int `json:"regionId" gorm:"index;comment:区域ID"`
	SchoolId int `json:"schoolId" gorm:"index;comment:学校ID"`
}

type BaseModel struct {
	common.Model
	TenantScope
	common.ControlBy
	common.ModelTime
}
