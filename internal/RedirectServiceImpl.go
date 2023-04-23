package internal

import (
	"errors"
	"time"

	errs "github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

var (
	ErrRedirectNotFound = errors.New("redirect not found")
	ErrRedirectInvalid  = errors.New("redirect is invalid")
)

type RedirectService struct {
	RedirectRepo RedirectRepositoryPort
}

func NewRedirectService(redirectRepo RedirectRepositoryPort) *RedirectService {
	return &RedirectService{
		RedirectRepo: redirectRepo,
	}
}

func (r *RedirectService) Find(code string) (*Redirect, error) {
	return r.RedirectRepo.Find(code)
}

func (r *RedirectService) Store(redirect *Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errs.Wrap(ErrRedirectInvalid, "service.Redirect.Store")
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.RedirectRepo.Store(redirect)
}
