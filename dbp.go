package dbp

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"KairosDrive/src/kairos.drive/elast"
	"KairosDrive/src/kairos.drive/utils/consts/ccomponents"
	server "KairosDrive/src/kairos.drive/utils/srv"
	_ "github.com/lib/pq"
)

type contextKey string

//DB contains database connection
var curSession = contextKey("curSession")

const (
	dbhost = "DBHOST"
	dbport = "DBPORT"
	dbuser = "DBUSER"
	dbpass = "DBPASS"
	dbname = "DBNAME"
)

//InitDB creates db connection and object
var InitDB = func() {
	srv := &server.KairosService
	config := dbConfig()
	var err error
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config[dbhost], config[dbport],
		config[dbuser], config[dbpass], config[dbname])

	srv.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Println(err)
	}
	err = srv.DB.Ping()
	if err != nil {
		srv.DB.Close()
		elast.LogSys(nil, "500", "Database is gone Kairos Drive shutting down for restart", nil, ccomponents.DBService)
		srv.HTTP.Close()
		panic(err)
	}
	elast.LogSys(nil, "200", "Successfully connected to database "+config[dbname]+" on "+config[dbhost]+":"+config[dbport], nil, ccomponents.DBService)
}

func dbConfig() map[string]string {
	conf := make(map[string]string)
	host, ok := os.LookupEnv("db_host")
	if !ok {
		log.Println("DBHOST environment variable required but not set")
	}
	port, ok := os.LookupEnv("db_port")
	if !ok {
		log.Println("DBPORT environment variable required but not set")
	}
	user, ok := os.LookupEnv("db_user")
	if !ok {
		log.Println("DBUSER environment variable required but not set")
	}
	password, ok := os.LookupEnv("db_pass")
	if !ok {
		log.Println("DBPASS environment variable required but not set")
	}
	name, ok := os.LookupEnv("db_name")
	if !ok {
		log.Println("DBNAME environment variable required but not set")
	}
	conf[dbhost] = host
	conf[dbport] = port
	conf[dbuser] = user
	conf[dbpass] = password
	conf[dbname] = name
	return conf
}

//DoQuery ... executes Query, sets application name
func DoQuery(r *http.Request, sqlStatement string, args ...interface{}) (rows *sql.Rows, err error) {
	srv := &server.KairosService
	err = srv.DB.Ping()
	if err != nil {
		srv.DB.Close()
		elast.LogSys(r, "500", "Database is gone Kairos Drive shutting down for restart", nil, ccomponents.DBService)
		srv.HTTP.Close()
		panic(err)
	}
	usrid := r.Context().Value("userID")
	str := fmt.Sprintf("%v", usrid)
	srv.DB.Exec("SET application_name TO 'Drive:/" + r.URL.Path + "|" + str + "';")
	rows, err = srv.DB.Query(sqlStatement, args...)
	return
}
