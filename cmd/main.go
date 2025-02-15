package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/judgeAPI"
	"github.com/Modalessi/nuha-api/internal/nuha-api"
	"github.com/Modalessi/nuha-api/internal/utils"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load()
	utils.Assert(err, "somthing went wrong while loading .env file")

	DB_URL := os.Getenv("DB_URL")
	utils.AssertOn(DB_URL != "", "something went wrong when reading 'DB_URL' env variable")

	ADDRESS := os.Getenv("ADDRESS")
	utils.AssertOn(ADDRESS != "", "something went wrong when reading 'ADDRESS' env variable")

	JWTSecret := os.Getenv("JWT_SECRET")
	utils.AssertOn(JWTSecret != "", "somethign went wrong when reading 'JWT_SECRET' env variable")

	db, err := sql.Open("postgres", DB_URL)
	utils.Assert(err, "error connecting to database")

	dbQueries := database.New(db)

	rapidAPIKey := os.Getenv("X_RAPIDAPI_KEY")
	utils.AssertOn(rapidAPIKey != "", "somethign went wrong when reading 'X_RAPIDAPI_KEY' env variable")

	rapidAPIHost := os.Getenv("X_RAPIDAPI_HOST")
	utils.AssertOn(rapidAPIKey != "", "somethign went wrong when reading 'X_RAPIDAPI_HOST' env variable")

	judgeAPI := judgeAPI.NewJudgeAPI(rapidAPIKey, rapidAPIHost)

	adminEmail := os.Getenv("ADMIN_EMAIL")
	utils.AssertOn(adminEmail != "", "somethign went wrong when reading 'ADMIN_EMAIL' env variable")

	nuhaServer := nuha.NewServer(judgeAPI, db, dbQueries, JWTSecret, adminEmail)

	// TODO: gracful shut down for SIGINT, SIGTERM
	fmt.Println("server is now running...")
	http.ListenAndServe("localhost"+ADDRESS, nuhaServer.GetServer())
}
