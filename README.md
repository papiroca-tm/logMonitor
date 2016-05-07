#logMonitor
log monitor for golang revel project with webix frontend framework
### todo list:
- реализовать вывод строки лога на бэке (функция getLine() уже написана)
 - поправить все запросы к бд
 - код создания таблицы с логами
- реализовать вывод строки лога на фронте
 - новая колонка
- leftSidebar menu на фронте
 - первый пункт меню будет просмотр логов
 - второй пункт меню будет просмотр статистики
- статистика в виде диаграмм:
 - по программе
 - по модулю
 - по функции
 - по типу лога
 - по коду ошибки

##instruction
необходимо переместить файлы проекта согласно следующей структуре:

**app**/
----**controllers**/
--------logMonitor.go
----**services**/
--------**logMonitor**/
------------config.json
------------SLogMonitor.go
**views**/
----**LogMonitor**/
--------Header.html
--------Index.html
**public**/
----**js**/
--------logMonitor.js
----**lib**/
--------**codebase**/ - *folder with webix framework*

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
        "DbUserPassword": "js110682msm",
        "DbHost": "localhost",
        "DbPort": "5432",
        "DbName": "workHelper",
        "DbShema": "logs",
        "DbTable": "logs",
        "Sslmode": "disable"
    }
   
используем в коде

    package controllers

	import (
		"github.com/revel/revel"
		logit "logMonitor/app/services/logMonitor"
	)
	
    func (c SomeController) Action() revel.Result {		
    	logit.TRACE("текст лога", "контекст", "")
    	logit.INFO("текст лога", "контекст", "")
    	logit.WARN("текст лога", "контекст", "")
    	
    	if err != nul {
			logit.ERROR("текст лога", "контекст", "код ошибки")
		}
    }