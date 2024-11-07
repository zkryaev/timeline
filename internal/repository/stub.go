package repository

import (
	"context"
	"timeline/internal/model"
)

type StubRepo struct{}

func (s StubRepo) SaveUser(ctx context.Context, user *model.User) (uint64, error) {
	return 0, nil
}
func (s StubRepo) User(ctx context.Context) (*model.User, error) {
	return nil, nil
}
func (s StubRepo) SaveOrg(ctx context.Context, org *model.Organization) (uint64, error) {
	return 0, nil
}
func (s StubRepo) Organization(ctx context.Context) (*model.Organization, error) {
	return nil, nil
}
