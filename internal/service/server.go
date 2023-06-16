package service

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"gitlab.com/wb-dynamics/wb-go-advert-bid-checker/internal/handler"
)

func Serve(ctx context.Context, port int) {
	router := http.NewServeMux()
	router.HandleFunc("/advert-bet-info", handler.HandleAdvertBetInfo)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func ()  {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("error starting server: %s\n", err)
		}				
	}()

	log.Infof("server started on port %d", port)

	<-ctx.Done()
	log.Infof("stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Fatalf("error shutting down server: %s\n", err)
	}
	log.Infof("server stopped")
}

