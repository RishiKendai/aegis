package stream

import (
	"strconv"

	"github.com/RishiKendai/aegis/internal/models"
)

type StreamMessage struct {
	ID     string
	Fields map[string]string
}

func ParseSubmission(msg *StreamMessage) (*models.Submission, error) {
	submission := &models.Submission{
		AttemptID:  msg.Fields["attemptId"],
		SourceCode: msg.Fields["sourceCode"],
		Language:   msg.Fields["language"],
		LangCode:   msg.Fields["langCode"],
		Email:      msg.Fields["email"],
		TestID:     msg.Fields["testId"],
		DriveID:    msg.Fields["driveId"],
		Difficulty: msg.Fields["difficulty"],
	}

	qidStr := msg.Fields["qId"]
	if qidStr == "" {
		return nil, ErrMissingField("qId")
	}
	qid, err := strconv.ParseInt(qidStr, 10, 64)
	if err != nil {
		return nil, err
	}
	submission.QID = qid

	// Validate required fields
	if submission.AttemptID == "" {
		return nil, ErrMissingField("attemptId")
	}
	if submission.SourceCode == "" {
		return nil, ErrMissingField("sourceCode")
	}
	if submission.Email == "" {
		return nil, ErrMissingField("email")
	}
	if submission.DriveID == "" {
		return nil, ErrMissingField("driveId")
	}

	return submission, nil
}

type ErrMissingField string

func (e ErrMissingField) Error() string {
	return "missing required field: " + string(e)
}
