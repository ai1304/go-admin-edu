package version

import (
	"go-admin/app/edu/models"
	"go-admin/cmd/migrate/migration"
	common "go-admin/common/models"
	"runtime"

	"gorm.io/gorm"
)

func init() {
	_, fileName, _, _ := runtime.Caller(0)
	migration.Migrate.SetVersion(migration.GetFilename(fileName), _2026052000040EduDemoSeed)
}

func _2026052000040EduDemoSeed(db *gorm.DB, version string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		categoryIDs, err := seedDemoCategories(tx)
		if err != nil {
			return err
		}
		if err := seedDemoRegionsAndSchools(tx); err != nil {
			return err
		}
		if err := seedDemoTags(tx); err != nil {
			return err
		}
		if err := seedDemoResources(tx, categoryIDs); err != nil {
			return err
		}
		if err := seedDemoCourses(tx, categoryIDs); err != nil {
			return err
		}
		if err := seedDemoActivities(tx); err != nil {
			return err
		}
		if err := seedDemoCases(tx); err != nil {
			return err
		}
		if err := seedDemoExperts(tx); err != nil {
			return err
		}
		if err := seedDemoNewsAndAI(tx); err != nil {
			return err
		}
		return tx.Create(&common.Migration{Version: version}).Error
	})
}

func seedDemoRegionsAndSchools(tx *gorm.DB) error {
	regions := []models.EduRegion{
		{Name: "华东区域", Code: "east-china", Sort: 1, Status: 1, Remark: "演示区域"},
		{Name: "华北区域", Code: "north-china", Sort: 2, Status: 1, Remark: "演示区域"},
	}
	for _, item := range regions {
		if err := tx.Where("code = ?", item.Code).FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}

	schools := []models.EduSchool{
		{Name: "示范特殊教育学院", Code: "senedu-demo", Address: "上海市浦东新区融合教育示范园", Contact: "教研办公室", Phone: "021-88886666", Status: 1, Remark: "special_edu 演示学校"},
		{Name: "阳光职业技术学院特教系", Code: "sunshine-vocational", Address: "杭州市滨江区阳光路 88 号", Contact: "特教系办公室", Phone: "0571-77778888", Status: 1, Remark: "special_edu 演示学校"},
	}
	for _, item := range schools {
		if err := tx.Where("code = ?", item.Code).FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedDemoCategories(tx *gorm.DB) (map[string]int, error) {
	categories := []models.EduResourceCategory{
		{Name: "专科", Code: "stage-college", Type: "stage", Sort: 1, Status: 1},
		{Name: "本科", Code: "stage-undergraduate", Type: "stage", Sort: 2, Status: 1},
		{Name: "研究生", Code: "stage-postgraduate", Type: "stage", Sort: 3, Status: 1},
		{Name: "听障", Code: "disability-hearing", Type: "disability", Sort: 1, Status: 1},
		{Name: "视障", Code: "disability-visual", Type: "disability", Sort: 2, Status: 1},
		{Name: "智力障碍", Code: "disability-intellectual", Type: "disability", Sort: 3, Status: 1},
		{Name: "言语障碍", Code: "disability-speech", Type: "disability", Sort: 4, Status: 1},
		{Name: "肢体障碍", Code: "disability-physical", Type: "disability", Sort: 5, Status: 1},
		{Name: "孤独症", Code: "disability-autism", Type: "disability", Sort: 6, Status: 1},
		{Name: "其他", Code: "disability-other", Type: "disability", Sort: 7, Status: 1},
		{Name: "认知", Code: "ability-cognition", Type: "ability_domain", Sort: 1, Status: 1},
		{Name: "语言", Code: "ability-language", Type: "ability_domain", Sort: 2, Status: 1},
		{Name: "社交", Code: "ability-social", Type: "ability_domain", Sort: 3, Status: 1},
		{Name: "运动技能", Code: "ability-motor", Type: "ability_domain", Sort: 4, Status: 1},
		{Name: "生活技能", Code: "ability-life", Type: "ability_domain", Sort: 5, Status: 1},
		{Name: "文档", Code: "resource-doc", Type: "resource_type", Sort: 1, Status: 1},
		{Name: "课件PPT", Code: "resource-ppt", Type: "resource_type", Sort: 2, Status: 1},
		{Name: "视频", Code: "resource-video", Type: "resource_type", Sort: 3, Status: 1},
		{Name: "微课", Code: "resource-micro", Type: "resource_type", Sort: 4, Status: 1},
		{Name: "融合教育", Code: "topic-inclusive", Type: "topic", Sort: 1, Status: 1},
		{Name: "个别化教育", Code: "topic-iep", Type: "topic", Sort: 2, Status: 1},
		{Name: "康复训练", Code: "topic-rehab", Type: "topic", Sort: 3, Status: 1},
		{Name: "辅助技术", Code: "topic-assistive", Type: "topic", Sort: 4, Status: 1},
		{Name: "职业教育", Code: "topic-vocational", Type: "topic", Sort: 5, Status: 1},
	}

	ids := make(map[string]int, len(categories))
	for _, item := range categories {
		if err := tx.Where("code = ?", item.Code).FirstOrCreate(&item).Error; err != nil {
			return nil, err
		}
		ids[item.Type+":"+item.Name] = item.Id
	}
	return ids, nil
}

func seedDemoTags(tx *gorm.DB) error {
	tags := []string{"听障", "融合教育", "IEP", "康复训练", "辅助技术", "孤独症", "生活技能", "教师发展"}
	for _, name := range tags {
		item := models.EduResourceTag{Name: name, Status: 1}
		if err := tx.Where("name = ?", name).FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedDemoResources(tx *gorm.DB, c map[string]int) error {
	resources := []models.EduResource{
		{Title: "听障大学生沟通技巧教学课件", Summary: "面向听障大学生的课堂沟通策略与手语辅助教学课件。", StageCategoryId: c["stage:本科"], DisabilityTypeId: c["disability:听障"], AbilityDomainId: c["ability_domain:社交"], TopicCategoryId: c["topic:融合教育"], ResourceTypeId: c["resource_type:课件PPT"], AuthorName: "张老师", Keywords: "听障,沟通,课件", Status: models.ResourceStatusPublished, ViewCount: 1280, DownloadCount: 356, FavoriteCount: 42},
		{Title: "视障学生信息无障碍操作指南", Summary: "介绍读屏软件、盲文点显器等辅助技术的使用方法。", StageCategoryId: c["stage:本科"], DisabilityTypeId: c["disability:视障"], AbilityDomainId: c["ability_domain:生活技能"], TopicCategoryId: c["topic:辅助技术"], ResourceTypeId: c["resource_type:文档"], AuthorName: "李梅", Keywords: "视障,无障碍,辅助技术", Status: models.ResourceStatusPublished, ViewCount: 980, DownloadCount: 233, FavoriteCount: 31},
		{Title: "孤独症谱系学生社交干预微课", Summary: "5分钟微课，讲解结构化社交训练的核心步骤。", StageCategoryId: c["stage:专科"], DisabilityTypeId: c["disability:孤独症"], AbilityDomainId: c["ability_domain:社交"], TopicCategoryId: c["topic:康复训练"], ResourceTypeId: c["resource_type:微课"], AuthorName: "张老师", Keywords: "孤独症,社交,微课", Status: models.ResourceStatusPublished, ViewCount: 2150, DownloadCount: 540, FavoriteCount: 88},
		{Title: "智力障碍学生生活技能训练视频", Summary: "生活自理能力训练示范视频，含洗漱、穿衣等模块。", StageCategoryId: c["stage:专科"], DisabilityTypeId: c["disability:智力障碍"], AbilityDomainId: c["ability_domain:生活技能"], TopicCategoryId: c["topic:康复训练"], ResourceTypeId: c["resource_type:视频"], AuthorName: "李梅", Keywords: "智力障碍,生活技能", Status: models.ResourceStatusPublished, ViewCount: 1760, DownloadCount: 410, FavoriteCount: 65},
		{Title: "个别化教育计划(IEP)编制范例", Summary: "完整的 IEP 文档模板与填写说明。", StageCategoryId: c["stage:研究生"], DisabilityTypeId: c["disability:其他"], AbilityDomainId: c["ability_domain:认知"], TopicCategoryId: c["topic:个别化教育"], ResourceTypeId: c["resource_type:文档"], AuthorName: "张老师", Keywords: "IEP,模板", Status: models.ResourceStatusPublished, ViewCount: 3050, DownloadCount: 1120, FavoriteCount: 156},
		{Title: "融合教育课堂管理策略", Summary: "普特融合课堂的差异化教学与课堂管理方法。", StageCategoryId: c["stage:本科"], DisabilityTypeId: c["disability:其他"], AbilityDomainId: c["ability_domain:认知"], TopicCategoryId: c["topic:融合教育"], ResourceTypeId: c["resource_type:文档"], AuthorName: "李梅", Keywords: "融合教育,课堂管理", Status: models.ResourceStatusPublished, ViewCount: 1430, DownloadCount: 305, FavoriteCount: 49},
	}
	for _, item := range resources {
		if err := tx.Where("title = ?", item.Title).FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedDemoCourses(tx *gorm.DB, c map[string]int) error {
	courses := []models.EduCourse{
		{Title: "特殊教育导论", Summary: "系统介绍特殊教育的基本理论、对象与发展历程。", StageCategoryId: c["stage:本科"], DisabilityTypeId: c["disability:其他"], Category: "专业基础", Difficulty: "EASY", TeacherName: "张老师", Objectives: "掌握特殊教育基本概念；了解各类特殊儿童的身心特点；建立融合教育理念。", Status: "published", ViewCount: 3200, LearnerCount: 860, Sort: 90},
		{Title: "个别化教育计划(IEP)实务", Summary: "从评估到目标制定，全流程讲解 IEP 编制与实施。", StageCategoryId: c["stage:研究生"], DisabilityTypeId: c["disability:其他"], Category: "教学技能", Difficulty: "MEDIUM", TeacherName: "李梅", Objectives: "能独立完成 IEP 评估与编制；掌握目标分解与成效追踪方法。", Status: "published", ViewCount: 2480, LearnerCount: 640, Sort: 80},
		{Title: "听障学生沟通与教学", Summary: "聚焦听障大学生的沟通策略与课堂教学方法。", StageCategoryId: c["stage:本科"], DisabilityTypeId: c["disability:听障"], Category: "教学技能", Difficulty: "MEDIUM", TeacherName: "张老师", Objectives: "掌握手语与口语并用的教学策略；提升听障课堂沟通效率。", Status: "published", ViewCount: 1760, LearnerCount: 420, Sort: 70},
		{Title: "孤独症儿童干预方法", Summary: "介绍 ABA、结构化教学等主流干预方法。", StageCategoryId: c["stage:专科"], DisabilityTypeId: c["disability:孤独症"], Category: "康复训练", Difficulty: "HARD", TeacherName: "李梅", Objectives: "理解孤独症核心障碍；掌握结构化教学与行为干预技术。", Status: "published", ViewCount: 2050, LearnerCount: 510, Sort: 60},
		{Title: "辅助技术与无障碍设计", Summary: "讲解特殊教育中的辅助技术与无障碍环境设计。", StageCategoryId: c["stage:本科"], DisabilityTypeId: c["disability:视障"], Category: "辅助技术", Difficulty: "MEDIUM", TeacherName: "王强", Objectives: "了解主流辅助技术产品；掌握无障碍教学环境设计原则。", Status: "published", ViewCount: 1320, LearnerCount: 310, Sort: 50},
	}
	for _, item := range courses {
		if err := tx.Where("title = ?", item.Title).FirstOrCreate(&item).Error; err != nil {
			return err
		}
		if err := seedDemoLessons(tx, item.Id, item.Title); err != nil {
			return err
		}
	}
	return nil
}

func seedDemoLessons(tx *gorm.DB, courseId int, courseTitle string) error {
	lessons := map[string][]models.EduCourseLesson{
		"特殊教育导论": {
			{CourseId: courseId, Title: "第1讲 特殊教育的概念与对象", DurationSeconds: 1820, Sort: 1, Status: 1},
			{CourseId: courseId, Title: "第2讲 特殊教育的发展历程", DurationSeconds: 1650, Sort: 2, Status: 1},
			{CourseId: courseId, Title: "第3讲 融合教育理念与实践", DurationSeconds: 2010, Sort: 3, Status: 1},
		},
		"个别化教育计划(IEP)实务": {
			{CourseId: courseId, Title: "第1讲 IEP 评估的方法与工具", DurationSeconds: 2200, Sort: 1, Status: 1},
			{CourseId: courseId, Title: "第2讲 教育目标的制定与分解", DurationSeconds: 1980, Sort: 2, Status: 1},
			{CourseId: courseId, Title: "第3讲 IEP 的实施与成效追踪", DurationSeconds: 2100, Sort: 3, Status: 1},
		},
		"听障学生沟通与教学": {
			{CourseId: courseId, Title: "第1讲 听障学生的语言特点", DurationSeconds: 1700, Sort: 1, Status: 1},
			{CourseId: courseId, Title: "第2讲 课堂沟通策略", DurationSeconds: 1850, Sort: 2, Status: 1},
		},
	}
	for _, item := range lessons[courseTitle] {
		if err := tx.Where("course_id = ? and title = ?", item.CourseId, item.Title).FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedDemoActivities(tx *gorm.DB) error {
	activities := []models.EduActivity{
		{Title: "特殊教育融合课堂教学创新研讨", Summary: "围绕融合课堂的差异化教学开展同课异构研讨。", School: "示范特殊教育学院", Teacher: "张老师", Track: "课程思政", Edition: "第五届", AwardLevel: "一等奖", SchoolType: "MINISTERIAL", TitleRank: "正高", StartTime: "2026-03-12 14:00:00", EndTime: "2026-03-12 17:00:00", Location: "学院报告厅", Organizer: "示范特殊教育学院", SignupCount: 86, ViewCount: 980, Status: "published", Sort: 90},
		{Title: "IEP 个别化教育计划编制工作坊", Summary: "面向一线特教教师的 IEP 编制实操工作坊。", School: "示范特殊教育学院", Teacher: "李梅", Track: "基础课程", Edition: "第五届", AwardLevel: "二等奖", SchoolType: "MINISTERIAL", TitleRank: "副高", StartTime: "2026-04-08 09:00:00", EndTime: "2026-04-08 16:00:00", Location: "实训中心", Organizer: "区域特教教研中心", SignupCount: 64, ViewCount: 820, Status: "published", Sort: 80},
		{Title: "孤独症儿童干预技术专题教研", Summary: "聚焦孤独症儿童的结构化教学与行为干预。", School: "阳光职业技术学院特教系", Teacher: "王强", Track: "新医科", Edition: "第四届", AwardLevel: "一等奖", SchoolType: "LOCAL", TitleRank: "中级及以下", StartTime: "2026-02-20 14:30:00", EndTime: "2026-02-20 17:00:00", Location: "阳光学院A栋", Organizer: "阳光职业技术学院", SignupCount: 52, ViewCount: 760, Status: "published", Sort: 70},
		{Title: "辅助技术赋能特教课堂观摩活动", Summary: "观摩 AAC、读屏等辅助技术的课堂应用。", School: "示范特殊教育学院", Teacher: "王强", Track: "新工科", Edition: "第四届", AwardLevel: "三等奖", SchoolType: "MINISTERIAL", TitleRank: "副高", StartTime: "2026-01-15 10:00:00", EndTime: "2026-01-15 12:00:00", Location: "智慧教室", Organizer: "示范特殊教育学院", SignupCount: 40, ViewCount: 540, Status: "published", Sort: 60},
	}
	for _, item := range activities {
		if err := tx.Where("title = ?", item.Title).FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}
	return nil
}

func seedDemoCases(tx *gorm.DB) error {
	cases := []models.EduCase{
		{Title: "一名听障大学生的融合学习支持案例", Summary: "通过个别化支持帮助听障学生顺利完成专业课程学习。", StudentName: "陈同学", StudentCode: "DEMO-HI-001", Gender: "女", Birthday: "2004-09-16", Stage: "本科", DisabilityType: "听障", AbilityDomain: "社交", CaseType: "融合教育", School: "示范特殊教育学院", Status: "published", ViewCount: 860, Sort: 90},
		{Title: "孤独症学生结构化教学干预个案", Summary: "运用结构化教学改善孤独症学生的课堂适应。", StudentName: "小明", StudentCode: "DEMO-ASD-001", Gender: "男", Birthday: "2010-05-21", Stage: "专科", DisabilityType: "孤独症", AbilityDomain: "认知", CaseType: "干预训练", School: "示范特殊教育学院", Status: "published", ViewCount: 1240, Sort: 80},
		{Title: "视障学生信息技术课程 IEP 管理案例", Summary: "为视障学生定制信息技术课程的个别化教育计划。", StudentName: "小华", StudentCode: "DEMO-VI-001", Gender: "女", Birthday: "2006-02-03", Stage: "本科", DisabilityType: "视障", AbilityDomain: "生活技能", CaseType: "IEP管理", School: "示范特殊教育学院", Status: "published", ViewCount: 980, Sort: 70},
		{Title: "智力障碍学生生活技能评估与干预", Summary: "生活技能评估驱动的个别化干预实践。", StudentName: "小芳", StudentCode: "DEMO-ID-001", Gender: "女", Birthday: "2011-11-08", Stage: "专科", DisabilityType: "智力障碍", AbilityDomain: "生活技能", CaseType: "评估诊断", School: "阳光职业技术学院特教系", Status: "published", ViewCount: 720, Sort: 60},
	}
	for _, item := range cases {
		if err := tx.Where("student_code = ?", item.StudentCode).FirstOrCreate(&item).Error; err != nil {
			return err
		}
		if err := seedDemoCaseDetails(tx, item.Id, item.Title); err != nil {
			return err
		}
	}
	return nil
}

func seedDemoCaseDetails(tx *gorm.DB, caseId int, caseTitle string) error {
	iep := models.EduCaseIEP{
		CaseId: caseId,
		Title:  caseTitle + " IEP 支持计划",
		Goal:   "提升学生课堂参与、沟通表达与独立完成任务的能力。",
		Plan:   "采用评估、目标拆解、课堂支持、家庭协同与阶段复盘的连续干预路径。",
		Status: "published",
	}
	if err := tx.Where("case_id = ? and title = ?", iep.CaseId, iep.Title).FirstOrCreate(&iep).Error; err != nil {
		return err
	}
	assessment := models.EduCaseAssessment{CaseId: caseId, ToolName: "综合能力观察记录", Result: "学生在社交沟通和任务完成方面需要结构化支持。", AssessedAt: "2026-03-05 10:00:00"}
	if err := tx.Where("case_id = ? and tool_name = ?", assessment.CaseId, assessment.ToolName).FirstOrCreate(&assessment).Error; err != nil {
		return err
	}
	intervention := models.EduCaseIntervention{CaseId: caseId, Title: "结构化课堂支持", Content: "使用视觉提示、任务分解和同伴支持提升课堂参与度。", StartDate: "2026-03-15", EndDate: "2026-06-30", Status: "active"}
	return tx.Where("case_id = ? and title = ?", intervention.CaseId, intervention.Title).FirstOrCreate(&intervention).Error
}

func seedDemoExperts(tx *gorm.DB) error {
	experts := []models.EduExpert{
		{Name: "张三", Title: "正高", Organization: "示范特殊教育学院", Specialties: "融合教育、特殊教育政策", Introduction: "长期从事特殊教育研究，主持多项国家级课题，在融合教育领域有丰富的理论与实践经验。", IsRecommended: 1, ViewCount: 2400, FavoriteCount: 120, ShareCount: 35, Sort: 90, Status: 1},
		{Name: "李四", Title: "副高", Organization: "示范特殊教育学院", Specialties: "孤独症干预、行为支持", Introduction: "专注孤独症儿童干预研究，擅长结构化教学与正向行为支持，培训特教教师千余人次。", IsRecommended: 1, ViewCount: 1980, FavoriteCount: 98, ShareCount: 28, Sort: 80, Status: 1},
		{Name: "王五", Title: "正高", Organization: "阳光职业技术学院特教系", Specialties: "辅助技术、无障碍设计", Introduction: "辅助技术领域专家，推动多所高校无障碍环境建设，研发多款特教辅助工具。", IsRecommended: 1, ViewCount: 1560, FavoriteCount: 74, ShareCount: 20, Sort: 70, Status: 1},
		{Name: "赵六", Title: "副高", Organization: "示范特殊教育学院", Specialties: "言语康复、构音矫正", Introduction: "言语病理学背景，擅长言语障碍评估与构音矫正，临床经验丰富。", IsRecommended: 0, ViewCount: 1120, FavoriteCount: 45, ShareCount: 12, Sort: 60, Status: 1},
	}
	for _, item := range experts {
		if err := tx.Where("name = ? and organization = ?", item.Name, item.Organization).FirstOrCreate(&item).Error; err != nil {
			return err
		}
		resources := []models.EduExpertResource{
			{ExpertId: item.Id, Title: item.Name + "专题讲座", Type: "lecture", Status: 1},
			{ExpertId: item.Id, Title: item.Name + "示范课程", Type: "course", Status: 1},
		}
		for _, resource := range resources {
			if err := tx.Where("expert_id = ? and title = ?", resource.ExpertId, resource.Title).FirstOrCreate(&resource).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func seedDemoNewsAndAI(tx *gorm.DB) error {
	news := []models.EduNews{
		{Title: "教育部印发新一轮特殊教育发展提升行动计划", ModuleType: "POLICY", Source: "教育部", Summary: "新一轮行动计划聚焦特殊教育普及与质量提升。", Content: "近日，教育部印发新一轮特殊教育发展提升行动计划，提出进一步扩大特殊教育资源供给、提升融合教育质量、加强特教师资队伍建设等重点举措。", PublishTime: "2026-03-05 09:00:00", ViewCount: 3200, LikeCount: 210, Status: models.NewsStatusPublished, IsTop: 1, Sort: 90},
		{Title: "高校特殊教育专业建设标准研讨会召开", ModuleType: "INDUSTRY", Source: "中国特殊教育学会", Summary: "研讨会就高校特教专业建设标准达成共识。", Content: "日前，高校特殊教育专业建设标准研讨会在京召开，与会专家围绕课程体系、实践教学、师资标准等议题展开深入交流。", PublishTime: "2026-02-28 09:00:00", ViewCount: 1680, LikeCount: 96, Status: models.NewsStatusPublished, Sort: 80},
		{Title: "人工智能辅助特殊教育评估研究取得新进展", ModuleType: "ACADEMIC", Source: "特殊教育研究期刊", Summary: "AI 技术在特教评估中的应用研究获新突破。", Content: "最新研究表明，人工智能技术在特殊儿童发展评估中展现出良好的辅助价值，可提升评估效率与客观性。", PublishTime: "2026-03-12 09:00:00", ViewCount: 2100, LikeCount: 145, Status: models.NewsStatusPublished, Sort: 70},
		{Title: "某高校融合教育实践入选全国优秀案例", ModuleType: "PRACTICE", Source: "示范特殊教育学院", Summary: "融合教育实践经验获全国推广。", Content: "该校构建了评估、支持、跟踪一体化融合教育支持体系，相关实践入选全国特殊教育优秀案例并向全国推广。", PublishTime: "2026-03-18 09:00:00", ViewCount: 1450, LikeCount: 88, Status: models.NewsStatusPublished, Sort: 60},
	}
	for _, item := range news {
		if err := tx.Where("title = ?", item.Title).FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}

	conversation := models.EduAIConversation{ClientKey: "demo-client", Title: "如何为孤独症学生设计 IEP", Mode: "offline"}
	if err := tx.Where("client_key = ? and title = ?", conversation.ClientKey, conversation.Title).FirstOrCreate(&conversation).Error; err != nil {
		return err
	}
	messages := []models.EduAIMessage{
		{ConversationId: conversation.Id, Role: "user", Content: "如何为孤独症学生设计个别化教育计划？"},
		{ConversationId: conversation.Id, Role: "assistant", Content: "建议先完成行为与社交水平评估，再结合课堂观察确定基线，制定可测量的阶段目标，并采用结构化教学、视觉提示和同伴支持等策略定期追踪成效。"},
	}
	for _, item := range messages {
		if err := tx.Where("conversation_id = ? and role = ? and content = ?", item.ConversationId, item.Role, item.Content).FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}
	return nil
}
