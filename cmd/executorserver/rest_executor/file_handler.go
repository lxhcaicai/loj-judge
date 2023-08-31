package restexecutor

import (
	"github.com/gin-gonic/gin"
	"github.com/lxhcaicai/loj-judge/filestore"
	"net/http"
)

type fileHandle struct {
	fs filestore.FileStore
}

func (f *fileHandle) fileGet(c *gin.Context) {
	ids := f.fs.List()
	c.JSON(http.StatusOK, ids)
}
