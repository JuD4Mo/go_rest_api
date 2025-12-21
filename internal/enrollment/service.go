package enrollment

import (
	"errors"
	"log"

	"github.com/JuD4Mo/go_rest_api/internal/course"
	"github.com/JuD4Mo/go_rest_api/internal/domain"
	"github.com/JuD4Mo/go_rest_api/internal/user"
)

type (
	Service interface {
		Create(userId, courseId string) (*domain.Enrollment, error)
	}

	service struct {
		log           *log.Logger
		userService   user.Service   //Inyección de dependencias
		courseService course.Service //Inyección de dependencias
		repo          Repository
	}
)

func NewService(log *log.Logger, repo Repository, userService user.Service, courseService course.Service) Service {
	return &service{
		log:           log,
		repo:          repo,
		userService:   userService,
		courseService: courseService,
	}
}

func (s service) Create(userId, courseId string) (*domain.Enrollment, error) {
	enroll := &domain.Enrollment{
		UserId:   userId,
		CourseId: courseId,
		Status:   "P",
	}

	_, err := s.userService.Get(userId)
	if err != nil {
		return nil, errors.New("user id does not exists")
	}

	_, err = s.courseService.Get(courseId)
	if err != nil {
		return nil, errors.New("course id does not exists")
	}

	err = s.repo.Create(enroll)
	if err != nil {
		return nil, err
	}

	return enroll, nil
}
