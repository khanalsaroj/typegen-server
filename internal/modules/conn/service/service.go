package service

import (
	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/modules/conn/repository"
	"time"
)

type Service interface {
	CheckConnection(req domain.DatabaseConnectionInfo) domain.ConnectResponse
	CheckHealth(results []domain.ConnectResponse) ([]domain.DatabaseConnectionHealth, error)
}

type service struct {
	repo repository.Repo
}

func New(repo repository.Repo) Service {
	return &service{repo: repo}
}

func (s *service) CheckConnection(req domain.DatabaseConnectionInfo) domain.ConnectResponse {
	pingMs, tables, sizeMb, err, tableInfo :=
		s.repo.PingDB(req)

	if err != nil {
		return domain.ConnectResponse{
			Message: err.Error(),
			Success: false,
		}
	}

	return domain.ConnectResponse{
		Message:     "Database connected successfully",
		Success:     true,
		PingMs:      pingMs,
		TablesFound: tables,
		SizeMb:      sizeMb,
		TableInfo:   tableInfo,
	}
}

func (s *service) CheckHealth(results []domain.ConnectResponse) ([]domain.DatabaseConnectionHealth, error) {
	health := make([]domain.DatabaseConnectionHealth, 0, len(results))
	now := time.Now().UTC()
	for _, r := range results {

		h := domain.DatabaseConnectionHealth{
			ConnectionID:  r.ConnectionId,
			Status:        "connected",
			LastCheckedAt: &now,
		}
		health = append(health, h)
	}

	return health, nil
}
