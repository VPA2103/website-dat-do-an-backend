package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	fs  FileStore
	rag RAG
	llm Gemini
}

type chatRequest struct {
	ThreadID    *string        `json:"threadId"`
	Message     string         `json:"message"`
	Preferences map[string]any `json:"preferences"`
}

type chatResponse struct {
	ThreadID         string `json:"threadId"`
	AssistantMessage string `json:"assistantMessage"`
}

func NewChatHandler(fs FileStore, rag RAG, llm Gemini) *ChatHandler {
	return &ChatHandler{fs: fs, rag: rag, llm: llm}
}

func (h *ChatHandler) Chat(c *gin.Context) {
	var req chatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"detail": "Invalid request"})
		return
	}

	threadID := ""
	if req.ThreadID != nil {
		threadID = *req.ThreadID
	}
	id, err := h.fs.EnsureThread(threadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to create thread"})
		return
	}

	if err := h.fs.AppendThreadMessage(id, "user", req.Message); err != nil {
		log.Println("❌ Append user message failed:", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"detail": err.Error(),
		})
		return
	}

	ctxText, err := h.rag.RetrieveContext(
		c.Request.Context(),
		req.Message,
		6,
	)
	//ctxText := "" // 👈 tạm thời không dùng RAG
	if err != nil {
		// RAG context is best-effort; keep chat available even if vector search or embeddings fail.
		log.Printf("rag RetrieveContext failed: %v", err)
		ctxText = ""
	}

	history, err := h.fs.GetThreadMessages(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to load thread history"})
		return
	}
	contents := buildChatContents(ctxText, history, req.Message)
	if req.Preferences != nil && len(req.Preferences) > 0 && len(contents) > 0 {
		if b, err := json.Marshal(req.Preferences); err == nil {
			contents[len(contents)-1] = contents[len(contents)-1] + "\nUser preferences (JSON): " + string(b)
		}
	}
	assistant, rateLimit, err := h.llm.Generate(c.Request.Context(), contents)
	if rateLimit != nil {
		if rateLimit.RetryAfterSeconds != nil {
			c.Header("Retry-After", strconv.Itoa(*rateLimit.RetryAfterSeconds))
		}
		c.JSON(http.StatusTooManyRequests, gin.H{"detail": gin.H{
			"error":               "RATE_LIMIT",
			"message":             rateLimit.Message,
			"retry_after_seconds": rateLimit.RetryAfterSeconds,
		}})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "LLM error"})
		return
	}

	if err := h.fs.AppendThreadMessage(id, "assistant", assistant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"detail": "Failed to persist message"})
		return
	}
	c.JSON(http.StatusOK, chatResponse{ThreadID: id, AssistantMessage: assistant})
}

func buildChatContents(contextText string, history []Message, message string) []string {
	contents := make([]string, 0, 1+len(history)+1)
	if contextText != "" {
		contents = append(contents, "CONTEXT (menu/nhà hàng):\n"+contextText)
	}

	start := 0
	if len(history) > 12 {
		start = len(history) - 12
	}
	for _, m := range history[start:] {
		if m.Content == "" {
			continue
		}
		switch m.Role {
		case "user":
			contents = append(contents, "User: "+m.Content)
		case "assistant":
			contents = append(contents, "Assistant: "+m.Content)
		}
	}

	contents = append(contents, "User: "+message)
	return contents
}
