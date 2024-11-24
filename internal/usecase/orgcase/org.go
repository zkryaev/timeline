package orgcase

import (
	"context"
	"errors"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/repository"
	"timeline/internal/repository/database/postgres"
	"timeline/internal/repository/mapper/orgmap"

	"go.uber.org/zap"
)

type OrgUseCase struct {
	user   repository.UserRepository
	org    repository.OrgRepository
	Logger *zap.Logger
}

func New(userRepo repository.UserRepository, orgRepo repository.OrgRepository, logger *zap.Logger) *OrgUseCase {
	return &OrgUseCase{
		user:   userRepo,
		org:    orgRepo,
		Logger: logger,
	}
}

func (o *OrgUseCase) OrgUpdate(ctx context.Context, org *orgdto.OrgUpdateReq) (*orgdto.OrgUpdateResp, error) {
	data, err := o.org.OrgUpdate(ctx, orgmap.UpdateToModel(org))
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return nil, err
		}
		o.Logger.Error(
			"failed org update",
			zap.Error(err),
		)
	}
	return orgmap.UpdateToDTO(data), nil
}
