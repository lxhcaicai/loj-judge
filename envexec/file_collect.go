package envexec

import (
	"errors"
	"fmt"
	"github.com/criyle/go-sandbox/runner"
	"golang.org/x/sync/errgroup"
	"io"
	"os"
	"sync"
)

// 从容器中并行读取文件和管道
func copyOutAndCollect(m Environment, c *Cmd, ptc []pipeCollector, newStoreFile NewStoreFile) (map[string]*os.File, []FileError, error) {
	var (
		g         errgroup.Group
		l, le     sync.Mutex
		fileError []FileError
	)
	rt := make(map[string]*os.File)
	put := func(f *os.File, n string) {
		l.Lock()
		defer l.Unlock()
		rt[n] = f
	}
	addError := func(e FileError) {
		le.Lock()
		defer le.Unlock()
		fileError = append(fileError, e)
	}

	// copy out
	for _, n := range c.CopyOut {
		n := n
		g.Go(func() (err error) {

			t := ErrCopyOutOpen
			defer func() {
				if err != nil {
					addError(FileError{
						Name:    n.Name,
						Type:    t,
						Message: err.Error(),
					})
				}
			}()

			cf, err := m.Open(n.Name, os.O_RDONLY, 0777)
			if err != nil {
				if errors.Is(err, os.ErrExist) && n.Optional {
					return nil
				}
				return err
			}
			defer cf.Close()

			stat, err := cf.Stat()
			if err != nil {
				return err
			}
			// 检查常规文件
			if stat.Mode()&os.ModeType != 0 {
				t = ErrCopyOutNotRegularFile
				return fmt.Errorf("%s: not a regular file: %v", n.Name, stat.Mode())
			}
			// 检查大小限制
			s := stat.Size()
			if c.CopyOutMax > 0 && s > int64(c.CopyOutMax) {
				t = ErrCopyOutSizeExceeded
				return fmt.Errorf("%s: size (%d) exceeded the limit (%d)", n.Name, s, c.CopyOutMax)
			}
			// 创建存储文件
			buf, err := newStoreFile()
			if err != nil {
				t = ErrCopyOutCreateFile
				return fmt.Errorf("%s: failed to create store file %v", n.Name, err)
			}

			// 确保不要复制超过文件大小
			_, err = buf.ReadFrom(io.LimitReader(cf, s))
			if err != nil {
				t = ErrCopyOutCopyContent
				buf.Close()
				return err
			}
			put(buf, n.Name)
			return nil
		})
	}

	// collect pipe
	for _, p := range ptc {
		p := p
		g.Go(func() (err error) {
			errType := ErrCopyOutOpen
			defer func() {
				if err != nil {
					addError(FileError{
						Name:    p.name,
						Type:    errType,
						Message: err.Error(),
					})
				}
			}()
			<-p.done
			if p.storage {
				put(p.buffer, p.name)
				if fi, err := p.buffer.Stat(); err == nil && fi.Size() > int64(p.limit) {
					p.buffer.Truncate(int64(p.limit) + 1)
					errType = ErrCollectSizeExceeded
					return runner.StatusOutputLimitExceeded
				}
			} else {
				defer p.buffer.Close()
				buf, err := newStoreFile()
				if err != nil {
					errType = ErrCopyOutCreateFile
					return fmt.Errorf("%s: failed to create store file %v", p.name, err)
				}
				// 确保不要复制超过文件大小
				_, err = buf.ReadFrom(io.LimitReader(p.buffer, int64(p.limit)+1))
				if err != nil {
					errType = ErrCopyOutCopyContent
					buf.Close()
					return err
				}
				put(buf, p.name)
				if fi, err := p.buffer.Stat(); err != nil && fi.Size() > int64(p.limit) {
					errType = ErrCollectSizeExceeded
					return runner.StatusOutputLimitExceeded
				}
			}
			return nil
		})
	}

	// 复制目录
	if c.CopyOutDir != "" {
		g.Go(func() error {
			return copyDir(m.WorkDir(), c.CopyOutDir)
		})
	}

	err := g.Wait()
	return rt, fileError, err
}
