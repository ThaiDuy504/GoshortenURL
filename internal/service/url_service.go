package service

import (
	"Go_shortenURL/internal/repository"
	"Go_shortenURL/pkg/shortener"
	"context"
	"errors"
	"net/http"
	"strings"
)

type URLService struct {
	URLRepository *repository.URLRepository
}

func NewURLService(urlRepository *repository.URLRepository) *URLService {
	return &URLService{URLRepository: urlRepository}
}

func (s *URLService) validateURL(url string) error {
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
	err := s.validateURL(url)
	if err != nil {
		return "", err
	}

	// Kiểm tra xem URL này đã được rút gọn lần nào chưa
	existingShortCode, err := s.URLRepository.GetShortCodeByOriginalURL(ctx, url)
	if err == nil {
		// URL đã tồn tại, trả về shortCode cũ
		return existingShortCode, nil
	}

	// URL chưa tồn tại, tạo mới với retry logic
	// Retry tối đa 5 lần nếu shortCode bị trùng
	for attempt := 0; attempt < 5; attempt++ {
		shortCode := shortener.Encode()
		err = s.URLRepository.SetURL(ctx, shortCode, url)

		// Nếu SetURL thành công, return shortCode
		if err == nil {
			return shortCode, nil
		}

		// Nếu có lỗi khác (không phải trùng lặp), return lỗi ngay
		if !strings.Contains(err.Error(), "unique") && !strings.Contains(err.Error(), "Duplicate") {
			return "", err
		}
		// Nếu lỗi trùng lặp, loop lại để retry
	}

	return "", errors.New("failed to generate unique short code after 5 attempts")
}

func (s *URLService) GetURL(ctx context.Context, shortCode string) (string, error) {
	url, err := s.URLRepository.GetURL(ctx, shortCode)
	if err != nil {
		return "", err
	}
	return url, nil
}
