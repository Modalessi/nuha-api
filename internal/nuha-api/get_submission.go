package nuha

import (
	"fmt"
	"net/http"

	"github.com/Modalessi/nuha-api/internal"
	"github.com/Modalessi/nuha-api/internal/judgeAPI"
	"github.com/google/uuid"
)

func getSubmission(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	submissionId := r.URL.Query().Get("submission_id")
	if submissionId == "" {
		respondWithError(w, 400, INVALID_QUERY_ERROR)
		return fmt.Errorf("error, submission_id query was not provided")
	}

	id, err := uuid.Parse(submissionId)
	if err != nil {
		respondWithError(w, 400, fmt.Errorf("invalid submission id"))
		return fmt.Errorf("invalid submssion id")
	}
	submissionDB, err := ns.DBQueries.GetSubmissionByID(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, EntityDoesNotExistError("PROBLEM"))
		return err
	}

	response := struct {
		ID        uuid.UUID `json:"id"`
		ProblemID uuid.UUID `json:"problem_id"`
		UserID    uuid.UUID `json:"user_id"`
		Status    string    `json:"status"`
		Language  string    `json:"language"`
		Code      string    `json:"code"`
		CreatedAT string    `json:"created_at"`
	}{
		ID:        submissionDB.ID,
		ProblemID: submissionDB.ProblemID,
		UserID:    submissionDB.UserID,
		Status:    submissionDB.Status,
		Language:  judgeAPI.JudgeLanguageDescription[judgeAPI.JudgeLanguage(submissionDB.Language)],
		Code:      submissionDB.SourceCode,
		CreatedAT: submissionDB.CreatedAt.String(),
	}

	respondWithJson(w, 200, &internal.JsonWrapper{Data: response})
	return nil
}
