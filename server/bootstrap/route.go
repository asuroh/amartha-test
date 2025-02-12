package bootstrap

import (
	"amartha-test/pkg/logruslogger"
	api "amartha-test/server/handler"

	chimiddleware "github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

// RegisterRoutes ...
func (boot *Bootup) RegisterRoutes() {
	handlerType := api.Handler{
		EnvConfig:  boot.EnvConfig,
		Validate:   boot.Validator,
		Translator: boot.Translator,
		ContractUC: &boot.ContractUC,
	}

	boot.R.Route("/api", func(r chi.Router) {
		// Define a limit rate to 1000 requests per IP per request.
		rate, _ := limiter.NewRateFromFormatted("1000-S")
		store, _ := sredis.NewStoreWithOptions(boot.ContractUC.Redis, limiter.StoreOptions{
			Prefix:   "limiter_global",
			MaxRetry: 3,
		})
		rateMiddleware := stdlib.NewMiddleware(limiter.New(store, rate, limiter.WithTrustForwardHeader(true)))
		r.Use(rateMiddleware.Handler)

		// Logging setup
		r.Use(chimiddleware.RequestID)
		r.Use(logruslogger.NewStructuredLogger(boot.EnvConfig["LOG_FILE_PATH"], boot.EnvConfig["LOG_DEFAULT"], boot.ContractUC.ReqID))
		r.Use(chimiddleware.Recoverer)

		// API
		r.Route("/v1", func(r chi.Router) {
			loanHandler := api.LoanHandler{Handler: handlerType}
			r.Route("/loan", func(r chi.Router) {
				r.Post("/create", loanHandler.CreateLoan)
			})

			paymentHandler := api.PaymentHandler{Handler: handlerType}
			r.Route("/payment", func(r chi.Router) {
				r.Post("/execute", paymentHandler.Execute)
				r.Get("/get-outstanding/{loanId}", paymentHandler.GetOutstanding)
				r.Get("/get-delinquent/{loanId}", paymentHandler.GetDelinquent)
			})
		})
	})
}
