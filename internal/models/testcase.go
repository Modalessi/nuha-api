package models

import (
	"encoding/json"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/utils"
)

type Testcase struct {
	Stdin          string `json:"stdin"`
	ExpectedOutput string `json:"expected_output"`
}

func NewTestCase(stdin string, expectedOutput string) *Testcase {
	return &Testcase{
		Stdin:          stdin,
		ExpectedOutput: expectedOutput,
	}
}

func TestCasesFromDBObjects(testcases []database.TestCase) []Testcase {

	tcs := make([]Testcase, len(testcases))
	for i, tc := range testcases {
		tcs[i] = *NewTestCase(tc.Stdin, tc.ExpectedOutput)
	}

	return tcs
}

func (t *Testcase) JSON() []byte {
	data, err := json.Marshal(t)
	utils.Assert(err, "error converting user object to json")
	return data
}
