package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/totorialman/go-task-avito/restapi"
	"github.com/totorialman/go-task-avito/restapi/operations"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/totorialman/go-task-avito/internal/middleware/acl"
	metricsmw "github.com/totorialman/go-task-avito/internal/middleware/metrics"
	authHandler "github.com/totorialman/go-task-avito/internal/pkg/auth/delivery/http"
	authRepo "github.com/totorialman/go-task-avito/internal/pkg/auth/repo"
	authUsecase "github.com/totorialman/go-task-avito/internal/pkg/auth/usecase"
	"github.com/totorialman/go-task-avito/internal/pkg/metrics"
)

func main() {
	logFile, err := os.OpenFile("main.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("ошибка открытия файла логов:", err)
		return
	}
	defer logFile.Close()

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(logFile, os.Stdout), &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	db, err := initDB(logger)
	if err != nil {
		logger.Error("Ошибка при подключении к PostgreSQL", slog.String("err", err.Error()))
		return
	}
	defer db.Close()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		logger.Error("JWT_SECRET не задан")
		return
	}
	mt0, err := metrics.NewProductMetrics()
	if err != nil {
		log.Fatal(err)
	}

	authRepo := authRepo.NewAuthRepo(db)
	authUsecase := authUsecase.NewAuthUsecase(authRepo)
	authHandler := authHandler.NewAuthHandler(authUsecase, mt0)
	
	mt, err := metrics.NewHttpMetrics()
	if err != nil {
		log.Fatal(err)
	}
	middl := metricsmw.CreateHttpMetricsMiddleware(mt)
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}
	api := operations.NewBackendServiceAPI(swaggerSpec)
	configureAPI(api, authHandler)

	handler := api.Serve(nil)
	wrapped := middl(acl.NewAclMiddleware(handler))

	server := restapi.NewServer(api)
	server.SetHandler(wrapped)

	defer server.Shutdown()

	r := mux.NewRouter()
	r.PathPrefix("/metrics").Handler(promhttp.Handler())

	http.Handle("/", r)
	httpSrv := http.Server{Handler: r, Addr: "0.0.0.0:9000"}
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			logger.Error("fail httpSrv.ListenAndServe")
		}
	}()
	
	

	server.Host = "0.0.0.0"
	server.Port = 8080
	if err := server.Serve(); err != nil {
		logger.Error("Ошибка запуска сервера", slog.String("err", err.Error()))
	}

}

func errorAsJSON(err error) []byte {
	//nolint:errchkjson
	b, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{err.Error()})
	return b
}

func configureAPI(api *operations.BackendServiceAPI, handler *authHandler.AuthHandler) {
	api.ServeError = func(rw http.ResponseWriter, r *http.Request, err error) {
		rw.Header().Set("Content-Type", "application/json")
		switch e := err.(type) {
		case *errors.CompositeError:
			rw.Write(errorAsJSON(e.Unwrap()[0]))
			rw.WriteHeader(http.StatusBadRequest)
		default:
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
	api.PostDummyLoginHandler = operations.PostDummyLoginHandlerFunc(handler.HandleDummyLogin)
	api.PostLoginHandler = operations.PostLoginHandlerFunc(handler.HandleLogin)
	api.PostRegisterHandler = operations.PostRegisterHandlerFunc(handler.HandleSignUp)
}

func initDB(logger *slog.Logger) (*pgxpool.Pool, error) {
	connStr := os.Getenv("POSTGRES_CONN")
	if connStr == "" {
		return nil, fmt.Errorf("POSTGRES_CONN не задан")
	}
	db, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(context.Background()); err != nil {
		return nil, err
	}
	logger.Info("Подключение к PostgreSQL успешно")
	return db, nil
}
