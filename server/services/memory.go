package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"server/config"
	"server/models"
	"time"
)

type MemoryService struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

func NewMemoryService() *MemoryService {
	apiKey := config.Config("MEM0_API_KEY")
	if apiKey == "" {
		panic("MEM0_API_KEY environment variable is required")
	}

	return &MemoryService{
		apiKey:  apiKey,
		baseURL: "https://api.mem0.ai",
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (ms *MemoryService) AddMemory(request models.Mem0AddRequest) ([]models.Memory, error) {
	url := fmt.Sprintf("%s/v1/memories/", ms.baseURL)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token %s", ms.apiKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := ms.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var memories []models.Memory
	if err := json.Unmarshal(body, &memories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return memories, nil
}

func (ms *MemoryService) SearchMemories(request models.Mem0SearchRequest) ([]models.Memory, error) {
	url := fmt.Sprintf("%s/v2/memories/search/", ms.baseURL)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token %s", ms.apiKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := ms.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var memories []models.Memory
	if err := json.Unmarshal(body, &memories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return memories, nil
}

func (ms *MemoryService) GetMemories(request models.Mem0GetRequest) ([]models.Memory, error) {
	url := fmt.Sprintf("%s/v2/memories/", ms.baseURL)

	jsonData, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token %s", ms.apiKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := ms.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var memories []models.Memory
	if err := json.Unmarshal(body, &memories); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return memories, nil
}

func (ms *MemoryService) GetMemory(memoryID string) (*models.Memory, error) {
	url := fmt.Sprintf("%s/v1/memories/%s/", ms.baseURL, memoryID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token %s", ms.apiKey))

	resp, err := ms.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	var memory models.Memory
	if err := json.Unmarshal(body, &memory); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &memory, nil
}

func (ms *MemoryService) DeleteMemory(memoryID string) error {
	url := fmt.Sprintf("%s/v1/memories/%s/", ms.baseURL, memoryID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token %s", ms.apiKey))

	resp, err := ms.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (ms *MemoryService) BatchDeleteMemories(memoryIDs []string) error {
	url := fmt.Sprintf("%s/v1/batch/", ms.baseURL)

	request := models.Mem0BatchDeleteRequest{
		MemoryIds: memoryIDs,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("DELETE", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Token %s", ms.apiKey))
	req.Header.Add("Content-Type", "application/json")

	resp, err := ms.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func (ms *MemoryService) GetUserMemories(userID string) ([]models.Memory, error) {
	request := models.Mem0GetRequest{
		Filters: map[string]interface{}{
			"user_id": userID,
		},
	}
	return ms.GetMemories(request)
}

func (ms *MemoryService) SearchUserMemories(userID, query string, topK int) ([]models.Memory, error) {
	request := models.Mem0SearchRequest{
		Query: query,
		Filters: map[string]interface{}{
			"user_id": userID,
		},
		TopK:   topK,
		Rerank: true,
	}
	return ms.SearchMemories(request)
}

func (ms *MemoryService) AddChatMemory(userID, sessionID, content, role string) ([]models.Memory, error) {
	request := models.Mem0AddRequest{
		Messages: []map[string]string{
			{
				"role":    role,
				"content": content,
			},
		},
		UserID: userID,
		RunID:  sessionID,
		Metadata: map[string]interface{}{
			"session_id": sessionID,
			"timestamp":  time.Now().Unix(),
		},
		Infer:   true,
		Version: "v2",
	}
	return ms.AddMemory(request)
}
