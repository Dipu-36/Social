package main

import (
	"log"

	"github.com/Dipu-36/social/internal/db"
	"github.com/Dipu-36/social/internal/env"
	"github.com/Dipu-36/social/internal/store"
	"github.com/joho/godotenv"
)
const version = "0.0.1"

func main() {
	// Load the .env file before accessing environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	} else {
		log.Println(".env file loaded successfully")
	}

	//Initializing the dependencies
	cfg := config{
		addr: env.GetString("ADDR", ":8081"),
		db: dbConfig{
			addr: env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),//This defines the maximum number of open connections to the DB at any given time. It limits how many simultaneous connection our application can have to avoid overwhelming the database
			maxidleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),// This specifies the maximum number of idle (unused) connections that can be kept in the pool. Idle connections are kept alive so they can be reused without needing to reconnect, which improves efficiency.
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"), //This defines how long an idle connection can remain open before it's closed. It helps release resources when they're no longer needed.
		},
		env : env.GetString("ENV", "development"),
	}
	log.Printf("Attempting connection to: %s", cfg.db.addr)
	//Establishing the database connection using the db package from the internal directory
	db, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxidleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Print(err)
	}

	defer db.Close()
	log.Println("database connection pool established")
	//passing the pool into the storage layer or store package, PostStore and UserStore use this pool for queries
	store := store.NewStorage(db)
	//passing the dependencies to hold it application wide
	app := &application{
		config: cfg,
		store:  store,
	}
	mux := app.mount()      //Injecting dependencies into the router setup
	log.Fatal(app.run(mux)) //Inject router into server
}
