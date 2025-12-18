package user

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	Create(user *User) error
	GetAll() ([]User, error)
	Get(id string) (*User, error)
	Delete(id string) error
	Update(id string, firstName, lastName, email, phone *string) error
}

type repo struct {
	log *log.Logger
	db  *gorm.DB
}

func NewRepo(log *log.Logger, db *gorm.DB) Repository {
	return &repo{
		log: log,
		db:  db,
	}
}

func (repo *repo) Create(user *User) error {
	user.ID = uuid.New().String()

	if err := repo.db.Create(user).Error; err != nil {
		repo.log.Println(err)
		return err
	}
	repo.log.Println("user created with id: ", user.ID)
	return nil
}

func (repo *repo) GetAll() ([]User, error) {
	var u []User

	result := repo.db.Model(&u).Order("created_at desc").Find(&u)

	if result.Error != nil {
		return nil, result.Error
	}

	return u, nil
}

func (repo *repo) Get(id string) (*User, error) {
	user := User{ID: id}

	result := repo.db.First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}

// Delete físico y completo
func (repo *repo) Delete(id string) error {
	user := User{
		ID: id,
	}
	result := repo.db.Delete(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (repo *repo) Update(id string, firstName, lastName, email, phone *string) error {
	values := make(map[string]interface{})

	if firstName != nil {
		values["firstName"] = firstName
	}
	if lastName != nil {
		values["lastName"] = lastName
	}
	if email != nil {
		values["email"] = email
	}
	if phone != nil {
		values["phone"] = phone
	}

	res := repo.db.Model(&User{}).Where("id = ?", id).Updates(values)
	if res.Error != nil {
		return res.Error
	}
	return nil
}
