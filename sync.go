package syncer

type Sync struct {
	PartnerID  string
	OtherPID   string
	UidPartner string
	UidOther   string
}

type SyncRepository interface {
	Store(sync *Sync) error
}
