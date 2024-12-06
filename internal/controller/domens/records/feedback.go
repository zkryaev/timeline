package records

import (
	"context"
	"net/http"
	"timeline/internal/controller/validation"
	"timeline/internal/entity/dto/recordto"

	"github.com/gorilla/mux"
)

type Feedback interface {
	Feedback(ctx context.Context, params *recordto.FeedbackParams) (*recordto.Feedback, error)
	FeedbackSet(ctx context.Context, feedback *recordto.Feedback) error
	FeedbackUpdate(ctx context.Context, feedback *recordto.Feedback) error
	FeedbackDelete(ctx context.Context, params *recordto.FeedbackParams) error
}

// @Summary Feedback info
// @Description Get feedbakc for specified record
// @Tags record / feedback
// @Param   recordID path int true "record_id"
// @Success 200 {object} recordto.Feedback
// @Failure 400
// @Failure 500
// @Router /records/feedbacks/{recordID} [get]
func (rec *RecordCtrl) Feedback(w http.ResponseWriter, r *http.Request) {
	params, err := validation.FetchSpecifiedID(mux.Vars(r), "recordID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req := &recordto.FeedbackParams{RecordID: params["recordID"]}
	if err := rec.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	data, err := rec.usecase.Feedback(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if rec.json.NewEncoder(w).Encode(&data) != nil {
		http.Error(w, "An error occurred while processing the response", http.StatusInternalServerError)
		return
	}
}

// @Summary Set feedback
// @Description Set feedback for specified record
// @Tags record / feedback
// @Accept  json
// @Param feedback body recordto.Feedback true "Feedback data"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /records/feedbacks [post]
func (rec *RecordCtrl) FeedbackSet(w http.ResponseWriter, r *http.Request) {
	req := &recordto.Feedback{}
	if rec.json.NewDecoder(r.Body).Decode(req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	if err := rec.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := rec.usecase.FeedbackSet(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Update feedback
// @Description Update feedback for specified record
// @Tags record / feedback
// @Accept  json
// @Param feedback body recordto.Feedback true "Feedback data"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /records/feedbacks [put]
func (rec *RecordCtrl) FeedbackUpdate(w http.ResponseWriter, r *http.Request) {
	req := &recordto.Feedback{}
	if rec.json.NewDecoder(r.Body).Decode(req) != nil {
		http.Error(w, "An error occurred while processing the request", http.StatusBadRequest)
		return
	}
	if err := rec.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := rec.usecase.FeedbackUpdate(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete feedback
// @Description Delete feedback for specified record
// @Tags record / feedback
// @Param   recordID path int true "record_id"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /records/feedbacks/{recordID} [delete]
func (rec *RecordCtrl) FeedbackDelete(w http.ResponseWriter, r *http.Request) {
	params, err := validation.FetchSpecifiedID(mux.Vars(r), "recordID")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req := &recordto.FeedbackParams{RecordID: params["recordID"]}
	if err := rec.validator.Struct(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := rec.usecase.FeedbackDelete(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
