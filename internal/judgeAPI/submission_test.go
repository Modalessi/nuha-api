package judgeAPI

import (
	"bytes"
	"testing"

	"github.com/Modalessi/nuha-api/internal/models"
)

func TestGenerateBatch(t *testing.T) {
	sourceCode := "a, b = map(int, input().split()\n"
	sourceCode += "print(a + b)"

	submission := NewSubmission(sourceCode, PYTHON_311)

	tc1 := models.NewTestCase("5 6", "11")
	tc2 := models.NewTestCase("1 2", "3")
	tc3 := models.NewTestCase("9 1", "10")

	batch := submission.GenerateBatchFromTestCases(*tc1, *tc2, *tc3)

	want1 := Submission{
		SourceCode:     sourceCode,
		LanguageID:     92,
		Stdin:          "5 6",
		ExpectedOutput: "11",
	}
	got1 := batch[0]

	if want1 != got1 {
		t.Fatalf("got %v, wanted %v", got1, want1)
	}

	want2 := Submission{
		SourceCode:     sourceCode,
		LanguageID:     92,
		Stdin:          "1 2",
		ExpectedOutput: "3",
	}
	got2 := batch[1]

	if want2 != got2 {
		t.Fatalf("got %v, wanted %v", got2, want2)
	}

	want3 := Submission{
		SourceCode:     sourceCode,
		LanguageID:     92,
		Stdin:          "9 1",
		ExpectedOutput: "10",
	}
	got3 := batch[2]

	if want3 != got3 {
		t.Fatalf("got %v, wanted %v", got3, want3)
	}

}

func TestGenerateBatchJson(t *testing.T) {
	sourceCode := "a, b = map(int, input().split()\n"
	sourceCode += "print(a + b)"

	submission := NewSubmission(sourceCode, PYTHON_311)

	tc1 := models.NewTestCase("5 6", "11")
	tc2 := models.NewTestCase("1 2", "3")
	tc3 := models.NewTestCase("9 1", "10")

	batch := submission.GenerateBatchFromTestCases(*tc1, *tc2, *tc3)

	batchJson := batch.JSON()

	want := `[`
	want += `{"source_code":"a, b = map(int, input().split()\nprint(a + b)","language_id":92,"stdin":"5 6","expected_output":"11"},`
	want += `{"source_code":"a, b = map(int, input().split()\nprint(a + b)","language_id":92,"stdin":"1 2","expected_output":"3"},`
	want += `{"source_code":"a, b = map(int, input().split()\nprint(a + b)","language_id":92,"stdin":"9 1","expected_output":"10"}`
	want += `]`

	if !bytes.Equal(batchJson, []byte(want)) {
		t.Fatalf("wanted %v\ngot %v", string(want), string(batchJson))
	}

}
