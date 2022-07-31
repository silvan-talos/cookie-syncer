package syncer

type Sync struct {
	PartnerID  string
	OtherPID   string
	UidPartner string
	UidOther   string
}

type SyncRepository interface {
	Store(sync *Sync) error
	GetAllForPartner(partnerID string) ([]Sync, error)
}
