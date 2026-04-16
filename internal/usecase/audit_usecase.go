package usecase

import (
	"context"
	"fmt"

	"github.com/sorrawichyooboon/go-audit-partition-purger/internal/domain"
)

type auditUsecase struct {
	repo AuditRepository
}

func NewAuditUsecase(repo AuditRepository) AuditUsecase {
	return &auditUsecase{
		repo: repo,
	}
}

func (u *auditUsecase) TrackAction(ctx context.Context, userID string, action string, details string) error {
	log := &domain.AuditLog{
		UserID:  userID,
		Action:  action,
		Details: details,
	}

	return u.repo.Save(ctx, log)
}

func (u *auditUsecase) ForcePurgeOldData(ctx context.Context, targetMonth string) error {
	partitionName := fmt.Sprintf("audit_logs_p%s", targetMonth)

	return u.repo.ManualPurgePartition(ctx, partitionName)
}
