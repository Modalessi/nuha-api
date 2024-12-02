package nuha

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Modalessi/nuha-api/internal"
	"github.com/Modalessi/nuha-api/internal/judgeAPI"
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

	submission := judgeAPI.NewSubmission(submissionData.Code, judgeAPI.JudgeLanguage(submissionData.Language))
	token, err := ns.JudgeAPI.PostSubmission(submission)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	response := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	respondWithJson(w, 201, &internal.JsonWrapper{Data: response})
	return nil
}
