package nuha

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/Modalessi/nuha-api/internal/models"
	"github.com/Modalessi/nuha-api/internal/repositories"
)

func addTestCases(ns *NuhaServer, w http.ResponseWriter, r *http.Request) error {
	problemId := r.URL.Query().Get("problem_id")
	if problemId == "" {
		respondWithError(w, 400, INVALID_QUERY_ERROR)
		return fmt.Errorf("error, problem_id query was not provided")
	}

	defer r.Body.Close()

	requestTestcases := []models.Testcase{}

	// check if file exist

	file, fileHeader, err := r.FormFile("testcases_file")
	if err == nil {

		requestTestcases, err = getTestCasesFromZipFile(&file, fileHeader.Size)
		if err != nil {
			respondWithError(w, 400, err)
			return err
		}

	} else {
		err = json.NewDecoder(r.Body).Decode(&requestTestcases)
		if err != nil {
			respondWithError(w, 400, INVALID_JSON_ERROR)
			return err
		}
	}

	pr := repositories.NewProblemRepository(ns.S3.Client, ns.DB, ns.DBQueries, r.Context(), ns.S3.BucketName)

	// check if problem exist
	problemDb, err := pr.GetProblemInfo(problemId)
	if err != nil {
		respondWithError(w, 400, err)
		return err
	}

	// store test cases
	err = pr.AddNewTestCases(problemDb.ID.String(), requestTestcases...)
	if err != nil {
		respondWithError(w, 500, SERVER_ERROR)
		return err
	}

	respondWithSuccess(w, 201, "test cases addes successfuly")
	return nil
}

func getTestCasesFromZipFile(file *multipart.File, size int64) ([]models.Testcase, error) {
	zipReader, err := zip.NewReader(*file, size)
	if err != nil {
		return nil, fmt.Errorf("error reading testcases zip file: %w", err)
	}

	testcases := []models.Testcase{}

	for i := range 500 {
		i += 1

		inFile, err := zipReader.Open(fmt.Sprintf("%d.in", i))
		if err != nil {
			log.Printf("inFile: %v", err)
			break
		}
		outFile, err := zipReader.Open(fmt.Sprintf("%d.out", i))
		if err != nil {
			return nil, fmt.Errorf(".in files and .out files are not matching")
		}

		inContent, err := io.ReadAll(inFile)
		if err != nil {
			return nil, fmt.Errorf("error reading input file: %v", err)
		}
		outContent, err := io.ReadAll(outFile)
		if err != nil {
			return nil, fmt.Errorf("error reading output file: %v", err)
		}

		tc := models.Testcase{
			Stdin:          string(inContent),
			ExpectedOutput: string(outContent),
		}

		testcases = append(testcases, tc)
	}

	return testcases, nil
}
