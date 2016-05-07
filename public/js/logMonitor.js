"use strict";

var config = {
    moduleName: 'logMonitor',
    moduleVer: '0.0.1',
    
    defaultHost: 'localhost',
    defaultPort: '9000',
    
    webixLocale: 'ru-RU',    
    
    toolBarID: 'mainToolbar',
    datePickerStartID: 'datePickerStartID',
    datePickerEndID: 'datePickerEndID',
    datePickerFormat: '%d.%m.%Y %H:%i:%s',
    
    dataListID: 'logDataList',
    dataListColums: [
        { id : "pk_id", header : "#", width : 40, sort : "int", adjust: "data" },
        { id : "time_stamp", header : "Дата и время", adjust: "data", editor:"text" },
        { id : "app_name", header : ["Программа", {content:"selectFilter"}], adjust: "data", editor:"text"},
        { id : "pkg_name", header : ["Пакет", {content:"selectFilter"}], adjust: "data", editor:"text"},
        { id : "module_name", header : ["Модуль", {content:"selectFilter"}], adjust: "data", editor:"text"},
        { id : "proc_name", header : ["функция", {content:"textFilter"}], adjust: "data", editor:"text"},
        { id : "log_context", header : ["Контекст", {content:"textFilter"}], fillspace : 1, editor:"text"},
        { id : "log_text", header : ["Текст лога", {content:"textFilter"}], fillspace : 3, editor:"text"},
        { id : "log_type", header : ["Тип лога", {content:"selectFilter"}], adjust: "header", editor:"text"},
        { id : "err_code", header : ["код ошибки", {content:"textFilter"}], adjust: "header", editor:"text"},
    ],
    
    btnRequestID: 'btnRequest',
    btnSqlID: "btnSql",
};

var LogMonitor = function () {    
    
    var THIS_NAME = config.moduleName;
    var THIS_VERSION = config.moduleVer;       
    this.getThisName = function () { return THIS_NAME; };        
    this.getThisVersion = function () { return THIS_VERSION; };   
    
    this.init = function () {        
        webix.i18n.setLocale(config.webixLocale);
        this.views = new Views();
        webix.ui(this.views.getCompletedDOM());               
		var dt = new Date();
		$$(config.datePickerEndID).setValue(dt);
		dt.setHours(dt.getHours() - 1);
		$$(config.datePickerStartID).setValue(dt);                
    };
};

var RestApiControler = function (master, host, port) {
    	
    var ctrl = this;    
    ctrl.host = (host !== undefined) ? host : config.defaultHost;
    ctrl.port = (port !== undefined) ? port : config.defaultPort;
        
    ctrl.get = function(url) {
        return $.ajax({
            type: 'GET',
            host: ctrl.host,
            port: ctrl.port,
            url: url,
            success: function (data, textStatus, jqXHR) {
                var dataList = master.views.getDataListID();
                $$(dataList).clearAll();
                $$(dataList).parse(data);                
            },
            error: function (jqXHR, textStatus, errorThown) {
                console.error(jqXHR);
				webix.alert(master.views.getErrView('error', jqXHR.responseText));
            }
        });
    }; 
          
};

var Views = function () {    
    
    this.mainToolbarView = {        
        view : "toolbar",
        id: config.toolBarID,
        cols : [
            { view : "label", label : config.moduleName, width : 150, align : "left" },
            { view : "label", label : "с:", width : 30, align : "left" },
            { 
                view : "datepicker",
                id: config.datePickerStartID,
                format: webix.Date.dateToStr(config.datePickerFormat), 
                timepicker : true, width : 200, align : "left" 
            },
            { view : "label", label : "по:", width : 30, align : "left" },
            { 
                view : "datepicker",
                id: config.datePickerEndID, 
                format: webix.Date.dateToStr(config.datePickerFormat), 
                timepicker : true, width : 200, align : "left" 
            },
            { view : "button", id : config.btnRequestID, value : "сформировать", inputWidth : 180, align: "left" },
            { view : "button", id : config.btnSqlID, value : "SQL-запрос", inputWidth : 180, align: "right" }
        ]
    };
    
    this.mainDatatableView =  {        
        view : "datatable",
        id: config.dataListID,
        resizeColumn: true, select : true, clipboard: true,
        delimiter:{
            rows:"\n", // the rows delimiter
            cols:"|"   // the columns delimiter
        },
        scrollX:true, scrollY:true, dragColumn:true, editable:true, editaction:"dblclick",
        columns : config.dataListColums,		
    };    
    
    this.getDataListID = function () { return config.dataListID; };
    
    this.getErrView = function (msgType, msgText) {
        return { type: msgType, width : 'auto', text : msgText }
    };
    
    this.getCompletedDOM = function () {
        return { rows : [this.mainToolbarView, this.mainDatatableView] };    
    };
        
};

var lm = new LogMonitor();
var restApiCtrl = new RestApiControler(lm);

function onDocumentReady () {    
    console.log(lm.getThisName() + ' ver:' + lm.getThisVersion());    
    lm.init();
    bildHandlers();
};

function btnRequestClick () {
    var dtStart = $$(config.datePickerStartID).getText();
    var dtEnd = $$(config.datePickerEndID).getText();
    var url = '/logmonitor/getlogs/' + '?dtStart='+ dtStart + '&dtEnd=' + dtEnd;
    restApiCtrl.get(url);
};
    
function bildHandlers () {
    $$(config.btnRequestID).attachEvent("onItemClick", btnRequestClick);    
};






