package judgeAPI

import (
	"bytes"
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

	payload := bytes.NewReader(s.JSON())
	req, err := http.NewRequest("POST", postSubmissionURL.String(), payload)
	if err != nil {
		return "", fmt.Errorf("error making submission request for judge zero: %w", err)
	}

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

func (j *JudgeAPI) PostBatchSubmission(bs *SubmissionBatch) ([]string, error) {
	postBatchSubmissionURL := j.baseURL.JoinPath("submissions/batch")

	payload := bytes.NewReader(bs.JSON())
	req, err := http.NewRequest("POST", postBatchSubmissionURL.String(), payload)
	if err != nil {
		return nil, fmt.Errorf("error making batch submission request for judge zero: %w", err)
	}

	req.Header.Add("x-rapidapi-key", j.apiKey)
	req.Header.Add("x-rapidapi-host", j.host)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("someting went wrong with judge zero sending batch submission request %w", err)
	}
	defer res.Body.Close()

	type response struct {
		Token string `json:"token"`
	}

	resBody := []response{}
	err = json.NewDecoder(res.Body).Decode(&resBody)
	if err != nil {
		return nil, fmt.Errorf("someting went wrong while decoding judge zero sending batch submission response %w", err)
	}

	tokens := make([]string, 0, len(resBody))
	for _, t := range resBody {
		tokens = append(tokens, t.Token)
	}

	return tokens, nil
}

func (j *JudgeAPI) GetBatchSubmissionsResult(tokens []string) ([]Submission, error) {
	tokensQuery := strings.Join(tokens, ",")
	postBatchSubmissionURL := j.baseURL.JoinPath("submissions/batch")
	query := postBatchSubmissionURL.Query()
	query.Add("tokens", tokensQuery)
	query.Add("base64_encoded", "false")
	query.Add("fields", "*")
	postBatchSubmissionURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", postBatchSubmissionURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error making get batch submission request for judge zero: %w", err)
	}

	req.Header.Add("x-rapidapi-key", j.apiKey)
	req.Header.Add("x-rapidapi-host", j.host)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("someting went wrong with judge zero sending get batch submission request %w", err)
	}
	defer res.Body.Close()

	resultData := struct {
		Submissions []Submission `json:"submissions"`
	}{}

	err = json.NewDecoder(res.Body).Decode(&resultData)
	if err != nil {
		return nil, fmt.Errorf("someting went wrong while decoding judge zero get batch submission response %w", err)
	}

	return resultData.Submissions, nil
}
