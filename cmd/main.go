package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

	"github.com/Modalessi/nuha-api/internal/database"
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

	dbConnection, err := sql.Open("postgres", DB_URL)
	utils.Assert(err, "error connecting to database")

	database := database.New(dbConnection)

	nuhaServer := nuha.NewServer(database)

	fmt.Println("server is now running...")
	http.ListenAndServe("localhost"+ADDRESS, nuhaServer.GetServer())
}
