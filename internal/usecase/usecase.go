package usecase

import (
	"github.com/patyukin/go-chat/internal/auth"
	"github.com/patyukin/go-chat/internal/cacher"
	"github.com/patyukin/go-chat/internal/config"
	"github.com/patyukin/go-chat/internal/db"
)

type UseCase struct {
	cfg        *config.Config
	registry   *db.Client
	authClient *auth.Client
	chr        *cacher.Cacher
}

func New(cfg *config.Config, registry *db.Client, authClient *auth.Client, chr *cacher.Cacher) *UseCase {
	return &UseCase{
		cfg:        cfg,
		registry:   registry,
		authClient: authClient,
		chr:        chr,
	}
}

func (u *UseCase) GetDomainUseCase() string {
	return u.cfg.CookieDomain
}
