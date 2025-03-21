package orgs

import (
	"context"
	"net/http"
	"strconv"
	"timeline/internal/controller/common"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/orgdto"

	"github.com/gorilla/mux"
)

type Timetable interface {
	Timetable(ctx context.Context, OrgID int) (*orgdto.Timetable, error)
	TimetableAdd(ctx context.Context, newTimetable *orgdto.Timetable) error
	TimetableUpdate(ctx context.Context, newTimetable *orgdto.Timetable) error
	TimetableDelete(ctx context.Context, orgID, weekday int) error
}

// @Summary Get timetable
// @Description Get organization timetable
// @Tags organization / timetables
// @Accept  json
// @Param   orgID path int true "org_id"
// @Success 200 {object} orgdto.Timetable
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/timetable [get]
func (o *OrgCtrl) Timetable(w http.ResponseWriter, r *http.Request) {
	path, err := validation.FetchPathID(mux.Vars(r), "orgID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := o.usecase.Timetable(r.Context(), path["orgID"])
	if err != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if common.WriteJSON(w, data) != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}

// @Summary Add timetable
// @Description Add organization timetable
// @Tags organization / timetables
// @Accept  json
// @Param   request body orgdto.Timetable true "New org info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/timetable [post]
func (o *OrgCtrl) TimetableAdd(w http.ResponseWriter, r *http.Request) {
	req := &orgdto.Timetable{}
	if common.DecodeAndValidate(r, req) != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if o.usecase.TimetableAdd(r.Context(), req) != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Update timetable
// @Description Update organization timetable
// @Tags organization / timetables
// @Accept  json
// @Param   request body orgdto.Timetable true "New org info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/timetable [put]
func (o *OrgCtrl) TimetableUpdate(w http.ResponseWriter, r *http.Request) {
	req := &orgdto.Timetable{}
	if common.DecodeAndValidate(r, req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	if o.usecase.TimetableUpdate(r.Context(), req) != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete timetable
// @Description Delete organization timetable. If weekday doesnt set then whole timetable will be deleted
// @Tags organization / timetables
// @Accept  json
// @Param orgID path int true "org_id"
// @Param weekday query int false "weekday"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/timetable [delete]
func (o *OrgCtrl) TimetableDelete(w http.ResponseWriter, r *http.Request) {
	params, err := validation.FetchPathID(mux.Vars(r), "orgID")
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
	if o.usecase.TimetableDelete(r.Context(), params["orgID"], weekday) != nil {
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
