package eventservice

import (
	"context"
	"sea-api/internal/errs"
	"sea-api/internal/models"
	"sea-api/internal/models/eventmodels"
	"sea-api/internal/utils/valid"
)

func (s *EventService) GetEventViewList(ctx context.Context, req *eventmodels.EventListRequest) (*eventmodels.EventViewListResponse, error) {
	count, err := s.repo.CountEvents()
	if err != nil {
		return nil, err
	}

	totalPages := valid.Limit(&req.ListRequest, count)

	events, err := s.repo.GetList(req)
	if err != nil {
		return nil, err
	}

	var list = []eventmodels.EventListItemResponse{}
	for _, e := range *events {
		url, err := s.s3.GenerateDownloadUrlByKey(ctx, e.BackgroundKey)
		if err != nil {
			return nil, err
		}

		list = append(list, eventmodels.EventListItemResponse{
			ID:            e.ID,
			Name:          e.Name,
			Description:   e.Description,
			BackgroundURL: url,
			Belonging:     e.Belonging,
			StartDate:     e.StartDate,
			EndDate:       e.EndDate,
		})
	}

	return &eventmodels.EventViewListResponse{
		List: list,
		ListResponse: models.ListResponse{
			TotalPages:  totalPages,
			CurrentPage: req.Page,
			Count:       count,
		},
	}, nil
}

func (s *EventService) GetEventView(ctx context.Context, id int64) (*eventmodels.EventViewResponse, error) {
	event, err := s.repo.GetView(id)
	if err != nil {
		return nil, errs.New(errs.NotFound, "event not found", nil)
	}

	url, err := s.s3.GenerateDownloadUrlByKey(ctx, event.BackgroundKey)
	if err != nil {
		return nil, err
	}

	return &eventmodels.EventViewResponse{
		ID:              event.ID,
		Name:            event.Name,
		Description:     event.Description,
		BackgroundURL:   url,
		Belonging:       event.Belonging,
		RequireApplying: event.RequireApplying,
		StartDate:       event.StartDate,
		EndDate:         event.EndDate,
	}, nil
}
