package connection

import (
	"context"
	"github.com/khanalsaroj/typegen-server/internal/domain"
	"github.com/khanalsaroj/typegen-server/internal/pkg/crypto"
)

type Service struct {
	repo   Repository
	crypto *crypto.Service
}

func NewService(
	repo Repository,
	cryptoSvc *crypto.Service,
) *Service {
	return &Service{
		repo:   repo,
		crypto: cryptoSvc,
	}
}

func (s *Service) Create(ctx context.Context, req *DatabaseConnectionsRequest) (*domain.DatabaseConnection, error) {
	ecn, err := s.crypto.Encrypt(req.Password)
	if err != nil {
		return nil, err
	}

	connection := &domain.DatabaseConnection{
		Name:         req.Name,
		DbType:       req.DbType,
		Host:         req.Host,
		Port:         req.Port,
		DatabaseName: req.DatabaseName,
		SchemaName:   req.SchemaName,
		Username:     req.Username,
		Password:     ecn,
	}

	if err := s.repo.Create(ctx, connection); err != nil {
		return nil, err
	}

	return connection, nil
}

func (s *Service) GetByID(ctx context.Context, id uint) (*domain.DatabaseConnection, error) {
	conn, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	dec, err := s.crypto.Decrypt(conn.Password)
	if err != nil {
		return nil, err
	}
	conn.Password = dec
	return conn, nil
}

func (s *Service) List(ctx context.Context, page, pageSize int) ([]*domain.DatabaseConnection, int64, error) {
	offset := (page - 1) * pageSize

	conn, val, err := s.repo.FindAll(ctx, offset, pageSize)

	for i := range conn {
		dec, err := s.crypto.Decrypt(conn[i].Password)
		if err != nil {
			return nil, 0, err
		}
		conn[i].Password = dec
	}
	return conn, val, err
}

func (s *Service) Update(ctx context.Context, id uint, req *DatabaseConnectionsRequest) (*domain.DatabaseConnection, error) {
	connection, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Host != "" {
		connection.Host = req.Host
	}
	if req.Port > 0 {
		connection.Port = req.Port
	}
	if req.DatabaseName != "" {
		connection.Name = req.Name
	}

	if req.DbType != "" {
		connection.DbType = req.DbType
	}

	if req.Password != "" {
		connection.Password = req.Password
	}

	if req.DatabaseName != "" {
		connection.Username = req.Username
	}

	if req.SchemaName != "" {
		connection.DatabaseName = req.SchemaName
	}

	if err := s.repo.Update(ctx, connection); err != nil {
		return nil, err
	}

	return connection, nil
}

func (s *Service) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}
