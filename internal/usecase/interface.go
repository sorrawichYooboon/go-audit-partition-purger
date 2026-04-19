package usecase

import (
	"context"

	"github.com/sorrawichyooboon/go-audit-partition-purger/internal/domain"
)

type AuditRepository interface {
	Save(ctx context.Context, log *domain.AuditLog) error
	ManualPurgePartition(ctx context.Context, partitionName string) error
}

type AuditUsecase interface {
	TrackAction(ctx context.Context, userID string, action string, details string) error
	ForcePurgeOldData(ctx context.Context, targetMonth string) error
}
