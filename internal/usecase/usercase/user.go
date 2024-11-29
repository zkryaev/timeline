package usercase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/general"
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

func (u *UserUseCase) UserUpdate(ctx context.Context, newUser *userdto.UserUpdateReq) (*userdto.UserUpdateReq, error) {
	if err := u.user.UserUpdate(ctx, usermap.UserUpdateToModel(newUser)); err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return nil, err
		}
		u.Logger.Error(
			"failed user update",
			zap.Error(err),
		)
	}
	return newUser, nil
}

func (u *UserUseCase) SearchOrgs(ctx context.Context, sreq *general.SearchReq) (*general.SearchResp, error) {
	sreq.Name = strings.TrimSpace(sreq.Name)
	data, found, err := u.org.OrgsBySearch(ctx, facemap.SearchToModel(sreq))
	if err != nil {
		if errors.Is(err, postgres.ErrOrgsNotFound) {
			return nil, ErrNoOrgs
		}
		u.Logger.Error("SearchOrgs", zap.Error(err))
		return nil, err
	}
	resp := &general.SearchResp{
		Pages: 1,
		Orgs:  make([]*orgdto.Organization, 0, len(data)),
	}
	if found > sreq.Limit {
		resp.Pages = (found + 1) / sreq.Limit
	}
	for _, v := range data {
		resp.Orgs = append(resp.Orgs, orgmap.OrgInfoToDTO(v))
	}
	return resp, nil
}

func (u *UserUseCase) OrgsInArea(ctx context.Context, area *general.OrgAreaReq) (*general.OrgAreaResp, error) {
	data, err := u.org.OrgsInArea(ctx, facemap.AreaToModel(area))
	if err != nil {
		if errors.Is(err, postgres.ErrOrgsNotFound) {
			return nil, ErrNoOrgs
		}
		u.Logger.Error("OrgsInArea", zap.Error(err))
		return nil, err
	}
	resp := &general.OrgAreaResp{
		Orgs: make([]*entity.MapOrgInfo, 0, len(data)),
	}
	for _, v := range data {
		resp.Orgs = append(resp.Orgs, orgmap.OrgSummaryToDTO(v))
	}
	return resp, nil
}
