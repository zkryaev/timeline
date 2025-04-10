package orgs

import (
	"context"
	"errors"
	"net/http"
	"strconv"
	"timeline/internal/controller/common"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/orgdto"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Timetable interface {
	Timetable(ctx context.Context, logger *zap.Logger, req orgdto.TimetableReq) (*orgdto.Timetable, error)
	TimetableAdd(ctx context.Context, logger *zap.Logger, newTimetable *orgdto.Timetable) error
	TimetableUpdate(ctx context.Context, logger *zap.Logger, newTimetable *orgdto.Timetable) error
	TimetableDelete(ctx context.Context, logger *zap.Logger, orgID, weekday int) error
}

// @Summary Get timetable
// @Description Get organization timetable
// @Tags organization / timetables
// @Accept  json
// @Param   orgID path int true "org_id"
// @Success 200 {object} orgdto.Timetable
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /orgs/{orgID}/timetable [get]
func (o *OrgCtrl) Timetable(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	path, err := validation.FetchPathID(mux.Vars(r), "orgID")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	token, _ := o.middleware.ExtractToken(r)
	tdata := common.GetTokenData(token.Claims)
	req := orgdto.TimetableReq{OrgID: path["orgID"]}
	if !tdata.IsOrg {
		req.UserID = tdata.ID
	}
	data, err := o.usecase.Timetable(r.Context(), logger, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Error("Timetable", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("Timetable", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	if err := common.WriteJSON(w, data); err != nil {
		logger.Error("WriteJSON", zap.Error(err))
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
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/timetable [post]
func (o *OrgCtrl) TimetableAdd(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	req := &orgdto.Timetable{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.TimetableAdd(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrTimeIncorrect):
			logger.Info("TimetableAdd", zap.Error(err))
			http.Error(w, common.ErrTimeIncorrect.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("TimetableAdd", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("TimetableAdd", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Update timetable
// @Description Update organization timetable
// @Tags organization / timetables
// @Accept  json
// @Param   request body orgdto.Timetable true "New org info"
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/timetable [put]
func (o *OrgCtrl) TimetableUpdate(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	req := &orgdto.Timetable{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.TimetableUpdate(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrTimeIncorrect):
			logger.Info("TimetableUpdate", zap.Error(err))
			http.Error(w, common.ErrTimeIncorrect.Error(), http.StatusBadRequest)
			return
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("TimetableUpdate", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("TimetableUpdate", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
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
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/timetable [delete]
func (o *OrgCtrl) TimetableDelete(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	params, err := validation.FetchPathID(mux.Vars(r), "orgID")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	var weekday int
	if r.URL.Query().Get("weekday") != "" {
		weekday, err = strconv.Atoi(r.URL.Query().Get("weekday"))
		if err != nil {
			logger.Error("Atoi", zap.Error(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if o.usecase.TimetableDelete(r.Context(), logger, params["orgID"], weekday) != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("TimetableDelete", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("TimetableDelete", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
