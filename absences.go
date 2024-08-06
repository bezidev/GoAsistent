package GoAsistent

import (
	"fmt"
	"time"
)

type AbsencesResponse struct {
	Summary struct {
		PendingHours      int `json:"pending_hours"`
		ExcusedHours      int `json:"excused_hours"`
		UnexcusedHours    int `json:"unexcused_hours"`
		UnmanagedAbsences int `json:"unmanaged_absences"`
	} `json:"summary"`
	Items []struct {
		Id                int     `json:"id"`
		Date              string  `json:"date"`
		MissingCount      int     `json:"missing_count"`
		ExcusedCount      int     `json:"excused_count"`
		NotExcusedCount   int     `json:"not_excused_count"`
		State             string  `json:"state"`
		Seen              bool    `json:"seen"`
		ExcuseSent        bool    `json:"excuse_sent"`
		ExcuseDescription *string `json:"excuse_description"`
		ExcuseWrittenDate any     `json:"excuse_written_date"`
		Hours             []struct {
			ClassName      string `json:"class_name"`
			ClassShortName string `json:"class_short_name"`
			Value          string `json:"value"`
			From           string `json:"from"`
			To             string `json:"to"`
			State          string `json:"state"`
		} `json:"hours"`
		Attachments []any `json:"attachments"`
	} `json:"items"`
}

func (s *sessionImpl) GetAbsences() (AbsencesResponse, error) {
	if s.TokenExpiration < int(time.Now().Unix()) {
		err := s.RefreshSession()
		if err != nil {
			return AbsencesResponse{}, err
		}
		defer s.RefreshTokenCallback(s.Username, s.RefreshToken)
	}

	res, err := s.Client.R().Get(fmt.Sprintf("%s/m/absences", EASISTENT_URL))
	if err != nil {
		return AbsencesResponse{}, err
	}
	var response AbsencesResponse
	err = res.UnmarshalJson(&response)
	if err != nil {
		return AbsencesResponse{}, err
	}
	return response, nil
}
