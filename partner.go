package syncer

type Partner struct {
	ID   string
	Name string
	URL  string
}

type PartnerRepository interface {
	Store(partner *Partner) error
	GetByID(id string) (*Partner, error)
	GetAll() []Partner
}
