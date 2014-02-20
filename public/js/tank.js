function Tank (game, player) {
  this.game = game;
  this.player = player;
}

Tank.prototype.getSizes = function() {
  this.playerY = this.game.y(this.player.Position.Y);
  this.playerX = this.game.x(this.player.Position.X);
  switch (this.player.Vector.Direction) {
    case this.game.Directions.up:
    case this.game.Directions.down:
      this.playerHeight = this.game.y(this.game.constants.PlayerSize.Y);
      this.playerWidth = this.game.x(this.game.constants.PlayerSize.X);
      break;
    case this.game.Directions.left:
    case this.game.Directions.right:
      this.playerHeight = this.game.y(this.game.constants.PlayerSize.X);
      this.playerWidth = this.game.x(this.game.constants.PlayerSize.Y);
      break;
  }
  // Lame, 2 switch statements to set this up?
  switch (this.player.Vector.Direction) {
  case this.game.Directions.down:
  case this.game.Directions.up:
    this.bodyHeight = this.playerHeight * .8;
    this.bodyWidth = this.playerWidth;
    this.barrelHeight = this.playerHeight * .6;
    this.barrelWidth = this.game.x(this.game.constants.BulletSize.X);
    this.bodyY = this.playerHeight - this.bodyHeight;
    this.bodyX = 0;
    this.barrelY = 0;
    this.barrelX = (this.playerWidth/2) - (this.barrelWidth/2);
    break;
  case this.game.Directions.right:
  case this.game.Directions.left:
    this.bodyHeight = this.playerHeight;
    this.bodyWidth = this.playerWidth * .8;
    this.barrelHeight = this.game.y(this.game.constants.BulletSize.Y);
    this.barrelWidth = this.playerWidth * .6;
    this.bodyY = 0;
    this.bodyX = this.playerWidth - this.bodyWidth;
    this.barrelY = (this.playerHeight/2) - (this.barrelHeight/2);
    this.barrelX = 0;
    break;
  }
  switch (this.player.Vector.Direction) {
  case this.game.Directions.down:
    this.bodyY = 0;
    this.barrelY = this.playerHeight - this.barrelHeight;
    break;
  case this.game.Directions.right:
    this.bodyX = 0;
    this.barrelX = this.playerWidth - this.barrelWidth;
    break;
  }
}
