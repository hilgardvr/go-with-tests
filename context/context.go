package context

import (
	"context"
	"fmt"
	"net/http"
)

type Store interface {
	Fetch(ctx context.Context) (string, error)
	Cancel()
}

func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := store.Fetch(r.Context())

		if err != nil {
			return
		}

		fmt.Fprint(w, data)

		// ctx := r.Context()
		// data := make(chan string, 1)

		// go func() {
		// 	data <- store.Fetch()
		// }()

		// select {
		// case d := <-data:
		// 	fmt.Fprint(w, d)
		// case <-ctx.Done():
		// 	store.Cancel()
		// }
	}
}
