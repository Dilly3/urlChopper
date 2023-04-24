package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dilly3/urlshortner/api"
	"github.com/dilly3/urlshortner/internal"
	"github.com/dilly3/urlshortner/repository"

	"github.com/go-chi/chi"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	repo := repository.ChooseRepo()
	service := internal.NewRedirectService(repo)
	handler := api.NewHandler(service)
	port := api.SelectPort()
	app := api.NewServer(handler, chi.NewRouter())

	errs := make(chan error, 2)

	go func() {
		fmt.Printf("listening on port => %v\n", port)
		errs <- http.ListenAndServe(port, app.Router)
	}()
	c := make(chan os.Signal, 1)
	go func() {
		signal.Notify(c, syscall.SIGINT)

	}()

	errs <- fmt.Errorf("%s", <-c)
	close(errs)
	fmt.Printf(" terminated ====>  ERR1 : %v\n", <-errs)
	fmt.Printf(" terminated ====>  ERR2 : %v\n", <-errs)

	fmt.Println("\nshutting down server")

}
