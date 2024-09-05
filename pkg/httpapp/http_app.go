// Package httpapp package provides functionality for http requests handling.
package httpapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"goboilerplate/pkg/apperr"
	"goboilerplate/pkg/httpapp/gql/generated"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql"
	gqlHandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const InternalServerErrorCode = "INTERNAL_SERVER_ERROR"

// App struct is instance of http application
type App struct {
	Server *echo.Echo
	Port   Port
}

type Port string

// Initialize instantiates http application with configured routes
func Initialize(
	resolver generated.ResolverRoot,
	handlers Handler,
	port Port,
) *App {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPatch,
			http.MethodPost,
			http.MethodDelete,
			http.MethodOptions,
		},
	}))
	e.HTTPErrorHandler = customHTTPErrorHandler
	gqlServer := gqlHandler.New(
		generated.NewExecutableSchema(generated.Config{Resolvers: resolver}),
	)
	gqlServer.Use(extension.Introspection{})
	gqlServer.AddTransport(transport.POST{})
	gqlServer.SetErrorPresenter(func(ctx context.Context, e error) *gqlerror.Error {
		err := graphql.DefaultErrorPresenter(ctx, e)
		var known apperr.Known
		if errors.As(e, &known) {
			err.Message = known.Error()
			err.Extensions = map[string]interface{}{
				"code": known.Code(),
			}
			err.Path = nil

			return err
		}

		reqCtx := graphql.GetOperationContext(ctx)
		er := errors.Errorf("%s, gqlQuery: %s | gqlVariables: %s | stacktrace: %s",
			err.Unwrap(),
			reqCtx.RawQuery,
			reqCtx.Variables,
			fmt.Sprintf("%.10000s", string(debug.Stack())),
		)
		log.Error(er)

		return &gqlerror.Error{
			Message: "Internal server error",
			Extensions: map[string]interface{}{
				"code": InternalServerErrorCode,
			},
		}
	})
	gqlServer.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		reqCtx := graphql.GetOperationContext(ctx)
		e := errors.Errorf("%s, gqlQuery: %s | gqlVariables: %s | stacktrace: %s",
			err,
			reqCtx.RawQuery,
			reqCtx.Variables,
			fmt.Sprintf("%.10000s", string(debug.Stack())),
		)

		return e
	})
	playgroundHandler := playground.Handler("GraphQL", "/graphql")
	e.POST("/graphql", func(c echo.Context) error {
		// Create newrelic transaction segment
		var params *graphql.RawParams
		buf, _ := io.ReadAll(c.Request().Body)
		rdr1 := io.NopCloser(bytes.NewBuffer(buf))
		rdr2 := io.NopCloser(bytes.NewBuffer(buf))
		dec := json.NewDecoder(rdr1)
		if err := dec.Decode(&params); err == nil {
			tx := newrelic.FromContext(c.Request().Context())
			seg := tx.StartSegment(params.OperationName)
			defer seg.End()
		}
		c.Request().Body = rdr2
		// Process graphql request
		gqlServer.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	e.POST("/upload", handlers.UploadFile)
	e.GET("/playground", func(c echo.Context) error {
		playgroundHandler.ServeHTTP(c.Response(), c.Request())
		return nil
	})
	e.GET("/health", func(c echo.Context) error {
		return nil
	})

	return &App{
		Server: e,
		Port:   port,
	}
}

// Run starts http Server
func (a *App) Run() {
	if err := a.Server.Start(fmt.Sprintf(":%v", a.Port)); err != nil &&
		err != http.ErrServerClosed {
		a.Server.Logger.Fatal("shutting down the server")
	}
}

// WaitForInterrupt waits for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
func (a *App) WaitForInterrupt() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Info("shutting down HTTP server", nil)

	if err := a.Server.Shutdown(ctx); err != nil {
		a.Server.Logger.Fatal(err)
	}
}
