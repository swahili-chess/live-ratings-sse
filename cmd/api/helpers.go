package main

import (
	"context"
	"database/sql"
	"time"

	"github.com/pingcap/log"
)

func (app application)InsertUsernames(list []PlayerMinDt, param interface{}, indicator string) {
	// get current usernames in db

	if indicator == "static"{
         	lichess_ids, err := app.slqxModels.GetallStatic()
	}

	if err != nil {
		log.Error("Failed to get usernames in DB")
		return
	}

	newPlayers := findNewPlayers(lichess_ids, list)

	for _, player := range newPlayers {

		err := sw.models.Lichess.Insert(player)

		if err != nil {
			log.Error("Failed to insert user", player)
		}

	}
}

func findNewPlayers(lichess_ids []string, players []PlayerMinDt) []PlayerMinDt {
	newPlayers := []PlayerMinDt{}
	elementSet := make(map[string]bool)

	for _, lichess_id := range lichess_ids {
		elementSet[lichess_id] = true
	}

	for _, dt := range players {
		if _, found := elementSet[dt.ID]; !found {
			newPlayers = append(newPlayers, dt)
		} else {
			delete(elementSet, dt.ID) // Remove common elements
		}
	}

	return newPlayers
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	duration, err := time.ParseDuration(cfg.db.maxIdleTime)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)

	if err != nil {
		return nil, err
	}

	return db, nil
}
