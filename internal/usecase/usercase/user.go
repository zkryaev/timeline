package usercase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/general"
	"timeline/internal/entity/dto/userdto"
	"timeline/internal/repository"
	"timeline/internal/repository/database/postgres"
	"timeline/internal/repository/mail"
	"timeline/internal/repository/mail/notify"
	"timeline/internal/repository/mapper/orgmap"
	"timeline/internal/repository/mapper/recordmap"
	"timeline/internal/repository/mapper/usermap"

	"go.uber.org/zap"
)

var (
	ErrNoOrgs = errors.New("There are no organizations")
)

type UserUseCase struct {
	user    repository.UserRepository
	org     repository.OrgRepository
	records repository.RecordRepository
	mail    notify.Mail
	Logger  *zap.Logger
}

func New(userRepo repository.UserRepository, orgRepo repository.OrgRepository, recRepo repository.RecordRepository, logger *zap.Logger) *UserUseCase {
	return &UserUseCase{
		user:   userRepo,
		org:    orgRepo,
		Logger: logger,
	}
}

func (u *UserUseCase) User(ctx context.Context, id int) (*entity.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("Id must be > 0")
	}
	data, err := u.user.UserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	resp := usermap.UserInfoToDTO(data)
	return resp, nil
}

func (u *UserUseCase) UserUpdate(ctx context.Context, newUser *userdto.UserUpdateReq) error {
	if err := u.user.UserUpdate(ctx, usermap.UserUpdateToModel(newUser)); err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			return err
		}
		u.Logger.Error(
			"failed user update",
			zap.Error(err),
		)
		return err
	}
	return nil
}

func (u *UserUseCase) SearchOrgs(ctx context.Context, sreq *general.SearchReq) (*general.SearchResp, error) {
	sreq.Name = strings.TrimSpace(sreq.Name)
	data, found, err := u.org.OrgsBySearch(ctx, orgmap.SearchToModel(sreq))
	if err != nil {
		if errors.Is(err, postgres.ErrOrgsNotFound) {
			return nil, ErrNoOrgs
		}
		u.Logger.Error("SearchOrgs", zap.Error(err))
		return nil, err
	}
	resp := &general.SearchResp{
		Found: found,
		Orgs:  make([]*entity.OrgsBySearch, 0, len(data)),
	}
	for _, v := range data {
		resp.Orgs = append(resp.Orgs, orgmap.OrgsBySearchToDTO(v))
	}
	return resp, nil
}

func (u *UserUseCase) OrgsInArea(ctx context.Context, area *general.OrgAreaReq) (*general.OrgAreaResp, error) {
	data, err := u.org.OrgsInArea(ctx, orgmap.AreaToModel(area))
	if err != nil {
		if errors.Is(err, postgres.ErrOrgsNotFound) {
			return nil, ErrNoOrgs
		}
		u.Logger.Error("OrgsInArea", zap.Error(err))
		return nil, err
	}
	resp := &general.OrgAreaResp{
		Found: len(data),
		Orgs:  make([]*entity.MapOrgInfo, 0, len(data)),
	}
	for _, v := range data {
		resp.Orgs = append(resp.Orgs, orgmap.OrgSummaryToDTO(v))
	}
	return resp, nil
}

func (u *UserUseCase) UserRecordReminder(ctx context.Context) error {
	data, err := u.records.UpcomingRecords(ctx)
	if err != nil {
		u.Logger.Error(
			"failed user update",
			zap.Error(err),
		)
		return err
	}
	for i := range data {
		msg := &mail.Message{
			Email: data[i].UserEmail,
			Type:  notify.ReminderType,
			Value: recordmap.RecordToReminder(data[i]),
		}
		if err := u.mail.SendMsg(msg); err != nil {
			u.Logger.Error(
				"failed to notify users",
				zap.Error(err),
			)
			return err
		}
	}
	return nil
}
