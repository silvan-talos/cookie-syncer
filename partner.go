package syncer

type Partner struct {
	ID   string
	Name string
}

type PartnerRepository interface {
	Store(partner *Partner) error
}
