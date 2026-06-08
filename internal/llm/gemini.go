package llm

import (
	"context"
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/vpa/quanlynhahang-backend/config"
	"github.com/vpa/quanlynhahang-backend/models/AIChatBot"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const systemInstruction = "Bạn là trợ lý tư vấn món ăn cho nhà hàng. " +
	"Phong cách: Lịch sự, niềm nở, sử dụng các từ như 'Dạ', 'Mời quý khách', 'Rất vui được tư vấn', 'Mình gợi ý', v.v. " +
	"Bạn phải tư vấn dựa trên dữ liệu menu/nhà hàng được cung cấp trong CONTEXT. " +
	"QUY TẮC BẮT BUỘC: Bạn CHỈ được gợi ý/đề xuất những món xuất hiện trong CONTEXT, mục 'Menu liên quan' (các dòng bắt đầu bằng '- Món: ...' hoặc có tên món trong các bullet). " +
	"KHÔNG được tự bịa món, không được đề xuất món ngoài menu. " +
	"QUY TẮC AN TOÀN (dị ứng/kiêng): Nếu khách nói 'dị ứng', 'tránh', 'không ăn/không dùng' một nguyên liệu/loại thịt (ví dụ: thịt bò), bạn TUYỆT ĐỐI không được gợi ý các món có chứa nguyên liệu đó. Hãy dựa vào thông tin trong CONTEXT, đặc biệt các phần 'Nguyên liệu:' và 'Dị ứng:' của từng món. Nếu CONTEXT không đủ thông tin để chắc chắn món có/không chứa nguyên liệu bị tránh, hãy hỏi lại để làm rõ và ưu tiên gợi ý các món mà bạn chắc chắn phù hợp. " +
	"Nếu CONTEXT không có mục 'Menu liên quan' hoặc danh sách menu trống/không có dữ liệu, hãy nói rõ hiện chưa có dữ liệu menu trong hệ thống và hỏi khách muốn mô tả khẩu vị/ngân sách để tư vấn chung, hoặc đề nghị admin cập nhật menu."

type Gemini struct {
	client *genai.Client

	models      []string
	embModels   []string
	temperature float32
}

func NewGemini(cfg config.GeminiConfig) (*Gemini, error) {
	if strings.TrimSpace(cfg.APIKey) == "" {
		return nil, errors.New("missing GEMINI_API_KEY")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := genai.NewClient(ctx, option.WithAPIKey(cfg.APIKey))
	if err != nil {
		return nil, err
	}

	preferred := strings.TrimSpace(cfg.Model)
	candidates := []string{preferred, "gemini-2.5-flash-lite", "gemini-2.5-flash", "gemini-3.1-flash-lite", "gemini-2.0-flash-lite"}
	models := dedupeNonEmpty(candidates)

	embPreferred := strings.TrimSpace(cfg.EmbeddingModel)
	embCandidates := []string{embPreferred, "gemini-embedding-2", "gemini-embedding-001", "text-embedding-004"}
	embModels := dedupeNonEmpty(embCandidates)

	return &Gemini{client: client, models: models, embModels: embModels, temperature: 0.4}, nil
}

func (g *Gemini) Embed(ctx context.Context, text string) ([]float32, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return []float32{}, nil
	}

	var lastErr error
	for _, modelName := range g.embModels {
		m := g.client.EmbeddingModel(modelName)
		resp, err := m.EmbedContent(ctx, genai.Text(text))
		if err != nil {
			lastErr = err
			continue
		}
		if resp == nil || resp.Embedding == nil {
			return []float32{}, nil
		}
		vals := resp.Embedding.Values
		out := make([]float32, 0, len(vals))
		for _, v := range vals {
			out = append(out, float32(v))
		}
		return out, nil
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return []float32{}, nil
}

func (g *Gemini) Generate(ctx context.Context, contents []string) (string, *core.RateLimitError, error) {
	prompt := strings.Join(contents, "\n")
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return "", nil, nil
	}

	var lastErr error
	for _, modelName := range g.models {
		m := g.client.GenerativeModel(modelName)
		m.Temperature = &g.temperature
		m.SystemInstruction = &genai.Content{Parts: []genai.Part{genai.Text(systemInstruction)}}

		resp, err := m.GenerateContent(ctx, genai.Text(prompt))
		if err != nil {
			// Best-effort detect 429-like errors by message.
			if isRateLimitErr(err) {
				retry := extractRetryAfterSeconds(err.Error())
				return "", &core.RateLimitError{Message: "Gemini API rate limit / quota exceeded", RetryAfterSeconds: retry}, nil
			}
			lastErr = err
			continue
		}
		text := extractText(resp)
		text = strings.TrimSpace(text)
		if text == "" {
			text = "Mình chưa có đủ dữ liệu để tư vấn. Bạn cho mình biết bạn thích vị nào (cay/ngọt), ngân sách, và bạn đang muốn ăn món gì?"
		}
		return text, nil, nil
	}
	if lastErr != nil {
		return "", nil, lastErr
	}
	return "", nil, errors.New("generate failed")
}

func dedupeNonEmpty(in []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, s := range in {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		if seen[s] {
			continue
		}
		seen[s] = true
		out = append(out, s)
	}
	return out
}

func extractText(resp *genai.GenerateContentResponse) string {
	if resp == nil || len(resp.Candidates) == 0 {
		return ""
	}
	c := resp.Candidates[0]
	if c.Content == nil {
		return ""
	}
	var b strings.Builder
	for _, p := range c.Content.Parts {
		if t, ok := p.(genai.Text); ok {
			b.WriteString(string(t))
		}
	}
	return b.String()
}

func isRateLimitErr(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "429") || strings.Contains(msg, "resource_exhausted") || strings.Contains(msg, "quota")
}

var retryRe = regexp.MustCompile(`(?i)(?:retry\s+in\s+)?([0-9]+(?:\.[0-9]+)?)s`)

func extractRetryAfterSeconds(msg string) *int {
	m := retryRe.FindStringSubmatch(msg)
	if len(m) < 2 {
		return nil
	}
	// parse integer seconds
	parts := strings.Split(m[1], ".")
	sec := 0
	for _, ch := range parts[0] {
		if ch < '0' || ch > '9' {
			return nil
		}
		sec = sec*10 + int(ch-'0')
	}
	if sec < 0 {
		return nil
	}
	return &sec
}
