package mysql

import (
	"database/sql"
	_ "embed"
	"log"

	syncer "github.com/silvan-talos/cookie-syncer"
)

type partnerRepository struct {
	DB *sql.DB
}

func NewPartnerRepository(db *sql.DB) syncer.PartnerRepository {
	return &partnerRepository{
		DB: db,
	}

}

//go:embed sqls/storePartner.sql
var storeScript string

func (pr *partnerRepository) Store(p *syncer.Partner) error {
	log.Println("SQL script:", storeScript)
	res, err := pr.DB.Exec(storeScript, p.ID, p.Name)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Println("failed to get last id")
		return err
	}
	log.Println("id:", id)
	return nil
}
