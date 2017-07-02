var MYGAME = {
    killStages: ["", "head", "body", "lefth", "righth", "leftf", "rightf"],
    hint:       '',
    curr:       [],
    missed:     [],
    game:       0, 
    auth:       '',
    execute:      function() {
        try {
            if (this.missed.length > 0) {
                document.getElementById(this.killStages[this.missed.length]).style.display = "block";
            }
        } catch(err) {
            newGame();
        }
    },
    reset:      function() {
        var sk = this.killStages.slice(1, this.killStages.length)
        for (i in sk) {
            document.getElementById(sk[i]).style.display = "none";
        }
    }
}

var wsocket;

document.getElementById("myButton").onclick = function() {
    guess = document.getElementById("letter").value;
    game = {Cmd: "P1T", Play: guess, Gid: MYGAME.game, Auth: MYGAME.auth};
    console.log("Sending: " + JSON.stringify(game));
    wsocket.send(JSON.stringify(game));
};

function finalGuess() {
    var answer;
    stf = function() { answer = prompt("Final guess: "); };
    do {         // Oh why do you do this iOS Safari?
        setTimeout(stf(), 0);
    } while(answer.length < 1);
    console.log(answer);
    game = {Cmd: "FIN", Play: answer, Gid: MYGAME.game, Auth: MYGAME.auth};
    console.log("Sending: " + JSON.stringify(game));
    wsocket.send(JSON.stringify(game));
};

$("#letter").keyup(function(event){
    if(event.keyCode == 13){
        $("#myButton").click();
    }
});

if ("WebSocket" in window) {
    url = "hangman.example.com/wshangman";
    //var isSafari = navigator.vendor.indexOf("Apple")==0 && /\sSafari\//.test(navigator.userAgent);
    //wsocket = isSafari ? new WebSocket("ws://" + url) : new WebSocket("wss://" + url);
    wsocket = new WebSocket("ws://" + url);
    wsocket.onopen = function() {
        newGame();
    };

    wsocket.onmessage = function (event) { 
        var msg = JSON.parse(event.data);
        console.log("Received: " + event.data);
        switch (msg.Cmd) {
            case "NEW":
                MYGAME.hint = msg.Hint;
                MYGAME.auth = msg.Cred; 
                MYGAME.game = msg.Game;
                document.getElementById("hint").innerHTML = MYGAME.hint;
                break;
        }
        MYGAME.curr = function() { 
            var tmp = [];
            for (i in msg.Curr) {
                tmp[i] = String.fromCharCode(msg.Curr[i])
            }
            return tmp;
        }();
        MYGAME.missed = msg.Missed;
        document.getElementById("abstract").innerHTML = function() {
            var tmp = [];
            for (i in MYGAME.curr) {
                if (MYGAME.curr[i] === "\00") {
                    tmp[i] = '\u2610';
                } else {
                    tmp[i] = MYGAME.curr[i];
                }
            }
            return tmp;
        }().join(' ');
        $("#abstract").effect("shake");
        $("#letter").val('');
        MYGAME.execute();
    };
    wsocket.onclose = function() { 
        console.log("Connection is closed..."); 
    };
} else {
    var wsocket = ''
}

function newGame() {
    wsocket.send(JSON.stringify({Cmd: "NEW"}));
    MYGAME.reset();
}

document.getElementById("newGame").onclick = newGame;
