package nuha

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/Modalessi/nuha-api/internal"
	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/judgeAPI"
	"github.com/google/uuid"
)

func getSubmission(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	submissionId := r.URL.Query().Get("submission_id")
	problemId := r.URL.Query().Get("problem_id")
	userId := r.URL.Query().Get("user_id")

	problemIdGiven := problemId != ""
	userIdGiven := userId != ""
	submissionIdGiven := submissionId != ""

	if !problemIdGiven && !userIdGiven && !submissionIdGiven {
		return getSubmissions(ns, w, r)
	}

	if problemIdGiven && !userIdGiven && !submissionIdGiven {
		return submissionsForProblem(ns, w, r)
	}

	if userIdGiven && !problemIdGiven && !submissionIdGiven {
		return submissionsForUser(ns, w, r)
	}

	if userIdGiven && problemIdGiven {
		return userSubmissionForProblem(ns, w, r)
	}

	id, err := uuid.Parse(submissionId)
	if err != nil {
		respondWithError(w, 400, fmt.Errorf("invalid submission id"))
		return err
	}
	submissionDB, err := ns.DBQueries.GetSubmissionByID(r.Context(), id)
	if err != nil {
		respondWithError(w, 404, EntityDoesNotExistError("SUBMISSION"))
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

func submissionsForProblem(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	problemID := r.URL.Query().Get("problem_id")
	if problemID == "" {
		respondWithError(w, 400, INVALID_QUERY_ERROR)
		return fmt.Errorf("error, problem_id query was not provided")
	}

	id, err := uuid.Parse(problemID)
	if err != nil {
		respondWithError(w, 400, fmt.Errorf("invalid problem id"))
		return fmt.Errorf("invalid problem id: %w", err)
	}

	_, err = ns.DBQueries.GetProblemByID(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, 404, EntityDoesNotExistError("PROBLEM"))
		return err
	}

	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	getSubmissionsParams := database.GetSubmissionsByProblemIDParams{
		ProblemID: id,
		Offset:    0, // TODO: add pagination
		Limit:     20,
	}
	submissionsDB, err := ns.DBQueries.GetSubmissionsByProblemID(r.Context(), getSubmissionsParams)

	if errors.Is(err, sql.ErrNoRows) {
		respondWithJson(w, 200, &internal.JsonWrapper{Data: []string{}})
		return nil
	}

	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	respondWithSubmissoins(w, submissionsDB)
	return nil
}

func submissionsForUser(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	userId := r.URL.Query().Get("user_id")
	if userId == "" {
		respondWithError(w, 400, INVALID_QUERY_ERROR)
		return fmt.Errorf("error user_id was not provided ")
	}

	id, err := uuid.Parse(userId)
	if err != nil {
		respondWithError(w, 400, INVALID_ID_ERROR)
		return err
	}

	_, err = ns.DBQueries.GetUserByID(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, 404, EntityDoesNotExistError("USER"))
		return err
	}
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	submissionsByUserIDParamas := database.GetSubmissionsByUserIDParams{
		UserID: id,
		Offset: 0, // TODO: pagination
		Limit:  20,
	}
	submissionsDB, err := ns.DBQueries.GetSubmissionsByUserID(r.Context(), submissionsByUserIDParamas)

	if errors.Is(err, sql.ErrNoRows) {
		respondWithJson(w, 200, &internal.JsonWrapper{Data: []string{}})
		return nil
	}

	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	respondWithSubmissoins(w, submissionsDB)
	return nil
}

func userSubmissionForProblem(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	userIdQuery := r.URL.Query().Get("user_id")
	if userIdQuery == "" {
		respondWithError(w, 400, INVALID_QUERY_ERROR)
		return fmt.Errorf("error user_id was not provided ")
	}

	userId, err := uuid.Parse(userIdQuery)
	if err != nil {
		respondWithError(w, 400, INVALID_ID_ERROR)
		return err
	}

	problemIDQuery := r.URL.Query().Get("problem_id")
	if problemIDQuery == "" {
		respondWithError(w, 400, INVALID_QUERY_ERROR)
		return fmt.Errorf("error, problem_id query was not provided")
	}

	problemId, err := uuid.Parse(problemIDQuery)
	if err != nil {
		respondWithError(w, 400, fmt.Errorf("invalid problem id"))
		return fmt.Errorf("invalid problem id: %w", err)
	}

	_, err = ns.DBQueries.GetUserByID(r.Context(), userId)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, 404, EntityDoesNotExistError("USER"))
		return err
	}
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	_, err = ns.DBQueries.GetProblemByID(r.Context(), problemId)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, 404, EntityDoesNotExistError("PROBLEM"))
		return err
	}

	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	userSubmissionsForProblemParams := database.GetUserSubmissionsForProblemParams{
		UserID:    userId,
		ProblemID: problemId,
		Offset:    0, // TODO: pagination
		Limit:     20,
	}

	submissionsDB, err := ns.DBQueries.GetUserSubmissionsForProblem(r.Context(), userSubmissionsForProblemParams)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithJson(w, 200, &internal.JsonWrapper{Data: []string{}})
		return nil
	}

	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	respondWithSubmissoins(w, submissionsDB)
	return nil
}

func getSubmissions(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	getSubmissionsParams := database.GetSubmissionsParams{
		Offset: 0,
		Limit:  20, // TODO: Pagination
	}
	submissionsDB, err := ns.DBQueries.GetSubmissions(r.Context(), getSubmissionsParams)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithJson(w, 200, &internal.JsonWrapper{Data: []string{}})
		return nil
	}

	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	respondWithSubmissoins(w, submissionsDB)
	return nil
}

func respondWithSubmissoins(w http.ResponseWriter, submissions []database.Submission) {
	type SubmissionDetails struct {
		ID        uuid.UUID `json:"id"`
		ProblemID uuid.UUID `json:"problem_id"`
		UserID    uuid.UUID `json:"user_id"`
		Status    string    `json:"status"`
		Language  string    `json:"language"`
		Code      string    `json:"code"`
		CreatedAT string    `json:"created_at"`
	}

	response := make([]SubmissionDetails, 0)
	for _, submissionDB := range submissions {
		response = append(response, SubmissionDetails{
			ID:        submissionDB.ID,
			ProblemID: submissionDB.ProblemID,
			UserID:    submissionDB.UserID,
			Status:    submissionDB.Status,
			Language:  judgeAPI.JudgeLanguageDescription[judgeAPI.JudgeLanguage(submissionDB.Language)],
			Code:      submissionDB.SourceCode,
			CreatedAT: submissionDB.CreatedAt.String(),
		})
	}

	respondWithJson(w, 200, &internal.JsonWrapper{Data: response})
}
