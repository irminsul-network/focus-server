package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

func registerHandlers(app *App) {

	http.HandleFunc("/points/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {

			ctx, cancel := context.WithCancel(r.Context())
			defer cancel()

			var point Point
			err := json.NewDecoder(r.Body).Decode(&point)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			_, err = app.db.ExecContext(ctx,
				"insert into point(content, urgency, created, archived, achieved) values($1, $2, $3, false, false)",
				point.Content, point.Urgency, time.Now().Unix())
			if err != nil {
				log.Print(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusCreated)
			write(w, "ok")

		} else if r.Method == "PATCH" {

			ctx, cancel := context.WithCancel(r.Context())
			defer cancel()

			var point Point
			decoder := json.NewDecoder(r.Body)
			decoder.DisallowUnknownFields()
			err := decoder.Decode(&point)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			id := strings.TrimPrefix(r.URL.Path, "/points/")
			if id == "" {
				http.Error(w, "ID required", http.StatusBadRequest)
			}
			res, err := app.db.ExecContext(ctx,
				"update point set content = $1, urgency = $2 where id = $3",
				point.Content,
				point.Urgency,
				id)

			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			rows, err := res.RowsAffected()
			if err != nil {
				log.Println(err)
				return
			}
			if rows == 1 {
				w.WriteHeader(http.StatusAccepted)
				write(w, "ok")
				return
			}

			if rows < 1 {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			if rows > 1 {
				log.Println("too many records updated for id: ", id, "records: ", rows)
				w.WriteHeader(http.StatusInternalServerError)
			}

		} else if r.Method == "DELETE" {

			ctx, cancel := context.WithCancel(r.Context())
			defer cancel()

			var point Point
			err := json.NewDecoder(r.Body).Decode(&point)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			id := strings.TrimPrefix(r.URL.Path, "/points/")
			if id == "" {
				http.Error(w, "ID required", http.StatusBadRequest)
			}
			res, err := app.db.ExecContext(ctx, "delete from point  where id = $1", id)

			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			rows, err := res.RowsAffected()

			if err != nil {
				log.Println(err)
				return
			}

			if rows > 1 {
				log.Println("too many records deleted for id: ", id, "records: ", rows)
				w.WriteHeader(http.StatusInternalServerError)
			}

			if rows == 1 {
				w.WriteHeader(http.StatusAccepted)
				write(w, "ok")

			} else if rows < 1 {
				w.WriteHeader(http.StatusNotFound)
				write(w, "not found")
			}
		} else if r.Method == "GET" {
			ctx, cancel := context.WithCancel(r.Context())
			defer cancel()

			rows, err := app.db.QueryContext(ctx, "select id, content, urgency, created, archived, achieved from point")
			if err != nil {
				log.Println(err)
				return
			}
			var points []Point = make([]Point, 0)

			for rows.Next() {
				var point Point
				err := rows.Scan(&point.Id, &point.Content, &point.Urgency, &point.Created, &point.Archived, &point.Achieved)
				if err != nil {
					log.Println(err)
					return
				}

				points = append(points, point)
			}

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			err = json.NewEncoder(w).Encode(points)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			http.Error(w, "Invalid Operation", http.StatusMethodNotAllowed)
		}
		w.WriteHeader(http.StatusBadRequest)

	})
}

type Point struct {
	Id       int64  `json:"id"`
	Content  string `json:"content"`
	Urgency  int    `json:"urgency"`
	Created  int64  `json:"created"`
	Achieved bool   `json:"achieved"`
	Archived bool   `json:"archived"`
}

func write(w http.ResponseWriter, data string) {
	_, err := w.Write([]byte(data))
	if err != nil {
		log.Println(err)
		return
	}
}
