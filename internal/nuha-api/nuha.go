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
	Server        http.Handler
	serverMux     *http.ServeMux
	JudgeAPI      *judgeAPI.JudgeAPI
	SubmissionsPL *submissionsPL.SubmissionsPipeline
	DB            *sql.DB
	DBQueries     *database.Queries
	Auth          *auth.AuthService
	UserRepo      *repositories.UserRespository
	JWTSecret     string
	AdminEmail    string
}

// this should be better
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		w.Header().Set("Access-Control-Allow-Credentials", "true")

		w.Header().Set("Access-Control-Max-Age", "300")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
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
		serverMux:     serverMux,
		JudgeAPI:      ja,
		SubmissionsPL: submissionsPipeline,
		DB:            db,
		DBQueries:     dbQuereis,
		Auth:          authService,
		UserRepo:      userRepo,
		JWTSecret:     jwtSecret,
		AdminEmail:    adminEmail,
	}

	corsHandler := CORSMiddleware(serverMux)

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

	ns.Server = corsHandler
	return &ns
}

func (ns *NuhaServer) GetServer() http.Handler {
	return ns.Server
}
