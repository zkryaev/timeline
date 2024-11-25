package usercase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/entity/dto/userdto"
	"timeline/internal/repository"
	"timeline/internal/repository/database/postgres"
	"timeline/internal/repository/mapper/facemap"
	"timeline/internal/repository/mapper/orgmap"
	"timeline/internal/repository/mapper/usermap"

	"go.uber.org/zap"
)

var (
	ErrNoOrgs = errors.New("There are no organizations")
)

type UserUseCase struct {
	user   repository.UserRepository
	org    repository.OrgRepository
	Logger *zap.Logger
}

func New(userRepo repository.UserRepository, orgRepo repository.OrgRepository, logger *zap.Logger) *UserUseCase {
	return &UserUseCase{
		user:   userRepo,
		org:    orgRepo,
		Logger: logger,
	}
}

func (u *UserUseCase) User(ctx context.Context, id int) (*userdto.UserGetResp, error) {
	if id <= 0 {
		return nil, fmt.Errorf("Id must be > 0")
	}
	data, err := u.user.UserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := usermap.UserInfoToGetResp(data)
	return resp, nil
}

func (u *UserUseCase) UserUpdate(ctx context.Context, user *userdto.UserUpdateReq) (*userdto.UserUpdateResp, error) {
	data, err := u.user.UserUpdate(ctx, usermap.UserUpdateToModel(user))
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return nil, err
		}
		u.Logger.Error(
			"failed user update",
			zap.Error(err),
		)
	}
	return usermap.UserUpdateToDTO(data), nil
}

func (u *UserUseCase) SearchOrgs(ctx context.Context, sreq *orgdto.SearchReq) (*orgdto.SearchResp, error) {
	sreq.Name = strings.TrimSpace(sreq.Name)
	data, err := u.org.OrgsBySearch(ctx, facemap.SearchToModel(sreq))
	if err != nil {
		if errors.Is(err, postgres.ErrOrgsNotFound) {
			return nil, ErrNoOrgs
		}
		u.Logger.Error("SearchOrgs", zap.Error(err))
		return nil, err
	}
	resp := &orgdto.SearchResp{
		Orgs: make([]*entity.Organization, 0, len(data)),
	}
	for _, v := range data {
		resp.Orgs = append(resp.Orgs, orgmap.OrgInfoToDTO(v))
	}
	return resp, nil
}

func (u *UserUseCase) OrgsInArea(ctx context.Context, area *orgdto.OrgAreaReq) (*orgdto.OrgAreaResp, error) {
	data, err := u.org.OrgsInArea(ctx, facemap.AreaToModel(area))
	if err != nil {
		if errors.Is(err, postgres.ErrOrgsNotFound) {
			return nil, ErrNoOrgs
		}
		u.Logger.Error("OrgsInArea", zap.Error(err))
		return nil, err
	}
	resp := &orgdto.OrgAreaResp{
		Orgs: make([]*entity.MapOrgInfo, 0, len(data)),
	}
	for _, v := range data {
		resp.Orgs = append(resp.Orgs, orgmap.OrgSummaryToDTO(v))
	}
	return resp, nil
}
