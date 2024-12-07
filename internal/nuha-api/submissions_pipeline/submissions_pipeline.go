package submissionsPL

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/judgeAPI"
	"github.com/Modalessi/nuha-api/internal/models"
	"github.com/google/uuid"
)

const (
	WAIT_TILL_CHANNEL_BUFFER_OPENS      = 5 // seconds
	SUBMISSIONS_PROCESSORS_COUNT        = 5
	RESULTS_PROCESSORS_COUNT            = 5
	DB_WRITER_COUNT                     = 1
	CHECK_WITH_JUDGEAPI_COUNT           = 5
	CHANNELS_BUFFER                     = 100
	PERIOD_BETWEEN_EACH_JUDGE_API_CHECK = 3 // seconds
)

type SubmissionJob struct {
	SubmissionID uuid.UUID
	Language     judgeAPI.JudgeLanguage
	Code         string
	Timelimit    float64
	MemoryLimit  float64
	ProblemID    uuid.UUID
	Testcases    []models.Testcase
}

type ResultTokens struct {
	SubmissionID uuid.UUID
	Tokens       []string
}

type DBUpdate struct {
	SubmissionID uuid.UUID
	Results      []judgeAPI.Submission
}

type TestCaseResult struct {
	SubmissionID   uuid.UUID
	Token          string
	Status         string
	Stdin          string
	Stdout         string
	ExpectedOutput string
	TimeUsed       float64
	MemoryUsed     float64
	JudgeResponse  []byte
}

type SubmissionsPipeline struct {
	submissionsChan chan *SubmissionJob
	resultsChan     chan *ResultTokens
	dbUpdateChan    chan *DBUpdate
	judgeAPI        *judgeAPI.JudgeAPI
	db              *sql.DB
	dbQueries       *database.Queries
	wg              sync.WaitGroup
	ctx             context.Context
	cancel          context.CancelFunc
}

func NewSubmissionPipeline(judgeAPI *judgeAPI.JudgeAPI, db *sql.DB, dbQueries *database.Queries) *SubmissionsPipeline {
	ctx, cancel := context.WithCancel(context.Background())

	return &SubmissionsPipeline{
		submissionsChan: make(chan *SubmissionJob, CHANNELS_BUFFER),
		resultsChan:     make(chan *ResultTokens, CHANNELS_BUFFER),
		dbUpdateChan:    make(chan *DBUpdate, CHANNELS_BUFFER),
		judgeAPI:        judgeAPI,
		db:              db,
		dbQueries:       dbQueries,
		ctx:             ctx,
		cancel:          cancel,
	}
}

func (sp *SubmissionsPipeline) Start() {
	for range SUBMISSIONS_PROCESSORS_COUNT {
		sp.wg.Add(1)
		go sp.submissionsProcessor()
	}

	for range RESULTS_PROCESSORS_COUNT {
		sp.wg.Add(1)
		go sp.resultsProcessor()
	}

	for range DB_WRITER_COUNT {
		sp.wg.Add(1)
		go sp.databaseUpdater()
	}

}

func (sp *SubmissionsPipeline) Submit(job *SubmissionJob) error {
	select {
	case sp.submissionsChan <- job:
		return nil
	case <-time.After(WAIT_TILL_CHANNEL_BUFFER_OPENS * time.Second):
		return fmt.Errorf("submission pipeline is full")
	}
}

func (sp *SubmissionsPipeline) Shutdown() {
	sp.cancel()
	sp.wg.Wait()
	close(sp.submissionsChan)
	close(sp.resultsChan)
	close(sp.dbUpdateChan)
}

func (sp *SubmissionsPipeline) submissionsProcessor() {
	defer sp.wg.Done()

	for {
		select {

		case job := <-sp.submissionsChan:
			submission := judgeAPI.NewSubmission(job.Code, job.Language)
			submission.SetCPUTimeLimit(job.Timelimit)
			submission.SetMemoryLimit(job.MemoryLimit)
			batch := submission.GenerateBatchFromTestCases(job.Testcases...)

			tokens, err := sp.judgeAPI.PostBatchSubmission(batch)
			if err != nil {
				log.Printf("Error submitting to judge0: %v", err)
				continue
			}

			resultTokens := ResultTokens{
				SubmissionID: job.SubmissionID,
				Tokens:       tokens,
			}

			sp.resultsChan <- &resultTokens

		case <-sp.ctx.Done():
			return
		}
	}

}

func (sp *SubmissionsPipeline) resultsProcessor() {
	defer sp.wg.Done()

	ticker := time.NewTicker(PERIOD_BETWEEN_EACH_JUDGE_API_CHECK * time.Second)
	defer ticker.Stop()

	type pendingSubmission struct {
		submissionID uuid.UUID
		tokens       []string
		checkCount   int
	}
	pendingSubmissions := make(map[uuid.UUID]pendingSubmission)

	for {
		select {
		case result := <-sp.resultsChan:
			pendingSubmissions[result.SubmissionID] = pendingSubmission{
				submissionID: result.SubmissionID,
				tokens:       result.Tokens,
				checkCount:   0,
			}

		case <-ticker.C:
			for id, pending := range pendingSubmissions {
				if pending.checkCount >= CHECK_WITH_JUDGEAPI_COUNT {
					log.Printf("Submission %v timed out after %d checks", id, CHECK_WITH_JUDGEAPI_COUNT)
					delete(pendingSubmissions, id)
					continue
				}

				submissions, err := sp.judgeAPI.GetBatchSubmissionsResult(pending.tokens)
				if err != nil {
					log.Printf("Error getting submissions from judge api: %v", err)
					continue
				}

				allDone := areAllSubmissionsDone(submissions)

				if allDone {
					dpUpdate := DBUpdate{
						SubmissionID: pending.submissionID,
						Results:      submissions,
					}
					sp.dbUpdateChan <- &dpUpdate
					delete(pendingSubmissions, id)
				} else {
					pending.checkCount += 1
					pendingSubmissions[id] = pending
				}
			}

		case <-sp.ctx.Done():
			return
		}
	}
}

func (sp *SubmissionsPipeline) databaseUpdater() {
	defer sp.wg.Done()

	for {
		select {
		case update := <-sp.dbUpdateChan:

			tx, err := sp.db.Begin()
			if err != nil {
				log.Printf("Error starting transaction: %v", err)
				continue
			}

			txq := sp.dbQueries.WithTx(tx)

			status := calculateSubmissionStatus(update.Results)

			tokens, stdins, stdouts, expectedOuts, statuses, times, memories, responses := getResultColumns(update.Results)
			createResultsParams := database.CreateSubmissionResultsParams{
				SubmissionID:    update.SubmissionID,
				Tokens:          tokens,
				Stdins:          stdins,
				Stdouts:         stdouts,
				Expectedoutputs: expectedOuts,
				Statuses:        statuses,
				Times:           times,
				Memories:        memories,
				Responses:       responses,
			}
			_, err = txq.CreateSubmissionResults(sp.ctx, createResultsParams)
			if err != nil {
				log.Printf("error creating submission results: %v", err)
				tx.Rollback()
				continue
			}

			updateSubmissionStatus := database.UpdateSubmissionStatusParams{
				ID:     update.SubmissionID,
				Status: string(status),
			}
			_, err = txq.UpdateSubmissionStatus(sp.ctx, updateSubmissionStatus)
			if err != nil {
				log.Printf("error updating submission status: %v", err)
				tx.Rollback()
				continue
			}

			err = tx.Commit()
			if err != nil {
				log.Printf("errror commitng transaction to add submissions resluts and update status: %v", err)
				tx.Rollback()
			}

		case <-sp.ctx.Done():
			return
		}
	}
}

func getResultColumns(results []judgeAPI.Submission) (
	tokens []string,
	stdins []string,
	stdouts []string,
	expectedOutputs []string,
	statuses []int32,
	timeUsed []string,
	memoryUsed []float64,
	judgeResponses [][]byte,
) {
	n := len(results)
	tokens = make([]string, n)
	stdins = make([]string, n)
	stdouts = make([]string, n)
	expectedOutputs = make([]string, n)
	statuses = make([]int32, n)
	timeUsed = make([]string, n)
	memoryUsed = make([]float64, n)
	judgeResponses = make([][]byte, n)

	for i, result := range results {
		tokens[i] = result.Token
		stdins[i] = result.Stdin
		stdouts[i] = result.Stdout
		expectedOutputs[i] = result.ExpectedOutput
		statuses[i] = int32(result.Status.ID)
		timeUsed[i] = result.Time
		memoryUsed[i] = result.Memory
		judgeResponses[i] = result.JSON()

	}

	return
}

func calculateSubmissionStatus(submissions []judgeAPI.Submission) models.SubmissionStatus {
	if len(submissions) == 0 {
		return models.PEDNING_SUBMISSION_STATUS
	}

	for _, s := range submissions {
		if s.Status.ID == judgeAPI.IN_QUEUE_STATUS || s.Status.ID == judgeAPI.PROCESSING_STATUS {
			return models.PEDNING_SUBMISSION_STATUS
		}
	}

	for _, s := range submissions {

		if s.Status.ID == judgeAPI.COMPILATION_ERROR_STATUS || s.Status.ID == judgeAPI.EXECFORMAT_ERROR_STATUS {
			return models.COMPILATION_ERROR_SUBMISSION_STATUS
		}

		if s.Status.ID == judgeAPI.INTERNAL_ERROR_STATUS {
			return models.SERVER_ERROR_SUBMISSION_STATUS
		}
	}

	for _, s := range submissions {
		if s.Status.ID == judgeAPI.RUNTIME_ERROR_NZEC_STATUS && s.Message == "Exited with error status 137" {
			return models.MEMORY_LIMIT_SUBMISSION_STATUS
		}
	}

	for _, s := range submissions {
		if s.Status.ID == judgeAPI.RUNTIME_ERROR_SIGSEGV_STATUS ||
			s.Status.ID == judgeAPI.RUNTIME_ERROR_SIGXFSZ_STATUS ||
			s.Status.ID == judgeAPI.RUNTIME_ERROR_SIGFPE_STATUS ||
			s.Status.ID == judgeAPI.RUNTIME_ERROR_SIGABRT_STATUS ||
			s.Status.ID == judgeAPI.RUNTIME_ERROR_NZEC_STATUS ||
			s.Status.ID == judgeAPI.RUNTIME_ERROR_OTHER_STATUS {
			return models.RUNTIME_ERROR_SUBMISSION_STATUS
		}
	}

	for _, s := range submissions {
		if s.Status.ID == judgeAPI.TIME_LIMIT_EXCEEDED_STATUS {
			return models.TIME_LIMIT_SUBMISSION_STATUS
		}
	}

	for _, s := range submissions {
		if s.Status.ID == judgeAPI.WRONG_ANSWER_STATUS {
			return models.WRONG_ANSWER_SUBMISSION_STATUS
		}
	}

	for _, s := range submissions {
		if s.Status.ID != judgeAPI.ACCEPTED_STATUS {
			return models.SERVER_ERROR_SUBMISSION_STATUS
		}
	}

	return models.ACCEPTED_SUBMISSION_STATUS
}

func areAllSubmissionsDone(submissions []judgeAPI.Submission) bool {
	submissionsCount := len(submissions)
	doneCount := 0

	for _, s := range submissions {
		if s.Status.ID != judgeAPI.IN_QUEUE_STATUS && s.Status.ID != judgeAPI.PROCESSING_STATUS {
			doneCount += 1
		}
	}

	return submissionsCount == doneCount
}
