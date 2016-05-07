package logMonitor

import (
	"github.com/revel/revel"
    "database/sql"
	"runtime"
	"strings"
	"strconv"
	"time"
	"os"
	"encoding/json"
	//
	_ "github.com/lib/pq"
)

var settings struct {
    DateTimeFormatString  string
	DbDateTimeFormatString  string
	StackLevel int
	DbDriver string
	DbUser string
    DbUserPassword string
    DbHost string
    DbPort string
    DbName string
	DbShema string
	DbTable string
    Sslmode string	
}

var db *sql.DB
var insertQuery string

// Config ...
func Config() {
	
	// чтение конфига подключения к базе данных из JSON файла конфига
	configFile, err := os.Open("app/services/logMonitor/config.json")
    checkErr(err)
    jsonParser := json.NewDecoder(configFile)
    err = jsonParser.Decode(&settings)
	checkErr(err)
	
	insertQuery= `INSERT INTO ` + settings.DbShema + `.` + settings.DbTable + ` (
						time_stamp, 
						app_name, 
						pkg_name, 
						module_name, 
						proc_name, 
						log_context, 
						log_text, 
						log_type,
						err_code
					) VALUES (NOW(), $1, $2, $3, $4, $5, $6, $7, $8);`
	
	openDB()
	defer closeDB()
	
	_, err = db.Query(
		`CREATE SCHEMA IF NOT EXISTS ` + settings.DbShema + ` 
			AUTHORIZATION ` + settings.DbUser + `;
			GRANT ALL ON SCHEMA ` + settings.DbShema + ` TO ` + settings.DbUser + `;
			COMMENT ON SCHEMA ` + settings.DbShema + `
			IS 'standard logs schema';`)
	checkErr(err)	
	
	_, err = db.Query(
		`CREATE TABLE IF NOT EXISTS ` + settings.DbShema + `.` + settings.DbTable + ` (
			pk_id serial,
			time_stamp timestamp without time zone,
			app_name character varying(50),
			pkg_name character varying(50),
			module_name character varying(50),
			proc_name character varying(50),
			log_context character varying(100),
			log_text text,
			log_type character varying(50),
			err_code character varying(50),
			CONSTRAINT "первичный ключ" PRIMARY KEY (pk_id)
		) WITH (OIDS=FALSE);`)
	checkErr(err)
}

func checkErr(err error) {
    if err != nil {        
		revel.ERROR.Println(err)
        //panic(err)
    }
}

func strToTime (s string) time.Time {
	t, err := time.Parse(settings.DateTimeFormatString, s)
	checkErr(err)
	return t
}

func timeToDbStr (t time.Time) string {
	return t.Format(settings.DbDateTimeFormatString)
}

func timeToStr (t time.Time) string {
	return t.Format(settings.DateTimeFormatString)
}
					
// Get ... todo
func Get(params map[string]interface{}) (data string) {
	
	var dttmStart, dttmEnd time.Time	
	
	dttmStart = strToTime(params["dtStart"].(string))
	dttmEnd = strToTime(params["dtEnd"].(string))
		
	openDB()
	query := `SELECT 
                    pk_id, 
                    to_char(time_stamp, 'DD.MM.YYYY HH24:MI:SS') As time_stamp, 
                    app_name, 
                    pkg_name, 
                    module_name, 
                    proc_name, 
                    log_context, 
                    log_text, 
                    log_type,
                    err_code 
              FROM ` + settings.DbShema + `.` + settings.DbTable + ` 
			  WHERE (time_stamp BETWEEN '` + timeToDbStr(dttmStart) + `' AND '` + timeToDbStr(dttmEnd) + `')`
                    
    rows, err := db.Query("SELECT array_to_json(ARRAY_AGG(row_to_json(row))) FROM (" + query + ") row")
	defer rows.Close()
	checkErr(err)
	defer closeDB()
	var row sql.NullString
	for rows.Next() {
		err = rows.Scan(&row)
		if err != nil {
			revel.ERROR.Println("QueryManager.QueryJson scan error", err)
			break
		}
		data = row.String
	}
	return
}

// INFO ...
func INFO(logText, logContext, errCode string) {
	revel.INFO.Println(logText, logContext, errCode)
	openDB()
	defer closeDB()
	query, err := db.Prepare(insertQuery)
	_, err = query.Exec(getAppName(), getPkgName(), getModuleName(), getFuncName(), logContext, logText, "INFO", "")	
	checkErr(err)
}

// TRACE ...
func TRACE(logText, logContext, errCode string) {
	revel.TRACE.Println(logText, logContext, errCode)
	openDB()
	defer closeDB()	
	query, err := db.Prepare(insertQuery)
	_, err = query.Exec(getAppName(), getPkgName(), getModuleName(), getFuncName(), logContext, logText, "TRACE", "")	
	checkErr(err)
}

// WARN ...
func WARN(logText, logContext, errCode string) {
	revel.WARN.Println(logText, logContext, errCode)
	openDB()
	defer closeDB()	
	query, err := db.Prepare(insertQuery)
	_, err = query.Exec(getAppName(), getPkgName(), getModuleName(), getFuncName(), logContext, logText, "WARN", "")	
	checkErr(err)
}
// ERROR ...
func ERROR(logText, logContext, errCode string) {
	revel.ERROR.Println(logText, logContext, errCode)
	openDB()
	defer closeDB()	
	query, err := db.Prepare(insertQuery)
	_, err = query.Exec(getAppName(), getPkgName(), getModuleName(), getFuncName(), logContext, logText, "ERROR", errCode)	
	checkErr(err)
}

func getAppName() string {
	return revel.Config.StringDefault("app.name", "")
}

func getPkgName() string {
	pc, _, _, _ := runtime.Caller(settings.StackLevel)
	functionObject := runtime.FuncForPC(pc)
	arr := strings.Split(functionObject.Name(), ".")
	sPkg := strings.Split(arr[0], "/")
	return sPkg[len(sPkg)-1]
}

func getModuleName() string {
	_, modulePathName, _, _ := runtime.Caller(settings.StackLevel)
	sModule := strings.Split(modulePathName, "/")
	return sModule[len(sModule)-1]
}

func getFuncName() string {
	pc, _, _, _ := runtime.Caller(settings.StackLevel)
	functionObject := runtime.FuncForPC(pc)
	arr := strings.Split(functionObject.Name(), ".")
	return arr[len(arr)-1]
}

func getLine() string {
	_, _, line, _ := runtime.Caller(settings.StackLevel)
	return strconv.Itoa(line)
}

func openDB() (err error) {
	driver := settings.DbDriver
	connectString := settings.DbDriver + "://"
	connectString += settings.DbUser + ":"
	connectString += settings.DbUserPassword + "@"
	connectString += settings.DbHost + ":"
	connectString += settings.DbPort + "/"
	connectString += settings.DbName
	connectString += "?sslmode=" + settings.Sslmode	
	db, err = sql.Open(driver, connectString)
	if err != nil {
		revel.ERROR.Println("DB open Error", err)
		return err
	}
	return nil
}

// closeDB ...
func closeDB() (err error) {
	err = db.Close()
	if err != nil {
		revel.ERROR.Println("DB close Error", err)
		return err
	}
	return nil
}