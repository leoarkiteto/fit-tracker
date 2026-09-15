// Package infrastructure implements the profile slice outbound ports on SQLite.
package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"fittracker-api/internal/profile"
)

// profileColumns is the column list the previous implementation selected.
const profileColumns = `id, name, age, height, current_weight, goal_weight,
	avatar_url, experience_level, available_days_per_week,
	preferred_workout_duration, equipment_type`

// ProfileRepository is the SQLite implementation of profile.Repository.
type ProfileRepository struct {
	db *sql.DB
}

var _ profile.Repository = (*ProfileRepository)(nil)

// NewProfileRepository builds the profile repository.
func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// GetAll returns every profile.
func (r *ProfileRepository) GetAll(ctx context.Context) ([]profile.Profile, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT `+profileColumns+` FROM UserProfiles`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []profile.Profile
	for rows.Next() {
		var p profile.Profile
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Age,
			&p.Height,
			&p.CurrentWeight,
			&p.GoalWeight,
			&p.AvatarURL,
			&p.ExperienceLevel,
			&p.AvailableDaysPerWeek,
			&p.PreferredWorkoutDuration,
			&p.EquipmentType,
		)
		if err != nil {
			return nil, err
		}
		profiles = append(profiles, p)
	}

	return profiles, rows.Err()
}

// GetByID returns a profile by id.
func (r *ProfileRepository) GetByID(ctx context.Context, id string) (*profile.Profile, bool, error) {
	var p profile.Profile

	err := r.db.QueryRowContext(
		ctx,
		`SELECT `+profileColumns+` FROM UserProfiles WHERE id = ?`,
		id,
	).Scan(
		&p.ID,
		&p.Name,
		&p.Age,
		&p.Height,
		&p.CurrentWeight,
		&p.GoalWeight,
		&p.AvatarURL,
		&p.ExperienceLevel,
		&p.AvailableDaysPerWeek,
		&p.PreferredWorkoutDuration,
		&p.EquipmentType,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, false, nil
		}
		return nil, false, err
	}

	return &p, true, nil
}

// Update stores every editable field of a profile.
func (r *ProfileRepository) Update(ctx context.Context, p *profile.Profile) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE UserProfiles SET name = ?, age = ?, height = ?, current_weight = ?, goal_weight = ?, avatar_url = ?, experience_level = ?, available_days_per_week = ?, preferred_workout_duration = ?, equipment_type = ?, updated_at = ? WHERE id = ?`,
		p.Name,
		p.Age,
		p.Height,
		p.CurrentWeight,
		p.GoalWeight,
		p.AvatarURL,
		p.ExperienceLevel,
		p.AvailableDaysPerWeek,
		p.PreferredWorkoutDuration,
		p.EquipmentType,
		time.Now(),
		p.ID,
	)
	return err
}

// Delete removes the profile and every row referencing it. Exercises are
// cascaded by the Workouts delete (ON DELETE CASCADE in the schema).
func (r *ProfileRepository) Delete(ctx context.Context, id string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	statements := []struct {
		query string
		arg   string
	}{
		{`DELETE FROM CompletedWorkouts WHERE user_profile_id = ?`, id},
		{`DELETE FROM WaterIntakeEntries WHERE user_profile_id = ?`, id},
		{`DELETE FROM BioimpedanceData WHERE user_profile_id = ?`, id},
		{`DELETE FROM Workouts WHERE user_profile_id = ?`, id},
		{`UPDATE Users SET user_profile_id = NULL WHERE user_profile_id = ?`, id},
		{`DELETE FROM UserProfiles WHERE id = ?`, id},
	}

	for _, stmt := range statements {
		if _, err = tx.ExecContext(ctx, stmt.query, stmt.arg); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}
