package models

import "time"

// Artifact represents a plagiarism artifact stored in MongoDB
type Artifact struct {
	Email            string        `bson:"email" json:"email"`
	AttemptID        string        `bson:"attemptId" json:"attemptId"`
	TestID           string        `bson:"testId" json:"testId"`
	DriveID          string        `bson:"driveId" json:"driveId"`
	Difficulty       string        `bson:"difficulty" json:"difficulty"`
	SourceCode       string        `bson:"sourceCode" json:"sourceCode"`
	QID              int64         `bson:"qId" json:"qId"`
	Language         string        `bson:"language" json:"language"`
	LangCode         string        `bson:"langCode" json:"langCode"`
	Tokens           []string      `bson:"tokens" json:"tokens"`
	NormalizedTokens []string      `bson:"normalizedTokens" json:"normalizedTokens"`
	AST              *ASTNode      `bson:"ast" json:"ast"`
	CFG              *CFG          `bson:"cfg" json:"cfg"`
	Fingerprints     *Fingerprints `bson:"fingerprints" json:"fingerprints"`
	CreatedAt        time.Time     `bson:"createdAt" json:"createdAt"`
}

// CandidateResult represents a candidate's plagiarism result
type CandidateResult struct {
	Email           string              `bson:"email" json:"email"`
	AttemptID       string              `bson:"attemptId" json:"attemptId"`
	DriveID         string              `bson:"driveId" json:"driveId"`
	Risk            string              `bson:"risk" json:"risk"` // clean, suspicious, highly suspicious, Near copy
	FlaggedQN       []string            `bson:"flagged_qn" json:"flagged_qn"`
	PlagiarismPeers map[string][]string `bson:"plagiarism_peers" json:"plagiarism_peers"` // qId -> []attemptId
	CodeSimilarity  int                 `bson:"code_similarity" json:"code_similarity"`
	AlgoSimilarity  int                 `bson:"algo_similarity" json:"algo_similarity"`
	Status          string              `bson:"status" json:"status"` // pending, completed, failed
	CreatedAt       time.Time           `bson:"createdAt" json:"createdAt"`
}

// TestReport represents an overall test plagiarism report
type TestReport struct {
	DriveID   string    `bson:"driveId" json:"driveId"`
	Risk      string    `bson:"risk" json:"risk"`     // Safe, Moderate, High, Critical
	Status    string    `bson:"status" json:"status"` // pending, completed, failed
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
	FlaggedQN []string  `bson:"flagged_qn" json:"flagged_qn"`
}

// ComputeRequest represents a request to compute plagiarism
type ComputeRequest struct {
	DriveID string `json:"driveId" binding:"required"`
}

// ComputeResponse represents the response from compute endpoint
type ComputeResponse struct {
	Message string `json:"message"`
	TestID  string `json:"test_id"`
}

// ErrorResponse represents a standard error response
type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}
