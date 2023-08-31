package restexecutor

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lxhcaicai/loj-judge/envexec"
	"github.com/lxhcaicai/loj-judge/filestore"
	"io"
	"mime"
	"net/http"
	"path"
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

func (f *fileHandle) fileIDGet(c *gin.Context) {
	type fileURI struct {
		FileID string `uri:"fid"`
	}
	var uri fileURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	name, file := f.fs.Get(uri.FileID)
	if file == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	r, err := envexec.FileToReader(file)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	defer r.Close()

	content, err := io.ReadAll(r)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	typ := mime.TypeByExtension(path.Ext(name))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", name))
	c.Data(http.StatusOK, typ, content)
}
