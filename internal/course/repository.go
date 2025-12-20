package course

import (
	"fmt"
	"log"
	"strings"

	"gorm.io/gorm"
)

type (
	Repository interface {
		Create(course *Course) error
		Get(id string) (*Course, error)
		GetAll(filters Filters, offset, limit int) ([]Course, error)
		Update()
		Delete()
		Count(filters Filters) (int, error)
	}

	repo struct {
		log *log.Logger
		db  *gorm.DB
	}
)

func NewRepo(db *gorm.DB, log *log.Logger) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(course *Course) error {
	result := repo.db.Create(course)
	if result.Error != nil {
		return result.Error
	}
	repo.log.Println("course created with id: ", course.ID)
	return nil
}

func (repo *repo) Get(id string) (*Course, error) {
	//Como el id está poblado y GORM detecta que es la PK filtra por ese parámetro
	course := Course{
		ID: id,
	}
	// result := repo.db.Model(&Course{}).Where("id = ?", id).First(&course)
	result := repo.db.First(&course)
	if result.Error != nil {
		return nil, result.Error
	}

	return &course, nil
}

func (repo *repo) GetAll(filters Filters, offset, limit int) ([]Course, error) {
	var courses []Course
	tx := repo.db.Model(&courses)
	tx = applyFilters(tx, filters)
	tx = tx.Limit(limit).Offset(offset)
	result := tx.Order("created_at DESC").Find(&courses)

	if result.Error != nil {
		return nil, result.Error
	}

	return courses, nil
}

func (repo *repo) Update() {

}

func (repo *repo) Delete() {

}

func (repo *repo) Count(filters Filters) (int, error) {
	var count int64
	tx := repo.db.Model(Course{})
	tx = applyFilters(tx, filters)

	result := tx.Count(&count)
	if result.Error != nil {
		return 0, result.Error
	}

	return int(count), nil
}

func applyFilters(tx *gorm.DB, filters Filters) *gorm.DB {
	if filters.Name != "" {
		filters.Name = fmt.Sprintf("%%%s%%", strings.ToLower(filters.Name))
		tx = tx.Where("lower(first_name) like ?", filters.Name)
	}

	return tx
}
