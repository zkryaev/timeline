package orgs

import (
	"context"
	"net/http"
	"strconv"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/orgdto"

	"github.com/gorilla/mux"
)

type Schedule interface {
	WorkerSchedule(ctx context.Context, params *orgdto.ScheduleParams) (*orgdto.ScheduleList, error)
	AddWorkerSchedule(ctx context.Context, schedule *orgdto.ScheduleList) error
	UpdateWorkerSchedule(ctx context.Context, schedule *orgdto.ScheduleList) error
	DeleteWorkerSchedule(ctx context.Context, params *orgdto.ScheduleParams) error
}

// @Summary Get worker schedule
// @Description Get specified worker schedule for specified org with weekday filter. If no weekday then all week will be returned
// @Tags organization/schedule
// @Produce json
// @Param   workerID path int true "worker_id"
// @Param   orgID path int true "org_id"
// @Param weekday query int false "weekday"
// @Success 200 {object} orgdto.ScheduleList
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/schedules/{workerID} [get]
func (o *OrgCtrl) WorkerSchedule(w http.ResponseWriter, r *http.Request) {
	params, err := validation.FetchSpecifiedID(mux.Vars(r), "orgID", "workerID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var weekday int
	if r.URL.Query().Get("weekday") != "" {
		weekday, err = strconv.Atoi(r.URL.Query().Get("weekday"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	req := &orgdto.ScheduleParams{
		WorkerID: params["workerID"],
		OrgID:    params["orgID"],
		Weekday:  weekday,
	}
	if err = o.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := o.usecase.WorkerSchedule(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if o.json.NewEncoder(w).Encode(data) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusInternalServerError)
		return
	}
}

// @Summary Delete worker schedule
// @Description Delete specified worker schedule for a specific organization and worker with an optional weekday filter
// @Tags organization/schedule
// @Param   workerID path int true "worker_id"
// @Param   orgID path int true "org_id"
// @Param   weekday query int false "weekday"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/schedules/{workerID} [delete]
func (o *OrgCtrl) DeleteWorkerSchedule(w http.ResponseWriter, r *http.Request) {
	params, err := validation.FetchSpecifiedID(mux.Vars(r), "orgID", "workerID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO: как то по умнее работать с query параметрами
	var weekday int
	if r.URL.Query().Get("weekday") != "" {
		weekday, err = strconv.Atoi(r.URL.Query().Get("weekday"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	req := &orgdto.ScheduleParams{
		WorkerID: params["workerID"],
		OrgID:    params["orgID"],
		Weekday:  weekday,
	}
	if err = o.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := o.usecase.DeleteWorkerSchedule(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Update worker schedule
// @Description Update the schedule for a specific worker in an organization
// @Tags organization/schedule
// @Accept json
// @Param   schedule body orgdto.ScheduleList true "Schedule data"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/schedules [put]
func (o *OrgCtrl) UpdateWorkerSchedule(w http.ResponseWriter, r *http.Request) {
	req := &orgdto.ScheduleList{}
	if o.json.NewDecoder(r.Body).Decode(req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	if err := o.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := o.usecase.UpdateWorkerSchedule(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Add worker schedule
// @Description Add a new schedule for a specific worker in an organization
// @Tags organization/schedule
// @Accept json
// @Produce json
// @Param   schedule body orgdto.ScheduleList true "Schedule data"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/schedules [post]
func (o *OrgCtrl) AddWorkerSchedule(w http.ResponseWriter, r *http.Request) {
	req := &orgdto.ScheduleList{}
	if err := o.json.NewDecoder(r.Body).Decode(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := o.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := o.usecase.AddWorkerSchedule(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}