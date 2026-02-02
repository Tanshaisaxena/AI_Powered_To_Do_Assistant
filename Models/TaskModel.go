package models

import "time"

type Task struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Tags        string    `json:"tags"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdat"`
	DueDate     time.Time `json:"duedate"`
	Priority    string    `json:"priority"`
	Notes       string    `json:"notes"`
	Embedding   []float64 `json:"-"`
}

type Ask struct {
	Question string `json:"question" binding:"required"`
}

type AskResponse struct {
	Tasks    []Task `json:"tasks"`
	Response string `json:"response"`
}

// bool uses memory
// struct{} uses ZERO bytes
var Stopwords = map[string]struct{}{
	"the":     {},
	"me":      {},
	"to":      {},
	"is":      {},
	"are":     {},
	"was":     {},
	"were":    {},
	"tell":    {},
	"related": {},
	"about":   {},
	"tasks":   {},
	"task":    {},
	"please":  {},
	"show":    {},
	"list":    {},
	"of":      {},
}

type ScoreTasks struct {
	Score int  `json:"score"`
	Task  Task `json:"task"`
}

type TaskSimilarityScore struct {
	SimilarityScore float64 `json:"similarityscore"`
	Task            Task    `json:"task"`
}
