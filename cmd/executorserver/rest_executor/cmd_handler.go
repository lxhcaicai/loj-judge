package restexecutor

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/lxhcaicai/loj-judge/cmd/executorserver/model"
	"github.com/lxhcaicai/loj-judge/filestore"
	"github.com/lxhcaicai/loj-judge/worker"
	"go.uber.org/zap"
	"net/http"
)

type Register interface {
	Register(engine *gin.Engine)
}

func New(worker worker.Worker, fs filestore.FileStore, srcPrefix []string, logger *zap.Logger) Register {
	return &handle{
		worker:     worker,
		fileHandle: fileHandle{fs: fs},
		srcPrefix:  srcPrefix,
		logger:     logger,
	}
}

type handle struct {
	worker worker.Worker
	fileHandle
	srcPrefix []string
	logger    *zap.Logger
}

func (h *handle) Register(r *gin.Engine) {

	r.POST("/run", h.handleRun)

	// File handle
	r.GET("/file", h.fileGet)
	r.POST("/file", h.filePost)
	r.GET("/file/:fid", h.fileIDGet)
	r.DELETE("/file/:fid", h.fileIDDelete)
}

func (h *handle) handleRun(c *gin.Context) {
	var req model.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	if len(req.Cmd) == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, "no cmd provided")
		return
	}
	r, err := model.ConvertRequest(&req, h.srcPrefix)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}
	h.logger.Sugar().Debugf("request: %+v", r)
	rtCh, _ := h.worker.Submit(c.Request.Context(), r)
	rt := <-rtCh
	h.logger.Sugar().Debugf("response: %+v", rt)
	if rt.Error != nil {
		c.Error(rt.Error)
		c.AbortWithStatusJSON(http.StatusInternalServerError, rt.Error.Error())
		return
	}

	// 直接编码json以避免分配
	c.Status(http.StatusOK)
	c.Header("Content-Type", "application/json; charset=utf-8")

	res, err := model.ConvertResponse(rt, true)
	if err != nil {
		c.Error(err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err.Error())
		return
	}
	defer res.Close()

	if err := json.NewEncoder(c.Writer).Encode(res.Results); err != nil {
		c.Error(err)
	}
}
