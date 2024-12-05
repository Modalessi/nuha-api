package nuha

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Modalessi/nuha-api/internal"
	"github.com/Modalessi/nuha-api/internal/judgeAPI"
	"github.com/Modalessi/nuha-api/internal/repositories"
)

func submitSolution(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	problemId := r.URL.Query().Get("problem_id")
	if problemId == "" {
		respondWithError(w, 400, INVALID_QUERY_ERROR)
		return fmt.Errorf("error, problem_id query was not provided")
	}

	type submissionSchema struct {
		Language int    `json:"language"`
		Code     string `json:"code"`
	}
	defer r.Body.Close()

	submissionData := submissionSchema{}
	err := json.NewDecoder(r.Body).Decode(&submissionData)
	if err != nil {
		respondWithError(w, 400, INVALID_JSON_ERROR)
		return err
	}

	pr := repositories.NewProblemRepository(ns.S3.Client, ns.DB, ns.DBQueries, r.Context(), ns.S3.BucketName)
	problem, err := pr.GetProblemInfo(problemId)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, 404, EntityDoesNotExistError("Problem"))
		return err
	}
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	testcases, err := pr.GetTestCases(problemId)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	submission := judgeAPI.NewSubmission(submissionData.Code, judgeAPI.JudgeLanguage(submissionData.Language))
	submission.SetCPUTimeLimit(problem.TimeLimit)
	submission.SetMemoryLimit(problem.MemoryLimit)

	submissionBatch := submission.GenerateBatchFromTestCases(testcases...)

	tokens, err := ns.JudgeAPI.PostBatchSubmission(submissionBatch)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	respondWithJson(w, 201, &internal.JsonWrapper{Data: tokens})
	return nil
}
