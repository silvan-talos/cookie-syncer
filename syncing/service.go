package syncing

import (
	"log"

	syncer "github.com/silvan-talos/cookie-syncer"
)

type Service interface {
	SyncUsers(partnerID, otherPID, uidPartner, uidOther string) error
}

type service struct {
	syncs syncer.SyncRepository
}

func NewService(syncs syncer.SyncRepository) Service {
	return &service{
		syncs: syncs,
	}
}

func (s *service) SyncUsers(partnerID, otherPID, uidPartner, uidOther string) error {
	sync := &syncer.Sync{
		PartnerID:  partnerID,
		OtherPID:   otherPID,
		UidPartner: uidPartner,
		UidOther:   uidOther,
	}
	err := s.syncs.Store(sync)
	if err != nil {
		log.Println("failed to store userSync:", err)
	}
	return err
}
