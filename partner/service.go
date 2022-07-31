package partner

import (
	"errors"

	"github.com/google/uuid"
	syncer "github.com/silvan-talos/cookie-syncer"
)

// ErrInvalidArgument is returned when one or more arguments are invalid.
var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	AddPartner(name, url string) (string, error)
	GetByID(id string) (*syncer.Partner, error)
	GetAll() []syncer.Partner
}

type service struct {
	partners syncer.PartnerRepository
}

func NewService(partners syncer.PartnerRepository) Service {
	return &service{
		partners: partners,
	}
}

func (s *service) AddPartner(name, url string) (string, error) {
	if name == "" || url == "" {
		return "", ErrInvalidArgument
	}

	p := &syncer.Partner{
		ID:   uuid.New().String(),
		Name: name,
		URL:  url,
	}

	if err := s.partners.Store(p); err != nil {
		return "", err
	}

	return p.ID, nil
}

func (s *service) GetByID(id string) (*syncer.Partner, error) {
	if id == "" {
		return nil, ErrInvalidArgument
	}
	return s.partners.GetByID(id)
}

func (s *service) GetAll() []syncer.Partner {
	return s.partners.GetAll()
}
