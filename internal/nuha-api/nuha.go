package nuha

import (
	"database/sql"
	"net/http"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/judgeAPI"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type NuhaServer struct {
	Server     *http.ServeMux
	JudgeAPI   *judgeAPI.JudgeAPI
	DB         *sql.DB
	DBQueries  *database.Queries
	S3         *S3Config
	JWTSecret  string
	AdminEmail string
}

type S3Config struct {
	Client     *s3.Client
	BucketName string
}

func NewServer(ja *judgeAPI.JudgeAPI, db *sql.DB, dbQuereis *database.Queries, s3config *S3Config, jwtSecret string, adminEmail string) *NuhaServer {
	serverMux := http.NewServeMux()

	ns := NuhaServer{
		Server:     serverMux,
		JudgeAPI:   ja,
		DB:         db,
		DBQueries:  dbQuereis,
		S3:         s3config,
		JWTSecret:  jwtSecret,
		AdminEmail: adminEmail,
	}

	serverMux.HandleFunc("GET /healthz", checkHealth)

	serverMux.HandleFunc("GET /login", withServer(&ns, login))
	serverMux.HandleFunc("POST /register", withServer(&ns, register))

	serverMux.HandleFunc("GET /protected", authorized(withServer(&ns, protected), ns.JWTSecret))

	serverMux.HandleFunc("POST /submit", authorized(withServer(&ns, submitSolution), ns.JWTSecret))

	serverMux.HandleFunc("POST /problem", authorized(adminOnly(withServer(&ns, createProblem), adminEmail), ns.JWTSecret))
	serverMux.HandleFunc("GET /problem", withServer(&ns, getProblem))
	serverMux.HandleFunc("DELETE /problem", authorized(adminOnly(withServer(&ns, deleteProblem), ns.AdminEmail), ns.JWTSecret))

	serverMux.HandleFunc("POST /testcase", authorized(adminOnly(withServer(&ns, addTestCases), adminEmail), ns.JWTSecret))

	return &ns

}

func NewS3Config(client *s3.Client, bucketName string) *S3Config {
	return &S3Config{
		Client:     client,
		BucketName: bucketName,
	}
}

func (ns *NuhaServer) GetServer() http.Handler {
	return ns.Server
}
