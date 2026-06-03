package models

type EduAIConversation struct {
	BaseModel
	UserId    int            `json:"userId" gorm:"index;comment:用户ID"`
	ClientKey string         `json:"clientKey" gorm:"size:128;index;comment:访客标识"`
	Title     string         `json:"title" gorm:"size:255;comment:会话标题"`
	Mode      string         `json:"mode" gorm:"size:32;default:offline;comment:应答模式"`
	Messages  []EduAIMessage `json:"messages" gorm:"-"`
}

func (*EduAIConversation) TableName() string {
	return "edu_ai_conversation"
}

type EduAIMessage struct {
	BaseModel
	ConversationId int    `json:"conversationId" gorm:"index;comment:会话ID"`
	Role           string `json:"role" gorm:"size:32;index;comment:角色"`
	Content        string `json:"content" gorm:"type:text;comment:消息内容"`
}

func (*EduAIMessage) TableName() string {
	return "edu_ai_message"
}
