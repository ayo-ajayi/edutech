package utils

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ApiResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func NewSuccessResponse(data interface{}, message string) ApiResponse {
	return ApiResponse{
		Success: true,
		Message: message,
		Data:    data,
	}
}

func NewErrorResponse(err interface{}, message string) ApiResponse {
	return ApiResponse{
		Success: false,
		Message: message,
		Error:   err,
	}
}

type StudentRegisteredTutorRes struct {
	TutorId      primitive.ObjectID `json:"tutor_id"`
	Email        string             `json:"email"`
	FirstName    string             `json:"first_name"`
	LastName     string             `json:"last_name"`
	SubjectId    primitive.ObjectID `json:"subject_id"`
	RegisteredAt time.Time          `json:"registered_at"`
}
