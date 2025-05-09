package usercase

import (
	"context"
	"errors"
	"fmt"
	"time"
	"timeline/internal/controller/scope"
	"timeline/internal/entity"
	"timeline/internal/entity/dto/general"
	"timeline/internal/entity/dto/userdto"
	"timeline/internal/infrastructure"
	"timeline/internal/infrastructure/database/postgres"
	"timeline/internal/infrastructure/mail"
	"timeline/internal/infrastructure/mapper/orgmap"
	"timeline/internal/infrastructure/mapper/recordmap"
	"timeline/internal/infrastructure/mapper/usermap"
	"timeline/internal/infrastructure/models"
	"timeline/internal/usecase/common"
	"timeline/internal/utils/loader"

	"go.uber.org/zap"
)

var (
	ErrNoOrgs = errors.New("no organizations was found")
)

type UserUseCase struct {
	user     infrastructure.UserRepository
	org      infrastructure.OrgRepository
	records  infrastructure.RecordRepository
	mail     infrastructure.Mail
	backdata *loader.BackData
	settings *scope.Settings
}

func New(userRepo infrastructure.UserRepository, orgRepo infrastructure.OrgRepository, recRepo infrastructure.RecordRepository,
	backdata *loader.BackData, settings *scope.Settings) *UserUseCase {
	return &UserUseCase{
		user:     userRepo,
		org:      orgRepo,
		records:  recRepo,
		backdata: backdata,
	}
}

func (u *UserUseCase) User(ctx context.Context, logger *zap.Logger, id int) (*entity.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("id < 0")
	}
	data, err := u.user.UserByID(ctx, id)
	if err != nil {
		if errors.Is(err, postgres.ErrUserNotFound) {
			logger.Info("UserByID", zap.Error(err))
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched user")
	resp := usermap.UserInfoToDTO(data)
	return resp, nil
}

func (u *UserUseCase) UserUpdate(ctx context.Context, logger *zap.Logger, newUser *userdto.UserUpdateReq) error {
	if err := u.user.UserUpdate(ctx, usermap.UserUpdateToModel(newUser)); err != nil {
		if errors.Is(err, postgres.ErrNoRowsAffected) {
			return common.ErrNothingChanged
		}
		return err
	}
	logger.Info("User has been updated")
	return nil
}

func (u *UserUseCase) SearchOrgs(ctx context.Context, logger *zap.Logger, sreq *general.SearchReq) (*general.SearchResp, error) {
	data, err := u.org.OrgsBySearch(ctx, orgmap.SearchToModel(sreq))
	if err != nil {
		if errors.Is(err, postgres.ErrOrgsNotFound) {
			return nil, common.ErrNotFound
		}
		return nil, err
	}
	logger.Info("Fetched organizations by search", zap.Int("Found", data.Found))
	tzid := u.backdata.Cities.GetCityTZ(data.UserCity)
	loc, err := time.LoadLocation(tzid)
	if err != nil {
		logger.Error("failed to load location, set UTC+03 (MSK)", zap.String("city-tzid", data.UserCity+"="+tzid), zap.Error(err))
		loc = time.Local // UTC+03 = MSK
	}
	resp := &general.SearchResp{
		Found: data.Found,
		Orgs:  make([]*entity.OrgsBySearch, 0, len(data.Data)),
	}
	for _, v := range data.Data {
		resp.Orgs = append(resp.Orgs, orgmap.OrgsBySearchToDTO(v, loc))
	}
	return resp, nil
}

func (u *UserUseCase) OrgsInArea(ctx context.Context, logger *zap.Logger, area *general.OrgAreaReq) (*general.OrgAreaResp, error) {
	data, err := u.org.OrgsInArea(ctx, orgmap.AreaToModel(area))
	if err != nil {
		if errors.Is(err, postgres.ErrOrgsNotFound) {
			return nil, common.ErrNotFound
		}
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
	if !u.settings.EnableMail {
		return nil
	}
	data, err := u.records.UpcomingRecords(ctx)
	if err != nil {
		if errors.Is(err, postgres.ErrRecordsNotFound) {
			return common.ErrNotFound
		}
		return err
	}
	logger.Info("Fetched user's upcoming records")
	if len(data) > 0 {
		for i := range data {
			tzid := u.backdata.Cities.GetCityTZ(data[0].UserCity)
			loc, err := time.LoadLocation(tzid)
			if err != nil {
				logger.Error("failed to load location, set UTC+03 (MSK)", zap.String("city-tzid", data[0].UserCity+"="+tzid), zap.Error(err))
				loc = time.Local // UTC+03 = MSK
			}
			msg := &models.Message{
				Email:    data[i].UserEmail,
				Type:     mail.ReminderType,
				Value:    recordmap.ReminderRecordToReminder(data[i], loc),
				IsAttach: false,
			}
			if err = u.mail.SendMsg(msg); err != nil {
				return err
			}
		}
		logger.Info("Notifications has been sent to user's email")
	}
	logger.Info("There are no records to be reminded of")
	return nil
}
