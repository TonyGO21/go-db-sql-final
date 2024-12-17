package main

import (
	"database/sql"
	"errors"
)

// ParcelStore представляет хранилище для посылок
type ParcelStore struct {
	db *sql.DB
}

// Parcel представляет данные о посылке.
type Parcel struct {
	Number    int
	Client    int
	Status    string
	Address   string
	CreatedAt string
}

// NewParcelStore создаёт новый экземпляр ParcelStore
func NewParcelStore(db *sql.DB) *ParcelStore {
	return &ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	query := `INSERT INTO parcel (client, status, address, created_at) 
	          VALUES (?, ?, ?, ?)`
	result, err := s.db.Exec(query, p.Client, p.Status, p.Address, p.CreatedAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	query := `SELECT number, client, status, address, created_at 
	          FROM parcel WHERE number = ?`
	row := s.db.QueryRow(query, number)

	p := Parcel{}
	err := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return p, nil
		}
		return p, err
	}
	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	query := `SELECT number, client, status, address, created_at 
	          FROM parcel WHERE client = ?`
	rows, err := s.db.Query(query, client)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var parcels []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		parcels = append(parcels, p)
	}
	return parcels, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	query := `UPDATE parcel SET status = ? WHERE number = ?`
	_, err := s.db.Exec(query, status, number)
	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	query := `UPDATE parcel SET address = ? WHERE number = ? AND status = 'registered'`
	result, err := s.db.Exec(query, address, number)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("address can only be changed if status is 'registered'")
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	query := `DELETE FROM parcel WHERE number = ? AND status = 'registered'`
	result, err := s.db.Exec(query, number)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("only parcels with status 'registered' can be deleted")
	}
	return nil
}

// ParcelStatusRegistered и ParcelStatusSent — статусы посылок.
const (
	ParcelStatusRegistered = "registered"
	ParcelStatusSent       = "sent"
)
