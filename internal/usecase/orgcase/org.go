package orgcase

import (
	"context"
	"errors"
	"fmt"
	"timeline/internal/entity/dto/orgdto"
	"timeline/internal/infrastructure"
	"timeline/internal/infrastructure/database/postgres"
	"timeline/internal/infrastructure/mapper/orgmap"
	"timeline/internal/usecase/common"

	"go.uber.org/zap"
)

type OrgUseCase struct {
	user infrastructure.UserRepository
	org  infrastructure.OrgRepository
}

func New(userRepo infrastructure.UserRepository, orgRepo infrastructure.OrgRepository) *OrgUseCase {
	return &OrgUseCase{
		user: userRepo,
		org:  orgRepo,
	}
}

func (o *OrgUseCase) Organization(ctx context.Context, logger *zap.Logger, id int) (*orgdto.Organization, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id must be > 0")
	}
	data, err := o.org.OrgByID(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrOrgNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched organization")
	return orgmap.OrganizationToDTO(data), nil
}

func (o *OrgUseCase) OrgUpdate(ctx context.Context, logger *zap.Logger, newOrg *orgdto.OrgUpdateReq) error {
	if err := o.org.OrgUpdate(ctx, orgmap.OrgUpdateToModel(newOrg)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("Organization has been updated")
	return nil
}
