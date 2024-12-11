package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/models"
	"github.com/google/uuid"
)

type ProblemRepository struct {
	db        *sql.DB
	dbQueries *database.Queries
	ctx       context.Context
}

func NewProblemRepository(db *sql.DB, dbQueries *database.Queries, ctx context.Context) *ProblemRepository {
	return &ProblemRepository{
		db:        db,
		dbQueries: dbQueries,
		ctx:       ctx,
	}
}

func (pr *ProblemRepository) StoreNewProblem(p *models.Problem) (*database.Problem, error) {

	newProblemParams := database.CreateProblemParams{
		Title:       p.Title,
		Difficulty:  p.Difficulty,
		Tags:        p.Tags,
		TimeLimit:   p.Timelimit,
		MemoryLimit: p.Memorylimit,
	}

	tx, err := pr.db.BeginTx(pr.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("starting transaction: %w", err)
	}
	defer tx.Rollback()

	qtx := pr.dbQueries.WithTx(tx)

	dbProblem, err := qtx.CreateProblem(pr.ctx, newProblemParams)
	if err != nil {
		return nil, err
	}

	addDescriptionParams := database.AddProblemDescriptionParams{
		ProblemID:   dbProblem.ID,
		Description: p.Description,
	}
	_, err = qtx.AddProblemDescription(pr.ctx, addDescriptionParams)
	if err != nil {
		return nil, err
	}

	testCasesStdins := make([]string, len(p.Testcases))
	testCasesExpectedOutputs := make([]string, len(p.Testcases))

	for i := range p.Testcases {
		testCasesStdins[i] = p.Testcases[i].Stdin
		testCasesExpectedOutputs[i] = p.Testcases[i].ExpectedOutput
	}

	addTestCasesParams := database.CreateTestCasesParams{
		ProblemID:       dbProblem.ID,
		Stdins:          testCasesStdins,
		ExpectedOutputs: testCasesExpectedOutputs,
	}
	_, err = qtx.CreateTestCases(pr.ctx, addTestCasesParams)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return &dbProblem, nil

}

func (pr *ProblemRepository) GetProblemInfo(problemID uuid.UUID) (*database.Problem, error) {

	problem, err := pr.dbQueries.GetProblemByID(pr.ctx, problemID)

	if err != nil {
		return nil, fmt.Errorf("database error checking problem %s: %w", problemID, err)
	}

	return &problem, nil
}

func (pr *ProblemRepository) GetProblemDescription(problemID uuid.UUID) (string, error) {

	dbDescription, err := pr.dbQueries.GetProblemDescription(pr.ctx, problemID)
	if err != nil {
		return "", err
	}

	return dbDescription.Description, nil
}

func (pr *ProblemRepository) GetProblems(offset int32, limit int32) ([]database.Problem, error) {
	getProblemsParams := database.GetProblemsParams{
		Offset: 0,
		Limit:  20,
	}
	problems, err := pr.dbQueries.GetProblems(pr.ctx, getProblemsParams)
	if err != nil {
		return nil, fmt.Errorf("database error getting problems: %w", err)
	}

	return problems, nil
}

func (pr *ProblemRepository) GetTestCases(problemId uuid.UUID) ([]database.TestCase, error) {

	dbTestCases, err := pr.dbQueries.GetTestCases(pr.ctx, problemId)
	if err != nil {
		return nil, err
	}
	return dbTestCases, nil
}

func (pr *ProblemRepository) AddNewTestCases(problemId uuid.UUID, testcases ...models.Testcase) error {

	testCasesStdins := make([]string, len(testcases))
	testCasesExpectedOutputs := make([]string, len(testcases))

	for i := range testcases {
		testCasesStdins[i] = testcases[i].Stdin
		testCasesExpectedOutputs[i] = testcases[i].ExpectedOutput
	}

	addTestCasesParams := database.CreateTestCasesParams{
		ProblemID:       problemId,
		Stdins:          testCasesStdins,
		ExpectedOutputs: testCasesExpectedOutputs,
	}
	_, err := pr.dbQueries.CreateTestCases(pr.ctx, addTestCasesParams)

	if err != nil {
		return err
	}

	return nil
}

func (pr *ProblemRepository) DeleteProblem(problemId uuid.UUID) (*database.Problem, error) {

	tx, err := pr.db.BeginTx(pr.ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}

	defer tx.Rollback()

	txq := pr.dbQueries.WithTx(tx)

	// delete description
	_, err = txq.DeleteProblemDescription(pr.ctx, problemId)
	if err != nil {
		return nil, err
	}

	// delete testcases
	_, err = txq.DeleteTestCases(pr.ctx, problemId)
	if err != nil {
		return nil, err
	}

	// delete problem
	deletedProblem, err := txq.DeleteProblem(pr.ctx, problemId)
	if err != nil {
		return nil, fmt.Errorf("error deleting problem from db with id %q: %w", problemId, err)
	}

	return &deletedProblem, nil
}

func (pr *ProblemRepository) UpdateProblem(problem *models.Problem) error {

	tx, err := pr.db.BeginTx(pr.ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	txq := pr.dbQueries.WithTx(tx)

	updateProblemParams := database.UpdateProblemParams{
		ID:          *problem.ID,
		Title:       problem.Title,
		Difficulty:  problem.Difficulty,
		Tags:        problem.Tags,
		TimeLimit:   problem.Timelimit,
		MemoryLimit: problem.Memorylimit,
	}
	_, err = txq.UpdateProblem(pr.ctx, updateProblemParams)
	if err != nil {
		return fmt.Errorf("error updating problem in database: %w", err)
	}

	if problem.ID == nil {
		return fmt.Errorf("problem has nil id, unacceptable")
	}

	updateProblemDescParams := database.UpdateProblemDescriptionParams{
		ProblemID:   *problem.ID,
		Description: problem.Description,
	}
	_, err = txq.UpdateProblemDescription(pr.ctx, updateProblemDescParams)
	if err != nil {
		return err
	}

	return nil
}
