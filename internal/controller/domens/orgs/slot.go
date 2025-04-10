package orgs

import (
	"context"
	"errors"
	"net/http"
	"timeline/internal/controller/common"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/orgdto"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Slots interface {
	Slots(ctx context.Context, logger *zap.Logger, req *orgdto.SlotReq) ([]*orgdto.SlotResp, error)
	UpdateSlot(ctx context.Context, logger *zap.Logger, req *orgdto.SlotUpdate) error
}

// @Summary Get slots
// @Description Get all slots for specified worker
// @Tags organization/slots
// @Produce json
// @Param   workerID path int true "worker_id"
// @Param   orgID path int true "org_id"
// @Success 200 {array} orgdto.SlotResp
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /orgs/{orgID}/slots/workers/{workerID} [get]
func (o *OrgCtrl) Slots(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	params, err := validation.FetchPathID(mux.Vars(r), "workerID", "orgID")
	if err != nil {
		logger.Error("FetchPathID", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	token, _ := o.middleware.ExtractToken(r)
	tdata := common.GetTokenData(token.Claims)
	req := &orgdto.SlotReq{WorkerID: params["workerID"], OrgID: params["orgID"]}
	if !tdata.IsOrg {
		req.UserID = tdata.ID
	}
	if err := common.Validate(req); err != nil {
		logger.Error("Validate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	data, err := o.usecase.Slots(r.Context(), logger, req)
	if err != nil {
		switch {
		case errors.Is(err, common.ErrNotFound):
			logger.Info("Slots", zap.Error(err))
			http.Error(w, "", http.StatusNotFound)
			return
		default:
			logger.Error("Slots", zap.Error(err))
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

// @Summary Update slots
// @Description Update specified slot for specified worker
// @Tags organization/slots
// @Accept json
// @Param   request body orgdto.SlotUpdate true "slots info"
// @Param   orgID path int true "org_id"
// @Success 200
// @Failure 304
// @Failure 400
// @Failure 500
// @Router /orgs/{orgID}/slots [put]
func (o *OrgCtrl) UpdateSlot(w http.ResponseWriter, r *http.Request) {
	uuid, _ := r.Context().Value("uuid").(string)
	logger := o.Logger.With(zap.String("uuid", uuid))
	req := &orgdto.SlotUpdate{}
	if err := common.DecodeAndValidate(r, req); err != nil {
		logger.Error("DecodeAndValidate", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	if err := o.usecase.UpdateSlot(r.Context(), logger, req); err != nil {
		switch {
		case errors.Is(err, common.ErrNothingChanged):
			logger.Info("UpdateSlot", zap.Error(err))
			http.Error(w, "", http.StatusNotModified)
			return
		default:
			logger.Error("UpdateSlot", zap.Error(err))
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
