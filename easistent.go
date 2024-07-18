package GoAsistent

import (
	"github.com/imroc/req/v3"
	"time"
)

const EASISTENT_URL = "https://www.easistent.com"

type sessionImpl struct {
	AuthToken            string
	RefreshToken         string
	ChildId              string
	TokenExpiration      int
	Username             string
	Name                 string
	DevMode              bool
	Client               *req.Client
	RefreshTokenCallback func(username string, refreshToken string)
}

type Session interface {
	RefreshSession() error
	RefreshWebSession() error
	GetGradings(past bool) ([]Grading, error)
	GetTimetable(startDate time.Time, endDate time.Time) (TimetableResponse, error)
	GetNotifications() ([]Notification, error)
	GetGrades() (GradesResponse, error)
	GetGradesForSubject(subjectId int) (SubjectGradesResponse, error)
	GetSessionData() *sessionImpl
}
