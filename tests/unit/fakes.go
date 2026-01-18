package unit

import (
	"context"
	"errors"
	"time"

	"github.com/ruslanuskembaev/geo-alerts-system/internal/domain"
)

type fakeIncidentRepo struct {
	createFn        func(context.Context, domain.CreateIncidentRequest) (*domain.Incident, error)
	getByIDFn       func(context.Context, string) (*domain.Incident, error)
	listFn          func(context.Context, int, int) ([]*domain.Incident, int, error)
	updateFn        func(context.Context, string, domain.UpdateIncidentRequest) (*domain.Incident, error)
	deactivateFn    func(context.Context, string) error
	listActiveFn    func(context.Context) ([]*domain.Incident, error)
	createCalls     int
	getByIDCalls    int
	listCalls       int
	updateCalls     int
	deactivateCalls int
	listActiveCalls int
}

func (f *fakeIncidentRepo) Create(ctx context.Context, req domain.CreateIncidentRequest) (*domain.Incident, error) {
	f.createCalls++
	if f.createFn != nil {
		return f.createFn(ctx, req)
	}
	return nil, errors.New("Create not implemented")
}

func (f *fakeIncidentRepo) GetByID(ctx context.Context, id string) (*domain.Incident, error) {
	f.getByIDCalls++
	if f.getByIDFn != nil {
		return f.getByIDFn(ctx, id)
	}
	return nil, errors.New("GetByID not implemented")
}

func (f *fakeIncidentRepo) List(ctx context.Context, limit, offset int) ([]*domain.Incident, int, error) {
	f.listCalls++
	if f.listFn != nil {
		return f.listFn(ctx, limit, offset)
	}
	return nil, 0, errors.New("List not implemented")
}

func (f *fakeIncidentRepo) Update(ctx context.Context, id string, req domain.UpdateIncidentRequest) (*domain.Incident, error) {
	f.updateCalls++
	if f.updateFn != nil {
		return f.updateFn(ctx, id, req)
	}
	return nil, errors.New("Update not implemented")
}

func (f *fakeIncidentRepo) Deactivate(ctx context.Context, id string) error {
	f.deactivateCalls++
	if f.deactivateFn != nil {
		return f.deactivateFn(ctx, id)
	}
	return errors.New("Deactivate not implemented")
}

func (f *fakeIncidentRepo) ListActive(ctx context.Context) ([]*domain.Incident, error) {
	f.listActiveCalls++
	if f.listActiveFn != nil {
		return f.listActiveFn(ctx)
	}
	return nil, errors.New("ListActive not implemented")
}

type fakeIncidentCache struct {
	getFn           func(context.Context) ([]*domain.Incident, bool, error)
	setFn           func(context.Context, []*domain.Incident) error
	invalidateFn    func(context.Context) error
	getCalls        int
	setCalls        int
	invalidateCalls int
}

func (f *fakeIncidentCache) GetActive(ctx context.Context) ([]*domain.Incident, bool, error) {
	f.getCalls++
	if f.getFn != nil {
		return f.getFn(ctx)
	}
	return nil, false, nil
}

func (f *fakeIncidentCache) SetActive(ctx context.Context, incidents []*domain.Incident) error {
	f.setCalls++
	if f.setFn != nil {
		return f.setFn(ctx, incidents)
	}
	return nil
}

func (f *fakeIncidentCache) Invalidate(ctx context.Context) error {
	f.invalidateCalls++
	if f.invalidateFn != nil {
		return f.invalidateFn(ctx)
	}
	return nil
}

type fakeCheckRepo struct {
	createFn        func(context.Context, domain.LocationCheck, []string) error
	statsFn         func(context.Context, time.Time) ([]domain.IncidentStats, error)
	createCalls     int
	statsCalls      int
	lastCheck       domain.LocationCheck
	lastIncidentIDs []string
}

func (f *fakeCheckRepo) Create(ctx context.Context, check domain.LocationCheck, incidentIDs []string) error {
	f.createCalls++
	f.lastCheck = check
	f.lastIncidentIDs = append([]string(nil), incidentIDs...)
	if f.createFn != nil {
		return f.createFn(ctx, check, incidentIDs)
	}
	return nil
}

func (f *fakeCheckRepo) StatsByIncident(ctx context.Context, since time.Time) ([]domain.IncidentStats, error) {
	f.statsCalls++
	if f.statsFn != nil {
		return f.statsFn(ctx, since)
	}
	return nil, nil
}

type fakeQueue struct {
	enqueueFn func(context.Context, domain.WebhookJob) error
	dequeueFn func(context.Context, time.Duration) (*domain.WebhookJob, bool, error)
	enqueued  []domain.WebhookJob
}

func (f *fakeQueue) Enqueue(ctx context.Context, job domain.WebhookJob) error {
	f.enqueued = append(f.enqueued, job)
	if f.enqueueFn != nil {
		return f.enqueueFn(ctx, job)
	}
	return nil
}

func (f *fakeQueue) Dequeue(ctx context.Context, timeout time.Duration) (*domain.WebhookJob, bool, error) {
	if f.dequeueFn != nil {
		return f.dequeueFn(ctx, timeout)
	}
	return nil, false, nil
}
