package nuha

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Modalessi/nuha-api/internal"
	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/judgeAPI"
	"github.com/Modalessi/nuha-api/internal/models"
	submissionsPL "github.com/Modalessi/nuha-api/internal/nuha-api/submissions_pipeline"
	"github.com/Modalessi/nuha-api/internal/repositories"
	"github.com/google/uuid"
)

func submitSolution(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	problemId := r.URL.Query().Get("problem_id")
	if problemId == "" {
		respondWithError(w, 400, INVALID_QUERY_ERROR)
		return fmt.Errorf("error, problem_id query was not provided")
	}

	id, err := uuid.Parse(problemId)
	if err != nil {
		respondWithError(w, 400, INVALID_ID_ERROR)
		return err
	}

	type submissionSchema struct {
		Language int    `json:"language"`
		Code     string `json:"code"`
	}
	defer r.Body.Close()

	submissionData := submissionSchema{}
	err = json.NewDecoder(r.Body).Decode(&submissionData)
	if err != nil {
		respondWithError(w, 400, INVALID_JSON_ERROR)
		return err
	}

	pr := repositories.NewProblemRepository(ns.DB, ns.DBQueries, r.Context())
	problem, err := pr.GetProblemInfo(id)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, 404, EntityDoesNotExistError("Problem"))
		return err
	}
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	// get user id
	ur := repositories.NewUserRespository(r.Context(), ns.DBQueries)
	userEmail, ok := r.Context().Value(userEmailKey).(string)
	if !ok {
		respondWithError(w, 500, SERVER_ERROR)
		return fmt.Errorf("error getting user email from context")
	}

	user, err := ur.GetUserByEmail(userEmail)
	if err != nil {
		respondWithError(w, 404, EntityDoesNotExistError("USER"))
		return err
	}

	testcases, err := pr.GetTestCases(id)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	// create the submiios here
	submission := models.NewSubmission(problem.ID, user.ID, submissionData.Language, submissionData.Code)

	// store it here
	createSubmissionParams := &database.CreateSubmissionParams{
		ProblemID:  problem.ID,
		UserID:     user.ID,
		Language:   int32(submission.LanguageID),
		SourceCode: submission.SourceCode,
		Status:     string(submission.Status),
	}
	submissionDB, err := ns.DBQueries.CreateSubmission(r.Context(), *createSubmissionParams)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return fmt.Errorf("error creating submission in database: %w", err)
	}

	// give it to submision piplie line here
	submissionJob := &submissionsPL.SubmissionJob{
		SubmissionID: submissionDB.ID,
		Language:     judgeAPI.JudgeLanguage(submissionDB.Language),
		Code:         submissionDB.SourceCode,
		Timelimit:    problem.TimeLimit,
		MemoryLimit:  problem.MemoryLimit,
		ProblemID:    problem.ID,
		Testcases:    models.TestCasesFromDBObjects(testcases),
	}
	ns.SubmissionsPL.Submit(submissionJob)

	response := struct {
		SubmissionID uuid.UUID `json:"submission_id"`
	}{
		SubmissionID: submissionDB.ID,
	}
	respondWithJson(w, 201, &internal.JsonWrapper{Data: response})
	return nil
}
