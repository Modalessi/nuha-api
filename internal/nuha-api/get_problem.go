package nuha

import (
	"fmt"
	"net/http"

	"github.com/Modalessi/nuha-api/internal"
	"github.com/Modalessi/nuha-api/internal/repositories"
)

func getProblem(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	problemId := r.URL.Query().Get("problem_id")
	if problemId == "" {
		return respondWithProblemList(ns, w, r)
	}

	respondWithError(w, 200, fmt.Errorf("we still did not implement getting one problem"))
	return nil
}

func respondWithProblemList(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {

	// TODO: Pagination

	pr := repositories.NewProblemRepository(ns.S3.Client, ns.DB, ns.DBQueries, r.Context(), ns.S3.BucketName)

	problemsDB, err := pr.GetProblems(0, 100)
	if err != nil {
		respondWithError(w, 500, err)
		return err
	}

	type responseProblem struct {
		Id         string   `json:"id"`
		Title      string   `json:"title"`
		Difficulty string   `json:"difficulty"`
		Tags       []string `json:"tags"`
	}

	responseProblems := []responseProblem{}

	for _, p := range problemsDB {
		rp := responseProblem{
			Id:         p.ID.String(),
			Title:      p.Title,
			Difficulty: p.Difficulty,
			Tags:       p.Tags,
		}
		responseProblems = append(responseProblems, rp)
	}

	respondWithJson(w, 200, &internal.JsonWrapper{Data: responseProblems})
	return nil
}
