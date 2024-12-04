package nuha

import (
	"encoding/json"
	"net/http"

	"github.com/Modalessi/nuha-api/internal"
	"github.com/Modalessi/nuha-api/internal/models"
	"github.com/Modalessi/nuha-api/internal/repositories"
)

func createProblem(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {

	type createProblemSchema struct {
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Difficulty  string   `json:"difficulty"`
		Tags        []string `json:"tags"`
		Timelimit   float64  `json:"timelimit,omitempty"`
		Memorylimit float64  `json:"memorylimit,omitempty"`
	}

	defer r.Body.Close()

	problemData := createProblemSchema{}
	err := json.NewDecoder(r.Body).Decode(&problemData)
	if err != nil {
		respondWithError(w, 400, INVALID_JSON_ERROR)
		return err
	}

	// create problem
	problem, err := models.CreateNewProblem(problemData.Title, problemData.Description, problemData.Difficulty, problemData.Tags)
	if err != nil {
		respondWithError(w, 400, err)
		return err
	}

	problem.AddTags(problemData.Tags)

	if problemData.Timelimit != 0 {
		problem.SetTimelimit(problemData.Timelimit)
	}
	if problemData.Memorylimit != 0 {
		problem.SetMemoryLimit(problem.Memorylimit)
	}

	// store problem
	pr := repositories.NewProblemRepository(ns.S3.Client, ns.DB, ns.DBQueries, r.Context(), ns.S3.BucketName)
	problemDB, err := pr.StoreNewProblem(problem)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	response := struct {
		ID    string
		Title string
	}{
		ID:    problemDB.ID.String(),
		Title: problem.Title,
	}

	respondWithJson(w, 201, &internal.JsonWrapper{Data: response})
	return nil
}
