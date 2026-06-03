package version

import (
	"go-admin/app/edu/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
	"runtime"
	"time"

	"gorm.io/gorm"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _2026060300010EduRecruitSeed)
}

func _2026060300010EduRecruitSeed(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		companies, err := seedRecruitCompanies(tx)
		if err != nil {
			return err
		}
		if err := seedRecruitJobs(tx, companies); err != nil {
			return err
		}
		return tx.Create(&common.Migration{Version: version}).Error
	})
}

func seedRecruitCompanies(tx *gorm.DB) (map[string]models.EduRecruitCompany, error) {
	now := time.Now()
	reviewedAt := "2026-05-20 10:00:00"
	companies := []models.EduRecruitCompany{
		{
			CompanyName: "星光特殊教育学校", CreditCode: "91320105MA1X2A3B4C", CompanyNature: "民办非企业单位", Industry: "特殊教育学校", CompanySize: "100-499人",
			Region: "南京市雨花台区", Address: "南京市雨花台区软件大道88号", Website: "https://www.xgtsxx.edu.cn", LogoUrl: "", ContactName: "张老师", ContactTitle: "人事负责人", ContactPhone: "138-0512-6688", ContactEmail: "hr@xgtsxx.com",
			Intro:           "星光特殊教育学校面向听障、智力障碍和孤独症学生提供义务教育、职业启蒙与融合支持服务，长期开展个别化教育计划和家校协同实践。",
			MainBusiness:    "特殊教育教学、融合课程建设、学生发展评估、职业适应训练和家校支持服务。",
			TalentNeeds:     "融合教育教师、特教班主任、言语康复教师、心理支持教师",
			Cooperation:     "希望与高校共建实习实践基地，联合开展 IEP 课程和融合教育课堂研究。",
			LicenseMaterial: "营业执照及民办学校办学许可证已归档", Qualification: "江苏省融合教育示范单位", AuthorizationNote: "联系人已获学校授权发布招聘信息", Tags: "特殊教育学校,融合教育,IEP",
			CertStatus: "approved", Status: "normal", ReviewStatus: "approved", ReviewedAt: reviewedAt,
		},
		{
			CompanyName: "启航融合教育中心", CreditCode: "91330106MA2Y3B4C5D", CompanyNature: "民办非企业单位", Industry: "融合教育机构", CompanySize: "50-99人",
			Region: "杭州市滨江区", Address: "杭州市滨江区江南大道3888号", Website: "https://www.qhfusion.org", ContactName: "李老师", ContactTitle: "运营负责人", ContactPhone: "139-5711-8899", ContactEmail: "talent@qhfusion.org",
			Intro:           "启航融合教育中心为普通学校和家庭提供融合教育评估、课堂支持、影子教师培训和家长赋能服务。",
			MainBusiness:    "融合教育咨询、课堂支持服务、教师培训、个案督导和家庭支持。",
			TalentNeeds:     "融合教育支持专员、影子教师督导、资源教室教师",
			Cooperation:     "联合高校开展融合教育实习、教师培训和个案督导项目。",
			LicenseMaterial: "社会服务机构登记证书已归档", Qualification: "融合教育专业服务机构", AuthorizationNote: "人事部门授权发布招聘信息", Tags: "融合教育,资源教室,影子教师",
			CertStatus: "approved", Status: "normal", ReviewStatus: "approved", ReviewedAt: reviewedAt,
		},
		{
			CompanyName: "康语儿童康复中心", CreditCode: "91330203MA2E4F5G6H", CompanyNature: "民办非企业单位", Industry: "康复机构", CompanySize: "50-99人",
			Region: "宁波市鄞州区", Address: "宁波市鄞州区首南街道学士路168号", Website: "https://www.kangyu-rehab.cn", ContactName: "王老师", ContactTitle: "招聘主管", ContactPhone: "137-7788-9900", ContactEmail: "hr@kangyu-rehab.cn",
			Intro:           "康语儿童康复中心专注言语语言、感统训练和孤独症儿童早期干预，提供评估、训练、家庭指导一体化服务。",
			MainBusiness:    "言语语言评估与训练、感统训练、孤独症干预、家庭康复指导。",
			TalentNeeds:     "言语治疗师、康复训练师、感统教师、个训教师",
			Cooperation:     "欢迎康复治疗、特殊教育和心理学专业学生实习实践。",
			LicenseMaterial: "民办非企业登记证书已归档", Qualification: "儿童康复服务备案机构", AuthorizationNote: "中心负责人授权发布招聘信息", Tags: "康复机构,言语康复,孤独症",
			CertStatus: "approved", Status: "normal", ReviewStatus: "approved", ReviewedAt: reviewedAt,
		},
		{
			CompanyName: "智汇辅助科技有限公司", CreditCode: "91320115MA3H7J8K9L", CompanyNature: "民营企业", Industry: "教育科技", CompanySize: "100-499人",
			Region: "深圳市南山区", Address: "深圳市南山区粤海街道科技园南区18栋", Website: "https://www.zhassistive.com", ContactName: "陈老师", ContactTitle: "人才发展经理", ContactPhone: "136-6555-7788", ContactEmail: "jobs@zhassistive.com",
			Intro:           "智汇辅助科技有限公司研发无障碍学习软件、AAC 沟通设备和特殊教育数据平台，服务学校、康复机构与家庭。",
			MainBusiness:    "辅助沟通设备、无障碍学习工具、资源教室设备和特教数据平台研发。",
			TalentNeeds:     "辅助技术支持工程师、特教产品运营、无障碍测试专员",
			Cooperation:     "与高校共建辅助技术实验室，提供产品实习和毕业设计课题。",
			LicenseMaterial: "营业执照已归档", Qualification: "国家高新技术企业", AuthorizationNote: "人力资源部授权发布招聘信息", Tags: "辅助技术,教育科技,AAC",
			CertStatus: "approved", Status: "normal", ReviewStatus: "approved", ReviewedAt: reviewedAt,
		},
		{
			CompanyName: "童心康复医疗有限公司", CreditCode: "91330108MA4K6L7M8N", CompanyNature: "民营企业", Industry: "医疗/康复", CompanySize: "100-499人",
			Region: "成都市武侯区", Address: "成都市武侯区人民南路四段66号", Website: "https://www.tongxinrehab.cn", ContactName: "赵老师", ContactTitle: "人力资源经理", ContactPhone: "139-1391-0000", ContactEmail: "recruit@tongxinrehab.cn",
			Intro:           "童心康复医疗有限公司提供儿童发育评估、作业治疗、言语治疗和心理支持服务，关注医教结合与家庭协同。",
			MainBusiness:    "儿童康复评估、OT/ST 训练、心理支持、家庭康复指导和医教结合服务。",
			TalentNeeds:     "康复治疗师、作业治疗师、心理支持专员",
			Cooperation:     "与高校合作开展康复实习、临床督导和专业培训。",
			LicenseMaterial: "医疗机构执业许可证已归档", Qualification: "儿童康复专科服务机构", AuthorizationNote: "招聘负责人授权发布岗位", Tags: "医疗康复,作业治疗,心理支持",
			CertStatus: "approved", Status: "normal", ReviewStatus: "approved", ReviewedAt: reviewedAt,
		},
		{
			CompanyName: "阳光融合学校", CreditCode: "52330200MJ9132456L", CompanyNature: "事业单位", Industry: "教育", CompanySize: "100-499人",
			Region: "上海市浦东新区", Address: "上海市浦东新区锦绣东路126号", Website: "https://www.sunfusion.edu.cn", ContactName: "刘老师", ContactTitle: "办公室主任", ContactPhone: "158-6745-2211", ContactEmail: "office@sunfusion.edu.cn",
			Intro:           "阳光融合学校探索普特融合课程、职业转衔与学生社会适应支持，建设资源教室和多学科协作团队。",
			MainBusiness:    "融合教育课程、资源教室支持、职业转衔训练和学生发展服务。",
			TalentNeeds:     "资源教室管理员、融合课程教师、职业转衔辅导员",
			Cooperation:     "开展融合教育示范课、实习实践和课题研究合作。",
			LicenseMaterial: "事业单位法人证书已归档", Qualification: "市级融合教育示范校", AuthorizationNote: "学校办公室授权发布岗位", Tags: "融合学校,资源教室,职业转衔",
			CertStatus: "approved", Status: "normal", ReviewStatus: "approved", ReviewedAt: reviewedAt,
		},
	}

	result := make(map[string]models.EduRecruitCompany, len(companies))
	for i := range companies {
		companies[i].CreatedAt = now.AddDate(0, 0, -i-8)
		companies[i].UpdatedAt = companies[i].CreatedAt
		if err := tx.Where("credit_code = ?", companies[i].CreditCode).Assign(companies[i]).FirstOrCreate(&companies[i]).Error; err != nil {
			return nil, err
		}
		result[companies[i].CompanyName] = companies[i]
	}
	return result, nil
}

func seedRecruitJobs(tx *gorm.DB, companies map[string]models.EduRecruitCompany) error {
	now := time.Now()
	reviewedAt := "2026-05-20 11:00:00"
	jobs := []models.EduRecruitJob{
		recruitJob(companies["星光特殊教育学校"], "融合教育教师", "全职", "特殊教育", "南京市雨花台区", 3, "8K-12K/月", "本科及以上", "特殊教育、教育学、心理学等相关专业", "应届生/社会招聘", "2026-05-20 09:00:00", "2026-08-31", "负责融合班级教学与 IEP 制定；开展学生发展评估；与班主任和康复教师协作；记录教学支持过程。", "持教师资格证；熟悉 IEP 流程；具备良好沟通能力和班级管理能力。", "周一至周五 8:00-16:30", "五险一金、寒暑假、专业培训、餐补、交通补贴", "新教师导师制、融合课堂观摩和课题参与机会", "融合教育,特殊教育,IEP", "张老师", "人事负责人", "138-0512-6688", "hr@xgtsxx.com"),
		recruitJob(companies["星光特殊教育学校"], "行为干预师", "全职", "心理支持", "南京市雨花台区", 2, "7K-10K/月", "本科及以上", "心理学、特殊教育、应用行为分析相关专业", "社会招聘", "2026-05-22 09:00:00", "2026-08-20", "为学生制定行为支持方案；开展课堂观察和数据记录；与家长及教师沟通干预进展。", "熟悉 ABA 或正向行为支持；能独立完成观察记录和阶段复盘。", "周一至周五 8:30-17:00", "五险一金、专业督导、节日福利", "提供行为干预个案督导和校本培训", "行为干预,心理支持,ABA", "张老师", "人事负责人", "138-0512-6688", "hr@xgtsxx.com"),
		recruitJob(companies["启航融合教育中心"], "融合教育支持专员", "全职", "融合教育", "杭州市滨江区", 2, "10K-15K/月", "本科及以上", "特殊教育、教育学、心理学相关专业", "社会招聘", "2026-05-19 09:00:00", "2026-08-15", "为普通学校提供课堂支持；参与个案评估和支持计划制定；开展教师和家长沟通。", "有融合教育或资源教室经验优先；能适应校区走访和项目协作。", "周一至周五 9:00-18:00", "五险一金、项目奖金、督导培训", "资深督导师带教，参与跨校融合教育项目", "融合教育,资源教室,课堂支持", "李老师", "运营负责人", "139-5711-8899", "talent@qhfusion.org"),
		recruitJob(companies["康语儿童康复中心"], "言语治疗师", "全职", "言语听觉", "宁波市鄞州区", 3, "9K-14K/月", "本科及以上", "言语听觉科学、康复治疗、特殊教育相关专业", "应届生/社会招聘", "2026-05-18 09:00:00", "2026-09-15", "负责儿童言语语言评估与训练；制定个别化训练计划；跟踪训练效果并与家长沟通。", "具备言语治疗或儿童康复实习经验优先；表达清晰，有耐心。", "周二至周六 9:00-18:00", "五险一金、带薪培训、绩效奖金", "提供言语治疗师成长路径和个案督导", "言语治疗,康复训练,儿童康复", "王老师", "招聘主管", "137-7788-9900", "hr@kangyu-rehab.cn"),
		recruitJob(companies["智汇辅助科技有限公司"], "辅助技术支持工程师", "全职", "辅助技术", "深圳市南山区", 2, "12K-18K/月", "本科及以上", "特殊教育、计算机、康复工程、教育技术相关专业", "社会招聘", "2026-05-17 09:00:00", "2026-09-30", "为学校和机构提供辅助技术方案；安装调试辅助设备；培训教师使用并跟踪应用效果。", "理解特殊教育场景；具备技术支持或产品培训经验；能接受短期出差。", "周一至周五 9:00-18:30", "五险一金、项目奖金、补充医疗、交通补贴", "产品、教研和客户成功联合培养", "辅助技术,AAC,无障碍", "陈老师", "人才发展经理", "136-6555-7788", "jobs@zhassistive.com"),
		recruitJob(companies["童心康复医疗有限公司"], "作业治疗师", "全职", "康复", "成都市武侯区", 2, "8K-13K/月", "本科及以上", "康复治疗学、作业治疗、特殊教育相关专业", "社会招聘", "2026-05-16 09:00:00", "2026-08-25", "开展儿童作业治疗评估与训练；设计家庭训练方案；参与多学科个案会议。", "具备 OT 训练经验；能完成评估记录和家庭指导；持相关资格优先。", "周二至周六 9:00-18:00", "五险一金、绩效奖金、专业督导", "多学科团队督导和外部培训机会", "作业治疗,儿童康复,医教结合", "赵老师", "人力资源经理", "139-1391-0000", "recruit@tongxinrehab.cn"),
		recruitJob(companies["阳光融合学校"], "资源教室管理员", "全职", "特殊教育", "上海市浦东新区", 1, "6K-8K/月", "大专及以上", "特殊教育、教育技术、心理学相关专业", "应届生/社会招聘", "2026-05-15 09:00:00", "2026-08-10", "负责资源教室日常管理；维护辅助器具和教学资源；协助教师开展学生支持记录。", "熟悉办公软件和资源管理；有耐心，能与教师和学生沟通协作。", "周一至周五 8:30-16:30", "五险一金、寒暑假、校内培训", "参与资源教室标准化建设和融合课程实践", "资源教室,辅助技术,融合教育", "刘老师", "办公室主任", "158-6745-2211", "office@sunfusion.edu.cn"),
		recruitJob(companies["阳光融合学校"], "特教实习辅导员", "实习", "特殊教育", "上海市浦东新区", 5, "3K-4K/月", "大专及以上", "特殊教育、学前教育、心理学相关专业", "在校生/应届生", "2026-05-14 09:00:00", "2026-07-31", "协助班级教师完成课堂支持；参与学生观察记录；协助组织融合活动。", "每周至少到岗 3 天；认真负责，愿意学习特殊教育支持方法。", "周一至周五可排班", "实习补贴、导师指导、实习证明", "配备校内指导教师，参与真实课堂支持", "实习,融合教育,课堂支持", "刘老师", "办公室主任", "158-6745-2211", "office@sunfusion.edu.cn"),
	}
	for i := range jobs {
		jobs[i].CreatedAt = now.AddDate(0, 0, -i-3)
		jobs[i].UpdatedAt = jobs[i].CreatedAt
		jobs[i].ReviewedAt = reviewedAt
		jobs[i].Status = "published"
		jobs[i].ReviewStatus = "approved"
		if err := tx.Where("company_id = ? and job_name = ?", jobs[i].CompanyId, jobs[i].JobName).Assign(jobs[i]).FirstOrCreate(&jobs[i]).Error; err != nil {
			return err
		}
	}
	return nil
}

func recruitJob(company models.EduRecruitCompany, name, jobType, direction, location string, headcount int, salary, education, major, target, publishTime, deadline, responsibilities, requirements, workTime, benefits, training, tags, contactName, contactTitle, contactPhone, contactEmail string) models.EduRecruitJob {
	return models.EduRecruitJob{
		CompanyId: company.Id, CompanyName: company.CompanyName, Industry: company.Industry, JobName: name, JobType: jobType, MajorDirection: direction, Location: location, Headcount: headcount, SalaryRange: salary, Education: education,
		MajorRequirement: major, RecruitTarget: target, PublishTime: publishTime, Deadline: deadline, Responsibilities: responsibilities, Requirements: requirements, WorkTime: workTime, Benefits: benefits, Training: training, Tags: tags,
		ContactName: contactName, ContactTitle: contactTitle, ContactPhone: contactPhone, ContactEmail: contactEmail, ContactAddress: company.Address, ExternalLink: company.Website,
	}
}
