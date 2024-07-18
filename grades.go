package GoAsistent

import (
	"fmt"
	"time"
)

/*type Grade struct {
	ID          int
	Course      string
	GradingName string
	Type        int
	Date        string
	Hour        int
	Test        bool
	Grade       int
}*/

type GradesResponse struct {
	Items []struct {
		Name       string `json:"name"`
		ShortName  string `json:"short_name"`
		ID         int    `json:"id"`
		GradeType  string `json:"grade_type"`
		IsExcused  bool   `json:"is_excused"`
		FinalGrade *struct {
			Value string `json:"value"`
		} `json:"final_grade"`
		AverageGrade string `json:"average_grade"`
		GradeRank    string `json:"grade_rank"`
		Semesters    []struct {
			ID         int `json:"id"`
			FinalGrade any `json:"final_grade"`
			Grades     []struct {
				TypeName     string  `json:"type_name"`
				Comment      *string `json:"comment"`
				ID           int     `json:"id"`
				Type         string  `json:"type"`
				OverridesIds any     `json:"overrides_ids"`
				Value        string  `json:"value"`
				Color        string  `json:"color"`
				Date         string  `json:"date"`
			} `json:"grades"`
		} `json:"semesters"`
	} `json:"items"`
}

type SubjectGradesResponse struct {
	Name       string `json:"name"`
	ShortName  string `json:"short_name"`
	ID         int    `json:"id"`
	GradeType  string `json:"grade_type"`
	IsExcused  bool   `json:"is_excused"`
	FinalGrade *struct {
		ID         int    `json:"id"`
		Value      string `json:"value"`
		Date       string `json:"date"`
		InsertedBy *struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"inserted_by"`
	} `json:"final_grade"`
	AverageGrade string `json:"average_grade"`
	GradeRank    string `json:"grade_rank"`
	Semesters    []struct {
		ID         int `json:"id"`
		FinalGrade any `json:"final_grade"`
		Grades     []struct {
			TypeName     string `json:"type_name"`
			Comment      any    `json:"comment"`
			ID           int    `json:"id"`
			OverridesIds any    `json:"overrides_ids"`
			Value        string `json:"value"`
			Subject      any    `json:"subject"`
			Date         string `json:"date"`
			NotifiedAt   any    `json:"notified_at"`
			InsertedAt   string `json:"inserted_at"`
			InsertedBy   *struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"inserted_by"`
			GradeRank string `json:"grade_rank"`
		} `json:"grades"`
	} `json:"semesters"`
}

func (s *sessionImpl) GetGrades() (GradesResponse, error) {
	if s.TokenExpiration < int(time.Now().Unix()) {
		err := s.RefreshSession()
		if err != nil {
			return GradesResponse{}, err
		}
		defer s.RefreshTokenCallback(s.Username, s.RefreshToken)
	}

	res, err := s.Client.R().Get(fmt.Sprintf("%s/m/grades", EASISTENT_URL))
	if err != nil {
		return GradesResponse{}, err
	}
	var response GradesResponse
	err = res.UnmarshalJson(&response)
	if err != nil {
		return GradesResponse{}, err
	}
	return response, nil
}

func (s *sessionImpl) GetGradesForSubject(subjectId int) (SubjectGradesResponse, error) {
	if s.TokenExpiration < int(time.Now().Unix()) {
		err := s.RefreshSession()
		if err != nil {
			return SubjectGradesResponse{}, err
		}
		defer s.RefreshTokenCallback(s.Username, s.RefreshToken)
	}

	res, err := s.Client.R().Get(fmt.Sprintf("%s/m/grades/classes/%d", EASISTENT_URL, subjectId))
	if err != nil {
		return SubjectGradesResponse{}, err
	}
	var response SubjectGradesResponse
	err = res.UnmarshalJson(&response)
	if err != nil {
		return SubjectGradesResponse{}, err
	}
	return response, nil
}
