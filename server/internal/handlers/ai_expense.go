package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type AIExpenseHandler struct{}

func NewAIExpenseHandler() *AIExpenseHandler {
	return &AIExpenseHandler{}
}

type ExtractExpenseRequest struct {
	InputData struct {
		Paragraph  string `json:"paragraph"`
		Categories []struct {
			CategoryID string `json:"category_id"`
			Name       string `json:"name"`
			IsDefault  bool   `json:"is_default"`
		} `json:"categories"`
	} `json:"input_data"`
	ConversationHistory []interface{} `json:"conversation_history"`
}

type ExtractExpenseResponse struct {
	Success    bool        `json:"success"`
	OutputData *OutputData `json:"output_data,omitempty"`
	Error      *string     `json:"error,omitempty"`
}

type OutputData struct {
	Expenses       []ExtractedExpense `json:"expenses"`
	TotalAmount    float64            `json:"total_amount"`
	Currency       string             `json:"currency"`
	CategoriesUsed []Category         `json:"categories_used"`
}

type ExtractedExpense struct {
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	CategoryID  string  `json:"category_id"`
	Description string  `json:"description"`
	Date        *string `json:"date"`
	Merchant    *string `json:"merchant"`
}

type Category struct {
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
	IsDefault  bool   `json:"is_default"`
}

// ExtractExpenses handles POST /expenses/extract
func (h *AIExpenseHandler) ExtractExpenses(c *fiber.Ctx) error {
	var req ExtractExpenseRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Invalid request body",
		})
	}

	// Validate input
	if req.InputData.Paragraph == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Paragraph is required",
		})
	}

	// Get API key from environment
	apiKey := os.Getenv("AI_EXPENSE_API_KEY")
	if apiKey == "" {
		c.Context().Logger().Printf("[ERROR] AI_EXPENSE_API_KEY environment variable not set")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "AI service not configured",
		})
	}
	c.Context().Logger().Printf("[INFO] AI Extraction request - API Key loaded, length: %d", len(apiKey))

	// Prepare request to AI service
	// Use environment variable if set, otherwise use production URL
	apiURL := os.Getenv("AI_EXPENSE_API_URL")
	if apiURL == "" {
		apiURL = "https://multi-service-chatbot.onrender.com/chat/expense_tracker/extract_expenses"
	}
	c.Context().Logger().Printf("[INFO] Preparing request to AI service: %s", apiURL)
	c.Context().Logger().Printf("[INFO] Paragraph length: %d, Categories count: %d", len(req.InputData.Paragraph), len(req.InputData.Categories))

	requestBody, err := json.Marshal(req)
	if err != nil {
		c.Context().Logger().Printf("[ERROR] Failed to marshal request: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to prepare request",
		})
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		c.Context().Logger().Printf("[ERROR] Failed to create HTTP request: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to create request",
		})
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-Key", apiKey)

	// Make request to AI service
	c.Context().Logger().Printf("[INFO] Sending request to AI service...")
	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		c.Context().Logger().Printf("[ERROR] Failed to connect to AI service: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to connect to AI service: " + err.Error(),
		})
	}
	defer resp.Body.Close()

	c.Context().Logger().Printf("[INFO] Received response from AI service - Status: %d", resp.StatusCode)

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Context().Logger().Printf("[ERROR] Failed to read response body: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to read response",
		})
	}

	c.Context().Logger().Printf("[INFO] Response body length: %d bytes", len(body))

	// Parse response
	var aiResponse ExtractExpenseResponse
	if err := json.Unmarshal(body, &aiResponse); err != nil {
		c.Context().Logger().Printf("[ERROR] Failed to parse AI response: %v. Body: %s", err, string(body))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to parse AI response",
		})
	}

	// Check if expenses were extracted
	expenseCount := 0
	if aiResponse.OutputData != nil && aiResponse.OutputData.Expenses != nil {
		expenseCount = len(aiResponse.OutputData.Expenses)
	}

	c.Context().Logger().Printf("[INFO] Returning AI response - Success: %v, Expenses count: %d", aiResponse.Success, expenseCount)

	// If no expenses extracted but success is true, it might be a mock/unconfigured service
	if aiResponse.Success && expenseCount == 0 {
		c.Context().Logger().Printf("[WARN] AI service returned success but no expenses. Response body: %s", string(body))
	}

	return c.Status(resp.StatusCode).JSON(aiResponse)
}
