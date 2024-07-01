package GoAsistent

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Grading struct {
	ID          int
	Course      string
	GradingName string
	Type        int
	Date        string
	Hour        int
	Test        bool
	Grade       int
}

type GradingResponse struct {
	Items []struct {
		ID       int    `json:"id"`
		Course   string `json:"course"`
		Subject  string `json:"subject"`
		Type     string `json:"type"`
		Date     string `json:"date"`
		Period   string `json:"period"`
		Test     bool   `json:"test"` // Kaj zaboga naj bi ta vrednost pomenila/povedala? Še en duplicate za pisno oceno???
		Grade    string `json:"grade"`
		TypeName string `json:"type_name"`
	} `json:"items"`
}

func (s *sessionImpl) GetGradings(past bool) ([]Grading, error) {
	if s.TokenExpiration < int(time.Now().Unix()) {
		err := s.RefreshSession()
		if err != nil {
			return nil, err
		}
	}

	filter := "future"
	if past {
		filter = "past"
	}

	res, err := s.Client.R().Get(fmt.Sprintf("%s/m/evaluations?filter=%s", EASISTENT_URL, filter))
	if err != nil {
		return nil, err
	}
	var response GradingResponse
	err = res.UnmarshalJson(&response)
	if err != nil {
		return nil, err
	}
	gradings := make([]Grading, 0)
	for _, v := range response.Items {
		t := DrugaOcena
		if v.Type == "written" {
			t = PisnaOcena
		} else if v.Type == "oral" {
			t = UstnaOcena
		}
		period, err := strconv.Atoi(strings.ReplaceAll(v.Period, ". ura", ""))
		if err != nil {
			return nil, err
		}
		grade, err := strconv.Atoi(v.Grade)
		if err != nil {
			grade = -1
		}
		gradings = append(gradings, Grading{
			ID:          v.ID,
			Course:      v.Course,
			GradingName: v.Subject,
			Type:        t,
			Date:        v.Date,
			Hour:        period,
			Test:        v.Test,
			Grade:       grade,
		})
	}
	return gradings, nil
}
