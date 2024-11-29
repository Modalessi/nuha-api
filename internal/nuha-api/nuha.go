package nuha

import (
	"net/http"

	"github.com/Modalessi/nuha-api/internal/database"
)

type NuhaServer struct {
	Server    *http.ServeMux
	DB        *database.Queries
	JWTSecret string
}

func NewServer(db *database.Queries, jwtSecret string) *NuhaServer {
	serverMux := http.NewServeMux()

	ns := NuhaServer{
		Server:    serverMux,
		DB:        db,
		JWTSecret: jwtSecret,
	}

	serverMux.HandleFunc("GET /healthz", checkHealth)

	serverMux.HandleFunc("GET /login", InjectNuhaServer(&ns, login))
	serverMux.HandleFunc("POST /register", InjectNuhaServer(&ns, register))

	return &ns

}

func (ns *NuhaServer) GetServer() http.Handler {
	return ns.Server
}
