package collegeapi

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
	"strconv"
	"encoding/json"

	"github.com/mrDisa/Raspy/backend/internal/model"
)

type Client struct {
	httpClient *http.Client
	baseURL string
}

func NewClient(baseURL string) *Client {
    return &Client{
        httpClient: &http.Client {
			Timeout: 5 * time.Second,
		},
        baseURL: baseURL,
    }
}

func (c *Client) GetSchedule(group string, startDate *time.Time) ([]model.ScheduleDay, error) {
	params := url.Values{}
	params.Set("group", group)
	
	if startDate != nil {
		timestamp := startDate.Unix()
		params.Set("start_date", strconv.FormatInt(timestamp, 10))
	}
	queryset := params.Encode()
	fullURL := c.baseURL + "get?" + queryset

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cannot get response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api returned status: %d", resp.StatusCode)
	}

	fmt.Printf("URL: %s, Status response: %d", fullURL, resp.StatusCode)

	var schedule []model.ScheduleDay
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&schedule); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return schedule, nil
}