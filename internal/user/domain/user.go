package userdomain

import (
	"errors"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type User struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"created_at"`
}

type Tag struct {
	Question    string
	Answer      string
	Explanation string
}

type TrainingResult struct {
	ID         string    `db:"id"`
	SessionID  string    `db:"session_id"`
	UserID     string    `db:"user_id"`
	ScenarioID string    `db:"scenario_id"`
	TotalRisk  int32     `db:"total_risk"`
	FinalGrade string    `db:"final_grade"`
	Tags       []Tag     `db:"tags"`
	CreatedAt  time.Time `db:"created_at"`
}

type TrainingHistoryItem struct {
	SessionID  string    `db:"session_id"`
	ScenarioID string    `db:"scenario_id"`
	TotalRisk  int32     `db:"total_risk"`
	FinalGrade string    `db:"final_grade"`
	Tags       []Tag     `db:"tags"`
	CreatedAt  time.Time `db:"created_at"`
}
