package version

import (
	"go-admin/cmd/migrate/migration"
	migrationModels "go-admin/cmd/migrate/migration/models"
	common "go-admin/common/models"
	"runtime"

	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"gorm.io/gorm"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _2026051700010EduButtonMenus)
}

func _2026051700010EduButtonMenus(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		now := pkg.GetCurrentTime()
		buttons := []migrationModels.SysMenu{
			eduButtonMenu(9101, 9001, 1, "EduStatsQuery", "查询数据概览", "edu:stats:query"),
			eduButtonMenu(9102, 9001, 2, "EduStatsExport", "导出数据概览", "edu:stats:export"),

			eduButtonMenu(9111, 9002, 1, "EduResourceQuery", "查询资源", "edu:resource:query"),
			eduButtonMenu(9112, 9002, 2, "EduResourceAdd", "新增资源", "edu:resource:add"),
			eduButtonMenu(9113, 9002, 3, "EduResourceEdit", "编辑资源", "edu:resource:edit"),
			eduButtonMenu(9114, 9002, 4, "EduResourceRemove", "删除资源", "edu:resource:remove"),
			eduButtonMenu(9115, 9002, 5, "EduResourceFiles", "资源附件", "edu:resource:files"),
			eduButtonMenu(9116, 9002, 6, "EduResourceComments", "资源评论", "edu:resource:comments"),
			eduButtonMenu(9117, 9002, 7, "EduResourceReview", "资源审核", "edu:resource:review"),
			eduButtonMenu(9118, 9002, 8, "EduResourceStatus", "资源上下架", "edu:resource:status"),
			eduButtonMenu(9119, 9002, 9, "EduResourceSearch", "同步资源搜索", "edu:resource:search"),

			eduButtonMenu(9121, 9003, 1, "EduCourseQuery", "查询课程", "edu:course:query"),
			eduButtonMenu(9122, 9003, 2, "EduCourseAdd", "新增课程", "edu:course:add"),
			eduButtonMenu(9123, 9003, 3, "EduCourseEdit", "编辑课程", "edu:course:edit"),
			eduButtonMenu(9124, 9003, 4, "EduCourseRemove", "删除课程", "edu:course:remove"),
			eduButtonMenu(9125, 9003, 5, "EduCourseManage", "课程业务管理", "edu:course:manage"),

			eduButtonMenu(9131, 9004, 1, "EduActivityQuery", "查询活动", "edu:activity:query"),
			eduButtonMenu(9132, 9004, 2, "EduActivityAdd", "新增活动", "edu:activity:add"),
			eduButtonMenu(9133, 9004, 3, "EduActivityEdit", "编辑活动", "edu:activity:edit"),
			eduButtonMenu(9134, 9004, 4, "EduActivityRemove", "删除活动", "edu:activity:remove"),
			eduButtonMenu(9135, 9004, 5, "EduActivityManage", "活动业务管理", "edu:activity:manage"),

			eduButtonMenu(9141, 9005, 1, "EduCaseQuery", "查询案例", "edu:case:query"),
			eduButtonMenu(9142, 9005, 2, "EduCaseAdd", "新增案例", "edu:case:add"),
			eduButtonMenu(9143, 9005, 3, "EduCaseEdit", "编辑案例", "edu:case:edit"),
			eduButtonMenu(9144, 9005, 4, "EduCaseRemove", "删除案例", "edu:case:remove"),
			eduButtonMenu(9145, 9005, 5, "EduCaseManage", "案例业务管理", "edu:case:manage"),
			eduButtonMenu(9146, 9005, 6, "EduCaseReview", "案例审核", "edu:case:review"),
			eduButtonMenu(9147, 9005, 7, "EduCaseAuthorization", "案例访问授权", "edu:case:authorization"),
			eduButtonMenu(9148, 9005, 8, "EduCaseAccessLog", "案例访问日志", "edu:case:accessLog"),
			eduButtonMenu(9149, 9005, 9, "EduCaseExport", "导出访问日志", "edu:case:export"),

			eduButtonMenu(9151, 9006, 1, "EduExpertQuery", "查询专家", "edu:expert:query"),
			eduButtonMenu(9152, 9006, 2, "EduExpertAdd", "新增专家", "edu:expert:add"),
			eduButtonMenu(9153, 9006, 3, "EduExpertEdit", "编辑专家", "edu:expert:edit"),
			eduButtonMenu(9154, 9006, 4, "EduExpertRemove", "删除专家", "edu:expert:remove"),
			eduButtonMenu(9155, 9006, 5, "EduExpertResources", "专家资源管理", "edu:expert:resources"),

			eduButtonMenu(9161, 9007, 1, "EduSchoolQuery", "查询学校", "edu:school:query"),
			eduButtonMenu(9162, 9007, 2, "EduSchoolAdd", "新增学校", "edu:school:add"),
			eduButtonMenu(9163, 9007, 3, "EduSchoolEdit", "编辑学校", "edu:school:edit"),
			eduButtonMenu(9164, 9007, 4, "EduSchoolRemove", "删除学校", "edu:school:remove"),

			eduButtonMenu(9171, 9008, 1, "EduRegionQuery", "查询区域", "edu:region:query"),
			eduButtonMenu(9172, 9008, 2, "EduRegionAdd", "新增区域", "edu:region:add"),
			eduButtonMenu(9173, 9008, 3, "EduRegionEdit", "编辑区域", "edu:region:edit"),
			eduButtonMenu(9174, 9008, 4, "EduRegionRemove", "删除区域", "edu:region:remove"),

			eduButtonMenu(9181, 9009, 1, "EduCategoryQuery", "查询分类", "edu:category:query"),
			eduButtonMenu(9182, 9009, 2, "EduCategoryAdd", "新增分类", "edu:category:add"),
			eduButtonMenu(9183, 9009, 3, "EduCategoryEdit", "编辑分类", "edu:category:edit"),
			eduButtonMenu(9184, 9009, 4, "EduCategoryRemove", "删除分类", "edu:category:remove"),

			eduButtonMenu(9191, 9010, 1, "EduTagQuery", "查询标签", "edu:tag:query"),
			eduButtonMenu(9192, 9010, 2, "EduTagAdd", "新增标签", "edu:tag:add"),
			eduButtonMenu(9193, 9010, 3, "EduTagEdit", "编辑标签", "edu:tag:edit"),
			eduButtonMenu(9194, 9010, 4, "EduTagRemove", "删除标签", "edu:tag:remove"),
		}

		for _, button := range buttons {
			button.CreatedAt = now
			button.UpdatedAt = now
			button.CreateBy = 1
			button.UpdateBy = 1

			var count int64
			if err := tx.Model(&migrationModels.SysMenu{}).
				Where("menu_id = ? OR permission = ?", button.MenuId, button.Permission).
				Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				continue
			}
			if err := tx.Create(&button).Error; err != nil {
				return err
			}
		}

		return tx.Create(&common.Migration{Version: version}).Error
	})
}

func eduButtonMenu(menuID, parentID, sort int, name, title, permission string) migrationModels.SysMenu {
	return migrationModels.SysMenu{
		MenuId:     menuID,
		MenuName:   name,
		Title:      title,
		Path:       "",
		Paths:      "/0/9000/" + pkg.IntToString(parentID) + "/" + pkg.IntToString(menuID),
		MenuType:   "F",
		Action:     "无",
		Permission: permission,
		ParentId:   parentID,
		NoCache:    false,
		Sort:       sort,
		Visible:    "0",
		IsFrame:    "1",
	}
}
