package usecase

import (
	"context"

	"github.com/sorrawichyooboon/go-audit-partition-purger/internal/domain"
)

// AuditRepository กำหนดสัญญาว่า Database Layer ต้องทำอะไรได้บ้าง
type AuditRepository interface {
	Save(ctx context.Context, log *domain.AuditLog) error
	
	// สคริปต์แมนนวลสำหรับบังคับลบข้อมูล หรือบังคับให้ pg_partman ทำงาน
	ManualPurgePartition(ctx context.Context, partitionName string) error
}

// AuditUsecase กำหนดสัญญาว่า Business Logic ของ Audit ต้องมีฟังก์ชันอะไรบ้าง
type AuditUsecase interface {
	TrackAction(ctx context.Context, userID string, action string, details string) error
	ForcePurgeOldData(ctx context.Context, targetMonth string) error
}
