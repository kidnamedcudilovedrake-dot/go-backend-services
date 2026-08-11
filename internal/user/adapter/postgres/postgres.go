package postgres

import (
	"context"
	"encoding/json"

	"antiscam-simulator/internal/user/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) Create(ctx context.Context, user userdomain.User) error {
	query := `INSERT INTO users (id , username , created_at) VALUES ($1 , $2 , $3)`
	_, err := r.pool.Exec(ctx, query, user.ID, user.Username, user.CreatedAt)
	return err
}

func (r *UserRepository) SaveTrainingResult(ctx context.Context, result *userdomain.TrainingResult) error {
	tagsJSON, err := json.Marshal(result.Tags)
	if err != nil {
		return err
	}
	query := `INSERT INTO training_results (session_id, user_id, scenario_id, total_risk, final_grade, tags) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (session_id) DO NOTHING`
	_, err = r.pool.Exec(ctx, query, result.SessionID, result.UserID, result.ScenarioID, result.TotalRisk, result.FinalGrade, tagsJSON)
	return err
}

func (r *UserRepository) UserExists(ctx context.Context, userID string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
	err := r.pool.QueryRow(ctx, query, userID).Scan(&exists)
	return exists, err
}

func (r *UserRepository) GetHistory(ctx context.Context, userID string) ([]userdomain.TrainingHistoryItem, error) {
	exists, err := r.UserExists(ctx, userID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, userdomain.ErrUserNotFound
	}

	query := `SELECT session_id, scenario_id, total_risk, final_grade, tags, created_at FROM training_results WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []userdomain.TrainingHistoryItem
	for rows.Next() {
		var item userdomain.TrainingHistoryItem
		var tagsBytes []byte
		if err := rows.Scan(&item.SessionID, &item.ScenarioID, &item.TotalRisk, &item.FinalGrade, &tagsBytes, &item.CreatedAt); err != nil {
			return nil, err
		}
		if tagsBytes != nil {
			_ = json.Unmarshal(tagsBytes, &item.Tags)
		}
		history = append(history, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if history == nil {
		history = []userdomain.TrainingHistoryItem{}
	}

	return history, nil
}
