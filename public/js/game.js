function Game (constants, container, score_container, socket) {
  this.Directions = {
    up: 1,
    down: 2,
    left: 3,
    right: 4,
  }
  this.RoundTypes = {
    bullet: 1,
  }
  // IO
  this.socket = socket;
  this.id = socket.socket.sessionid;

  // Container element
  this.$container = container;
  this.$score_container = score_container;

  // World Portion of the Game State
  this.world = null;
  // Constants
  this.constants = constants;

  // Drawing stuffs
  this.svg = null;
  this.map = null;
  this.players = null;
  this.x = null;
  this.y = null;

  this.bindEvents();
}

Game.prototype.bindEvents = function() {
  var self = this;
  $(document).keypress(function(evt) {
    if ( event.keyCode == 'w'.charCodeAt(0) ) {
      self.socket.emit('command', 'direction', ""+self.Directions.up);
      console.log('up');
    } else if ( event.keyCode == 's'.charCodeAt(0) ) {
      self.socket.emit('command', 'direction', ""+self.Directions.down);
      console.log('down');
    } else if ( event.keyCode == 'a'.charCodeAt(0) ) {
      self.socket.emit('command', 'direction', ""+self.Directions.left);
      console.log('left');
    } else if ( event.keyCode == 'd'.charCodeAt(0) ) {
      self.socket.emit('command', 'direction', ""+self.Directions.right);
      console.log('right');
    } else if ( event.keyCode == 32 ) {
      self.socket.emit('command', 'fire', ""+self.RoundTypes.bullet);
      console.log('fire!');
    }
  });
}

Game.prototype.buildMap = function(GameState) {
  this.world = GameState.World
  if (!this.svg) {
    this.svg = d3.select(this.$container[0]).append('svg')
  }
  if (!this.map) {
    this.map = this.svg.append("rect")
      .attr('x', 0)
      .attr('y', 0)
      .attr('fill', '#fefefe')
      .attr('stroke', '#999999')
  }
  if (!this.players) {
    this.players = this.svg.append('g')
      .attr('class', 'players')
  }
  // Always update the sizes
  height = this.$container.innerHeight();
  this.y = d3.scale.linear().range([0,height])
  this.y.domain([0, this.world.Size.Y])
  width = this.$container.innerWidth();
  this.x = d3.scale.linear().range([0,width])
  this.x.domain([0, this.world.Size.X])
  this.svg
      .attr('height', height)
      .attr('width', width)
  this.map
      .attr('height', height)
      .attr('width', width)
  return this.map;
};

Game.prototype.draw = function(GameState) {
  var self = this;
  this.buildMap(GameState);
  // compute the tank for each player
  var tank = null;
  // clean the board
  this.players.selectAll('*').remove()
  // rebuild it
  this.players
    .selectAll('rect')
    .data(d3.entries(GameState.Players))
    .enter()
    .append('g').attr('class', 'player')
    .each(function(d) {
      //console.log("player: ", d);
      playerG = d3.select(this);
      tankG = playerG.append('g').attr('class', 'tank')
        .attr('transform', function(d){
          // New take for this player
          tank = new Tank(self, d.value);
          tank.getSizes();
          return 'translate(' + self.x(d.value.Position.X) + ',' + self.y(d.value.Position.Y) + ')';
        });
      tankG.append('rect')
        .attr('class', 'body')
          .attr('fill', function(d){
            return '#' + d.value.Color;
          })
          .attr('x', tank.bodyX)
          .attr('y', tank.bodyY)
          .attr('height', tank.bodyHeight)
          .attr('width', tank.bodyWidth)
      tankG.append('rect')
        .attr('class', 'barrel')
          .attr('fill', function(d){
            return '#999999';// + d.value.Color;
          })
          .attr('x', tank.barrelX)
          .attr('y', tank.barrelY)
          .attr('height', tank.barrelHeight)
          .attr('width', tank.barrelWidth);
      // For each bullet
      for (i in d.value.Bullets) {
        var bullet = d.value.Bullets[i];
        playerG.append('rect')
          .attr('class', 'bullet')
          .attr('fill', function(d){
            return '#' + d.value.Color;
          })
          .attr('x', self.x(bullet.Position.X))
          .attr('y', self.y(bullet.Position.Y))
          .attr('height', self.x(self.constants.BulletSize.Y))
          .attr('width', self.y(self.constants.BulletSize.X));
      }
    });
}

Game.prototype.score = function(GameState) {
  var self = this;
  var $c = this.$score_container;
  $c.empty();
  for (i in GameState.Players) {
    var player = GameState.Players[i];
    var font_weight = (i == self.id) ? 800 : 400;
    $c.append($('<div>')
      .css({ color: '#' + player.Color, 'font-weight': font_weight })
      .text(player.Name + ' | score: ' + player.Points)
    )
  }
}

