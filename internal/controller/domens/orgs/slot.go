package orgs

import (
	"context"
	"net/http"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/orgdto"

	"github.com/gorilla/mux"
)

type Slots interface {
	Slots(ctx context.Context, req *orgdto.SlotReq) ([]*orgdto.SlotResp, error)
	UpdateSlot(ctx context.Context, req *orgdto.SlotUpdate) error
}

// @Summary Get slots
// @Description Get all slots for specified worker
// @Tags organization/slots
// @Produce json
// @Param   workerID path int true "worker_id"
// @Success 200 {array} orgdto.SlotResp
// @Failure 400
// @Failure 500
// @Router /orgs/slots/{workerID} [get]
func (o *OrgCtrl) Slots(w http.ResponseWriter, r *http.Request) {
	params, err := validation.FetchSpecifiedID(mux.Vars(r), "workerID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req := &orgdto.SlotReq{WorkerID: params["workerID"]}
	if err = o.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := o.usecase.Slots(r.Context(), req)
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

// @Summary Update slots
// @Description Update specified slot for specified worker
// @Tags organization/slots
// @Accept json
// @Param   request body orgdto.SlotUpdate true "slots info"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /orgs/slots [put]
func (o *OrgCtrl) UpdateSlot(w http.ResponseWriter, r *http.Request) {
	req := &orgdto.SlotUpdate{}
	if o.json.NewDecoder(r.Body).Decode(req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	if err := o.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := o.usecase.UpdateSlot(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
