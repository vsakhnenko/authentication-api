package main

import (
	"authentication/internal/server"
	"errors"
	"fmt"
	"net/http"
)

func main() {
	serverInstance := server.NewServer()

	err := serverInstance.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}
