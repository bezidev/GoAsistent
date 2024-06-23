package GoAsistent

import "slices"

type Subject struct {
	ID       string
	Name     string
	Teachers []string
}

func ExtractSubjectsFromTimetable(timetable *TimetableResponse) (map[string]*Subject, error) {
	t := *timetable
	subjects := make(map[string]*Subject)
	for _, v := range t.SchoolHourEvents {
		_, exists := subjects[v.Subject.ID]
		if !exists {
			subjects[v.Subject.ID] = &Subject{
				ID:       v.Subject.ID,
				Name:     v.Subject.Name,
				Teachers: make([]string, 0),
			}
		}
		for _, k := range v.Teachers {
			if slices.Contains(subjects[v.Subject.ID].Teachers, k.Name) {
				continue
			}
			subjects[v.Subject.ID].Teachers = append(subjects[v.Subject.ID].Teachers, k.Name)
		}
	}
	return subjects, nil
}
