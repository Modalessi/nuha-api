package nuha

import (
	"net/http"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/handlers"
)

type NuhaServer struct {
	Server *http.ServeMux
	DB     *database.Queries
}

func NewServer(db *database.Queries) *NuhaServer {
	serverMux := http.NewServeMux()

	addRoutes(serverMux)

	return &NuhaServer{
		Server: serverMux,
		DB:     db,
	}
}

func addRoutes(serv *http.ServeMux) {

	// all the app routes
	serv.HandleFunc("GET /healthz", handlers.CheckHealth)

}

func (ns *NuhaServer) GetHandler() http.Handler {
	return ns.Server
}
