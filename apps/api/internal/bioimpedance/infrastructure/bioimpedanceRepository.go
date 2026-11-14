// Package infrastructure implements the bioimpedance slice outbound ports on
// SQLite.
package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"fittracker-api/internal/bioimpedance"
)

// measurementColumns is the column list the previous implementation selected.
const measurementColumns = `id, date, weight, body_fat_percentage, muscle_mass, bone_mass, water_percentage, visceral_fat, bmr, metabolic_age, notes`

// MeasurementRepository is the SQLite implementation of bioimpedance.Repository.
type MeasurementRepository struct {
	db *sql.DB
}

var _ bioimpedance.Repository = (*MeasurementRepository)(nil)

// NewMeasurementRepository builds the bioimpedance repository.
func NewMeasurementRepository(db *sql.DB) *MeasurementRepository {
	return &MeasurementRepository{db: db}
}

// GetByProfile returns every measurement of a profile, newest first.
func (r *MeasurementRepository) GetByProfile(
	ctx context.Context,
	profileID string,
) ([]bioimpedance.Measurement, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT `+measurementColumns+` FROM BioimpedanceData WHERE user_profile_id = ? ORDER BY date DESC`,
		profileID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []bioimpedance.Measurement
	for rows.Next() {
		var m bioimpedance.Measurement
		err := rows.Scan(
			&m.ID,
			&m.Date,
			&m.Weight,
			&m.BodyFatPercentage,
			&m.MuscleMass,
			&m.BoneMass,
			&m.WaterPercentage,
			&m.VisceralFat,
			&m.BMR,
			&m.MetabolicAge,
			&m.Notes,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, m)
	}

	return history, rows.Err()
}

// GetLatest returns the most recent measurement of a profile.
func (r *MeasurementRepository) GetLatest(
	ctx context.Context,
	profileID string,
) (*bioimpedance.Measurement, bool, error) {
	var m bioimpedance.Measurement

	err := r.db.QueryRowContext(
		ctx,
		`SELECT `+measurementColumns+` FROM BioimpedanceData WHERE user_profile_id = ? ORDER BY date DESC LIMIT 1`,
		profileID,
	).Scan(
		&m.ID,
		&m.Date,
		&m.Weight,
		&m.BodyFatPercentage,
		&m.MuscleMass,
		&m.BoneMass,
		&m.WaterPercentage,
		&m.VisceralFat,
		&m.BMR,
		&m.MetabolicAge,
		&m.Notes,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &m, true, nil
}

// Save stores a measurement.
func (r *MeasurementRepository) Save(ctx context.Context, m *bioimpedance.Measurement) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO BioimpedanceData (id, date, weight, body_fat_percentage, muscle_mass, bone_mass, water_percentage, visceral_fat, bmr, metabolic_age, notes, user_profile_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.ID,
		m.Date,
		m.Weight,
		m.BodyFatPercentage,
		m.MuscleMass,
		m.BoneMass,
		m.WaterPercentage,
		m.VisceralFat,
		m.BMR,
		m.MetabolicAge,
		m.Notes,
		m.UserProfileID,
	)
	return err
}

// Delete removes a measurement belonging to a profile.
func (r *MeasurementRepository) Delete(ctx context.Context, id, profileID string) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM BioimpedanceData WHERE id = ? AND user_profile_id = ?`,
		id,
		profileID,
	)
	return err
}
