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
var storePartner string

func (pr *partnerRepository) Store(p *syncer.Partner) error {
	log.Println("SQL script:", storePartner)
	_, err := pr.DB.Exec(storePartner, p.ID, p.Name, p.URL)
	return err
}

//go:embed sqls/getPartnerByID.sql
var getPartnerByID string

func (pr *partnerRepository) GetByID(id string) (*syncer.Partner, error) {
	p := syncer.Partner{}
	err := pr.DB.QueryRow(getPartnerByID, id).Scan(&p.ID, &p.Name, &p.URL)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

//go:embed sqls/getAllPartners.sql
var getAllPartners string

func (pr *partnerRepository) GetAll() []syncer.Partner {
	res, err := pr.DB.Query(getAllPartners)
	if err != nil {
		return []syncer.Partner{}
	}
	defer res.Close()

	partners := []syncer.Partner{}
	for res.Next() {
		p := syncer.Partner{}
		err := res.Scan(&p.ID, &p.Name, &p.URL)
		if err != nil {
			return []syncer.Partner{}
		}
		partners = append(partners, p)
	}
	return partners
}
