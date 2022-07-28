package partner

import (
	"errors"

	"github.com/google/uuid"
	syncer "github.com/silvan-talos/cookie-syncer"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	AddPartner(name string) (string, error)
}

type service struct {
	partners syncer.PartnerRepository
}

func NewService(partners syncer.PartnerRepository) Service {
	return &service{
		partners: partners,
	}
}

func (s *service) AddPartner(name string) (string, error) {
	if name == "" {
		return "", ErrInvalidArgument
	}

	p := &syncer.Partner{
		ID:   uuid.New().String(),
		Name: name,
	}

	if err := s.partners.Store(p); err != nil {
		return "", err
	}

	return p.ID, nil
}
