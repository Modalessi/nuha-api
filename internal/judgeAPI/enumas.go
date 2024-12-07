package judgeAPI

type JudgeLanguage int

const (
	PYTHON_312    JudgeLanguage = 100
	PYTHON_311    JudgeLanguage = 92
	C_CLANG18     JudgeLanguage = 104
	ASSEMBLY_NASM JudgeLanguage = 45
	BASH          JudgeLanguage = 46
	BASIC_FBC     JudgeLanguage = 47
	C_GCC_7       JudgeLanguage = 48
	CPP_GCC_7     JudgeLanguage = 52
	C_GCC_8       JudgeLanguage = 49
	CPP_GCC_8     JudgeLanguage = 53
	C_GCC_9       JudgeLanguage = 50
	CPP_GCC_9     JudgeLanguage = 54
	CSHARP_MONO   JudgeLanguage = 51
	COMMON_LISP   JudgeLanguage = 55
	D_DMD         JudgeLanguage = 56
	ELIXIR        JudgeLanguage = 57
	ERLANG        JudgeLanguage = 58
	EXECUTABLE    JudgeLanguage = 44
	FORTRAN       JudgeLanguage = 59
	GO            JudgeLanguage = 60
	HASKELL       JudgeLanguage = 61
	JAVA          JudgeLanguage = 62
	JAVASCRIPT    JudgeLanguage = 63
	LUA           JudgeLanguage = 64
	OCAML         JudgeLanguage = 65
	OCTAVE        JudgeLanguage = 66
	PASCAL        JudgeLanguage = 67
	PHP           JudgeLanguage = 68
	PLAIN_TEXT    JudgeLanguage = 43
	PROLOG        JudgeLanguage = 69
	PYTHON_2      JudgeLanguage = 70
	PYTHON_3      JudgeLanguage = 71
	RUBY          JudgeLanguage = 72
	RUST          JudgeLanguage = 73
	TYPESCRIPT    JudgeLanguage = 74
)

var JudgeLanguageDescription = map[JudgeLanguage]string{
	PYTHON_312:    "Python (3.12)",
	PYTHON_311:    "Python (3.11)",
	C_CLANG18:     "C (Clang 18)",
	ASSEMBLY_NASM: "Assembly (NASM 2.14.02)",
	BASH:          "Bash (5.0.0)",
	BASIC_FBC:     "Basic (FBC 1.07.1)",
	C_GCC_7:       "C (GCC 7.4.0)",
	CPP_GCC_7:     "C++ (GCC 7.4.0)",
	C_GCC_8:       "C (GCC 8.3.0)",
	CPP_GCC_8:     "C++ (GCC 8.3.0)",
	C_GCC_9:       "C (GCC 9.2.0)",
	CPP_GCC_9:     "C++ (GCC 9.2.0)",
	CSHARP_MONO:   "C# (Mono 6.6.0.161)",
	COMMON_LISP:   "Common Lisp (SBCL 2.0.0)",
	D_DMD:         "D (DMD 2.089.1)",
	ELIXIR:        "Elixir (1.9.4)",
	ERLANG:        "Erlang (OTP 22.2)",
	EXECUTABLE:    "Executable",
	FORTRAN:       "Fortran (GFortran 9.2.0)",
	GO:            "Go (1.13.5)",
	HASKELL:       "Haskell (GHC 8.8.1)",
	JAVA:          "Java (OpenJDK 13.0.1)",
	JAVASCRIPT:    "JavaScript (Node.js 12.14.0)",
	LUA:           "Lua (5.3.5)",
	OCAML:         "OCaml (4.09.0)",
	OCTAVE:        "Octave (5.1.0)",
	PASCAL:        "Pascal (FPC 3.0.4)",
	PHP:           "PHP (7.4.1)",
	PLAIN_TEXT:    "Plain Text",
	PROLOG:        "Prolog (GNU Prolog 1.4.5)",
	PYTHON_2:      "Python (2.7.17)",
	PYTHON_3:      "Python (3.8.1)",
	RUBY:          "Ruby (2.7.0)",
	RUST:          "Rust (1.40.0)",
	TYPESCRIPT:    "TypeScript (3.7.4)",
}

type JudgeSubmissionStatusID int

const (
	IN_QUEUE_STATUS              JudgeSubmissionStatusID = 1
	PROCESSING_STATUS            JudgeSubmissionStatusID = 2
	ACCEPTED_STATUS              JudgeSubmissionStatusID = 3
	WRONG_ANSWER_STATUS          JudgeSubmissionStatusID = 4
	TIME_LIMIT_EXCEEDED_STATUS   JudgeSubmissionStatusID = 5
	COMPILATION_ERROR_STATUS     JudgeSubmissionStatusID = 6
	RUNTIME_ERROR_SIGSEGV_STATUS JudgeSubmissionStatusID = 7
	RUNTIME_ERROR_SIGXFSZ_STATUS JudgeSubmissionStatusID = 8
	RUNTIME_ERROR_SIGFPE_STATUS  JudgeSubmissionStatusID = 9
	RUNTIME_ERROR_SIGABRT_STATUS JudgeSubmissionStatusID = 10
	RUNTIME_ERROR_NZEC_STATUS    JudgeSubmissionStatusID = 11
	RUNTIME_ERROR_OTHER_STATUS   JudgeSubmissionStatusID = 12
	INTERNAL_ERROR_STATUS        JudgeSubmissionStatusID = 13
	EXECFORMAT_ERROR_STATUS      JudgeSubmissionStatusID = 14
)

var JudgeSubmissionStatusDescription = map[JudgeSubmissionStatusID]string{
	IN_QUEUE_STATUS:              "In Queue",
	PROCESSING_STATUS:            "Processing",
	ACCEPTED_STATUS:              "Accepted",
	WRONG_ANSWER_STATUS:          "Wrong Answer",
	TIME_LIMIT_EXCEEDED_STATUS:   "Time Limit Exceeded",
	COMPILATION_ERROR_STATUS:     "Compilation Error",
	RUNTIME_ERROR_SIGSEGV_STATUS: "Runtime Error (SIGSEGV)",
	RUNTIME_ERROR_SIGXFSZ_STATUS: "Runtime Error (SIGXFSZ)",
	RUNTIME_ERROR_SIGFPE_STATUS:  "Runtime Error (SIGFPE)",
	RUNTIME_ERROR_SIGABRT_STATUS: "Runtime Error (SIGABRT)",
	RUNTIME_ERROR_NZEC_STATUS:    "Runtime Error (NZEC)",
	RUNTIME_ERROR_OTHER_STATUS:   "Runtime Error (Other)",
	INTERNAL_ERROR_STATUS:        "Internal Error",
	EXECFORMAT_ERROR_STATUS:      "Exec Format Error",
}
