package course

import (
	"log"
	"time"
)

type (
	Service interface {
		Create(name string, startDate, endDate string) (*Course, error)
		Get(id string) (*Course, error)
		GetAll(filters Filters, offset, limit int) ([]Course, error)
		Update(id string, name, startDate, endDate *string) error
		Delete(id string) error
		Count(filters Filters) (int, error)
	}

	service struct {
		log  *log.Logger
		repo Repository
	}

	Filters struct {
		Name string
	}
)

func NewService(log *log.Logger, repo Repository) Service {
	return &service{
		log:  log,
		repo: repo,
	}
}

func (s service) Create(name string, startDate, endDate string) (*Course, error) {

	startDateParsed, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return nil, err
	}

	endDateParsed, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return nil, err
	}

	course := &Course{
		Name:      name,
		StartDate: startDateParsed,
		EndDate:   endDateParsed,
	}
	err = s.repo.Create(course)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (s service) Get(id string) (*Course, error) {
	course, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return course, nil
}

func (s service) GetAll(filters Filters, offset, limit int) ([]Course, error) {
	courses, err := s.repo.GetAll(filters, offset, limit)
	if err != nil {
		return nil, err
	}

	return courses, nil
}

func (s service) Update(id string, name, startDate, endDate *string) error {

	var startDateParsed *time.Time
	var endDateParsed *time.Time

	if startDate != nil && *startDate != "" {
		date, err := time.Parse("2006-01-02", *startDate)
		if err != nil {
			return err
		}
		startDateParsed = &date
	}

	if endDate != nil && *endDate != "" {
		date, err := time.Parse("2006-01-02", *endDate)
		if err != nil {
			return err
		}
		endDateParsed = &date
	}

	return s.repo.Update(id, name, startDateParsed, endDateParsed)
}

func (s service) Delete(id string) error {
	return s.repo.Delete(id)
}

func (s service) Count(filters Filters) (int, error) {
	num, err := s.repo.Count(filters)
	if err != nil {
		return 0, err
	}
	return num, nil
}
