package api

import (
	"context"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/oakeshq/go-starter/config"
	"github.com/oakeshq/go-starter/internal/healthcheck"
	"github.com/oakeshq/go-starter/internal/user"
	"github.com/oakeshq/go-starter/pkg/httperr"
	"github.com/oakeshq/go-starter/pkg/logs"
	gmiddleware "github.com/oakeshq/go-starter/pkg/middleware"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"

	"github.com/oakeshq/go-starter/pkg/router"
)

// API exposes the integral struct
type API struct {
	handler    http.Handler
	r          *router.Router
	config     *config.Config
	db     *gorm.DB
}

// NewAPI instantiates a new REST API.
func NewAPI(
	config *config.Config,
	r *router.Router,
	db *gorm.DB,
) *API {

	api := &API{
		r:          r,
		config:     config,
		db:     db,
	}

	ctx := context.Background()
	r.Chi.Use(middleware.RealIP)
	r.Use(gmiddleware.RequestIDCtx)
	r.Use(httperr.Recoverer)
	r.UseBypass(logs.NewStructuredLogger(logrus.StandardLogger()))

	corsHandler := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link", "X-Total-Count"},
		AllowCredentials: true,
	})

	healthcheck.RegisterHandlers(r)

	r.Route("/v1", func(r *router.Router) {
		user.RegisterHandlers(r, db, config)
	})

	api.handler = corsHandler.Handler(chi.ServerBaseContext(ctx, r))
	return api
}