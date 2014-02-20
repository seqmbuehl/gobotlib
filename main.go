// 2014, Nathan Tsoi
// Note: This needs massive cleanup

package main

import (
  "os"
  "fmt"
  "github.com/googollee/go-socket.io"
  "log"
  "net/http"
  "time"
  "encoding/json"
  "strings"
  "strconv"
  "math/rand"
  "sync"
)
type Color string
const (
  colorless Color = ""
  orange Color = "ebbc8d"
  green Color = "aad98d"
  blue Color = "8daceb"
  red Color = "eb8d8d"
  yellow Color = "d9d682"
  turquoise Color = "95d1cb"
  purple Color = "af93c9"
)

type Action int
const (
  fire Action = iota+1
)

type Direction int
const (
  up Direction = 1+iota
  down
  left
  right
)

type RoundType int
const (
  bullet RoundType = iota+1
  missle
  bomb
)

type Vector2 struct {
  X int
  Y int
}

type Vector struct {
  Direction Direction
  Velocity int
}

type Bullet struct {
  RoundType RoundType
  Vector *Vector
  Position *Vector2
}

type Player struct {
  Name string
  Points int
  Color Color
  Vector *Vector
  Position *Vector2
  Bullets map[string]Bullet
  ns *socketio.NameSpace
}

type World struct {
  Size *Vector2
}

type Constants struct {
  BulletSize Vector2
  PlayerSize Vector2
}

type Game struct {
  Players map[string]Player
  World *World
  Constants Constants
}

////
// Globals
////
var Colors = [...]Color{orange, green, blue, red, yellow, turquoise, purple}
var randSource = rand.New(rand.NewSource(time.Now().UnixNano()))
// Vertically Oriented Player Size, bullets should always be square
var constants = Constants{PlayerSize: Vector2{X: 40, Y: 64}, BulletSize: Vector2{X: 10, Y: 10}}
var maxPlayers = 7
var playerVelocity = 3
var bulletVelocity = 6
var maxBullets = 5
var playersMutex = sync.Mutex{}

////
// Game Methods
////
func (g *Game) addPlayer(id string, name string) *Player {
  //fmt.Println("adding ", name, " ", strconv.Itoa(len(g.Players)), " already exist")
  player := newPlayer(g, name)
  //fmt.Printf(" new player: %#v\n", player)
  if len(g.Players) > 0 {
    colors_avail := make(map[Color]bool)
    for _, c := range Colors {
      colors_avail[c] = true
    }
    for _, p := range g.Players {
      delete(colors_avail, p.Color)
    }
    for c := range colors_avail {
      player.Color = c
      break
    }
  } else {
    player.Color = Colors[0]
  }
  g.Players[id] = player
  //fmt.Printf("Adding '%s'\n", name)
  return &player
}
func (g *Game) respawnPlayer(id string) {
  if player, ok := g.Players[id]; ok {
    player.Vector.Direction = up
    player.Position = g.newPosition()
    g.Players[id] = player
  }
}
func (g *Game) randPosition() *Vector2 {
  return &Vector2{randSource.Intn(g.World.Size.X-constants.PlayerSize.X), randSource.Intn(g.World.Size.Y-constants.PlayerSize.Y)}
}
func (g *Game) newPosition() *Vector2 {
  intersection := false
  v := g.randPosition()
  p := Player{Position: v, Vector: &Vector{Direction: up}}
  for {
    intersection = false
    for _, player := range g.Players {
      // find a new spot if this player overlaps any other player
      if player.IntersectsPlayer(&p) {
        //fmt.Printf("%#v intersects %#v", player, p)
        intersection = true
        break
      }
      if !intersection {
        for _, bullet := range player.Bullets {
          // find a new spot if this player overlaps any other player's bullets
          if p.IntersectsBullet(&bullet) {
            //fmt.Printf("%#v intersects %#v", p, bullet)
            intersection = true
            break
          }
        }
      }
    }
    if !intersection {
      break
    }
    v = g.randPosition()
    p.Position = v
  }
  return v;
}


////
// Player Methods
////
// Get the player's right edge x coord
func (p *Player) X2() int {
  switch p.Vector.Direction {
  case up:
  case down:
    return constants.PlayerSize.X + p.Position.X
  }
  return constants.PlayerSize.Y + p.Position.X
}
// Get the player's bottom edge y coord
func (p *Player) Y2() int {
  switch p.Vector.Direction {
  case up:
  case down:
    return constants.PlayerSize.X + p.Position.Y
  }
  return constants.PlayerSize.Y + p.Position.Y
}

// There is probably a good way to combine the next two methods
func (a *Player) IntersectsPlayer(b *Player) bool {
  // non-intersection is
  // a is to the left of b or
  // a is to the right of b or
  // a is below b or
  // a is above b
  //  which yields
  // !(a.X2() < b.Position.X || a.Position.X > b.X2() || a.Position.Y > b.Y2() || a.Y2() < b.Position.Y)
  // demorgan's law yields
  return a.X2() > b.Position.X && a.Position.X < b.X2() && a.Position.Y < b.Y2() && a.Y2() > b.Position.Y
}
func (a *Player) IntersectsBullet(b *Bullet) bool {
  // non-intersection is
  // a is to the left of b or
  // a is to the right of b or
  // a is below b or
  // a is above b
  //  which yields
  // !(a.X2() < b.Position.X || a.Position.X > b.X2() || a.Position.Y > b.Y2() || a.Y2() < b.Position.Y)
  // demorgan's law yields
  return a.X2() > b.Position.X && a.Position.X < b.X2() && a.Position.Y < b.Y2() && a.Y2() > b.Position.Y
}

// Get the bullets's right side x coord or bottom edge y coord
// Assume bullets are square
func (b *Bullet) X2() int {
  return constants.BulletSize.X + b.Position.X
}
func (b *Bullet) Y2() int {
  return constants.BulletSize.Y + b.Position.Y
}


////
// Maker methods
////
func newBullet() Bullet {
  return Bullet{RoundType: bullet, Vector: &Vector{Velocity: bulletVelocity}}
}

func newPlayer(g *Game, name string) Player {
  return Player{Name: name, Points: 0, Vector: &Vector{Direction: up, Velocity: playerVelocity}, Position: g.newPosition(), Bullets: make(map[string]Bullet)}
}

func newWorld() *World {
  return &World{Size: &Vector2{X: 500, Y: 500}}
}

func NewGame() *Game {
  return &Game{Players: make(map[string]Player), World: newWorld(), Constants: constants}
}


////
// Main
////
func main() {
  sock_config := &socketio.Config{}
  sock_config.HeartbeatTimeout = 2
  sock_config.ClosingTimeout = 4

  game := NewGame()
  sio := socketio.NewSocketIOServer(sock_config)
  gameLoop := func () {
    for {
      //playersMutex.Lock()
      // Move each player and their bullets
      for _, player := range game.Players {
        //fmt.Printf("player: %#v\n", player)
        // Get the current width or height depending on orientation
        currentXSize := constants.PlayerSize.X
        currentYSize := constants.PlayerSize.Y
        // Move
        switch player.Vector.Direction {
          case up:
            player.Position.Y = player.Position.Y - player.Vector.Velocity
          case down:
            player.Position.Y = player.Position.Y + player.Vector.Velocity
          case left:
            player.Position.X = player.Position.X - player.Vector.Velocity
            currentXSize = constants.PlayerSize.Y
            currentYSize = constants.PlayerSize.X
          case right:
            player.Position.X = player.Position.X + player.Vector.Velocity
            currentXSize = constants.PlayerSize.Y
            currentYSize = constants.PlayerSize.X
        }
        // lame, no min/max for int
        if player.Position.Y < 0 {
          player.Position.Y = 0
        } else if player.Position.Y + currentYSize > game.World.Size.Y {
          player.Position.Y = game.World.Size.Y - currentYSize
        }
        if player.Position.X < 0 {
          player.Position.X = 0
        } else if player.Position.X + currentXSize > game.World.Size.X {
          player.Position.X = game.World.Size.X - currentXSize
        }
        // Bullets
        for id, bullet := range player.Bullets {
          // If the bullet is new
          if bullet.Position == nil {
            x := 0
            y := 0
            // Set the position to the tip of the barrel
            switch player.Vector.Direction {
            case up:
              x = player.Position.X + (constants.PlayerSize.X / 2) - (constants.BulletSize.X / 2)
              y = player.Position.Y - constants.BulletSize.Y
              break;
            case down:
              x = player.Position.X + (constants.PlayerSize.X / 2) - (constants.BulletSize.X / 2)
              y = player.Position.Y + constants.PlayerSize.Y
              break;
            case left:
              x = player.Position.X - constants.BulletSize.X
              y = player.Position.Y + (constants.PlayerSize.X / 2) - (constants.BulletSize.X / 2)
              break;
            case right:
              x = player.Position.X + constants.PlayerSize.Y
              y = player.Position.Y + (constants.PlayerSize.X / 2) - (constants.BulletSize.X / 2)
              break;
            }
            player.Bullets[id] = Bullet{RoundType: bullet.RoundType, Position: &Vector2{X: x, Y: y}, Vector: &Vector{Velocity: bullet.Vector.Velocity, Direction: player.Vector.Direction}}
          // Else, move the bullet
          } else {
            switch bullet.Vector.Direction {
              case up:
                bullet.Position.Y = bullet.Position.Y - bullet.Vector.Velocity
              case down:
                bullet.Position.Y = bullet.Position.Y + bullet.Vector.Velocity
              case left:
                bullet.Position.X = bullet.Position.X - bullet.Vector.Velocity
              case right:
                bullet.Position.X = bullet.Position.X + bullet.Vector.Velocity
            }
            // discard any bullets off screen
            if !(bullet.Position.Y > 0 && bullet.Position.Y + constants.BulletSize.Y < game.World.Size.Y && bullet.Position.X > 0 && bullet.Position.X + constants.BulletSize.X < game.World.Size.X) {
              delete(player.Bullets, id)
            }
          }
        }
      }
      respawnPlayers := make(map[string]Player)
      // Check for collisions
      for playerId, player := range game.Players {
        if _, ok := respawnPlayers[playerId]; ok {
          continue
        }
        // if this player overlaps with any other player
        for player2Id, player2 := range game.Players {
          if _, ok := respawnPlayers[player2Id]; ok  {
            continue
          }
          if playerId == player2Id {
            continue
          }
          if player.IntersectsPlayer(&player2) {
            // respawn the players
            respawnPlayers[playerId] = player
            respawnPlayers[player2Id] = player2
            // report a crash
            //fmt.Println(player.Name, " and ", player2.Name, "crashed")
            json_bytes, _ := json.Marshal(map[string]interface{}{
              "player1": player,
              "player2": player2,
            })
            sio.Broadcast("gameevent", "crash", string(json_bytes))
          }
        }
        // For every bullet
        for id, bullet := range player.Bullets {
          // For every other player
          for player2Id, player2 := range game.Players {
            if _, ok := respawnPlayers[player2Id]; ok {
              continue
            }
            if playerId == player2Id{
              continue
            }
            // if this bullet overlaps the other player
            if player2.IntersectsBullet(&bullet) {
              // report a kill
              //fmt.Println(player.Name, " killed ", player2.Name)
              json_bytes, _ := json.Marshal(map[string]interface{}{
                "player1": player,
                "player2": player2,
              })
              sio.Broadcast("gameevent", "kill", string(json_bytes))
              // remove the bullet
              delete(player.Bullets, id)
              // respawn the player
              respawnPlayers[player2Id] = player2
              // increment the bullet's player's score by 1 point
              player.Points += 1
              game.Players[playerId] = player
            }
          }
        }
      }
      for id, _ := range respawnPlayers {
        game.respawnPlayer(id)
      }
      //playersMutex.Unlock()

      // Send out the new game state
      json_gamestate, err := json.Marshal(game)
      if err != nil {
        fmt.Println("error:", err)
      }
      sio.Broadcast("gamestate", string(json_gamestate))

      // Run at 30 fps
      time.Sleep((1000 / 30) * time.Millisecond)

      // Debug mode
      //time.Sleep(200 * time.Millisecond)
      //fmt.Printf("\n%s\n", json_gamestate)
    }
  }
  go gameLoop()

  // Handler for new connections, also adds socket.io event handlers
  sio.On("connect", func (ns *socketio.NameSpace) {
    numPlayers := len(game.Players)
    //fmt.Println("numPlayers: ", numPlayers)
    if numPlayers >= maxPlayers {
      ns.Emit("failure", "game is full")
      return
    }
    if p, ok := game.Players[ns.Id()]; ok {
      fmt.Println(p.Name, " reconnected")
    } else {
      name := strings.Join([]string{"Player", strconv.Itoa(len(game.Players) + 1)}, " ")
      //playersMutex.Lock()
      game.addPlayer(ns.Id(), name)
      //playersMutex.Unlock()
    }
    fmt.Println(game.Players[ns.Id()].Name, " connected:", ns.Id(), " in channel ", ns.Endpoint())
  })
  sio.On("disconnect", func (ns *socketio.NameSpace) {
    fmt.Println(game.Players[ns.Id()].Name, " disconnected:", ns.Id(), " in channel ", ns.Endpoint())
    //playersMutex.Lock()
    delete(game.Players, ns.Id())
    //playersMutex.Unlock()
  })
  sio.On("command", func (ns *socketio.NameSpace, command string, value string) {
    //playersMutex.Lock()
    player := game.Players[ns.Id()]
    //fmt.Printf("%s (%s) '%s': '%s'\n", player.Name, ns.Id(), command, value)
    switch command {
    case `set_name`:
      indexed_name := value
      if len(game.Players) > 1 {
        clear := false
        for i := 0; (i < maxPlayers) && !clear; i++ {
          for _, p := range game.Players {
            if p.Name == indexed_name {
              indexed_name = strings.Join([]string{value, strconv.Itoa(i+1)}, " ")
              break
            }
            clear = true
          }
        }
      }
      //fmt.Println(" resolved name: ", indexed_name)
      player.Name = indexed_name
      game.Players[ns.Id()] = player
    case `direction`:
      direction, _ := strconv.ParseInt(value, 10, 34)
      player.Vector.Direction = Direction(direction)
    case `fire`:
      if len(player.Bullets) < maxBullets {
        for i := 0; i < maxBullets; i++ {
          k := strconv.Itoa(i)
          if _, ok := player.Bullets[k]; ok {
            continue
          }
          player.Bullets[k] = newBullet()
          break
        }
      }
    }
    game.Players[ns.Id()] = player
    //playersMutex.Unlock()
  })
  //this will serve a http static file server
  sio.Handle("/", http.FileServer(http.Dir("./public/")))
  //startup the server
  port := os.Getenv("PORT")
  if port == "" {
    port = "80"
  }
  port = ":" + port
  log.Fatal(http.ListenAndServe(port, sio))
}
