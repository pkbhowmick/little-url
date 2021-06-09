package apis

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func Serve()  {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello"))
	})

	log.Println("server is listening on port 3000")

	http.ListenAndServe(":3000", r)
}
