# glowstone 🌟
### Minecraft server proxy 💫

> :4013 is a port you'd like to listen on for incoming connections

> :25565, :25566 are downstream servers, if multiple addresses are provided traffic will be load balanced between servers (round robin)
```
proxy := glowstone.NewProxy(":4013", ":25565", ":25566")
proxy.Listen()
```

> Add new server

![](https://i.imgur.com/YrBW4bP.png)

> Your server should be visible and repspond to client's healthchecks

![](https://i.imgur.com/Itc5rZv.png)

>If your downstream servers are reachable, you should be able to play 🍾

![](https://i.imgur.com/HZhy7HJ.png)


![](https://i.imgur.com/e7xKx8A.png)
