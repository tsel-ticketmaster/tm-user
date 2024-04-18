package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/cors"
	"github.com/tsel-ticketmaster/tm-user/config"
	"github.com/tsel-ticketmaster/tm-user/internal/module/adminapp/admin"
	"github.com/tsel-ticketmaster/tm-user/internal/pkg/jwt"
	internalMiddleare "github.com/tsel-ticketmaster/tm-user/internal/pkg/middleware"
	"github.com/tsel-ticketmaster/tm-user/internal/pkg/session"
	"github.com/tsel-ticketmaster/tm-user/pkg/applogger"
	"github.com/tsel-ticketmaster/tm-user/pkg/middleware"
	"github.com/tsel-ticketmaster/tm-user/pkg/monitoring"
	"github.com/tsel-ticketmaster/tm-user/pkg/postgresql"
	"github.com/tsel-ticketmaster/tm-user/pkg/redis"
	"github.com/tsel-ticketmaster/tm-user/pkg/server"
	"github.com/tsel-ticketmaster/tm-user/pkg/validator"
)

func main() {
	c := config.Get()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := applogger.GetLogrus()

	mon := monitoring.NewOpenTelemetry(
		c.Application.Name,
		c.Application.Environment,
		c.OpenTelemetry.Collector.Endpoint,
	)

	mon.Start(ctx)

	validate := validator.Get()

	jsonWebToken := jwt.NewJSONWebToken(c.JWT.PrivateKey, c.JWT.PublicKey)

	psqldb := postgresql.GetDatabase()
	if err := psqldb.Ping(); err != nil {
		logger.WithContext(ctx).WithError(err).Error()
	}

	rc := redis.GetClient()
	if err := rc.Ping(context.Background()).Err(); err != nil {
		logger.WithContext(ctx).WithError(err).Error()
	}

	session := session.NewRedisSessionStore(logger, rc)

	adminSessionMiddleware := internalMiddleare.NewAdminSessionMiddleware(jsonWebToken, session)

	router := http.NewServeMux()

	// admin's app
	adminappAdminRepository := admin.NewAdminRepository(logger, psqldb)
	adminappAdminUseCase := admin.NewAdminUseCase(admin.AdminUseCaseProperty{
		Logger:          logger,
		Timeout:         c.Application.Timeout,
		JSONWebToken:    jsonWebToken,
		Session:         session,
		AdminRepository: adminappAdminRepository,
	})
	admin.InitHTTPHandler(router, adminSessionMiddleware, validate, adminappAdminUseCase)

	handler := middleware.SetChain(
		router,
		middleware.HTTPOpenTelemetryTracer,
		middleware.HTTPResponseTraceInjection,
		middleware.NewHTTPRequestLogger(logger, c.Application.Debug).Middleware,
		cors.New(cors.Options{
			AllowedOrigins:   c.CORS.AllowedOrigins,
			AllowedMethods:   c.CORS.AllowedMethods,
			AllowedHeaders:   c.CORS.AllowedHeaders,
			ExposedHeaders:   c.CORS.ExposedHeaders,
			MaxAge:           c.CORS.MaxAge,
			AllowCredentials: c.CORS.AllowCredentials,
		}).Handler,
	)

	srv := &server.Server{
		Server: http.Server{
			Addr:    fmt.Sprintf(":%d", c.Application.Port),
			Handler: handler,
		},
		Logger: logger,
	}

	go func() {
		srv.ListenAndServe()
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	<-sigterm

	srv.Shutdown(ctx)
	psqldb.Close()
	rc.Close()
	mon.Stop(ctx)
}
