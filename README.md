# cookie-syncer

`cookie-syncer` is a demonstration cookie syncing system written in Go.

## Running

After cloning the repository, run the following:
```shell
docker compose pull
docker compose up
```
Open [sync.html](sync.html) to see the effects.


## Architecture

`cookie-syncer` components are packaged by dependencies:
* `server` - http server and routing
* `partner` - the service used for managing partners
* `syncing` - the service used for syncing
* `mysql` - repositories implementation for MySQL database
* `cmd` - entry point; ties up all the dependencies
* domain models/entities are stored in the root package called `syncer`

## Database

`cookie-syncer` uses a MySQL version 8 database.

### Database schema
![database schema](https://github.com/silvan-talos/cookie-syncer/blob/main/docs/db.png?raw=true)

## Documentation

Documentation can be generated using godoc. This topic is WIP.

## Testing

### Unit tests

The whole application should be testable. However, this is WIP.

### Integration tests

Integration tests using `Postman` are available [here](https://github.com/silvan-talos/cookie-syncer/blob/main/docs/postman_tests.json)
