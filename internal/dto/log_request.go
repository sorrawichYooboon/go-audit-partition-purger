package dto

type CreateAuditLogRequest struct {
	UserID  string `json:"user_id" binding:"required" example:"USR-001"`
	Action  string `json:"action" binding:"required" example:"LOGIN_SUCCESS"`
	Details string `json:"details" binding:"required" example:"{\"ip\":\"192.168.1.1\"}"`
}
