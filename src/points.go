package main

import (
	"context"
	"database/sql"
	"github.com/irminsul-network/focus-server/src/errors"
	"log"
	"time"
)

func createPoint(p CreatePoint) {

}

type PointManager interface {
	createPoint(ctx context.Context, p CreatePoint) (int64, error)
}

type DbPointManager struct {
	db *sql.DB
}

func (pm DbPointManager) createPoint(ctx context.Context, p CreatePoint) (int64, error) {

	now := time.Now().Unix()

	res, err := pm.db.ExecContext(ctx,
		`
			insert into points(content, created_on, archived, achieved) values($1, $2, false, false)
		`,
		p.Content, now)
	if err != nil {
		log.Print(err)
		return 0, errors.NewServerError("unable to insert point", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Print(err)
		return 0, errors.NewServerError("failed fetching inserted id", err)
	}

	_, err = pm.db.ExecContext(ctx,
		`
			insert into point_encounters(point_id, encountered_on, conquered, urgency) values($1, $2, $3, $4)
		`,
		id, now, p.Conquered, p.Urgency)
	if err != nil {
		log.Print(err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return 0, errors.NewServerError("unable to insert encounter to db", err)
	}

	return id, nil
}

type Point struct {
	Id       int64  `json:"id"`
	Content  string `json:"content"`
	Created  int64  `json:"created"`
	Achieved bool   `json:"achieved"`
	Archived bool   `json:"archived"`
}

type CreatePoint struct {
	Point
	Conquered int64 `json:"conquered"`
	Urgency   int64 `json:"urgency"`
}

type PointEncounter struct {
	Id            int64 `json:"id"`
	PointId       int64 `json:"point_id"`
	EncounteredOn int64 `json:"encountered_on"`
	Conquered     int64 `json:"conquered"`
	Urgency       int64 `json:"urgency"`
}
