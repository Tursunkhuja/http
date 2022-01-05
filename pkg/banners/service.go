package banners

import (
	"context"
	"errors"
	"log"
	"sync"
)

//service to control banner
type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

//create service
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
}

//Banner
type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
}

func (s *Service) All(ctx context.Context) ([]*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.items != nil {
		return s.items, nil
	}

	return nil, errors.New("no banners")
}

func (s *Service) ByID(ctx context.Context, id int64) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, banner := range s.items {
		if banner.ID == id {
			return banner, nil
		}
	}
	return nil, errors.New("item not found")
}

var newID int64

func (s *Service) Save(ctx context.Context, item *Banner) (*Banner, error) {

	if item.ID == 0 {
		newID++
		newBnner := &Banner{
			ID:      newID,
			Title:   item.Title,
			Content: item.Content,
			Button:  item.Button,
			Link:    item.Link,
		}
		s.items = append(s.items, newBnner)
		return newBnner, nil
	}

	ExistB, err := s.ByID(ctx, item.ID)

	if err != nil {
		log.Print(err)
		return nil, errors.New("item not found")
	}

	ExistB.Button = item.Button
	ExistB.Title = item.Title
	ExistB.Content = item.Content
	ExistB.Link = item.Link

	return ExistB, nil
}

func (s *Service) RemoveByID(ctx context.Context, id int64) (*Banner, error) {
	sBanner, err := s.ByID(ctx, id)
	if err != nil {
		log.Print(err)
		return nil, errors.New("item not found")
	}
	for i, banner := range s.items {
		if banner.ID == id {
			s.items = append(s.items[:i], s.items[i+1:]...)
			break
		}
	}

	return sBanner, nil
}
