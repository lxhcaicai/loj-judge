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

func (f *fileHandle) filePost(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	fi, err := fh.Open()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	sf, err := f.fs.New()
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer sf.Close()

	if _, err := sf.ReadFrom(fi); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	id, err := f.fs.Add(fh.Filename, sf.Name())
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.JSON(http.StatusOK, id)
}
