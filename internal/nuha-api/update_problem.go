package nuha

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Modalessi/nuha-api/internal/models"
	"github.com/Modalessi/nuha-api/internal/repositories"
)

func updateProblem(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {

	problemId := r.URL.Query().Get("problem_id")
	if problemId == "" {
		respondWithError(w, 400, INVALID_QUERY_ERROR)
		return fmt.Errorf("error, problem_id query was not provided")
	}

	type updateProblemSchema struct {
		Title       *string  `json:"title,omitempty"`
		Description *string  `json:"description,omitempty"`
		Difficulty  *string  `json:"difficulty,omitempty"`
		Tags        []string `json:"tags,omitempty"`
		TimeLimit   *float64 `json:"time_limit,omitempty"`
		MemoryLimit *float64 `json:"memory_limit,omitempty"`
	}

	defer r.Body.Close()
	updateData := updateProblemSchema{}
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		respondWithError(w, 400, INVALID_JSON_ERROR)
		return err
	}

	// Check if at least one field is provided
	if updateData.Title == nil &&
		updateData.Description == nil &&
		updateData.Difficulty == nil &&
		len(updateData.Tags) == 0 &&
		updateData.TimeLimit == nil &&
		updateData.MemoryLimit == nil {
		respondWithError(w, 400, INVALID_JSON_ERROR)
		return fmt.Errorf("at least one field must be provided for update")
	}

	// Validate difficulty if provided

	pr := repositories.NewProblemRepository(ns.S3.Client, ns.DB, ns.DBQueries, r.Context(), ns.S3.BucketName)
	problemDB, err := pr.GetProblemInfo(problemId)
	if err != nil {
		respondWithError(w, 404, EntityDoesNotExistError("Problem"))
		return err
	}

	description, err := pr.GetProblemDescription(problemId)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	problem := models.ProblemFromDBObject(problemDB)
	problem.Description = description

	if updateData.Title != nil {
		problem.SetTitle(*updateData.Title)
	}
	if updateData.Difficulty != nil {
		err := problem.SetDifficulty(*updateData.Difficulty)
		if err != nil {
			respondWithError(w, 400, err)
			return err
		}
	}
	if updateData.Description != nil {
		problem.SetDescription(*updateData.Description)
	}
	if updateData.Difficulty != nil {
		problem.Difficulty = *updateData.Difficulty
	}
	if len(updateData.Tags) > 0 {
		problem.Tags = updateData.Tags
	}
	if updateData.TimeLimit != nil {
		problem.SetTimelimit(*updateData.TimeLimit)
	}
	if updateData.MemoryLimit != nil {
		problem.SetMemoryLimit(*updateData.MemoryLimit)
	}

	err = pr.UpdateProblem(problem)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	respondWithSuccess(w, 200, fmt.Sprintf("problem with id %s has been updated successfully", problemId))
	return nil
}
