$(function(){
  var ui = new Ui($('.game-container'), $('.message-container'));
  var game = null;
  var socket = io.connect();

  socket.on("connect", function(){
    console.log('connected');
    socket.emit('command', 'set_name', 'Nathan');
  });
  socket.on("gamestate", function(state){
    var GameState = JSON.parse(state);
    ui.resizeContainer(GameState);
    if (!game) {
      game = new Game(GameState.Constants, ui.getGameContainer(), $('.score-container'), socket);
    }
    game.draw(GameState);
    game.score(GameState);
  });
  socket.on("disconnect", function() {
    console.log("disconnected");
  });
  socket.on("gameevent", function(type, info){
    var json = JSON.parse(info);
    console.log("gameevent", type, JSON.stringify(json));
    switch (type) {
      case 'crash':
        ui.showMessage(json.player1.Name + ' and ' + json.player2.Name + ' had a crash');
        break;
      case 'kill':
        ui.showMessage(json.player1.Name + ' took out ' + json.player2.Name);
        break;
    }
  });
});