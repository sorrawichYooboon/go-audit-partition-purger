package database

import (
	"context"
	"fmt"
	"regexp"

	"github.com/sorrawichyooboon/go-audit-partition-purger/internal/domain"
	"github.com/sorrawichyooboon/go-audit-partition-purger/internal/usecase"
	"gorm.io/gorm"
)

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) usecase.AuditRepository {
	return &auditRepository{db: db}
}

func (r *auditRepository) Save(ctx context.Context, log *domain.AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

var auditPartitionNameRE = regexp.MustCompile(`^audit_logs_p\d{4}_\d{2}$`)

func (r *auditRepository) ManualPurgePartition(ctx context.Context, partitionName string) error {
	if !auditPartitionNameRE.MatchString(partitionName) {
		return fmt.Errorf("invalid partition name: %q", partitionName)
	}

	if err := r.db.WithContext(ctx).
		Exec(fmt.Sprintf("ALTER TABLE audit_logs DETACH PARTITION %s", partitionName)).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).
		Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", partitionName)).Error; err != nil {
		return err
	}

	return nil
}
