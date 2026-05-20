package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-admin/app/edu/models"
	"go-admin/common/dto"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"gorm.io/gorm"
)

type EduAI struct {
	api.Api
}

type aiChatReq struct {
	ConversationId int    `json:"conversationId"`
	UserId         int    `json:"userId"`
	ClientKey      string `json:"clientKey"`
	Message        string `json:"message"`
}

type aiConversationQuery struct {
	dto.Pagination
	UserId    int    `form:"userId"`
	ClientKey string `form:"clientKey"`
	Keyword   string `form:"keyword"`
}

func (e EduAI) Chat(c *gin.Context) {
	req := aiChatReq{}
	if err := e.MakeContext(c).MakeOrm().Bind(&req, binding.JSON).Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	req.Message = strings.TrimSpace(req.Message)
	if req.Message == "" {
		e.Error(400, nil, "message is required")
		return
	}

	conversation := models.EduAIConversation{}
	if req.ConversationId != 0 {
		if err := e.Orm.First(&conversation, req.ConversationId).Error; err != nil {
			e.Error(404, err, "conversation not found")
			return
		}
	} else {
		title := req.Message
		if len([]rune(title)) > 24 {
			title = string([]rune(title)[:24])
		}
		conversation = models.EduAIConversation{
			UserId:    req.UserId,
			ClientKey: req.ClientKey,
			Title:     title,
			Mode:      "offline",
		}
		conversation.SetCreateBy(req.UserId)
		if err := e.Orm.Create(&conversation).Error; err != nil {
			e.Error(500, err, "create conversation failed")
			return
		}
	}

	userMsg := models.EduAIMessage{ConversationId: conversation.Id, Role: "user", Content: req.Message}
	userMsg.SetCreateBy(req.UserId)
	if err := e.Orm.Create(&userMsg).Error; err != nil {
		e.Error(500, err, "save message failed")
		return
	}

	history := make([]models.EduAIMessage, 0)
	_ = e.Orm.Where("conversation_id = ?", conversation.Id).Order("id asc").Find(&history).Error
	reply, online := generateAIReply(history)
	mode := "offline"
	if online {
		mode = "online"
	}
	assistantMsg := models.EduAIMessage{ConversationId: conversation.Id, Role: "assistant", Content: reply}
	if err := e.Orm.Create(&assistantMsg).Error; err != nil {
		e.Error(500, err, "save reply failed")
		return
	}
	_ = e.Orm.Model(&models.EduAIConversation{}).Where("id = ?", conversation.Id).Updates(map[string]interface{}{"mode": mode}).Error

	e.OK(gin.H{"conversationId": conversation.Id, "reply": reply, "mode": mode}, "chat success")
}

func (e EduAI) MyConversations(c *gin.Context) {
	req := aiConversationQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduAIConversation, 0)
	db := e.Orm.Model(&models.EduAIConversation{})
	if req.UserId != 0 {
		db = db.Where("user_id = ?", req.UserId)
	}
	if req.ClientKey != "" {
		db = db.Where("client_key = ?", req.ClientKey)
	}
	if req.Keyword != "" {
		db = db.Where("title like ?", "%"+req.Keyword+"%")
	}
	if err := db.Order("updated_at desc,id desc").Limit(50).Find(&list).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	e.OK(list, "query success")
}

func (e EduAI) ConversationDetail(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var data models.EduAIConversation
	if err := e.Orm.First(&data, c.Param("id")).Error; err != nil {
		e.Error(404, err, "conversation not found")
		return
	}
	messages := make([]models.EduAIMessage, 0)
	_ = e.Orm.Where("conversation_id = ?", data.Id).Order("id asc").Find(&messages).Error
	data.Messages = messages
	e.OK(data, "query success")
}

func (e EduAI) DeleteConversation(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	id := parsePathId(c.Param("id"))
	if err := e.Orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("conversation_id = ?", id).Delete(&models.EduAIMessage{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.EduAIConversation{}, id).Error
	}); err != nil {
		e.Error(500, err, "delete failed")
		return
	}
	e.OK(id, "delete success")
}

func (e EduAI) PublicDeleteConversation(c *gin.Context) {
	req := aiConversationQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	if req.UserId == 0 && req.ClientKey == "" {
		e.Error(400, nil, "missing conversation identity")
		return
	}
	id := parsePathId(c.Param("id"))
	db := e.Orm.Model(&models.EduAIConversation{}).Where("id = ?", id)
	if req.UserId != 0 {
		db = db.Where("user_id = ?", req.UserId)
	} else {
		db = db.Where("client_key = ?", req.ClientKey)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	if count == 0 {
		e.Error(404, nil, "conversation not found")
		return
	}
	if err := e.Orm.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("conversation_id = ?", id).Delete(&models.EduAIMessage{}).Error; err != nil {
			return err
		}
		return tx.Delete(&models.EduAIConversation{}, id).Error
	}); err != nil {
		e.Error(500, err, "delete failed")
		return
	}
	e.OK(id, "delete success")
}

func (e EduAI) AdminConversations(c *gin.Context) {
	req := aiConversationQuery{}
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	_ = c.ShouldBindQuery(&req)
	list := make([]models.EduAIConversation, 0)
	db := applyEduUserScope(c, e.Orm.Model(&models.EduAIConversation{}))
	if req.Keyword != "" {
		db = db.Where("title like ? or client_key like ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	if err := db.Order("updated_at desc,id desc").Limit(req.GetPageSize()).Offset((req.GetPageIndex() - 1) * req.GetPageSize()).Find(&list).Error; err != nil {
		e.Error(500, err, "query failed")
		return
	}
	e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "query success")
}

func (e EduAI) AdminStats(c *gin.Context) {
	if err := e.MakeContext(c).MakeOrm().Errors; err != nil {
		e.Error(500, err, err.Error())
		return
	}
	var conversations int64
	var messages int64
	_ = e.Orm.Model(&models.EduAIConversation{}).Count(&conversations).Error
	_ = e.Orm.Model(&models.EduAIMessage{}).Count(&messages).Error
	e.OK(gin.H{
		"conversationCount": conversations,
		"messageCount":      messages,
		"enabled":           aiEnabled(),
		"model":             aiModel(),
	}, "query success")
}

func generateAIReply(history []models.EduAIMessage) (string, bool) {
	if aiEnabled() {
		if reply, err := callOpenAICompatible(history); err == nil && strings.TrimSpace(reply) != "" {
			return reply, true
		}
	}
	last := ""
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Role == "user" {
			last = history[i].Content
			break
		}
	}
	return fallbackAIReply(last), false
}

func aiBaseURL() string {
	if value := strings.TrimRight(os.Getenv("APP_AI_BASE_URL"), "/"); value != "" {
		return value
	}
	return "https://api.deepseek.com"
}

func aiModel() string {
	if value := os.Getenv("APP_AI_MODEL"); value != "" {
		return value
	}
	return "deepseek-chat"
}

func aiEnabled() bool {
	return os.Getenv("APP_AI_API_KEY") != ""
}

func callOpenAICompatible(history []models.EduAIMessage) (string, error) {
	type chatMessage struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	messages := []chatMessage{{
		Role:    "system",
		Content: "你是高校特殊教育资源平台的智能助手，回答应专业、温和、可操作，重点覆盖 IEP、评估、干预训练、融合教育、资源建设和教师备课。",
	}}
	start := len(history) - 12
	if start < 0 {
		start = 0
	}
	for _, item := range history[start:] {
		messages = append(messages, chatMessage{Role: item.Role, Content: item.Content})
	}
	body := map[string]interface{}{
		"model":       aiModel(),
		"messages":    messages,
		"temperature": 0.7,
	}
	payload, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, aiBaseURL()+"/chat/completions", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("APP_AI_API_KEY"))
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("ai status %d", resp.StatusCode)
	}
	var decoded struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return "", err
	}
	if len(decoded.Choices) == 0 {
		return "", fmt.Errorf("empty choices")
	}
	return decoded.Choices[0].Message.Content, nil
}

func fallbackAIReply(question string) string {
	q := strings.TrimSpace(question)
	switch {
	case strings.Contains(q, "IEP") || strings.Contains(q, "个别化"):
		return "IEP 编制建议：先做现状评估，明确学生优势、需求和支持条件；再把长期目标拆成可观察、可测量的短期目标；最后同步记录教学调整、家庭协同和阶段性评估结果。目标建议使用“条件 + 行为 + 标准”的句式，便于后续追踪。"
	case strings.Contains(q, "孤独症") || strings.Contains(q, "自闭"):
		return "孤独症学生支持可以从结构化环境、视觉提示、任务分解和正向行为支持入手。先降低环境不确定性，再用清晰流程帮助学生理解任务，同时把社交目标放进真实课堂活动中逐步练习。"
	case strings.Contains(q, "评估") || strings.Contains(q, "诊断"):
		return "特殊学生评估建议采用多元证据：标准化量表、课堂观察、作品分析、访谈记录和情境任务表现都要纳入。评估结论不只描述不足，也要标出可利用的优势和下一步干预入口。"
	case strings.Contains(q, "融合"):
		return "融合教育实施可从课程调整、同伴支持、无障碍环境和教师协作四方面推进。建议先确定学生参与普通课堂的关键障碍，再为每个障碍配置具体支持，而不是只做笼统陪读。"
	case strings.Contains(q, "备课") || strings.Contains(q, "教学"):
		return "特教备课建议同时准备核心目标、替代目标和拓展目标；材料尽量多感官呈现；任务难度设置阶梯；课堂中保留即时反馈和重复练习空间。这样更容易兼顾不同能力水平的学生。"
	default:
		return "你好，我是特教智能助手。你可以向我询问 IEP 编制、学生评估、干预训练、融合教育、课堂备课和资源建设等问题。当前未配置在线大模型 Key 时，我会使用内置专业规则进行回答。"
	}
}
