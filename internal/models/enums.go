package models

type SubmissionStatus string

const (
	PEDNING_SUBMISSION_STATUS           SubmissionStatus = "PENDING"
	ACCEPTED_SUBMISSION_STATUS          SubmissionStatus = "ACCEPTED"
	WRONG_ANSWER_SUBMISSION_STATUS      SubmissionStatus = "WRONG ANSWER"
	TIME_LIMIT_SUBMISSION_STATUS        SubmissionStatus = "TIME LIMIT EXCEEDED"
	MEMORY_LIMIT_SUBMISSION_STATUS      SubmissionStatus = "MEMORY LIMIT EXCEEDED"
	COMPILATION_ERROR_SUBMISSION_STATUS SubmissionStatus = "COMPILATION ERROR"
	RUNTIME_ERROR_SUBMISSION_STATUS     SubmissionStatus = "RUNTIME ERROR"
	SERVER_ERROR_SUBMISSION_STATUS      SubmissionStatus = "SERVER ERROR"
)
