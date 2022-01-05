package banners

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"sync"
)

// this is for managing banners
type Service struct {
	mu    sync.RWMutex
	items []*Banner
}

//Banner
type Banner struct {
	ID      int64
	Title   string
	Content string
	Button  string
	Link    string
	Image   string
}

//create service
func NewService() *Service {
	return &Service{items: make([]*Banner, 0)}
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

func (s *Service) Save(ctx context.Context, item *Banner, file multipart.File) (*Banner, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if item.ID == 0 {
		newID++
		item.ID = newID
		if item.Image != "" {
			item.Image = fmt.Sprint(item.ID) + "." + item.Image

			data, err := ioutil.ReadAll(file)
			if err != nil {
				return nil, errors.New("not readible data")
			}

			err = ioutil.WriteFile("./web/banners/"+item.Image, data, 0666)
			if err != nil {
				return nil, err
			}
		}

		s.items = append(s.items, item)
		return item, nil
	}

	for k, v := range s.items {
		if v.ID == item.ID {
			if item.Image != "" {
				item.Image = fmt.Sprint(item.ID) + "." + item.Image

				data, err := ioutil.ReadAll(file)
				if err != nil {
					return nil, errors.New("not readible data")
				}

				err = ioutil.WriteFile("./web/banners/"+item.Image, data, 0666)
				if err != nil {
					return nil, err
				}
			} else {
				item.Image = s.items[k].Image
			}

			s.items[k] = item
			return item, nil
		}
	}

	return nil, errors.New("item not found")
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
