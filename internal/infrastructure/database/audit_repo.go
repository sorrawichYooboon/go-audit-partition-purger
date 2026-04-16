package database

import (
	"context"
	"fmt"

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

func (r *auditRepository) ManualPurgePartition(ctx context.Context, partitionName string) error {

	// THE TRAP (INTENTIONAL SQL INJECTION)
	// โค้ดส่วนนี้จงใจเจตนาต่อ String ตรงๆ ซึ่งมีความเสี่ยงสูงที่จะเกิด SQL Injection ได้

	query := fmt.Sprintf("ALTER TABLE %s DETACH PARTITION; DROP TABLE IF EXISTS %s;", partitionName, partitionName)

	// สั่ง Execute ผ่านการต่อ String ห้วน
	if err := r.db.WithContext(ctx).Exec(query).Error; err != nil {
		return err
	}

	return nil
}
