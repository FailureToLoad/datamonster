package user

import (
	"context"
	"datamonster/user/repo"
	"log"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repo.PostGres
}

func NewService(repo *repo.PostGres) *UserService {
	return &UserService{repo: repo}
}

func (s UserService) ValidateCredentials(ctx context.Context, username, password string) (userId int, err error) {
	user, fetchErr := s.repo.Get(ctx, username)
	if fetchErr != nil {
		log.Default().Printf("Error fetching user: %v\n", fetchErr)
		return 0, fetchErr
	}

	validateErr := bcrypt.CompareHashAndPassword(user.Hash, []byte(password))
	if validateErr != nil {
		log.Default().Printf("Error validating creds: %v\n", validateErr)
		return 0, validateErr
	}
	return user.Id, nil
}

func (s UserService) Register(ctx context.Context, username, password string) (userId int, err error) {
	hashedPassword, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if hashErr != nil {
		return 0, hashErr
	}

	user := repo.User{
		Username: username,
		Hash:     hashedPassword,
	}
	return s.repo.Insert(ctx, user)
}
