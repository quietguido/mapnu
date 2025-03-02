# MapNu

### DB dump:

```
pg_dump --no-owner -Fc -U postgres map_nu -f ./map_nu.custom
```

### DB restore:

```
dropdb -U postgres map_nu
createdb -U postgres map_nu
pg_restore --no-owner -d map_nu -U postgres ./map_nu.custom
```

### Install `migrate` command-tool:

https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

### Create new migration:

```
migrate create -ext sql -dir migrations mg_name
```

### Apply migration:

```
migrate -path migrations -database "postgres://localhost:5432/map_nu?sslmode=disable" up
```
