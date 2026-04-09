package corehandler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type groqMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type groqRequest struct {
	Model       string        `json:"model"`
	Messages    []groqMessage `json:"messages"`
	Temperature float64       `json:"temperature"`
}

type groqChoice struct {
	Message groqMessage `json:"message"`
}

type groqResponse struct {
	Choices []groqChoice `json:"choices"`
}

const systemPrompt = `Você é um assistente especializado em criar orçamentos de desenvolvimento de software.
Com base no prompt do usuário, gere um orçamento estruturado em JSON com o seguinte formato exato:

{
  "clientName": "nome do cliente",
  "scope": "descrição detalhada do escopo em HTML (pode usar <strong>, <em>, <p>)",
  "conditions": "condições comerciais padrão em HTML",
  "hourlyRate": 0,
  "items": [
    {
      "id": "uuid-1",
      "name": "nome da atividade",
      "quantity": 1,
      "unitPrice": 0,
      "estimateHours": 2,
      "itemStatus": "Aguard. Aprovação"
    }
  ]
}

Regras:
- Extraia o nome do cliente do prompt se mencionado
- Divida o trabalho em atividades específicas e realistas
- Use as horas estimadas mencionadas pelo usuário, distribuindo entre as atividades
- Se o usuário mencionar valor por hora (ex: "R$ 50/hora", "cobrar 90 a hora", "50 reais por hora"), extraia e coloque em hourlyRate como número
- Se não mencionar valor por hora, deixe hourlyRate como 0
- Responda APENAS com o JSON, sem texto adicional, sem markdown, sem blocos de código`

func GenerateQuote(c *gin.Context) {
	apiKey := os.Getenv("GROQ_API_KEY")
	if apiKey == "" {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "IA não configurada"})
		return
	}

	var body struct {
		Prompt string `json:"prompt" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prompt é obrigatório"})
		return
	}

	payload := groqRequest{
		Model:       "llama-3.3-70b-versatile",
		Temperature: 0.3,
		Messages: []groqMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: body.Prompt},
		},
	}

	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "https://api.groq.com/openai/v1/chat/completions", bytes.NewBuffer(payloadBytes))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao chamar IA"})
		return
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("groq error: %s", string(respBytes))})
		return
	}

	var groqResp groqResponse
	if err := json.Unmarshal(respBytes, &groqResp); err != nil || len(groqResp.Choices) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "resposta inválida da IA"})
		return
	}

	content := groqResp.Choices[0].Message.Content

	// Validate it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "IA retornou formato inválido", "raw": content})
		return
	}

	c.JSON(http.StatusOK, result)
}
