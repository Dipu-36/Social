package db

import (
	"context"
	"database/sql"
	"time"
	_"github.com/lib/pq"
)

//New is responsible for creation of the databse connection pool for PostgresSQL database which is necessary to handle concurrrentt reuests in a web application 
func New(addr string, maxOpenConns, maxIdleConns int , maxIdleTime string) (*sql.DB, error){
	//Creates a pool of connections 
	db, err := sql.Open("postgres", addr)
	if err!=nil{
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)//Limits simultaneous active connections
	db.SetMaxIdleConns(maxIdleConns)//Keeps warm connections for reuse

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	//
	db.SetConnMaxIdleTime(duration)

	//if it takes more than 5 seconds to connect we wil have a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	//To verify if the connection to the database is still alive, and establishing a connection if necessary 
	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}
	return db, nil
}