package infrastructure

import (
	"context"
	"fmt"
	"time"
	"timeline/internal/config"
	"timeline/internal/infrastructure/database/postgres"
	"timeline/internal/infrastructure/models"
	"timeline/internal/infrastructure/models/orgmodel"
	"timeline/internal/infrastructure/models/recordmodel"
	"timeline/internal/infrastructure/models/usermodel"
	"timeline/internal/utils/loader/objects"

	"go.uber.org/zap"
)

type Database interface {
	Open() error
	Close()
	infrastructure
}

type infrastructure interface {
	CodeRepository
	UserRepository
	OrgRepository
	RecordRepository
	BackgroundDataStore
}

type CodeRepository interface {
	SaveVerifyCode(ctx context.Context, info *models.CodeInfo) error
	VerifyCode(ctx context.Context, info *models.CodeInfo) (time.Time, error)
	ActivateAccount(ctx context.Context, id int, isOrg bool) error
	AccountExpiration(ctx context.Context, email string, isOrg bool) (*models.ExpInfo, error)
	DeleteExpiredCodes(ctx context.Context) error
}

type UserRepository interface {
	UserUpdate(ctx context.Context, upd *usermodel.UserInfo) error
	UserSave(ctx context.Context, user *usermodel.UserRegister) (int, error)
	UserByID(ctx context.Context, userID int) (*usermodel.UserInfo, error)
	UserDelete(ctx context.Context, userID int) error
	UserSoftDelete(ctx context.Context, userID int) error
	UserDeleteExpired(ctx context.Context) error
	UserUUID(ctx context.Context, userID int) (string, error)
	UserSetUUID(ctx context.Context, userID int, NewUUID string) error
	UserDeleteURL(ctx context.Context, userID int, URL string) error
}

type OrgRepository interface {
	OrgSave(ctx context.Context, org *orgmodel.OrgRegister) (int, error)
	OrgUpdate(ctx context.Context, upd *orgmodel.Organization) error
	OrgByID(ctx context.Context, id int) (*orgmodel.Organization, error)
	OrgsBySearch(ctx context.Context, params *orgmodel.SearchParams) (*orgmodel.RespData, error)
	OrgsInArea(ctx context.Context, area *orgmodel.AreaParams) ([]*orgmodel.OrgByArea, error)
	OrgDelete(ctx context.Context, orgID int) error
	OrgSoftDelete(ctx context.Context, orgID int) error
	OrgDeleteExpired(ctx context.Context) error
	OrgSaveShowcaseImageURL(ctx context.Context, meta *models.ImageMeta) error
	OrgUUID(ctx context.Context, orgID int) (string, error)
	OrgSetUUID(ctx context.Context, orgID int, NewUUID string) error
	OrgDeleteURL(ctx context.Context, meta *models.ImageMeta) error
	TimetableRepository
	WorkerRepository
	ServiceRepository
	SlotRepository
	ScheduleRepository
}

type RecordRepository interface {
	Record(ctx context.Context, param recordmodel.RecordParam) (*recordmodel.RecordScrap, error)
	RecordList(ctx context.Context, req *recordmodel.RecordListParams) ([]*recordmodel.RecordScrap, int, error)
	RecordAdd(ctx context.Context, req *recordmodel.Record) (*recordmodel.ReminderRecord, int, error)
	RecordDelete(ctx context.Context, recordID int) error
	RecordCancel(ctx context.Context, req *recordmodel.RecordCancelation) error
	UpcomingRecords(ctx context.Context) ([]*recordmodel.ReminderRecord, error)
	FeedbackRepository
}

type TimetableRepository interface {
	Timetable(ctx context.Context, orgID int, tdata models.TokenData) ([]*orgmodel.OpenHours, string, error)
	TimetableAdd(ctx context.Context, orgID int, newTime []*orgmodel.OpenHours) error
	TimetableUpdate(ctx context.Context, orgID int, upd []*orgmodel.OpenHours) error
	TimetableDelete(ctx context.Context, orgID, weekday int) error
}

type WorkerRepository interface {
	Worker(ctx context.Context, WorkerID, OrgID int) (*orgmodel.Worker, error)
	WorkerAdd(ctx context.Context, worker *orgmodel.Worker) (int, error)
	WorkerUpdate(ctx context.Context, worker *orgmodel.Worker) error
	WorkerAssignService(ctx context.Context, assignInfo *orgmodel.WorkerAssign) error
	WorkerUnAssignService(ctx context.Context, assignInfo *orgmodel.WorkerAssign) error
	WorkerList(ctx context.Context, OrgID, Limit, Offset int) ([]*orgmodel.Worker, int, error)
	WorkerDelete(ctx context.Context, WorkerID, OrgID int) error
	WorkerSoftDelete(ctx context.Context, WorkerID, OrgID int) error
	WorkerUUID(ctx context.Context, orgID, workerID int) (string, error)
	WorkerSetUUID(ctx context.Context, orgID, workerID int, NewUUID string) error
	WorkerDeleteURL(ctx context.Context, orgID int, URL string) error
}

type ServiceRepository interface {
	Service(ctx context.Context, ServiceID, OrgID int) (*orgmodel.Service, error)
	ServiceWorkerList(ctx context.Context, ServiceID, OrgID int) ([]*orgmodel.Worker, error)
	ServiceAdd(ctx context.Context, service *orgmodel.Service) (int, error)
	ServiceUpdate(ctx context.Context, service *orgmodel.Service) error
	ServiceList(ctx context.Context, OrgID int, Limit, Offset int) ([]*orgmodel.Service, int, error)
	ServiceDelete(ctx context.Context, ServiceID, OrgID int) error
	ServiceSoftDelete(ctx context.Context, ServiceID, OrgID int) error
}

type SlotRepository interface {
	GenerateSlots(ctx context.Context) error
	DeleteExpiredSlots(ctx context.Context) error
	Slots(ctx context.Context, params *orgmodel.SlotsReq) ([]*orgmodel.Slot, string, error)
}

type ScheduleRepository interface {
	WorkerSchedule(ctx context.Context, params *orgmodel.ScheduleParams) (*orgmodel.ScheduleList, error)
	AddWorkerSchedule(ctx context.Context, Schedule *orgmodel.WorkerSchedule) error
	UpdateWorkerSchedule(ctx context.Context, Schedule *orgmodel.WorkerSchedule) error
	DeleteWorkerSchedule(ctx context.Context, metainfo *orgmodel.ScheduleParams) error
	SoftDeleteWorkerSchedule(ctx context.Context, metainfo *orgmodel.ScheduleParams) error
}

type FeedbackRepository interface {
	FeedbackList(ctx context.Context, params *recordmodel.FeedbackParams) ([]*recordmodel.Feedback, int, error)
	FeedbackSet(ctx context.Context, feedback *recordmodel.Feedback) error
	FeedbackUpdate(ctx context.Context, feedback *recordmodel.Feedback) error
	FeedbackDelete(ctx context.Context, params *recordmodel.FeedbackParams) error
}

type BackgroundDataStore interface {
	SaveCities(ctx context.Context, logger *zap.Logger, cities objects.Cities) error
	PreLoadCities(ctx context.Context) (objects.Cities, error)
	IsCitiesLoad(ctx context.Context) (bool, error)
}

// Паттерн фабричный метод, чтобы не завязываться на конкретной БД
func GetDB(name string, cfg *config.Database) (Database, error) {
	switch name {
	case "postgres":
		return postgres.New(cfg, "", true), nil
	default:
		return nil, fmt.Errorf("unexpected db name")
	}
}
