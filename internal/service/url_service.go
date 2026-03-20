package service

import (
	"Go_shortenURL/internal/repository"
	"Go_shortenURL/pkg/shortener"
)

type URLService struct {
	URLRepository *repository.URLRepository
}

func NewURLService(urlRepository *repository.URLRepository) *URLService {
	return &URLService{URLRepository: urlRepository}
}

func (s *URLService) ShortenURL(url string) (string, error) {
	shortCode := shortener.Encode(url)
	err := s.URLRepository.SetURL(shortCode, url)
	if err != nil {
		return "", err
	}
	return shortCode, nil
}

func (s *URLService) GetURL(shortCode string) (string, error) {
	url, err := s.URLRepository.GetURL(shortCode)
	if err != nil {
		return "", err
	}
	return url, nil
}

