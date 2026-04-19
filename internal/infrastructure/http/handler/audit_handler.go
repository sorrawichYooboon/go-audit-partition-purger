package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sorrawichyooboon/go-audit-partition-purger/internal/dto"
	"github.com/sorrawichyooboon/go-audit-partition-purger/internal/usecase"
)

type AuditHandler struct {
	usecase usecase.AuditUsecase
}

func NewAuditHandler(u usecase.AuditUsecase) *AuditHandler {
	return &AuditHandler{usecase: u}
}

// TrackAuditLog godoc
// @Summary      Create a new audit log
// @Description  Stores audit log into partitioned table via pg_partman
// @Tags         audit
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateAuditLogRequest true "Audit Log Data"
// @Success      201  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/audit-logs [post]
func (h *AuditHandler) TrackAuditLog(c *gin.Context) {
	var req dto.CreateAuditLogRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.TrackAction(c.Request.Context(), req.UserID, req.Action, req.Details); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Audit log saved successfully"})
}
