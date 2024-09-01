package weekly_report_app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"wbapp/internal"
	"wbapp/internal/config"
	"wbapp/internal/entities"
)

func Run() error {
	file, err := os.Open("/Users/aizitm02/Education/wbapp/configs/config.yaml")
	if err != nil {
		return err
	}
	defer file.Close()

	config, err := config.ParseConfig(file)
	if err != nil {
		return err
	}

	wr := NewWeeklyReport(config.WeeklyReportConfig)
	fmt.Println(wr.GetListOfCategories())
	return nil
}

type WeeklyReport struct {
	scanner *internal.WbScanner
	methods []config.Method

	authToken string
}

func (wr *WeeklyReport) GetListOfCategories() ([]entities.Category, error) {
	// Prepare request parameters
	route := wr.methods[0].Route
	method := "GET"
	var payload bytes.Buffer
	queryParams := url.Values{}
	queryParams.Add("locale", "ru")
	headers := map[string]string{
		"Authorization": wr.authToken,
	}

	responseBody, err := wr.scanner.Request(route, method, payload, queryParams, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}

	var categoryResponse entities.CategoriesResponse
	err = json.Unmarshal(responseBody, &categoryResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal categories: %w", err)
	}

	return categoryResponse.Data.Categories, nil
}

func NewWeeklyReport(config config.WeeklyReportConfig) *WeeklyReport {
	wb := internal.NewWbScanner(config.Host, nil)
	token, exists := os.LookupEnv("AUTH_TOKEN")
	if !exists {
		fmt.Println("Environment variable AUTH_TOKEN is not set")
		return nil
	}
	return &WeeklyReport{
		wb, config.Methods, token,
	}
}
