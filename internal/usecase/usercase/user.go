package usercase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/general"
	"timeline/internal/entity/dto/userdto"
	"timeline/internal/infrastructure"
	"timeline/internal/infrastructure/mail"
	"timeline/internal/infrastructure/mapper/orgmap"
	"timeline/internal/infrastructure/mapper/recordmap"
	"timeline/internal/infrastructure/mapper/usermap"
	"timeline/internal/infrastructure/models"

	"go.uber.org/zap"
)

var (
	ErrNoOrgs = errors.New("no organizations was found")
)

type UserUseCase struct {
	user    infrastructure.UserRepository
	org     infrastructure.OrgRepository
	records infrastructure.RecordRepository
	mail    infrastructure.Mail
}

func New(userRepo infrastructure.UserRepository, orgRepo infrastructure.OrgRepository, recRepo infrastructure.RecordRepository) *UserUseCase {
	return &UserUseCase{
		user:    userRepo,
		org:     orgRepo,
		records: recRepo,
	}
}

func (u *UserUseCase) User(ctx context.Context, logger *zap.Logger, id int) (*entity.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id < 0")
	}
	data, err := u.user.UserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched user")
	resp := usermap.UserInfoToDTO(data)
	return resp, nil
}

func (u *UserUseCase) UserUpdate(ctx context.Context, logger *zap.Logger, newUser *userdto.UserUpdateReq) error {
	if err := u.user.UserUpdate(ctx, usermap.UserUpdateToModel(newUser)); err != nil {
		return err
	}
	logger.Info("User has been updated")
	return nil
}

func (u *UserUseCase) SearchOrgs(ctx context.Context, logger *zap.Logger, sreq *general.SearchReq) (*general.SearchResp, error) {
	sreq.Name = strings.TrimSpace(sreq.Name)
	data, found, err := u.org.OrgsBySearch(ctx, orgmap.SearchToModel(sreq))
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched organizations by search", zap.Int("Found", found))
	resp := &general.SearchResp{
		Found: found,
		Orgs:  make([]*entity.OrgsBySearch, 0, len(data)),
	}
	for _, v := range data {
		resp.Orgs = append(resp.Orgs, orgmap.OrgsBySearchToDTO(v))
	}
	return resp, nil
}

func (u *UserUseCase) OrgsInArea(ctx context.Context, logger *zap.Logger, area *general.OrgAreaReq) (*general.OrgAreaResp, error) {
	data, err := u.org.OrgsInArea(ctx, orgmap.AreaToModel(area))
	if err != nil {
		return nil, err
	}
	logger.Info("Fetched organizations in specified area", zap.Int("Found", len(data)))
	resp := &general.OrgAreaResp{
		Found: len(data),
		Orgs:  make([]*entity.MapOrgInfo, 0, len(data)),
	}
	for _, v := range data {
		resp.Orgs = append(resp.Orgs, orgmap.OrgSummaryToDTO(v))
	}
	return resp, nil
}

func (u *UserUseCase) UserRecordReminder(ctx context.Context, logger *zap.Logger) error {
	data, err := u.records.UpcomingRecords(ctx)
	if err != nil {
		return err
	}
	logger.Info("Fetched user's upcoming records")
	for i := range data {
		msg := &models.Message{
			Email:    data[i].UserEmail,
			Type:     mail.ReminderType,
			Value:    recordmap.RecordToReminder(data[i]),
			IsAttach: false,
		}
		if err = u.mail.SendMsg(msg); err != nil {
			return err
		}
		logger.Info("Notification has been sent to user's email")
	}
	return nil
}
