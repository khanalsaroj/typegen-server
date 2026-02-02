package health

import (
	"context"
	"github.com/khanalsaroj/typegen-server/internal/domain"
	"time"

	"gorm.io/gorm"
)

type Service struct {
	db            *gorm.DB
	startTime     time.Time
	version       string
	degradedAfter time.Duration
}

func NewService(db *gorm.DB, version string) *Service {
	return &Service{
		db:            db,
		startTime:     time.Now(),
		version:       version,
		degradedAfter: 500 * time.Millisecond,
	}
}

func (s *Service) CheckDatabase(ctx context.Context) domain.DatabaseHealth {
	start := time.Now()

	sqlDB, err := s.db.DB()
	if err != nil {
		return domain.DatabaseHealth{
			Connected: false,
			LatencyMS: 0,
		}
	}

	err = sqlDB.PingContext(ctx)
	latency := time.Since(start).Milliseconds()

	if err != nil {
		return domain.DatabaseHealth{
			Connected: false,
			LatencyMS: latency,
		}
	}

	return domain.DatabaseHealth{
		Connected: true,
		LatencyMS: latency,
	}
}

func (s *Service) DeriveStatus(db domain.DatabaseHealth) domain.HealthStatus {
	if !db.Connected {
		return domain.Unhealthy
	}

	if time.Duration(db.LatencyMS)*time.Millisecond > s.degradedAfter {
		return domain.Degraded
	}

	return domain.Healthy
}

func (s *Service) GetHealth(ctx context.Context) domain.HealthResponse {
	dbHealth := s.CheckDatabase(ctx)

	return domain.HealthResponse{
		Status:   s.DeriveStatus(dbHealth),
		Uptime:   int64(time.Since(s.startTime).Seconds()),
		Version:  s.version,
		Database: dbHealth,
	}
}
