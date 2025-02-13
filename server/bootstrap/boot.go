package bootstrap

import (
	"amartha-test/usecase"
	"database/sql"

	"github.com/go-chi/chi"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-redis/redis"
	validator "gopkg.in/go-playground/validator.v9"
)

// Bootup ...
type Bootup struct {
	R          *chi.Mux
	CorsDomain []string
	EnvConfig  map[string]string
	DB         *sql.DB
	Redis      *redis.Client
	Validator  *validator.Validate
	Translator ut.Translator
	ContractUC usecase.ContractUC
}
