package controllers

import (
	"github.com/revel/revel"
    "database/sql"
	"runtime"
	"strings"
	//_ "github.com/lib/pq"
)

var db *sql.DB
var insertQuery = `INSERT INTO logs (
						time_stamp, 
						app_name, 
						pkg_name, 
						module_name, 
						proc_name, 
						context, 
						log_text, 
						type
					) VALUES (NOW(), $1, $2, $3, $4, $5, $6, $7);`

// LogMonitor ...
type LogMonitor struct {
	*revel.Controller
}

// Index ...
func (c LogMonitor) Index() revel.Result {
	return c.Render()
}

// GetLogs ...
func (c LogMonitor) GetLogs() revel.Result {
	//return c.RenderJson("")
    return c.RenderJson(Get())
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}
					
// Get ...
func Get() (data string) {
	initDb()
	query := "SELECT pk_id, to_char(time_stamp, 'DD.MM.YYYY HH24:MI:SS') As time_stamp, app_name, pkg_name, module_name, proc_name, context, log_text, type FROM public.logs"
    rows, err := db.Query("SELECT array_to_json(ARRAY_AGG(row_to_json(row))) FROM (" + query + ") row")
	defer rows.Close()
	checkErr(err)
	defer closeDB()
	var row sql.NullString
	for rows.Next() {
		err = rows.Scan(&row)
		if err != nil {
			revel.INFO.Println("QueryManager.QueryJson scan error", err)
			break
		}
		data = row.String
	}
	return
}

// Info ...
func Info(context string, text string) (err error) {
	revel.INFO.Println(context, text)
	initDb()
	query, err := db.Prepare(insertQuery)
	_, err = query.Exec(
		getAppName(),
		getPkgName(),
		getModuleName(),
		getFuncName(),
		context,
		text,
		"INFO",
	)
	if err != nil {
		revel.ERROR.Println("db Error", err)
		return err
	}
	defer closeDB()
	return nil
}

// Trace ...
func Trace(context string, text string) (err error) {
	revel.TRACE.Println(context, text)
	initDb()
	query, err := db.Prepare(insertQuery)
	_, err = query.Exec(
		getAppName(),
		getPkgName(),
		getModuleName(),
		getFuncName(),
		context,
		text,
		"TRACE",
	)
	if err != nil {
		revel.ERROR.Println("DB Error", err)
		return err
	}
	defer closeDB()
	return nil
}

// Warn ...
func Warn(context string, text string) (err error) {
	revel.WARN.Println(context, text)
	initDb()
	query, err := db.Prepare(insertQuery)
	_, err = query.Exec(
		getAppName(),
		getPkgName(),
		getModuleName(),
		getFuncName(),
		context,
		text,
		"WARN",
	)
	if err != nil {
		revel.ERROR.Println("DB Error", err)
		return err
	}
	defer closeDB()
	return nil
}

// Error ...
func Error(context string, text string) (err error) {
	revel.ERROR.Println(context, text)
	initDb()
	query, err := db.Prepare(insertQuery)
	_, err = query.Exec(
		getAppName(),
		getPkgName(),
		getModuleName(),
		getFuncName(),
		context,
		text,
		"ERROR",
	)
	if err != nil {
		revel.ERROR.Println("DB Error", err)
		return err
	}
	defer closeDB()
	return nil
}

func getAppName() string {
	return revel.Config.StringDefault("app.name", "")
}

func getPkgName() string {
	pc, _, _, _ := runtime.Caller(2)
	functionObject := runtime.FuncForPC(pc)
	arr := strings.Split(functionObject.Name(), ".")
	sPkg := strings.Split(arr[0], "/")
	return sPkg[len(sPkg)-1]
}

func getModuleName() string {
	_, modulePathName, _, _ := runtime.Caller(2)
	sModule := strings.Split(modulePathName, "/")
	return sModule[len(sModule)-1]
}

func getFuncName() string {
	pc, _, _, _ := runtime.Caller(2)
	functionObject := runtime.FuncForPC(pc)
	arr := strings.Split(functionObject.Name(), ".")
	return arr[len(arr)-1]
}

func initDb() (err error) {
	driver := revel.Config.StringDefault("db.user", "postgres")
	connectString := revel.Config.StringDefault("db.spec", "")
	db, err = sql.Open(driver, connectString)
	if err != nil {
		revel.ERROR.Println("DB open Error", err)
		return err
	}
	return nil
}

func closeDB() (err error) {
	err = db.Close()
	if err != nil {
		revel.ERROR.Println("DB close Error", err)
		return err
	}
	return nil
}

/*
CREATE TABLE public.logs
(
  pk_id integer NOT NULL DEFAULT nextval('logs_pk_id_seq'::regclass), -- первичный ключ
  time_stamp timestamp without time zone, -- дата и время
  app_name character varying(50), -- имя приложения
  pkg_name character varying(50), -- имя пакета
  module_name character varying(50), -- имя модуля
  proc_name character varying(50), -- имя функции/процедуры
  context character varying(100), -- контекст бизнес-логики
  log_text text, -- текст лога
  type character varying(50), -- тип сообщения
  CONSTRAINT "первичный ключ" PRIMARY KEY (pk_id)
)
WITH (
  OIDS=FALSE
);
ALTER TABLE public.logs
  OWNER TO postgres;
COMMENT ON TABLE public.logs
  IS 'таблица логов';
COMMENT ON COLUMN public.logs.pk_id IS 'первичный ключ';
COMMENT ON COLUMN public.logs.time_stamp IS 'дата и время';
COMMENT ON COLUMN public.logs.app_name IS 'имя приложения';
COMMENT ON COLUMN public.logs.pkg_name IS 'имя пакета';
COMMENT ON COLUMN public.logs.module_name IS 'имя модуля';
COMMENT ON COLUMN public.logs.proc_name IS 'имя функции/процедуры';
COMMENT ON COLUMN public.logs.context IS 'контекст бизнес-логики';
COMMENT ON COLUMN public.logs.log_text IS 'текст лога';
COMMENT ON COLUMN public.logs.type IS 'тип сообщения';
*/