package nuha

import (
	"net/http"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/judgeAPI"
)

type NuhaServer struct {
	Server    *http.ServeMux
	JudgeAPI  *judgeAPI.JudgeAPI
	DB        *database.Queries
	JWTSecret string
}

func NewServer(ja *judgeAPI.JudgeAPI, db *database.Queries, jwtSecret string) *NuhaServer {
	serverMux := http.NewServeMux()

	ns := NuhaServer{
		Server:    serverMux,
		JudgeAPI:  ja,
		DB:        db,
		JWTSecret: jwtSecret,
	}

	serverMux.HandleFunc("GET /healthz", checkHealth)

	serverMux.HandleFunc("GET /login", InjectNuhaServer(&ns, login))
	serverMux.HandleFunc("POST /register", InjectNuhaServer(&ns, register))

	serverMux.HandleFunc("GET /protected", Authorized(InjectNuhaServer(&ns, protected), ns.JWTSecret))
	serverMux.HandleFunc("POST /submit", Authorized(InjectNuhaServer(&ns, submitSolution), ns.JWTSecret))

	return &ns

}

func (ns *NuhaServer) GetServer() http.Handler {
	return ns.Server
}
