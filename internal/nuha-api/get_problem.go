package nuha

import (
	"net/http"

	"github.com/Modalessi/nuha-api/internal"
	"github.com/Modalessi/nuha-api/internal/repositories"
)

func getProblem(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	problemId := r.URL.Query().Get("problem_id")
	if problemId == "" {
		return respondWithProblemList(ns, w, r)
	}

	pr := repositories.NewProblemRepository(ns.S3.Client, ns.DB, ns.DBQueries, r.Context(), ns.S3.BucketName)

	problemDB, err := pr.GetProblemInfo(problemId)
	if err != nil {
		respondWithError(w, 404, EntityDoesNotExistError("Problem"))
		return err
	}

	problemDescription, err := pr.GetProblemDescription(problemId)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	type responeProblem struct {
		Id          string   `json:"id"`
		Title       string   `json:"title"`
		Difficulty  string   `json:"difficulty"`
		Discription string   `json:"discription"`
		Tags        []string `json:"tags"`
		TimeLimit   float64  `json:"time_limit"`
		MemoryLimit float64  `json:"memory_limit"`
	}

	response := responeProblem{
		Id:          problemDB.ID.String(),
		Title:       problemDB.Title,
		Difficulty:  problemDB.Difficulty,
		Discription: problemDescription,
		Tags:        problemDB.Tags,
		TimeLimit:   problemDB.TimeLimit,
		MemoryLimit: problemDB.MemoryLimit,
	}

	respondWithJson(w, 200, &internal.JsonWrapper{Data: response})
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
		Id          string   `json:"id"`
		Title       string   `json:"title"`
		Difficulty  string   `json:"difficulty"`
		Tags        []string `json:"tags"`
		TimeLimit   float64  `json:"time_limit"`
		MemoryLimit float64  `json:"memory_limit"`
	}

	responseProblems := []responseProblem{}

	for _, p := range problemsDB {
		rp := responseProblem{
			Id:          p.ID.String(),
			Title:       p.Title,
			Difficulty:  p.Difficulty,
			Tags:        p.Tags,
			TimeLimit:   p.TimeLimit,
			MemoryLimit: p.MemoryLimit,
		}
		responseProblems = append(responseProblems, rp)
	}

	respondWithJson(w, 200, &internal.JsonWrapper{Data: responseProblems})
	return nil
}
