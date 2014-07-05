function Wall (host, name, onCloseHandler, onMessageHandler) {
    this.host = host;
    this.conn = null;
    this.name = name;
    this.onClose = onCloseHandler;
    this.onMessage = onMessageHandler;
}

Wall.prototype = {
    constructor: Wall,
    initSocket: function() {
        console.log("init");
        this.conn = new WebSocket("ws://" + this.host + "/ws/" + this.name);
        this.conn.onopen = function(){ console.log("ws connected."); }
        this.conn.onclose = this.onClose;
        this.conn.onmessage = this.onMessage;
    },
    sendCommand: function(cmdJson) { 
        if(!this.conn) return false;
        this.conn.send(cmdJson);
    },
    addCard: function(text, name) { 
        console.log("addCard");
        this.sendCommand(JSON.stringify({
            "Cmd": "add_card",
            "Data": {"Text": text, "Name": name}
        })); 
    },
    plusCard: function(id, name) { 
        console.log("plusCard");
        this.sendCommand(JSON.stringify({
            "Cmd": "plus_card",
            "Data": {"Id": "" + id, "Name": name}
        })); 
    },
    removeCard: function(id, name) { },
    moveCard: function(id, name, x, y) { 
        console.log("moveCard");
        this.sendCommand(JSON.stringify({
            "Cmd": "move_card",
            "Data": {"Id": "" + id, "Name": name, "X": x, "Y": y}
        })); 
    }
}

