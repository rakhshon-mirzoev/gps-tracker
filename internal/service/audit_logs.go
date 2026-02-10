package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/model"
	"github.com/rakhshon-mirzoev/gps-tracker/internal/repository"
)

type AuditLogs interface {
	Create(ctx context.Context, al *model.AuditLogs) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.AuditLogs, error)
	List(ctx context.Context, actorID uuid.UUID, action, targetTable, dateFrom, dateTo string, limit, offset int) ([]model.AuditLogs, int, error)
	Update(ctx context.Context, id uuid.UUID, dto model.UpdateAuditLogs) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type auditLogsService struct {
	repo repository.AuditLogsRepository
}

func NewAuditLogsService(repo repository.AuditLogsRepository) AuditLogs {
	return &auditLogsService{repo: repo}
}

func (s *auditLogsService) Create(ctx context.Context, al *model.AuditLogs) error {
	return s.repo.Create(ctx, al)
}

func (s *auditLogsService) GetByID(ctx context.Context, id uuid.UUID) (*model.AuditLogs, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *auditLogsService) List(ctx context.Context, actorID uuid.UUID, action, targetTable, dateFrom, dateTo string, limit, offset int) ([]model.AuditLogs, int, error) {
	return s.repo.List(ctx, actorID, action, targetTable, dateFrom, dateTo, limit, offset)
}

func (s *auditLogsService) Update(ctx context.Context, id uuid.UUID, dto model.UpdateAuditLogs) error {
	return s.repo.Update(ctx, id, dto)
}

func (s *auditLogsService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}
