package syncing

import (
	"log"

	syncer "github.com/silvan-talos/cookie-syncer"
)

type Service interface {
	SyncUsers(partnerID, otherPID, uidPartner, uidOther string) error
	GetStatus(partnerID string) ([]Status, error)
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

type Status struct {
	PartnerID  string `json:"partner_id"`
	UidPartner string `json:"partner_user_id"`
	OtherPID   string `json:"second_partner_id"`
	UidOther   string `json:"second_partner_user_id"`
}

func (s *service) GetStatus(partnerID string) ([]Status, error) {
	syncs, err := s.syncs.GetAllForPartner(partnerID)
	if err != nil {
		return []Status{}, err
	}
	status := []Status{}
	for _, sync := range syncs {
		st := Status{
			PartnerID:  sync.PartnerID,
			UidPartner: sync.UidPartner,
			OtherPID:   sync.OtherPID,
			UidOther:   sync.UidOther,
		}
		status = append(status, st)
	}
	return status, nil
}
