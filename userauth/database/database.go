package database

import (
	"database/sql"
	"fmt"

	logger "github.com/micro/micro/v3/service/logger"

	_ "github.com/lib/pq"
)

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

// DB is a global variable to hold db connection
var DB *sql.DB

// type DB struct {
// 	*sql.DB // specify Database Client
// }

func InitDb() {

	// var con *sql.DB
	config := dbConfig()
	// var Err error

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		msg := fmt.Sprintf("ERROR => Error creating Database...!\n%s", err.Error())
		logger.Info(msg)
		// panic(print_message.PrintMessage(msg))
	}
	err = db.Ping()
	if err != nil {
		msg := fmt.Sprintf("ERROR => Unable to connect Database...!\n%s", err.Error())
		logger.Info(msg)
		// panic(print_message.PrintMessage(msg))
	}
	fmt.Println("Successfully connected!")
	DB = db
	// return db
}

//.......................................................................
func dbConfig() map[string]string {

	// Err := godotenv.Load("userauth/.env")

	// if Err != nil {
	// 	log.Fatalf("Error loading .env file")
	// }

	conf := make(map[string]string)
	// host, ok := os.LookupEnv(dbhost)
	// if !ok {
	// 	panic("DBHOST environment variable required but not set")
	// }
	// port, ok := os.LookupEnv(dbport)
	// if !ok {
	// 	panic("DBPORT environment variable required but not set")
	// }
	// user, ok := os.LookupEnv(dbuser)
	// if !ok {
	// 	panic("DBUSER environment variable required but not set")
	// }
	// password, ok := os.LookupEnv(dbpass)
	// if !ok {
	// 	panic("DBPASS environment variable required but not set")
	// }
	// name, ok := os.LookupEnv(dbname)
	// if !ok {
	// 	panic("DBNAME environment variable required but not set")
	// }
	conf[dbhost] = "localhost"   //host
	conf[dbport] = "5432"        //port
	conf[dbuser] = "priadarshan" //user
	conf[dbpass] = ""            //password
	conf[dbname] = "priadarshan" // name
	return conf
}
