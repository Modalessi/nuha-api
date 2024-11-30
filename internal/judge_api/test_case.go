package judge_api

type Testcase struct {
	stdin          string
	expectedOutput string
}

func NewTestCase(stdin string, expectedOutput string) *Testcase {
	return &Testcase{
		stdin:          stdin,
		expectedOutput: expectedOutput,
	}
}
