package service

import (
	"Go_shortenURL/internal/repository"
	"Go_shortenURL/pkg/shortener"
	"context"
	"errors"
	"net/http"
)

type URLService struct {
	URLRepository *repository.URLRepository
}

func NewURLService(urlRepository *repository.URLRepository) *URLService {
	return &URLService{URLRepository: urlRepository}
}
func (s *URLService) ValidateURL(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("URL is not reachable")
	}
	return nil
}

func (s *URLService) ShortenURL(ctx context.Context, url string) (string, error) {
	shortCode := shortener.Encode()
	err := s.URLRepository.SetURL(ctx, shortCode, url)
	if err != nil {
		return "", err
	}
	return shortCode, nil
}

func (s *URLService) GetURL(ctx context.Context, shortCode string) (string, error) {
	url, err := s.URLRepository.GetURL(ctx, shortCode)
	if err != nil {
		return "", err
	}
	return url, nil
}

