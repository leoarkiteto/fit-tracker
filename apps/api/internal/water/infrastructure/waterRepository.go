// Package infrastructure implements the hydration slice outbound ports on
// SQLite.
package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"fittracker-api/internal/water"
)

// EntryRepository is the SQLite implementation of water.Repository.
type EntryRepository struct {
	db *sql.DB
}

var _ water.Repository = (*EntryRepository)(nil)

// NewEntryRepository builds the hydration repository.
func NewEntryRepository(db *sql.DB) *EntryRepository {
	return &EntryRepository{db: db}
}

// GetByDay returns the entries of a profile logged in [start, end).
func (r *EntryRepository) GetByDay(
	ctx context.Context,
	profileID string,
	start, end time.Time,
) ([]water.Entry, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, amount_ml, consumed_at, note FROM WaterIntakeEntries WHERE user_profile_id = ? AND consumed_at >= ? AND consumed_at < ?`,
		profileID,
		start,
		end,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []water.Entry
	for rows.Next() {
		var e water.Entry
		err := rows.Scan(&e.ID, &e.AmountMl, &e.ConsumedAt, &e.Note)
		if err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, rows.Err()
}

// Save stores an entry.
func (r *EntryRepository) Save(ctx context.Context, e *water.Entry) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO WaterIntakeEntries (id, user_profile_id, amount_ml, consumed_at) VALUES (?, ?, ?, ?)`,
		e.ID,
		e.UserProfileID,
		e.AmountMl,
		e.ConsumedAt,
	)
	return err
}

// Delete removes an entry belonging to a profile.
func (r *EntryRepository) Delete(ctx context.Context, id, profileID string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM WaterIntakeEntries WHERE id = ? AND user_profile_id = ?`,
		id,
		profileID,
	)
	return err
}

// CurrentWeightKg returns the recorded weight of a profile. A NULL weight
// already resolves to the default inside the query, so found only reports
// whether the profile exists.
func (r *EntryRepository) CurrentWeightKg(ctx context.Context, profileID string) (float64, bool, error) {
	var weight float64

	err := r.db.QueryRowContext(
		ctx,
		`SELECT COALESCE(current_weight, 70.0) FROM UserProfiles WHERE id = ?`,
		profileID,
	).Scan(&weight)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, nil
		}
		return 0, false, err
	}

	return weight, true, nil
}
