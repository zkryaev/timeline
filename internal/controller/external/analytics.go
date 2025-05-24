package external

import (
	"net/http"
	"strconv"
	"timeline/internal/config"
	"timeline/internal/controller/auth/middleware"
	"timeline/internal/controller/common"
	"timeline/internal/controller/query"
	"timeline/internal/controller/scope"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
)

var pathPool = buffer.NewPool()

type AnalitycsClient struct {
	address  string
	logger   *zap.Logger
	settings *scope.Settings
}

func NewAnalyticsClient(cfg config.AnalyticsService, logger *zap.Logger, settings *scope.Settings) *AnalitycsClient {
	return &AnalitycsClient{
		address:  "http://" + cfg.Host + ":" + cfg.Port,
		logger:   logger,
		settings: settings,
	}
}

func (a *AnalitycsClient) Retranslate(w http.ResponseWriter, r *http.Request) {
	logger := common.LoggerWithUUID(a.settings, a.logger, r.Context())
	tdata, err := middleware.GetTokenDataFromCtx(a.settings, r.Context())
	if err != nil {
		logger.Info("GetTokenDataFromCtx", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	orgID := query.NewParamString(scope.ORG_ID, false)
	if !a.settings.EnableAuthorization {
		orgID.Required = true
	} else {
		orgID.Val = strconv.Itoa(tdata.ID)
	}
	uri := query.NewParamString(scope.ANALYTICS_URI, true)
	params := query.NewParams(a.settings, orgID, uri)
	if err := params.Parse(r.URL.Query()); err != nil {
		logger.Error("param.Parse", zap.Error(err))
		http.Error(w, "", http.StatusBadRequest)
		return
	}
	a.logger.Info("Analytic redirector", zap.String("report", uri.Val))
	switch uri.Val {
	case "ai/feedbacks":
	case "cancellations":
	case "summary":
	case "workload":
	case "distribution/bookings":
	case "distribution/income":
	default:
		http.Error(w, "error: unknown report", http.StatusBadRequest)
		return
	}
	path := pathPool.Get()
	defer path.Free()
	path.AppendString(a.address + "/analytics" + "/" + uri.Val + "?" + scope.ORG_ID + "=" + orgID.Val)
	resp, err := http.Get(path.String())
	if err != nil {
		a.logger.Error("failed get request:", zap.String("Analytics service:", err.Error()))
		http.Error(w, "error: unknown report", http.StatusBadRequest)
		return
	}
	if err := resp.Write(w); err != nil {
		a.logger.Error("Analytic redirector: resp.Write error", zap.Error(err))
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
	a.logger.Info("Getting analytics completed")
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}
