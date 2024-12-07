package models

import (
	"time"

	"github.com/google/uuid"
)

type Submission struct {
	ID         *uuid.UUID       `json:"id,omitempty"`
	ProblemID  uuid.UUID        `json:"problem_id"`
	UserID     uuid.UUID        `json:"user_id"`
	LanguageID int              `json:"language_id"`
	SourceCode string           `json:"source_code"`
	Status     SubmissionStatus `json:"status"`
	Updated_at *time.Time       `json:"updated_at,omitempty"`
	Created_at *time.Time       `json:"created_at,omitempty"`
}

func NewSubmission(problemID uuid.UUID, userID uuid.UUID, languageID int, sourceCode string) *Submission {
	return &Submission{
		ProblemID:  problemID,
		UserID:     userID,
		LanguageID: languageID,
		SourceCode: sourceCode,
		Status:     PEDNING_SUBMISSION_STATUS,
	}
}
