"use strict";

var config = {
    moduleName: 'logMonitor',
    moduleVer: '0.0.1',    
    defaultHost: 'localhost',
    defaultPort: '9000',    
    webixLocale: 'ru-RU',
    datePickerFormat: '%d.%m.%Y %H:%i:%s',
    dataListColums: [
        { id : "pk_id", header : "#", width : 40, sort : "int", adjust: "data" },
        { id : "time_stamp", header : "Дата и время", adjust: "data", editor:"text" },
        { id : "app_name", header : ["Программа", {content:"selectFilter"}], adjust: "data", editor:"text"},
        { id : "pkg_name", header : ["Пакет", {content:"selectFilter"}], adjust: "data", editor:"text"},
        { id : "module_name", header : ["Модуль", {content:"selectFilter"}], adjust: "data", editor:"text"},
        { id : "proc_name", header : ["функция", {content:"textFilter"}], adjust: "data", editor:"text"},
        { id : "module_line", header : "строка", adjust: "header", editor:"text"},
        { id : "log_context", header : ["Контекст", {content:"textFilter"}], fillspace : 1, editor:"text"},
        { id : "log_text", header : ["Текст лога", {content:"textFilter"}], fillspace : 3, editor:"text"},
        { id : "log_type", header : ["Тип лога", {content:"selectFilter"}], adjust: "header", editor:"text"},
        { id : "err_code", header : ["код ошибки", {content:"textFilter"}], adjust: "header", editor:"text"},
    ]
};

var Views = function () {
    this.mainToolbarView = {
        container: 'mt',        
        view : "toolbar",
        id: 'mainToolbar',        
        cols : [
            { view : "label", label : config.moduleName, width : 150, align : "left" },            
            {
                view: "combo",
                id:"mainMenu",
                width : 160,
                value: "Логи", 
                options: ["Логи", "Статистика"]
            },            
            { view : "label", label : "с:", width : 30, align : "left" },
            { 
                view : "datepicker",
                id: 'datePickerStartID',
                format: webix.Date.dateToStr(config.datePickerFormat), 
                timepicker : true, width : 180, align : "left" 
            },
            { view : "label", label : "по:", width : 30, align : "left" },
            { 
                view : "datepicker",
                id: 'datePickerEndID', 
                format: webix.Date.dateToStr(config.datePickerFormat), 
                timepicker : true, width : 180, align : "left" 
            },
            { view : "button", id : 'btnRequest', value : "сформировать", inputWidth : 180, align: "left" }
            //{ view : "button", id : 'btnSql', value : "SQL-запрос", inputWidth : 180, align: "right" }
        ]
    }; 
    this.mainDatatableView =  {        
        view : "datatable",
        id: 'logDataList',
        container: 'content',
        //autoheight:true,
        //hidden: true,
        resizeColumn: true, select : true, clipboard: true,
        delimiter:{ rows:"\n", cols:"|" },
        scrollX: false, scrollY: true,
        //dragColumn:true, 
        editable:true, editaction:"dblclick",
        columns : config.dataListColums,		
    };    
    this.getDataListID = function () { return 'logDataList'; };    
    this.getErrView = function (msgType, msgText) {
        return { type: msgType, width : 'auto', text : msgText }
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
                $$(dataList).refresh();
                $('#content').height($(window).height() - ($('#mt').height() + 10));
                $$(dataList).adjust();              
            },
            error: function (jqXHR, textStatus, errorThown) {
                console.log('error'); 
                console.error(jqXHR);
				webix.alert(master.views.getErrView('error', jqXHR.responseText));
            }
        });
    };           
};

var Handler = function (master) {
    this.btnRequestClick = function () {
        var dtStart = $$('datePickerStartID').getText();
        var dtEnd = $$('datePickerEndID').getText();
        var url = '/logmonitor/getlogs/' + '?dtStart='+ dtStart + '&dtEnd=' + dtEnd;
        console.log(url);
        master.restApiCtrl.get(url);
    };
    this.menuOnChange = function (newv, oldv) {
        console.log(newv);                    
        switch (newv) {
            case 'Логи':                
                $$('logDataList').show();
                break
            case 'Статистика':
                $$('logDataList').hide();
                break
        }
    };
};

var LogMonitor = function () {
    var self = this;   
    var THIS_NAME = config.moduleName;
    var THIS_VERSION = config.moduleVer;       
    this.getThisName = function () { return THIS_NAME; };        
    this.getThisVersion = function () { return THIS_VERSION; };
    this.views = new Views();
    this.restApiCtrl = new RestApiControler(this);
    this.handler = new Handler(this);
    this.bildHandlers = function () {
        $$('btnRequest').attachEvent("onItemClick", self.handler.btnRequestClick);
        $$('mainMenu').attachEvent("onChange", function(newv, oldv){
            self.handler.menuOnChange(newv, oldv)
        });
    };
    this.init = function () {
        webix.i18n.setLocale(config.webixLocale);
        webix.ui(this.views.mainToolbarView);
        webix.ui( this.views.mainDatatableView );             
        var dt = new Date();
        $$('datePickerEndID').setValue(dt);
        dt.setHours(dt.getHours() - 1);
        $$('datePickerStartID').setValue(dt);
    };    
};

var lm = new LogMonitor();

function onDocumentReady () {    
    console.log(lm.getThisName() + ' ver:' + lm.getThisVersion());    
    lm.init();
    lm.bildHandlers();
};








