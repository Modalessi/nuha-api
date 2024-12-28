package nuha

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/Modalessi/nuha-api/internal/auth"
	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/email-service"
	"github.com/Modalessi/nuha-api/internal/judgeAPI"
	submissionsPL "github.com/Modalessi/nuha-api/internal/nuha-api/submissions_pipeline"
	"github.com/Modalessi/nuha-api/internal/repositories"
)

type NuhaServer struct {
	Server        *http.ServeMux
	JudgeAPI      *judgeAPI.JudgeAPI
	SubmissionsPL *submissionsPL.SubmissionsPipeline
	DB            *sql.DB
	DBQueries     *database.Queries
	Auth          *auth.AuthService
	UserRepo      *repositories.UserRespository
	JWTSecret     string
	AdminEmail    string
}

func NewServer(ja *judgeAPI.JudgeAPI, db *sql.DB, dbQuereis *database.Queries, jwtSecret string, adminEmail string) *NuhaServer {
	serverMux := http.NewServeMux()

	submissionsPipeline := submissionsPL.NewSubmissionPipeline(ja, db, dbQuereis)

	authConfig := auth.AuthServiceConfig{
		JWTSecretKey:              jwtSecret,
		TokensExpirationsDuration: (time.Hour * 24),
	}

	authService := auth.NewAuthService(db, dbQuereis, &email.EmailService{}, authConfig)
	userRepo := repositories.NewUserRespository(db, dbQuereis)

	ns := NuhaServer{
		Server:        serverMux,
		JudgeAPI:      ja,
		SubmissionsPL: submissionsPipeline,
		DB:            db,
		DBQueries:     dbQuereis,
		Auth:          authService,
		UserRepo:      userRepo,
		JWTSecret:     jwtSecret,
		AdminEmail:    adminEmail,
	}

	serverMux.HandleFunc("GET /healthz", checkHealth)

	serverMux.HandleFunc("POST /login", withServer(&ns, login))
	serverMux.HandleFunc("POST /logout", authorized(withServer(&ns, logout), ns.Auth))
	serverMux.HandleFunc("POST /register", withServer(&ns, register))
	serverMux.HandleFunc("GET /verify/", withServer(&ns, verifyEmail))

	serverMux.HandleFunc("GET /protected", authorized(withServer(&ns, protected), ns.Auth))

	serverMux.HandleFunc("POST /submit", authorized(withServer(&ns, submitSolution), ns.Auth))
	serverMux.HandleFunc("GET /submit", authorized(withServer(&ns, getSubmission), ns.Auth))

	serverMux.HandleFunc("POST /problem", authorized(adminOnly(withServer(&ns, createProblem), adminEmail), ns.Auth))
	serverMux.HandleFunc("GET /problem", withServer(&ns, getProblem))
	serverMux.HandleFunc("DELETE /problem", authorized(adminOnly(withServer(&ns, deleteProblem), ns.AdminEmail), ns.Auth))
	serverMux.HandleFunc("PUT /problem", authorized(adminOnly(withServer(&ns, updateProblem), ns.AdminEmail), ns.Auth))

	serverMux.HandleFunc("POST /testcase", authorized(adminOnly(withServer(&ns, addTestCases), adminEmail), ns.Auth))

	ns.SubmissionsPL.Start()
	return &ns

}

func (ns *NuhaServer) GetServer() http.Handler {
	return ns.Server
}
