package nuha

import (
	"net/http"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/judgeAPI"
)

type NuhaServer struct {
	Server     *http.ServeMux
	JudgeAPI   *judgeAPI.JudgeAPI
	DB         *database.Queries
	JWTSecret  string
	AdminEmail string
}

func NewServer(ja *judgeAPI.JudgeAPI, db *database.Queries, jwtSecret string, adminEmail string) *NuhaServer {
	serverMux := http.NewServeMux()

	ns := NuhaServer{
		Server:     serverMux,
		JudgeAPI:   ja,
		DB:         db,
		JWTSecret:  jwtSecret,
		AdminEmail: adminEmail,
	}

	serverMux.HandleFunc("GET /healthz", checkHealth)

	serverMux.HandleFunc("GET /login", injectNuhaServer(&ns, login))
	serverMux.HandleFunc("POST /register", injectNuhaServer(&ns, register))

	serverMux.HandleFunc("GET /protected", authorized(injectNuhaServer(&ns, protected), ns.JWTSecret))

	serverMux.HandleFunc("POST /submit", authorized(injectNuhaServer(&ns, submitSolution), ns.JWTSecret))

	serverMux.HandleFunc("POST /problem", authorized(adminOnly(injectNuhaServer(&ns, createProblem), adminEmail), ns.JWTSecret))

	return &ns

}

func (ns *NuhaServer) GetServer() http.Handler {
	return ns.Server
}
