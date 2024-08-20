var gui

//после загрузки страницы
window.addEventListener('load',function(){
    address = document.getElementById("address").value
    gui = new Guiws(address);

    //---заполнение словаря действий
    gui.functions.set("setValue",function(mes){
        let id= mes[0]
        let value = mes[1]
        let elem = document.getElementById(id)
        if(elem.hasAttribute("value"))
            elem.value=value
        else
            elem.innerText = value
    });
    gui.functions.set("setAttribute",function(mes){
        let id= mes[0]
        let attr = mes[1]
        let value = mes[2]
        let elem = document.getElementById(id)
        elem.setAttribute(attr,value)
    });

    gui.functions.set("getValue",function(mes){
        let id= mes[0]
        let value = document.getElementById(id).value
        gui.go("getValue", id, value)
    });   

    gui.functions.set("getAttribute",function(mes){
        let id= mes[0]
        let attr = mes[1]
        let value = document.getElementById(id).getAttribute(attr)
        gui.go("getAttribute", id, attr, value)
    });
    
    //---подписки на события
    document.addEventListener('click',function(event){
        let id=event.currentTarget.activeElement.id
        if (id!=null)
            gui.go('event', id + '_click');
    });    
    document.addEventListener('change',function(event){
        let id=event.currentTarget.activeElement.id
        if (id!=null)
            gui.go('event', id + '_change');
    });
});

//перед закрытием(обновлением)
window.addEventListener('beforeunload',()=>{
    gui.go("quit")
});


//--- Класс реализующий взаимодействие с gui
class Guiws {
    functions //словарь действий
    ws

    constructor(address){      
        this.functions = new Map()
        this.initChannel(address)
     }

    //---фукционирование вебсокетов
    initChannel(address){
        if(address==''||address==undefined)address="ws://127.0.0.1:8080/gui"
        this.ws = new WebSocket(address)

        if(this.ws != undefined){
            this.ws.addEventListener('message',(mess)=>{
                let data =JSON.parse(mess.data)
                this.functions.get(data.action)(data.arguments)
            })
            this.ws.addEventListener('open',()=>{
                document.getElementById('errorClose').hidden=true
            })
            this.ws.addEventListener('close',()=>{
                document.getElementById('errorClose').hidden=false
            })            
        }
    }

    //---фукция формирования структуры сообщения и отправки на сервер
    go(funcname_and_arguments){
        if(arguments.length<1) return;
        if(this.ws.readyState!=1) return;
        let mess ={}
        mess.action= arguments[0]
        mess.arguments=[]
        for(let i=1;i<arguments.length;i++){
            mess.arguments[i-1]=arguments[i]
        }
        this.ws.send(JSON.stringify(mess))
    }
}