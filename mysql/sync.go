package mysql

import (
	"database/sql"
	_ "embed"
	"log"

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
	log.Println("SQL script:", storeSync)
	_, err := sr.DB.Exec(storeSync, s.PartnerID, s.UidPartner, s.OtherPID, s.UidOther)
	return err
}
