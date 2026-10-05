package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type MoodleAccessClient struct {
	url        string
	token      string
	httpClient *http.Client
}

func NewMoodleAccessClient(cfg Config) *MoodleAccessClient {
	return &MoodleAccessClient{
		url:   cfg.MoodleWebServiceURL,
		token: cfg.MoodleWebServiceToken,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *MoodleAccessClient) Capabilities(ctx context.Context, dodID string) ([]string, error) {
	form := url.Values{
		"wstoken":            {c.token},
		"wsfunction":         {"local_otasignconnector_get_access"},
		"moodlewsrestformat": {"json"},
		"dodid":              {dodID},
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.url, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("Moodle authorization service returned %s", response.Status)
	}

	var result struct {
		Capabilities []string `json:"capabilities"`
		Exception    string   `json:"exception"`
		ErrorCode    string   `json:"errorcode"`
		Message      string   `json:"message"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, err
	}
	if result.Exception != "" || result.ErrorCode != "" {
		return nil, fmt.Errorf("Moodle authorization service rejected request: %s", result.Message)
	}

	allowed := map[string]bool{
		"viewown": true, "viewunit": true, "signascommander": true, "configure": true,
	}
	capabilities := make([]string, 0, len(result.Capabilities))
	for _, capability := range result.Capabilities {
		capability = strings.ToLower(strings.TrimSpace(capability))
		if allowed[capability] {
			capabilities = append(capabilities, capability)
		}
	}
	if !hasAny(capabilities, "viewown") {
		return nil, errors.New("Moodle does not grant OTA Sign access")
	}
	return capabilities, nil
}
