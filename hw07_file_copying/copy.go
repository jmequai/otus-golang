package main

import (
	"errors"
	"io"
	"os"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")

	fBufSize int64 = 1024
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	f, fErrOpen := os.Open(fromPath)
	if fErrOpen != nil {
		return ErrUnsupportedFile
	}

	defer f.Close()

	fstat, fErrStat := f.Stat()
	if fErrStat != nil {
		return ErrUnsupportedFile
	}

	fsize := fstat.Size()

	if offset > fsize {
		return ErrOffsetExceedsFileSize
	}

	if limit == 0 || (fsize > 0 && offset+limit > fsize) {
		limit = fsize - offset
	}

	_, fErrSeek := f.Seek(offset, 0)
	if fErrSeek != nil {
		return fErrSeek
	}

	t, tErrOpen := os.OpenFile(toPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if tErrOpen != nil {
		return tErrOpen
	}

	defer t.Close()

	fProgress := NewProgress(limit)
	defer fProgress.Finish()

	fCopySize := int64(0)

	for {
		if fCopySize >= limit {
			break
		}

		currentBufSize := limit - fCopySize

		if currentBufSize > fBufSize {
			currentBufSize = fBufSize
		}

		n, err := io.CopyN(t, f, currentBufSize)

		fCopySize += n
		fProgress.Add(n)
		// time.Sleep(time.Second)

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return err
		}
	}

	return nil
}
