package services

import (
	"backend/internal/config"
	"backend/internal/repositories"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type N8NService struct {
	webhookURL string

	httpClient *http.Client
}

type N8NReadinessInput struct {
	ClientName      string `json:"client_name"`
	ClientEmail     string `json:"client_email"`
	ServicePackage  string `json:"service_package"`
	ProjectGoal     string `json:"project_goal"`
	DesiredTimeline string `json:"desired_timeline"`
	AssetsProvided  string `json:"assets_provided"`
}

type N8NReadinessResult struct {
	ReadinessStatus       string   `json:"readiness_status"`
	MissingItems          []string `json:"missing_items"`
	AISummary             string   `json:"ai_summary"`
	RecommendedNextAction string   `json:"recommended_next_action"`
}

func NewN8NService(cfg config.Config) *N8NService {
	return &N8NService{
		webhookURL: strings.TrimSpace(cfg.N8NWebhookURL),
		httpClient: &http.Client{
			Timeout: time.Duration(30) * time.Second,
		},
	}
}

func (s *N8NService) CheckReadiness(ctx context.Context, input N8NReadinessInput) (*N8NReadinessResult, error) {
	if s.webhookURL == "" {
		return nil, fmt.Errorf("webhookURL is missing")
	}

	body, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("json marshal failed")
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, s.webhookURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("n8n request failed")
	}
	request.Header.Set("Content-Type", "application/json")

	request.Header.Set("Accept", "application/json")

	response, err := s.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("call n8n webhook")
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {

		responseBody, _ := io.ReadAll(io.LimitReader(response.Body, 2048))

		return nil, fmt.Errorf("n8n returned status %d: %s", response.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var result N8NReadinessResult

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode n8n response failed here: %w", err)
	}

	return &result, nil
}

func ReadinessToCreateSubmissionInput(userID string, form N8NReadinessInput, result *N8NReadinessResult) repositories.CreateSubmissionInput {
	return repositories.CreateSubmissionInput{
		UserID:                userID,
		ClientName:            form.ClientName,
		ClientEmail:           form.ClientEmail,
		ServicePackage:        form.ServicePackage,
		ProjectGoal:           form.ProjectGoal,
		DesiredTimeline:       form.DesiredTimeline,
		AssetsProvided:        form.AssetsProvided,
		ReadinessStatus:       result.ReadinessStatus,
		MissingItems:          result.MissingItems,
		AISummary:             result.AISummary,
		RecommendedNextAction: result.RecommendedNextAction,
	}
}
