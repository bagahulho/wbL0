package service

import "wbL0/pkg/repository"

type Service struct {
	Repository *repository.Repository
}

func NewService(repos *repository.Repository) *Service {
	return &Service{Repository: repos}
}
