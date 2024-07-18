package GoAsistent

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

type SubstitutionMetadata struct {
	Date       string `json:"date"`
	ScheduleID int    `json:"schedule_id"`
	EventSlug  string `json:"event_slug"`
	Link       string `json:"link"`
}

type GradeMetadata struct {
	GradeID   int    `json:"gradeId"`
	SubjectID int    `json:"subjectId"`
	Link      string `json:"link"`
}

type CommunicationMetadata struct {
	ChannelID   string `json:"channelId"`
	MessageID   int    `json:"messageId"`
	UserID      int    `json:"userId"`
	ChannelType string `json:"channelType"`
	Link        string `json:"link"`
}

type NotificationsResponse struct {
	Status    any   `json:"status"`
	Message   any   `json:"message"`
	Errfields []any `json:"errfields"`
	Data      []struct {
		Title     string          `json:"title"`
		Message   string          `json:"message"`
		ID        int             `json:"id"`
		CreatedAt string          `json:"created_at"`
		Seen      bool            `json:"seen"`
		Type      string          `json:"type"`
		Metadata  json.RawMessage `json:"meta_data"`
	} `json:"data"`
}

// Notification je javen izhod iz funkcije GetNotifications
type Notification struct {
	Type      int
	GradeData struct {
		GradeID     int
		SubjectID   int
		SubjectName string
		Grade       int
		GradeType   int
	}
	SubstitutionData  SubstitutionMetadata
	CommunicationData CommunicationMetadata
}

const (
	GradeNotification         = iota
	SubstitutionNotification  = iota
	CommunicationNotification = iota
)

const (
	PisnaOcena = iota
	UstnaOcena = iota
	DrugaOcena = iota
)

func (s *sessionImpl) GetNotifications() ([]Notification, error) {
	if s.TokenExpiration < int(time.Now().Unix()) {
		err := s.RefreshSession()
		if err != nil {
			return nil, err
		}
		defer s.RefreshTokenCallback(s.Username, s.RefreshToken)
	}

	res, err := s.Client.R().Get(fmt.Sprintf("%s/notifications/ajax_web_notifications_get", EASISTENT_URL))
	if err != nil {
		return nil, err
	}
	var response NotificationsResponse
	err = res.UnmarshalJson(&response)
	if err != nil {
		return nil, err
	}

	notifications := make([]Notification, 0)

	for _, v := range response.Data {
		if v.Type == "ocena" {
			var grade GradeMetadata
			err := json.Unmarshal(v.Metadata, &grade)
			if err != nil {
				return nil, err
			}
			r := regexp.MustCompile(`(?P<Grade>[1-5]) - (?P<Subject>.*), (?P<GradeType>.*)`)
			submatches := r.FindStringSubmatch(v.Message)
			if len(submatches) != 4 {
				return nil, errors.New(fmt.Sprintf("invalid grade – cannot parse grade notification: %s", v.Message))
			}
			g, err := strconv.Atoi(submatches[1])
			if err != nil {
				return nil, err
			}
			gt := DrugaOcena
			if submatches[3] == "Pisna ocena" {
				gt = PisnaOcena
			} else if submatches[3] == "Ustna ocena" {
				gt = UstnaOcena
			}
			n := Notification{
				Type: GradeNotification,
				GradeData: struct {
					GradeID     int
					SubjectID   int
					SubjectName string
					Grade       int
					GradeType   int
				}{
					GradeID:     grade.GradeID,
					SubjectID:   grade.SubjectID,
					SubjectName: submatches[2],
					Grade:       g,
					GradeType:   gt,
				},
			}
			notifications = append(notifications, n)
			continue
		}
		if v.Type == "substitution_students" {
			var substitution SubstitutionMetadata
			err := json.Unmarshal(v.Metadata, &substitution)
			if err != nil {
				return nil, err
			}
			n := Notification{
				Type:             SubstitutionNotification,
				SubstitutionData: substitution,
			}
			notifications = append(notifications, n)
			continue
		}
		if v.Type == "komunikacija" {
			var communication CommunicationMetadata
			err := json.Unmarshal(v.Metadata, &communication)
			if err != nil {
				return nil, err
			}
			n := Notification{
				Type:              CommunicationNotification,
				CommunicationData: communication,
			}
			notifications = append(notifications, n)
			continue
		}
	}
	return notifications, nil
}
