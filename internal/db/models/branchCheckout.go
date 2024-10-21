package models

import (
	"database/sql"
	"time"
)

type BranchCheckout struct {
	ID               int64
	Name             string
	CheckoutCount    int64
	LastCheckedOutAt time.Time
	LatestCommitMsg  string
}

type BranchCheckoutRepository struct {
	db *sql.DB
}

func NewBranchCheckoutRepository(db *sql.DB) *BranchCheckoutRepository {
	return &BranchCheckoutRepository{db: db}
}

func (r *BranchCheckoutRepository) Upsert(b *BranchCheckout) error {
	query := `
        INSERT INTO branch_checkouts (name, checkout_count, last_checked_out_at, latest_commit_msg)
        VALUES (?, 1, ?, ?)
        ON CONFLICT(name) DO UPDATE SET
            checkout_count = checkout_count + 1,
            last_checked_out_at = ?,
            latest_commit_msg = ?
    `

	now := time.Now()
	_, err := r.db.Exec(query, b.Name, now, b.LatestCommitMsg, now, b.LatestCommitMsg)
	return err
}

func (r *BranchCheckoutRepository) GetRecent(limit int, isReverse bool) ([]BranchCheckout, error) {
	// Get the most recently checked out branches
	query := `SELECT * FROM branch_checkouts ORDER BY last_checked_out_at DESC LIMIT ?`
	if isReverse {
		query = `SELECT * FROM branch_checkouts ORDER BY last_checked_out_at ASC LIMIT ?`
	}

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []BranchCheckout
	for rows.Next() {
		var b BranchCheckout
		err := rows.Scan(&b.ID, &b.Name, &b.CheckoutCount, &b.LastCheckedOutAt, &b.LatestCommitMsg)
		if err != nil {
			return nil, err
		}
		branches = append(branches, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return branches, nil
}

func (r *BranchCheckoutRepository) GetFrequent(limit int, isReverse bool) ([]BranchCheckout, error) {
	// Get the most frequently checked out branches
	query := `SELECT * FROM branch_checkouts ORDER BY checkout_count DESC LIMIT ?`
	if isReverse {
		query = `SELECT * FROM branch_checkouts ORDER BY checkout_count ASC LIMIT ?`
	}

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var branches []BranchCheckout
	for rows.Next() {
		var b BranchCheckout
		err := rows.Scan(&b.ID, &b.Name, &b.CheckoutCount, &b.LastCheckedOutAt, &b.LatestCommitMsg)
		if err != nil {
			return nil, err
		}
		branches = append(branches, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return branches, nil
}
