package GoAsistent

import (
	"fmt"
	"time"
)

type AllDayEvent struct {
	ID       int    `json:"id"`
	Date     string `json:"date"`
	Location struct {
		Name      any `json:"name"`
		ShortName any `json:"short_name"`
	} `json:"location"`
	// Teachers: TBD
	Teachers []any  `json:"teachers"`
	Name     string `json:"name"`
	// EventType je:
	// 4 – prazniki/počitnice
	// 6 – celodnevne aktivnosti, ki jih določi šola
	// ostalo – ????
	EventType int `json:"event_type"`
}

type TimetableHour struct {
	Time struct {
		FromID int    `json:"from_id"`
		ToID   int    `json:"to_id"`
		Date   string `json:"date"`
	} `json:"time"`
	// Videokonferenca je useless
	Videokonferenca struct {
		ID      any `json:"id"`
		Link    any `json:"link"`
		Zacetek any `json:"zacetek"`
		Opomba  any `json:"opomba"`
	} `json:"videokonferenca"`
	EventID int    `json:"event_id"`
	Color   string `json:"color"`
	Subject struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"subject"`
	Completed bool `json:"completed"`
	// HourSpecialType je lahko
	// - exam
	// - substitution
	// - pre-exam (preverjanje)
	// - cancelled (odpadla)
	// - prazen string (eAsistent pošlje null)
	//
	// Zanimivo, da eAsistent ne podpira exama in substitutiona hkrati ... slab design choice ig
	HourSpecialType string `json:"hour_special_type"`
	Departments     []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"departments"`
	Classroom struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"classroom"`
	Teachers []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"teachers"`
	Groups []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"groups"`
	Info []any `json:"info"`
}

type TimetableResponse struct {
	// TimeTable nam pove vse ure zahtevanega tedna (in še preveč???). To je torej najbolj levi stolpec v eAsistent sistemu.
	TimeTable []struct {
		ID        int    `json:"id"`
		Name      string `json:"name"`
		NameShort string `json:"name_short"`
		Time      struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"time"`
		Type string `json:"type"`
	} `json:"time_table"`
	// DayTable nam pove vse datume zahtevanega tedna. Torej gornja vrstica v eAsistent sistemu.
	DayTable []struct {
		Name      string `json:"name"`
		ShortName string `json:"short_name"`
		Date      string `json:"date"`
	} `json:"day_table"`
	// SchoolHourEvents nam pove dejanske ure (predmete) na urniku. Končno nekaj uporabnega.
	SchoolHourEvents []TimetableHour `json:"school_hour_events"`
	// Events: TBD
	Events       []any         `json:"events"`
	AllDayEvents []AllDayEvent `json:"all_day_events"`
}

func (s *sessionImpl) GetTimetable(startDate time.Time, endDate time.Time) (TimetableResponse, error) {
	if s.TokenExpiration < int(time.Now().Unix()) {
		err := s.RefreshSession()
		if err != nil {
			return TimetableResponse{}, err
		}
	}

	res, err := s.Client.R().Get(fmt.Sprintf("%s/m/timetable/weekly?from=%s&to=%s", EASISTENT_URL, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")))
	if err != nil {
		return TimetableResponse{}, err
	}
	var response TimetableResponse
	err = res.UnmarshalJson(&response)
	if err != nil {
		return TimetableResponse{}, err
	}
	return response, nil
}
