package judge_api

type JudgeAPI struct {
	apiKey string
	host   string
}

func NewJudgeAPI(apiKey string, host string) *JudgeAPI {
	return &JudgeAPI{
		apiKey: apiKey,
		host:   host,
	}
}
