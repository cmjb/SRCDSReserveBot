# SRCDSReserveBot

This is a bot written in Go that will let you reserve game servers (Speficially HLDS and SRCDS servers) based on the Steam Ids added via discord servers.

Configuration in `config.yml`:

```yml
token: "" # Discord Token
channel: "" # Discord Channel
lobby_request_channel: "" # Discord Request Channel
lobby_status_channel: "" # Discord Debug Channel
moderator_id: "" # Discord Admin ID
guild_id: "" # Discord Server ID
rcon_password: "" #SRCDS Password
```

Database configuration can be added via `pg.go`, under [`Get_DB`](https://github.com/cmjb/SRCDSReserveBot/blob/master/pg.go#L48)

```go
func Get_DB() *pg.DB {
    db := pg.Connect(&pg.Options{
        Addr: "10.0.0.10:5432",
        User: "postgres",
        Password: "example",
        Database: "schedulebot",
    })
    return db
}
```
