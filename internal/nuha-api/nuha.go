package nuha

import (
	"database/sql"
	"net/http"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/judgeAPI"
	submissionsPL "github.com/Modalessi/nuha-api/internal/nuha-api/submissions_pipeline"
)

type NuhaServer struct {
	Server        *http.ServeMux
	JudgeAPI      *judgeAPI.JudgeAPI
	SubmissionsPL *submissionsPL.SubmissionsPipeline
	DB            *sql.DB
	DBQueries     *database.Queries
	JWTSecret     string
	AdminEmail    string
}

func NewServer(ja *judgeAPI.JudgeAPI, db *sql.DB, dbQuereis *database.Queries, jwtSecret string, adminEmail string) *NuhaServer {
	serverMux := http.NewServeMux()

	submissionsPipeline := submissionsPL.NewSubmissionPipeline(ja, db, dbQuereis)

	ns := NuhaServer{
		Server:        serverMux,
		JudgeAPI:      ja,
		SubmissionsPL: submissionsPipeline,
		DB:            db,
		DBQueries:     dbQuereis,
		JWTSecret:     jwtSecret,
		AdminEmail:    adminEmail,
	}

	serverMux.HandleFunc("GET /healthz", checkHealth)

	serverMux.HandleFunc("POST /login", withServer(&ns, login))
	serverMux.HandleFunc("POST /register", withServer(&ns, register))

	serverMux.HandleFunc("GET /protected", authorized(withServer(&ns, protected), ns.JWTSecret))

	serverMux.HandleFunc("POST /submit", authorized(withServer(&ns, submitSolution), ns.JWTSecret))
	serverMux.HandleFunc("GET /submit", authorized(withServer(&ns, getSubmission), ns.JWTSecret))

	serverMux.HandleFunc("POST /problem", authorized(adminOnly(withServer(&ns, createProblem), adminEmail), ns.JWTSecret))
	serverMux.HandleFunc("GET /problem", withServer(&ns, getProblem))
	serverMux.HandleFunc("DELETE /problem", authorized(adminOnly(withServer(&ns, deleteProblem), ns.AdminEmail), ns.JWTSecret))
	serverMux.HandleFunc("PUT /problem", authorized(adminOnly(withServer(&ns, updateProblem), ns.AdminEmail), ns.JWTSecret))

	serverMux.HandleFunc("POST /testcase", authorized(adminOnly(withServer(&ns, addTestCases), adminEmail), ns.JWTSecret))

	ns.SubmissionsPL.Start()
	return &ns

}

func (ns *NuhaServer) GetServer() http.Handler {
	return ns.Server
}
