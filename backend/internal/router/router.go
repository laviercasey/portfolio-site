package router

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/portfolio/backend/internal/handler"
	"github.com/portfolio/backend/internal/middleware"
	"github.com/portfolio/backend/internal/service"
)

type Handlers struct {
	Health    *handler.HealthHandler
	Auth      *handler.AuthHandler
	Project   *handler.ProjectHandler
	Content   *handler.ContentHandler
	Career    *handler.CareerHandler
	Inquiry   *handler.InquiryHandler
	Upload    *handler.UploadHandler
	Analytics *handler.AnalyticsHandler
}

type Deps struct {
	Ctx         context.Context
	Handlers    Handlers
	AuthService *service.AuthService
	CORSOrigins []string
	Logger      *slog.Logger
}

func New(deps Deps) chi.Router {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(middleware.Logger(deps.Logger))
	r.Use(middleware.CORS(deps.CORSOrigins))
	r.Use(chimw.Recoverer)

	r.Get("/healthz", deps.Handlers.Health.Check)

	ctx := deps.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	inquiryLimiter := middleware.NewRateLimiter(ctx, 1, 5)
	loginLimiter := middleware.NewRateLimiter(ctx, 0.1, 5)
	analyticsLimiter := middleware.NewRateLimiter(ctx, 0.5, 30)

	r.Route("/api", func(r chi.Router) {
		r.With(loginLimiter.Limit).Post("/auth/login", deps.Handlers.Auth.Login)

		r.Get("/projects", deps.Handlers.Project.List)
		r.Get("/projects/{slug}", deps.Handlers.Project.GetBySlug)

		r.Get("/content", deps.Handlers.Content.List)

		r.Get("/career", deps.Handlers.Career.GetAll)

		r.With(inquiryLimiter.Limit).Post("/inquiries", deps.Handlers.Inquiry.Create)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(deps.AuthService))

			r.Post("/projects", deps.Handlers.Project.Create)
			r.Put("/projects/{id}", deps.Handlers.Project.Update)
			r.Delete("/projects/{id}", deps.Handlers.Project.Delete)

			r.Put("/content/{section}", deps.Handlers.Content.Update)

			r.Post("/career/{type}", deps.Handlers.Career.Create)
			r.Put("/career/{type}/{id}", deps.Handlers.Career.Update)
			r.Delete("/career/{type}/{id}", deps.Handlers.Career.Delete)

			r.Get("/inquiries", deps.Handlers.Inquiry.List)
			r.Get("/inquiries/{id}", deps.Handlers.Inquiry.GetByID)
			r.Patch("/inquiries/{id}", deps.Handlers.Inquiry.UpdateStatus)

			r.Post("/upload", deps.Handlers.Upload.Upload)
			r.Get("/media", deps.Handlers.Upload.List)
			r.Delete("/media/{id}", deps.Handlers.Upload.Delete)

			r.Route("/analytics", func(r chi.Router) {
				r.Use(analyticsLimiter.Limit)
				r.Get("/summary", deps.Handlers.Analytics.Summary)
				r.Get("/top-pages", deps.Handlers.Analytics.TopPages)
				r.Get("/top-referrers", deps.Handlers.Analytics.TopReferrers)
				r.Get("/top-countries", deps.Handlers.Analytics.TopCountries)
				r.Get("/top-utm", deps.Handlers.Analytics.TopUTM)
				r.Get("/timeseries", deps.Handlers.Analytics.Timeseries)
			})
		})
	})

	return r
}
