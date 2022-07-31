package mysql

import (
	"database/sql"
	_ "embed"

	syncer "github.com/silvan-talos/cookie-syncer"
)

type syncRepository struct {
	DB *sql.DB
}

func NewSyncRepository(db *sql.DB) syncer.SyncRepository {
	return &syncRepository{
		DB: db,
	}

}

//go:embed sqls/storeSync.sql
var storeSync string

func (sr *syncRepository) Store(s *syncer.Sync) error {
	_, err := sr.DB.Exec(storeSync, s.PartnerID, s.UidPartner, s.OtherPID, s.UidOther)
	return err
}

//go:embed sqls/getAllForPartner.sql
var getAllForPartner string

func (sr *syncRepository) GetAllForPartner(partnerID string) ([]syncer.Sync, error) {
	res, err := sr.DB.Query(getAllForPartner, partnerID)
	if err != nil {
		return []syncer.Sync{}, err
	}
	defer res.Close()

	syncs := []syncer.Sync{}
	for res.Next() {
		s := syncer.Sync{}
		err := res.Scan(&s.PartnerID, &s.UidPartner, &s.OtherPID, &s.UidOther)
		if err != nil {
			return []syncer.Sync{}, err
		}
		syncs = append(syncs, s)
	}
	return syncs, nil
}
