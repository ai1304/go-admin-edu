<template>
  <PortalLayout>
    <section class="ai-page">
      <aside class="conversation-panel">
        <a-button type="primary" long @click="newChat">新建对话</a-button>
        <div class="conversation-list">
          <button
            v-for="item in conversations"
            :key="item.id"
            class="conversation-item"
            :class="{ active: currentId === item.id }"
            @click="openConversation(item.id)"
          >
            <span>{{ item.title || "新对话" }}</span>
            <small>{{ item.mode === "online" ? "在线" : "离线" }}</small>
          </button>
          <a-empty v-if="!conversations.length" description="暂无对话" />
        </div>
      </aside>

      <main class="chat-panel">
        <header>
          <div>
            <h1>特教智能助手</h1>
            <p>{{ mode === "online" ? "已连接在线大模型" : "当前使用内置专业规则应答" }}</p>
          </div>
          <a-tag :color="mode === 'online' ? 'green' : 'gray'">{{ mode === "online" ? "在线模型" : "离线规则" }}</a-tag>
        </header>

        <div ref="messageBox" class="message-box">
          <div v-if="!messages.length" class="welcome">
            <h2>可以从这些问题开始</h2>
            <div class="suggestions">
              <button v-for="item in suggestions" :key="item" @click="quickAsk(item)">{{ item }}</button>
            </div>
          </div>
          <div v-for="(message, index) in messages" :key="index" class="message" :class="message.role">
            <div class="bubble">{{ message.content }}</div>
          </div>
          <div v-if="sending" class="message assistant">
            <div class="bubble">思考中...</div>
          </div>
        </div>

        <footer class="chat-input">
          <a-textarea v-model="input" :auto-size="{ minRows: 2, maxRows: 4 }" placeholder="输入你的问题，Enter 发送，Shift+Enter 换行" @keydown.enter.prevent="send" />
          <a-button type="primary" :loading="sending" @click="send">发送</a-button>
        </footer>
      </main>
    </section>
  </PortalLayout>
</template>

<script setup>
import { nextTick, onMounted, ref } from "vue";
import PortalLayout from "@/layouts/PortalLayout.vue";
import { chatWithAI, getAIConversation, getAIConversations } from "@/api/ai";
import { useSessionStore } from "@/stores/session";

const session = useSessionStore();
const conversations = ref([]);
const messages = ref([]);
const currentId = ref(0);
const input = ref("");
const sending = ref(false);
const mode = ref("offline");
const messageBox = ref(null);
const clientKey = getClientKey();
const suggestions = [
  "如何为孤独症学生编制 IEP？",
  "听障大学生课堂沟通有哪些支持策略？",
  "特殊学生评估应该关注哪些方面？",
  "融合教育课堂如何做差异化教学？"
];

function getClientKey() {
  let value = localStorage.getItem("portalClientKey");
  if (!value) {
    value = `visitor-${Date.now()}-${Math.random().toString(16).slice(2)}`;
    localStorage.setItem("portalClientKey", value);
  }
  return value;
}

function identity() {
  return {
    userId: session.user?.id || session.user?.userId || 0,
    clientKey
  };
}

async function loadConversations() {
  const res = await getAIConversations(identity());
  conversations.value = res.data || [];
}

async function openConversation(id) {
  currentId.value = id;
  const res = await getAIConversation(id);
  const data = res.data || {};
  messages.value = data.messages || [];
  mode.value = data.mode || "offline";
  scrollBottom();
}

function newChat() {
  currentId.value = 0;
  messages.value = [];
  mode.value = "offline";
}

function quickAsk(text) {
  input.value = text;
  send();
}

async function send(event) {
  if (event?.shiftKey) {
    input.value += "\n";
    return;
  }
  const text = input.value.trim();
  if (!text || sending.value) return;
  messages.value.push({ role: "user", content: text });
  input.value = "";
  sending.value = true;
  scrollBottom();
  try {
    const res = await chatWithAI({ ...identity(), conversationId: currentId.value || undefined, message: text });
    const data = res.data || {};
    currentId.value = data.conversationId;
    mode.value = data.mode || "offline";
    messages.value.push({ role: "assistant", content: data.reply });
    await loadConversations();
  } finally {
    sending.value = false;
    scrollBottom();
  }
}

function scrollBottom() {
  nextTick(() => {
    if (messageBox.value) {
      messageBox.value.scrollTop = messageBox.value.scrollHeight;
    }
  });
}

onMounted(loadConversations);
</script>

<style scoped>
.ai-page {
  display: grid;
  grid-template-columns: 260px minmax(0, 1fr);
  gap: 16px;
  height: calc(100vh - 128px);
  min-height: 620px;
  padding-top: 24px;
}

.conversation-panel,
.chat-panel {
  background: #fff;
  border: 1px solid #e5e6eb;
  border-radius: 8px;
}

.conversation-panel {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 14px;
  min-width: 0;
}

.conversation-list {
  display: grid;
  gap: 8px;
  overflow: auto;
}

.conversation-item {
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 8px;
  align-items: center;
  width: 100%;
  padding: 10px;
  color: #1d2129;
  text-align: left;
  background: #f7f8fa;
  border: 1px solid transparent;
  border-radius: 8px;
  cursor: pointer;
}

.conversation-item span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.conversation-item small {
  color: #86909c;
}

.conversation-item.active {
  color: #176fd6;
  background: #eef6ff;
  border-color: #bedaff;
}

.chat-panel {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.chat-panel header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  padding: 18px 22px;
  border-bottom: 1px solid #e5e6eb;
}

.chat-panel h1 {
  margin: 0;
  font-size: 24px;
}

.chat-panel p {
  margin: 4px 0 0;
  color: #86909c;
}

.message-box {
  flex: 1;
  overflow: auto;
  padding: 22px;
  background: #fbfdff;
}

.welcome {
  max-width: 720px;
  margin: 80px auto;
  text-align: center;
}

.welcome h2 {
  margin-bottom: 18px;
  font-size: 24px;
}

.suggestions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  justify-content: center;
}

.suggestions button {
  padding: 9px 14px;
  color: #176fd6;
  background: #eef6ff;
  border: 1px solid #bedaff;
  border-radius: 999px;
  cursor: pointer;
}

.message {
  display: flex;
  margin-bottom: 14px;
}

.message.user {
  justify-content: flex-end;
}

.bubble {
  max-width: min(720px, 78%);
  padding: 12px 15px;
  line-height: 1.8;
  white-space: pre-wrap;
  word-break: break-word;
  background: #f1f4f8;
  border-radius: 8px;
}

.message.user .bubble {
  color: #fff;
  background: #176fd6;
}

.chat-input {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 96px;
  gap: 12px;
  padding: 14px 18px;
  border-top: 1px solid #e5e6eb;
}

@media (max-width: 900px) {
  .ai-page {
    grid-template-columns: 1fr;
    height: auto;
  }

  .conversation-panel {
    max-height: 240px;
  }
}
</style>
