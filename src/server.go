package main

import (
	"context"
	"encoding/json"
	"errors"
	errs "github.com/irminsul-network/focus-server/src/errors"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func registerHandlers(app *App) {

	pm := &DbPointManager{
		app.db,
	}
	http.HandleFunc("/points/", func(w http.ResponseWriter, r *http.Request) {
		HandleGlobalErrors(w, func() error {

			if r.Method == "POST" {
				subPath := strings.TrimPrefix(r.URL.Path, "/points/")
				if subPath == "" {
					return postPoints(pm, r, w)
				}
				//} else if strings.HasSuffix(subPath, "/bump") {
				//	id := strings.TrimSuffix(subPath, "/bump")
				//	bump(r.Context(), app, id, r.Body, w)
				//}
			}

			//} else if r.Method == "DELETE" {
			//	if deletePoint(w, r, app) {
			//		return nil
			//	}
			//} else if r.Method == "GET" {
			//	if getPoints(w, r, app) {
			//		return nil
			//	}
			//} else {
			//	http.Error(w, "Invalid Operation", http.StatusMethodNotAllowed)
			//}
			return errs.NewUserError("bad method", nil)
		})

	})
}

func HandleGlobalErrors(w http.ResponseWriter, f func() error) {
	err := f()
	if err == nil {
		return
	}
	var UserError errs.UserError
	if errors.As(err, &UserError) {
		log.Printf("User Error: %s%n", err)
		http.Error(w, UserError.Error(), http.StatusBadRequest)
	}
	log.Printf("Server Error: %s%n", err)
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

//func getPoints(w http.ResponseWriter, r *http.Request, app *App) bool {
//	ctx, cancel := context.WithCancel(r.Context())
//	defer cancel()
//
//	rows, err := app.db.QueryContext(ctx, `
//		WITH encounter_counts AS (
//		SELECT
//			point_id,
//			COUNT(*) AS total_encountered_count,
//			SUM(CASE WHEN conquered = 1 THEN 1 ELSE 0 END) AS total_conquered_count,
//			COUNT(CASE WHEN encountered_on >= strftime('%s', 'now', '-7 days') THEN 1 END) AS recent_encountered_count,
//			SUM(CASE WHEN conquered = 1 AND encountered_on >= strftime('%s', 'now', '-7 days') THEN 1 ELSE 0 END) AS recent_conquered_count
//		FROM point_encounters
//		GROUP BY point_id
//		)
//		SELECT
//			p.id,
//			p.content,
//			p.achieved,
//			p.archived,
//			ec.recent_encountered_count,
//			ec.recent_conquered_count,
//			ec.total_encountered_count,
//			ec.total_conquered_count
//		FROM
//			point p
//
//		inner JOIN
//			encounter_counts ec ON p.id = ec.point_id;
//	`)
//	if err != nil {
//		log.Println(err)
//		return true
//	}
//	var points = make([]EvaluatedPoint, 0)
//
//	for rows.Next() {
//		var point Point
//		err := rows.Scan(&point.Id, &point.Content, &point.Encountered, &point.Conquered, &point.Created, &point.Archived, &point.Achieved)
//		if err != nil {
//			log.Println(err)
//			return true
//		}
//
//		points = append(points, point)
//	}
//
//	w.WriteHeader(http.StatusOK)
//	w.Header().Set("Content-Type", "application/json")
//	err = json.NewEncoder(w).Encode(points)
//	if err != nil {
//		log.Println(err)
//		return true
//	}
//	return false
//}

//func deletePoint(w http.ResponseWriter, r *http.Request, app *App) bool {
//	ctx, cancel := context.WithCancel(r.Context())
//	defer cancel()
//
//	var point Point
//	err := json.NewDecoder(r.Body).Decode(&point)
//	if err != nil {
//		log.Println(err)
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return true
//	}
//
//	id := strings.TrimPrefix(r.URL.Path, "/points/")
//	if id == "" {
//		http.Error(w, "ID required", http.StatusBadRequest)
//	}
//	res, err := app.db.ExecContext(ctx, "delete from point  where id = $1", id)
//
//	if err != nil {
//		log.Println(err)
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//	}
//
//	rows, err := res.RowsAffected()
//
//	if err != nil {
//		log.Println(err)
//		return true
//	}
//
//	if rows > 1 {
//		log.Println("too many records deleted for id: ", id, "records: ", rows)
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//
//	if rows == 1 {
//		w.WriteHeader(http.StatusAccepted)
//		write(w, "ok")
//
//	} else if rows < 1 {
//		w.WriteHeader(http.StatusNotFound)
//		write(w, "not found")
//	}
//	return false
//}

//func bump(ctx context.Context, app *App, id string, body io.Reader, w http.ResponseWriter) {
//
//	var point Point
//	ctx, cancel := context.WithCancel(ctx)
//	defer cancel()
//
//	err := json.NewDecoder(body).Decode(&point)
//	if err != nil {
//		log.Println(err)
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	res, err := app.db.ExecContext(ctx, "update point set encountered = encountered + $1, conquered = conquered + $2 where id = $3",
//		point.Encountered, point.Conquered, id)
//	if err != nil {
//		log.Println(err)
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	rows, err := res.RowsAffected()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//
//	if rows < 1 {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//
//	if rows > 1 {
//		log.Println("too many records updated for id: ", id, "records: ", rows)
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//
//	w.WriteHeader(http.StatusAccepted)
//	write(w, "ok")
//	return
//
//}

func postPoints(pm PointManager, r *http.Request, w http.ResponseWriter) error {

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	var point CreatePoint
	err := json.NewDecoder(r.Body).Decode(&point)
	res, err := pm.createPoint(ctx, point)
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		write(w, strconv.FormatInt(res, 10))
	}
	return err
}

//func patchPoints(app *App, r *http.Request, w http.ResponseWriter) {
//
//	ctx, cancel := context.WithCancel(r.Context())
//	defer cancel()
//
//	var point Point
//	decoder := json.NewDecoder(r.Body)
//	decoder.DisallowUnknownFields()
//	err := decoder.Decode(&point)
//	if err != nil {
//		log.Println(err)
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//
//	id := strings.TrimPrefix(r.URL.Path, "/points/")
//	if id == "" {
//		http.Error(w, "ID required", http.StatusBadRequest)
//	}
//	res, err := app.db.ExecContext(ctx,
//		"update point set content = $1, encountered = $2 where id = $3",
//		point.Content,
//		point.Encountered,
//		id)
//
//	if err != nil {
//		log.Println(err)
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	rows, err := res.RowsAffected()
//	if err != nil {
//		log.Println(err)
//		return
//	}
//	if rows == 1 {
//		w.WriteHeader(http.StatusAccepted)
//		write(w, "ok")
//		return
//	}
//
//	if rows < 1 {
//		w.WriteHeader(http.StatusNotFound)
//		return
//	}
//
//	if rows > 1 {
//		log.Println("too many records updated for id: ", id, "records: ", rows)
//		w.WriteHeader(http.StatusInternalServerError)
//	}
//}

func write(w http.ResponseWriter, data string) {
	_, err := w.Write([]byte(data))
	if err != nil {
		log.Println(err)
		return
	}
}
