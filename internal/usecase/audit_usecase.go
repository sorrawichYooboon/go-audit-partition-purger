package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"

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
	if details != "" && !json.Valid([]byte(details)) {
		return errors.New("invalid json format for details field")
	}

	log := &domain.AuditLog{
		UserID:  userID,
		Action:  action,
		Details: details,
	}

	return u.repo.Save(ctx, log)
}

func (u *auditUsecase) ForcePurgeOldData(ctx context.Context, targetMonth string) error {
	monthRegex := regexp.MustCompile(`^\d{4}_\d{2}$`)
	if !monthRegex.MatchString(targetMonth) {
		return errors.New("invalid targetMonth format. expected YYYY_MM")
	}

	partitionName := fmt.Sprintf("audit_logs_p%s", targetMonth)

	return u.repo.ManualPurgePartition(ctx, partitionName)
}
