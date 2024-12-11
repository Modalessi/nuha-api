package nuha

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/Modalessi/nuha-api/internal/repositories"
	"github.com/google/uuid"
)

func deleteProblem(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
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
	defer r.Body.Close()

	pr := repositories.NewProblemRepository(ns.DB, ns.DBQueries, r.Context())
	deletedProblem, err := pr.DeleteProblem(id)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, 404, EntityDoesNotExistError("Problem"))
		return err
	}

	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	responseMSG := fmt.Sprintf("problem with id %s deleted successfully", deletedProblem.ID.String())

	respondWithSuccess(w, 200, responseMSG)
	return nil
}
