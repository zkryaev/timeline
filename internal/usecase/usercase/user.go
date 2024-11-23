package usercase

import (
	"context"
	"errors"
	"strings"
	"timeline/internal/entity"
	"timeline/internal/entity/dto"
	"timeline/internal/repository"
	"timeline/internal/repository/database/postgres"
	"timeline/internal/repository/mapper/facemap"
	"timeline/internal/repository/mapper/orgmap"

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
		user: userRepo,
		org:  orgRepo,
	}
}

func (u *UserUseCase) SearchOrgs(ctx context.Context, sreq *dto.SearchReq) (*dto.SearchResp, error) {
	sreq.Name = strings.TrimSpace(sreq.Name)
	data, err := u.org.OrgsSearch(ctx, facemap.SearchToModel(sreq))
	if err != nil {
		if errors.Is(err, postgres.ErrOrgsNotFound) {
			return nil, ErrNoOrgs
		}
		u.Logger.Error("SearchOrgs", zap.Error(err))
		return nil, err
	}
	resp := &dto.SearchResp{
		Orgs: make([]*entity.Organization, 0, len(data)),
	}
	for _, v := range data {
		resp.Orgs = append(resp.Orgs, orgmap.OrgInfoToDTO(v))
	}
	return resp, nil
}

func (u *UserUseCase) OrgsInArea(ctx context.Context, area *dto.OrgAreaReq) (*dto.OrgAreaResp, error) {
	data, err := u.org.OrgsInArea(ctx, facemap.AreaToModel(area))
	if err != nil {
		if errors.Is(err, postgres.ErrOrgsNotFound) {
			return nil, ErrNoOrgs
		}
		u.Logger.Error("OrgsInArea", zap.Error(err))
		return nil, err
	}
	resp := &dto.OrgAreaResp{
		Orgs: make([]*entity.MapOrgInfo, 0, len(data)),
	}
	for _, v := range data {
		resp.Orgs = append(resp.Orgs, orgmap.OrgSummaryToDTO(v))
	}
	return resp, nil
}
