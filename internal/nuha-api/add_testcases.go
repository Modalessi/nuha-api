package nuha

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Modalessi/nuha-api/internal/models"
	"github.com/Modalessi/nuha-api/internal/repositories"
)

func addTestCases(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	problemId := r.URL.Query().Get("problem_id")
	if problemId == "" {
		respondWithError(w, 400, INVALID_QUERY_ERROR)
		return fmt.Errorf("error, problem_id query was not provided")
	}

	defer r.Body.Close()

	requestTestcases := []models.Testcase{}
	err := json.NewDecoder(r.Body).Decode(&requestTestcases)
	if err != nil {
		respondWithError(w, 400, INVALID_JSON_ERROR)
		return err
	}

	pr := repositories.NewProblemRepository(ns.S3.Client, ns.DB, ns.DBQueries, r.Context(), ns.S3.BucketName)

	// check if problem exist
	problemDb, err := pr.GetProblemInfo(problemId)
	if err != nil {
		respondWithError(w, 400, err)
		return err
	}

	// store test cases
	err = pr.AddNewTestCases(problemDb.ID.String(), requestTestcases...)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	respondWithSuccess(w, 201, "test cases addes successfuly")
	return nil
}
