# glowstone 🌟

### glowstone allows other people to join your localhost games 🎢

### if you're familiar with services like https://aternos.org/ you know that queues can get ridicolous
### with glowstone you can just host your own server and let other people join it via glowstone proxy/tcp mux
### for now you have to host the proxy yourself, but even smallest vm will suffice which is not true for mc servers

## Try it yourself
#### Client, this app should be deployed on the same machine as your minecraft server
```
m := glowstone.NewMux()

//Tries to connect with the other mux on port 8000
if err := m.Dial(":8000"); err != nil {
  log.Fatal(err.Error())
}

//When connection is successfull, we can start receiving packets
if err := m.Recv(); err != nil {
  log.Fatal(err.Error())
}
```
#### Server, this usually sits somewhere in the clould or on any other machine exposed to external internet traffic
```
m := glowstone.NewMux()

//Listens for other mux to connect, 
if err := m.ListenMux(":8000"); err != nil {
  log.Fatal(err.Error())
}
//Listens for other players to connect, after that it handles each connection and forwards it to your minecraft server
if err := m.Listen(":9000"); err != nil {
  log.Fatal(err.Error())
}
```

## Architecture 🏗
![](https://i.imgur.com/8oID1nK.png)
