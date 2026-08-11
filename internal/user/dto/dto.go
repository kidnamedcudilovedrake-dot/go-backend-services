package userdto

import (
	"time"

	"antiscam-simulator/internal/user/domain"
)

type RegisterRequest struct {
	Username string `json:"username"`
}

type RegisterResponse struct {
	UserID string `json:"user_id"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type TagDto struct {
	Question    string `json:"question"`
	Answer      string `json:"answer"`
	Explanation string `json:"explanation"`
}


type TrainingHistoryItemDto struct {
	SessionID  string    `json:"session_id"`
	ScenarioID string    `json:"scenario_id"`
	TotalRisk  int32     `json:"total_risk"`
	FinalGrade string    `json:"final_grade"`
	Tags       []TagDto  `json:"tags"`
	CreatedAt  time.Time `json:"created_at"`
}

type UserHistoryResponse struct {
	UserID  string                   `json:"user_id"`
	History []TrainingHistoryItemDto `json:"history"`
}

func MapHistoryFromDomain(items []userdomain.TrainingHistoryItem) []TrainingHistoryItemDto {
	dtoItems := make([]TrainingHistoryItemDto, 0, len(items))
	for _, item := range items {
		dtoTags := make([]TagDto, 0, len(item.Tags))
		for _, tag := range item.Tags {
			dtoTags = append(dtoTags, TagDto{
				Question:    tag.Question,
				Answer:      tag.Answer,
				Explanation: tag.Explanation,
			})
		}
		dtoItems = append(dtoItems, TrainingHistoryItemDto{
			SessionID:  item.SessionID,
			ScenarioID: item.ScenarioID,
			TotalRisk:  item.TotalRisk,
			FinalGrade: item.FinalGrade,
			Tags:       dtoTags,
			CreatedAt:  item.CreatedAt,
		})
	}
	if dtoItems == nil {
		dtoItems = []TrainingHistoryItemDto{}
	}
	return dtoItems
}
