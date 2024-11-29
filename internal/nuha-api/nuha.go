package nuha

import (
	"net/http"

	"github.com/Modalessi/nuha-api/internal/database"
)

type NuhaServer struct {
	Server *http.ServeMux
	DB     *database.Queries
}

func NewServer(db *database.Queries) *NuhaServer {
	serverMux := http.NewServeMux()

	ns := NuhaServer{
		Server: serverMux,
		DB:     db,
	}

	serverMux.HandleFunc("GET /healthz", checkHealth)
	serverMux.HandleFunc("POST /register", InjectNuhaServer(&ns, register))

	return &ns

}

func (ns *NuhaServer) GetServer() http.Handler {
	return ns.Server
}
