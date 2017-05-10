package store

import (
	"time"

	"database/sql"

	"github.com/jmoiron/sqlx"
	uuid "github.com/satori/go.uuid"
)

type ContactID struct {
	sql.NullInt64
}

type Contact struct {
	Org  OrgID     `db:"org_id"`
	ID   ContactID `db:"id"`
	UUID string    `db:"uuid"`
	Name string    `db:"name"`

	URN ContactURNID `db:"urn_id"`

	CreatedOn  time.Time `db:"created_on"`
	ModifiedOn time.Time `db:"modified_on"`

	CreatedBy  int `db:"created_by_id"`
	ModifiedBy int `db:"modified_by_id"`
}

const lookupContactFromURNSQL = `
SELECT c.org_id, c.id, c.uuid, c.name, u.id as "urn_id"
FROM contacts_contact AS c, contacts_contacturn AS u 
WHERE u.urn = $1 AND u.contact_id = c.id AND u.org_id = $2 AND c.is_active = TRUE AND c.is_test = FALSE
`

const insertContactSQL = `
INSERT INTO contacts_contact(org_id, is_active, is_blocked, is_test, is_stopped, uuid, created_on, modified_on, created_by_id, modified_by_id, name) 
VALUES(:org_id, TRUE, FALSE, FALSE, FALSE, :uuid, :created_on, :modified_on, :created_by_id, :modified_by_id, :name)
RETURNING id
`

// InsertContact inserts the passed in contact, the id field will be populated with the result on success
func InsertContact(db *sqlx.DB, contact *Contact) error {
	rows, err := db.NamedQuery(insertContactSQL, contact)
	if err != nil {
		return err
	}
	if rows.Next() {
		rows.Scan(&contact.ID)
	}
	return err
}

// ContactForURN first tries to look up a contact for the passed in URN, if not finding one then creating one
func ContactForURN(db *sqlx.DB, org OrgID, channel ChannelID, urn URN, name string) (*Contact, error) {
	// try to look up our contact by URN
	var contact Contact
	err := db.Get(&contact, lookupContactFromURNSQL, urn, org)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// we found it, return it
	if err != sql.ErrNoRows {
		return &contact, nil
	}

	// didn't find it, we need to create it instead
	contact.Org = org
	contact.UUID = uuid.NewV4().String()
	contact.Name = name

	// TODO: Set these to a system user
	contact.CreatedBy = 1
	contact.ModifiedBy = 1

	// Insert it
	err = InsertContact(db, &contact)
	if err != nil {
		return nil, err
	}

	// now find our URN
	contactURN, err := ContactURNForURN(db, org, channel, contact.ID, urn)
	if err != nil {
		return nil, err
	}

	// save this URN on our contact
	contact.URN = contactURN.ID

	// and return it
	return &contact, err
}
