#logMonitor
log monitor for golang revel project with webix frontend framework
### todo list:
- сделать отдельные инпуты под выбор времени
- статистика в виде диаграмм:
 - по программе
 - по модулю
 - по функции
 - по типу лога
 - по коду ошибки
- решить проблему пути к конфигурационному файлу при запуске revel run из разных директорий

##instruction
необходимо переместить файлы проекта согласно следующей структуре:

**app/controllers/**logMonitor.go

**app/services/logMonitor/**config.json

**app/services/logMonitor/**SLogMonitor.go

**views/LogMonitor/**Header.html

**views/LogMonitor/**Index.html

**public/js/**logMonitor.js

**public/lib/codebase/** - folder with webix framework



Прописываем роуты в **conf/routes**

	GET	/logmonitor	                LogMonitor.Index
    GET	/logmonitor/getlogs/     	LogMonitor.GetLogs

настраиваем **app/services/logMonitor/config.json**

    {
        "DateTimeFormatString": "02.01.2006 15:04:05",
    	"DbDateTimeFormatString": "2006-01-02 15:04:05.999999999",    
    	"StackLevel": 2,    
    	"DbDriver": "postgres",
        "DbUser": "postgres",
        "DbUserPassword": "password",
        "DbHost": "localhost",
        "DbPort": "5432",
        "DbName": "dbName",
        "DbShema": "logs",
        "DbTable": "logs",
        "Sslmode": "disable"
    }
   
используем в коде

    package controllers

	import (
		"github.com/revel/revel"
		logit "appName/app/services/logMonitor"
	)
	
    func (c SomeController) Action() revel.Result {		
    	logit.TRACE("текст лога", "контекст", "")
    	logit.INFO("текст лога", "контекст", "")
    	logit.WARN("текст лога", "контекст", "")
    	
    	if err != nul {
			logit.ERROR("текст лога", "контекст", "код ошибки")
		}
    }
