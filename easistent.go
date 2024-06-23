package GoAsistent

import "time"

const EASISTENT_URL = "https://www.easistent.com"

type sessionImpl struct {
	AuthToken       string
	RefreshToken    string
	ChildId         string
	TokenExpiration int
	Username        string
	Name            string
}

type Session interface {
	RefreshSession() error
	GetGradings(past bool) ([]Grading, error)
	GetTimetable(startDate time.Time, endDate time.Time) (TimetableResponse, error)
	GetNotifications() ([]Notification, error)
	GetGrades() (GradesResponse, error)
	GetGradesForSubject(subjectId int) (SubjectGradesResponse, error)
	GetSessionData() *sessionImpl
}
