# rssator
Backend Web Server + RSS Aggregator (practicing go based of Boot dev's tutorial)

Extennding project: <br/>
    Diferent date formats that diferent rss feeds use , handle that (add different rss feeds and handle issues as they come) <br/>
    Add better logging for unhappy path <br/>
    Go to boot dot dev to see additionnal info for expanding the project <br/>


# Starting the project:

### To launch postgres and redis use:
    docker compose -f docker-compose.yml  up

### To shut down and remove postgres and redis use:
    docker compose -f docker-compose.yml down

### Make sure that you have the following variables in your .ev file:
    PORT=8080
    DB_URL=postgres://admin:admin@localhost:5432/postgres?sslmode=disable

### After starting the docker you can run the migrations
Go to ./sql/schema and run the following commannd i the termial:
    goose postgres postgres://admin:admin@localhost:5432/postgres

### Ten you can build and start the app by running:
    go build ; ./rssator



