package reporting

import (
	"context"
	"time"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetDailyTill(ctx context.Context, dateStr string) (*DailyTillReport, error) {
	// Parse date, default to today if empty
	var date time.Time
	var err error
	if dateStr == "" {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, err
		}
	}
	return s.repo.GetDailyTill(ctx, date)
}

func (s *Service) GetSalesSummary(ctx context.Context, startStr, endStr string) (*SalesSummaryReport, error) {
	now := time.Now()
	start := now.AddDate(0, 0, -30) // Default last 30 days
	end := now

	var err error
	if startStr != "" {
		start, err = time.Parse("2006-01-02", startStr)
		if err != nil {
			return nil, err
		}
	}
	if endStr != "" {
		end, err = time.Parse("2006-01-02", endStr)
		if err != nil {
			return nil, err
		}
		// Move end to end of day
		end = end.Add(24*time.Hour - time.Nanosecond)
	}

	return s.repo.GetSalesSummary(ctx, start, end)
}
