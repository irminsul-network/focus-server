package main

import (
	"log"
	"net/http"
)

func registerHandlers(app *App) {

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		//ctx, cancel := context.WithCancel(r.Context())
		//defer cancel()

		var bytes []byte = make([]byte, 1024)
		n, err := r.Body.Read(bytes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		log.Printf("read %d bytes", n)

		log.Println(string(bytes[:n]))

		//_, err = app.db.ExecContext(context.Background(), "insert into point(content, urgency) values('first', 1)")
		//if err != nil {
		//	log.Print(err)
		//	return
		//}
	})
}
