package service

import (
	"cruzes.co/junedc/entain/sport/db"
	"cruzes.co/junedc/entain/sport/proto/sport"

	"golang.org/x/net/context"
)

type Sport interface {
	// ListEvents will return a collection of races.
	ListEvents(ctx context.Context, in *sport.ListEventsRequest) (*sport.ListEventsResponse, error)

	GetEvent(ctx context.Context, in *sport.GetEventRequest) (*sport.Event, error)
}

// sportService implements the Sport interface.
type sportService struct {
	eventsRepo db.EventsRepo
}

// NewSportService instantiates and returns a new SportService.
func NewSportService(eventsRepo db.EventsRepo) Sport {
	return &sportService{eventsRepo}
}

func (s *sportService) ListEvents(ctx context.Context, in *sport.ListEventsRequest) (*sport.ListEventsResponse, error) {
	events, err := s.eventsRepo.List(in)
	if err != nil {
		return nil, err
	}

	return &sport.ListEventsResponse{Event: events}, nil
}

func (s *sportService) GetEvent(ctx context.Context, in *sport.GetEventRequest) (*sport.Event, error) {
	race, err := s.eventsRepo.Get(in)
	if err != nil {
		return nil, err
	}

	return race, nil
}
