package cronjob

import (
	"context"
	"timeline/internal/repository"

	gocron "github.com/go-co-op/gocron/v2"
)

// Database:
//   - slots: генерирует и удаляет стухшие
//   - user_codes, org_codes: удаляет стухшие
//   - users, orgs: удаляет стухшие
func InitCronScheduler(db repository.Repository) gocron.Scheduler {
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err.Error())
	}
	s.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(gocron.NewAtTime(16, 00, 00)),
		),
		gocron.NewTask(
			func(slots repository.SlotRepository) {
				ctx := context.Background()
				slots.DeleteExpiredSlots(ctx)
				slots.GenerateSlots(ctx)
			},
			db,
		),
		gocron.WithName("Database > Slots > Delete expired and Generate"),
	)
	s.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(gocron.NewAtTime(00, 00, 00)),
		),
		gocron.NewTask(
			func(codes repository.CodeRepository) {
				ctx := context.Background()
				codes.DeleteExpiredCodes(ctx)
			},
			db,
		),
		gocron.WithName("Database > Codes > Delete expired"),
	)
	s.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(gocron.NewAtTime(00, 00, 00)),
		),
		gocron.NewTask(
			func(users repository.UserRepository, orgs repository.OrgRepository) {
				ctx := context.Background()
				users.UserDeleteExpired(ctx)
				orgs.OrgDeleteExpired(ctx)

			},
			db, db,
		),
		gocron.WithName("Database > Accounts > Delete expired"),
	)
	return s
}
