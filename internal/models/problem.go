package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Modalessi/nuha-api/internal/database"
	"github.com/Modalessi/nuha-api/internal/utils"
	"github.com/google/uuid"
)

type Problem struct {
	ID          uuid.UUID
	Title       string
	Description string
	Difficulty  string
	Tags        []string
	Testcases   []Testcase
	Timelimit   float64
	Memorylimit float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func ProblemFromDBObject(p *database.Problem) *Problem {
	return &Problem{
		ID:          p.ID,
		Title:       p.Title,
		Difficulty:  p.Difficulty,
		Tags:        p.Tags,
		Timelimit:   p.TimeLimit,
		Memorylimit: p.MemoryLimit,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func CreateNewProblem(title string, description string, difficulty string, tags []string) (*Problem, error) {

	if difficulty != "HARD" && difficulty != "MEDIUM" && difficulty != "EASY" {
		return nil, fmt.Errorf("difficutly must be one of these (HARD, MEDIUM, EASY)")
	}

	return &Problem{
		ID:          uuid.New(),
		Title:       title,
		Description: description,
		Difficulty:  difficulty,
		Tags:        []string{},
		Testcases:   []Testcase{},
		Timelimit:   1,
		Memorylimit: 128000,
	}, nil
}

func (p *Problem) SetTitle(title string) {
	p.Title = title
}

func (p *Problem) SetDifficulty(difficulty string) error {
	if difficulty != "HARD" && difficulty != "MEDIUM" && difficulty != "EASY" {
		return fmt.Errorf("difficutly must be one of these (HARD, MEDIUM, EASY)")
	}

	p.Difficulty = difficulty
	return nil
}

func (p *Problem) SetDescription(description string) {
	p.Description = description
}

func (p *Problem) AddTags(tags []string) {
	p.Tags = append(p.Tags, tags...)
}

func (p *Problem) SetTimelimit(timelimt float64) {
	p.Timelimit = timelimt
}

func (p *Problem) SetMemoryLimit(memorylimit float64) {
	p.Memorylimit = memorylimit
}

func (p *Problem) JSON() []byte {
	data, err := json.Marshal(p)
	utils.Assert(err, "error converting problem object to json")
	return data
}
