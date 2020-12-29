# glowstone
Minecraft server proxy

```
//Where :4013 is a port you'd like to listen on for incoming connections
//And :25565, :25566 are downstream servers, if multiple address are provided traffic will be load balanced between servers (round robin)
proxy := glowstone.NewProxy(":4013", ":25565", ":25566")
proxy.Listen()
```
