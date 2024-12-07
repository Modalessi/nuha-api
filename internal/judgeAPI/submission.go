package judgeAPI

import (
	"encoding/json"
	"strconv"

	"github.com/Modalessi/nuha-api/internal/models"
	"github.com/Modalessi/nuha-api/internal/utils"
)

type SubmissionBatch []Submission

type Submission struct {
	SourceCode                           string           `json:"source_code,omitempty"`
	LanguageID                           int              `json:"language_id,omitempty"`
	CompilerOptions                      string           `json:"compiler_options,omitempty"`
	CommandLineArguments                 string           `json:"command_line_arguments,omitempty"`
	Stdin                                string           `json:"stdin,omitempty"`
	ExpectedOutput                       string           `json:"expected_output,omitempty"`
	CPUTimeLimit                         string           `json:"cpu_time_limit,omitempty"`
	CPUExtraTime                         string           `json:"cpu_extra_time,omitempty"`
	WallTimeLimit                        string           `json:"wall_time_limit,omitempty"`
	MemoryLimit                          float64          `json:"memory_limit,omitempty"`
	StackLimit                           int              `json:"stack_limit,omitempty"`
	MaxProcessesAndOrThreads             int              `json:"max_processes_and_or_threads,omitempty"`
	EnablePerProcessAndThreadTimeLimit   bool             `json:"enable_per_process_and_thread_time_limit,omitempty"`
	EnablePerProcessAndThreadMemoryLimit bool             `json:"enable_per_process_and_thread_memory_limit,omitempty"`
	MaxFileSize                          int              `json:"max_file_size,omitempty"`
	RedirectStderrToStdout               bool             `json:"redirect_stderr_to_stdout,omitempty"`
	EnableNetwork                        bool             `json:"enable_network,omitempty"`
	NumberOfRuns                         int              `json:"number_of_runs,omitempty"`
	AdditionalFiles                      string           `json:"additional_files,omitempty"`
	CallbackURL                          string           `json:"callback_url,omitempty"`
	Stdout                               string           `json:"stdout,omitempty"`
	Stderr                               string           `json:"stderr,omitempty"`
	CompileOutput                        string           `json:"compile_output,omitempty"`
	Message                              string           `json:"message,omitempty"`
	ExitCode                             int              `json:"exit_code,omitempty"`
	ExitSignal                           int              `json:"exit_signal,omitempty"`
	Status                               SubmissionStatus `json:"status,omitempty"`
	CreatedAt                            string           `json:"created_at,omitempty"`
	FinishedAt                           string           `json:"finished_at,omitempty"`
	Token                                string           `json:"token,omitempty"`
	Time                                 string           `json:"time,omitempty"`
	WallTime                             string           `json:"wall_time,omitempty"`
	Memory                               float64          `json:"memory,omitempty"`
}

type SubmissionStatus struct {
	ID          JudgeSubmissionStatusID `json:"id"`          // Status ID corresponding to our JudgeSubmissionStatus constants
	Description string                  `json:"description"` // Human readable description of the status
}

func NewSubmission(sourceCode string, language JudgeLanguage) *Submission {
	return &Submission{
		SourceCode: sourceCode,
		LanguageID: int(language),
	}
}

func (s *Submission) SetCompilerOptions(options string) {
	s.CompilerOptions = options
}

func (s *Submission) SetCommandLineArguments(args string) {
	s.CommandLineArguments = args
}

func (s *Submission) SetStdin(stdin string) {
	s.Stdin = stdin
}

func (s *Submission) SetExpectedOutput(expectedOutput string) {
	s.ExpectedOutput = expectedOutput
}

func (s *Submission) SetCPUTimeLimit(limit float64) {
	s.CPUTimeLimit = strconv.FormatFloat(limit, 'f', -1, 64)
}

func (s *Submission) SetCPUExtraTime(time float64) {
	s.CPUExtraTime = strconv.FormatFloat(time, 'f', -1, 64)
}

func (s *Submission) SetWallTimeLimit(limit float64) {
	s.WallTimeLimit = strconv.FormatFloat(limit, 'f', -1, 64)
}

func (s *Submission) SetMemoryLimit(limit float64) {
	s.MemoryLimit = limit
}

func (s *Submission) SetStackLimit(limit int) {
	s.StackLimit = limit
}

func (s *Submission) SetMaxProcessesAndThreads(max int) {
	s.MaxProcessesAndOrThreads = max
}

func (s *Submission) SetEnablePerProcessTimeLimit(enable bool) {
	s.EnablePerProcessAndThreadTimeLimit = enable
}

func (s *Submission) SetEnablePerProcessMemoryLimit(enable bool) {
	s.EnablePerProcessAndThreadMemoryLimit = enable
}

func (s *Submission) SetMaxFileSize(size int) {
	s.MaxFileSize = size
}

func (s *Submission) SetRedirectStderr(redirect bool) {
	s.RedirectStderrToStdout = redirect
}

func (s *Submission) SetEnableNetwork(enable bool) {
	s.EnableNetwork = enable
}

func (s *Submission) SetNumberOfRuns(runs int) {
	s.NumberOfRuns = runs
}

func (s *Submission) SetAdditionalFiles(files string) {
	s.AdditionalFiles = files
}

func (s *Submission) SetCallbackURL(url string) {
	s.CallbackURL = url
}

func (s *Submission) GenerateBatchFromTestCases(testcases ...models.Testcase) *SubmissionBatch {
	batch := SubmissionBatch{}

	for _, tc := range testcases {
		ns := *s
		ns.SetStdin(tc.Stdin)
		ns.SetExpectedOutput(tc.ExpectedOutput)

		batch = append(batch, ns)
	}

	return &batch
}

func (s *Submission) JSON() []byte {
	data, err := json.Marshal(s)
	utils.Assert(err, "error converting submission struct to json")
	return data
}

func (sb *SubmissionBatch) JSON() []byte {
	wrapper := struct {
		Submissions []Submission `json:"submissions"`
	}{
		Submissions: *sb,
	}

	data, err := json.Marshal(wrapper)
	utils.Assert(err, "error converting submission struct to json")
	return data
}
