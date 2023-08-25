package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	dbQueries "github.com/ChessSwahili/live-ratings-sse/internal/sqlc"
	"github.com/julienschmidt/httprouter"
)

type application struct {
	slqxModels *dbQueries.Queries
}

type Broker struct {
	clients map[chan string]bool

	newClients chan chan string

	defunctClients chan chan string

	messages chan string
}

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

func (b *Broker) Start() {

	go func() {

		for {

			select {

			case s := <-b.newClients:

				b.clients[s] = true

			case s := <-b.defunctClients:
				delete(b.clients, s)
				close(s)

			case msg := <-b.messages:
				for s := range b.clients {
					s <- msg
				}
			}
		}
	}()
}

func main() {

	var cfg config

	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("Invalid Port Number")
	}

	flag.IntVar(&cfg.port, "port", port, "API server port")
	flag.StringVar(&cfg.env, "env", "production", "Environment (development|production")
	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("LICHESS_DSN"), "PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max ilde connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgreSQL max connection  connections")

	flag.Parse()

	con, err := openDB(cfg)

	if err != nil {
		log.Fatal(err)
	}

	app := application{
		slqxModels: dbQueries.New(con),
	}

	ctx := context.Background()

	res_static, err := app.slqxModels.CheckEntriesStatic(ctx)

	res_dynamic, err := app.slqxModels.CheckEntriesDynamic(ctx)

    players := FetchTeamPlayers()

    switch {

	case res_static == 0:
		//insert in static

	case res_dynamic == 0:
		//insert in dynamic

	case res_dynamic > res_static:
        // insert in static

	case res_static > res_dynamic:
		//it cant happen too sad left members
	}


	



	// base_ratings :=
	// // if err != nil {
	// // 	log.Fatal("Fail to start the database")
	// // }

	b := &Broker{
		make(map[chan string]bool),
		make(chan (chan string)),
		make(chan (chan string)),
		make(chan string),
	}

	b.Start()

	router := httprouter.New()

	router.Handler(http.MethodGet, "/events", b)

	messages := []string{"Start Now"}
	go func() {
		for i := 0; ; i++ {
			b.messages <- messages[i%len(messages)]
			time.Sleep(3e9)

		}
	}()

	http.ListenAndServe(":7667", enableCORS(router))
}
