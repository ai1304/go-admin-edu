package apis

import (
	"go-admin/app/edu/models"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
)

type EduStats struct {
	api.Api
}

func (e EduStats) Overview(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	counts := gin.H{}
	counts["regions"] = e.countModel(&models.EduRegion{})
	counts["schools"] = e.countModel(&models.EduSchool{})
	counts["resources"] = e.countModel(&models.EduResource{})
	counts["publishedResources"] = e.countWhere(&models.EduResource{}, "status = ?", "published")
	counts["courses"] = e.countModel(&models.EduCourse{})
	counts["activities"] = e.countModel(&models.EduActivity{})
	counts["cases"] = e.countModel(&models.EduCase{})
	counts["experts"] = e.countModel(&models.EduExpert{})
	e.OK(counts, "查询成功")
}

func (e EduStats) countModel(model interface{}) int64 {
	var count int64
	_ = e.Orm.Model(model).Count(&count).Error
	return count
}

func (e EduStats) countWhere(model interface{}, query string, args ...interface{}) int64 {
	var count int64
	_ = e.Orm.Model(model).Where(query, args...).Count(&count).Error
	return count
}
