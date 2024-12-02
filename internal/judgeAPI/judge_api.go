package judgeAPI

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Modalessi/nuha-api/internal/utils"
)

type JudgeAPI struct {
	baseURL *url.URL
	apiKey  string
	host    string
}

func NewJudgeAPI(apiKey string, host string) *JudgeAPI {

	const BASE_URL_STRING = "https://judge0-ce.p.rapidapi.com"
	baseURL, err := url.Parse(BASE_URL_STRING)
	utils.Assert(err, "failed parsing Judge api base url")

	return &JudgeAPI{
		baseURL: baseURL,
		apiKey:  apiKey,
		host:    host,
	}
}

func (j *JudgeAPI) PostSubmission(s *Submission) (string, error) {
	postSubmissionURL := j.baseURL.JoinPath("submissions")

	payload := strings.NewReader(string(s.JSON()))
	req, err := http.NewRequest("POST", postSubmissionURL.String(), payload)
	utils.Assert(err, "error making submission request for judge zero")

	req.Header.Add("x-rapidapi-key", j.apiKey)
	req.Header.Add("x-rapidapi-host", j.host)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("someting went wrong with judge zero sending submission request %v", err)
	}
	defer res.Body.Close()

	type response struct {
		Token string `json:"token"`
	}

	resBody := response{}

	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return "", fmt.Errorf("someting went wrong while decoding judge zero sending submission response %v", err)
	}

	return resBody.Token, nil
}
