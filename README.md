# glowstone 🌟
### Minecraft server proxy 💫

> :4013 is a port you'd like to listen on for incoming connections

> :25565, :25566 are downstream servers, if multiple addresses are provided traffic will be load balanced between servers (round robin)
```
proxy := glowstone.NewProxy(":4013", ":25565", ":25566")
proxy.Listen()
```
