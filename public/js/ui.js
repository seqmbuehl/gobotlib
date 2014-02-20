function Ui (game_container, message_container) {
  this.game_container = game_container;
  this.message_container = message_container;
  this.messages = [];
}

Ui.prototype.getGameContainer = function() {
  return this.game_container;
}


Ui.prototype.resizeContainer = function(GameState) {
  var height_to_width = GameState.World.Size.Y / parseFloat(GameState.World.Size.X);
  // set our game container to the biggest size possible in the same porportions as the server's game
  var height = null;
  var width = null;
  if (GameState.World.Size.Y == GameState.World.Size.X) {
    height = width = Math.min($('body').innerWidth(), $('body').innerHeight())
  } else {
    throw 'TODO: handle portrait and landscape mode';
  }
  // take into account other game elements here
  this.game_container.css({
    width: width,
    height: height
  })
}

Ui.prototype.showMessage = function(message) {
  var $message = $('<div class="message">').text(message).hide();
  this.message_container.prepend($message);
  $message.fadeIn();
  window.setTimeout(function() {
    $message.fadeOut(function() {
      $message.remove();
    });
  }, 3 * 1000);
}
