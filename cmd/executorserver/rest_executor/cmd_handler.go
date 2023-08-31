package restexecutor

import (
	"github.com/gin-gonic/gin"
	"github.com/lxhcaicai/loj-judge/filestore"
	"go.uber.org/zap"
)

type Register interface {
	Register(engine *gin.Engine)
}

func New(fs filestore.FileStore, srcPrefix []string, logger *zap.Logger) Register {
	return &handle{
		fileHandle: fileHandle{fs: fs},
		srcPrefix:  srcPrefix,
		logger:     logger,
	}
}

type handle struct {
	fileHandle
	srcPrefix []string
	logger    *zap.Logger
}

func (h *handle) Register(r *gin.Engine) {

	// File handle
	r.GET("/file", h.fileGet)
	r.POST("/file", h.filePost)
	r.GET("/file/:fid", h.fileIDGet)
}
