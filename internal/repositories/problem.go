package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/models"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type ProblemRepository struct {
	db         *sql.DB
	dbQueries  *database.Queries
	ctx        context.Context
	s3Client   *s3.Client
	bucketName string
}

func NewProblemRepository(s3client *s3.Client, db *sql.DB, dbQueries *database.Queries, ctx context.Context, bucketName string) *ProblemRepository {
	return &ProblemRepository{
		db:         db,
		dbQueries:  dbQueries,
		s3Client:   s3client,
		ctx:        ctx,
		bucketName: bucketName,
	}
}

func (pr *ProblemRepository) StoreNewProblem(p *models.Problem) (*database.Problem, error) {

	descriptionPath := fmt.Sprintf("problems/%s/description.md", p.ID.String())
	testcasesPath := fmt.Sprintf("problems/%s/testcases", p.ID.String())

	newProblemParams := database.CreateProblemParams{
		ID:              p.ID,
		Title:           p.Title,
		DescriptionPath: descriptionPath,
		TestcasesPath:   testcasesPath,
		Tags:            p.Tags,
		TimeLimit:       p.Timelimit,
		MemoryLimit:     p.Memorylimit,
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

	// add description md to s3
	newObjectParams := &s3.PutObjectInput{
		Bucket: aws.String(pr.bucketName),
		Key:    aws.String(descriptionPath),
		Body:   strings.NewReader(p.Description),
	}
	_, err = pr.s3Client.PutObject(pr.ctx, newObjectParams)
	if err != nil {
		return nil, fmt.Errorf("uploading description to S3: %w", err)
	}

	// add test cases .in .out to s3
	for i, tc := range p.Testcases {
		inputPath := fmt.Sprintf("%s/%d.in", testcasesPath, i+1)
		newObjectParams = &s3.PutObjectInput{
			Bucket: aws.String(pr.bucketName),
			Key:    aws.String(inputPath),
			Body:   strings.NewReader(tc.Stdin),
		}
		_, err := pr.s3Client.PutObject(pr.ctx, newObjectParams)
		if err != nil {
			return nil, fmt.Errorf("uploading testcase input %d to S3: %w", i+1, err)
		}

		outputPath := fmt.Sprintf("%s/%d.out", testcasesPath, i+1)
		newObjectParams = &s3.PutObjectInput{
			Bucket: aws.String(pr.bucketName),
			Key:    aws.String(outputPath),
			Body:   strings.NewReader(tc.ExpectedOutput),
		}
		_, err = pr.s3Client.PutObject(pr.ctx, newObjectParams)
		if err != nil {
			return nil, fmt.Errorf("uploading testcase output %d to S3: %w", i+1, err)
		}

	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("committing transaction: %w", err)
	}

	return &dbProblem, nil
}
