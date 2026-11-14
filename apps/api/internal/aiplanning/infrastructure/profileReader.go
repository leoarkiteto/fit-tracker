package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"fittracker-api/internal/aiplanning"
)

// ProfileReader reads the profile columns a plan is generated from.
type ProfileReader struct {
	db *sql.DB
}

var _ aiplanning.ProfileRepository = (*ProfileReader)(nil)

// NewProfileReader builds the profile reader.
func NewProfileReader(db *sql.DB) *ProfileReader {
	return &ProfileReader{db: db}
}

// GetProfile returns the planning-relevant fields of a profile.
func (r *ProfileReader) GetProfile(
	ctx context.Context,
	profileID string,
) (*aiplanning.Profile, bool, error) {
	var profile aiplanning.Profile

	err := r.db.QueryRowContext(
		ctx,
		`SELECT name, age, current_weight FROM UserProfiles WHERE id = ?`,
		profileID,
	).Scan(&profile.Name, &profile.Age, &profile.CurrentWeight)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &profile, true, nil
}
