package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("/dev/urandom", func(t *testing.T) {
		tmp, _ := os.CreateTemp("/tmp", "test-")
		defer os.Remove(tmp.Name())

		require.FileExists(t, tmp.Name())

		err := Copy("/dev/urandom", tmp.Name(), 0, 256)

		f, _ := os.Open(tmp.Name())
		fstat, _ := f.Stat()

		require.NoError(t, err)
		require.Equal(t, int64(256), fstat.Size())
	})

	t.Run("/dev/urandom", func(t *testing.T) {
		tmp, _ := os.CreateTemp("/tmp", "test-")
		defer os.Remove(tmp.Name())

		require.FileExists(t, tmp.Name())

		err := Copy("/dev/urandom", tmp.Name(), 0, 0)

		f, _ := os.Open(tmp.Name())
		fstat, _ := f.Stat()

		require.NoError(t, err)
		require.Equal(t, int64(0), fstat.Size())
	})
}

func TestCopyFailed(t *testing.T) {
	defer os.Remove("/tmp/out.txt")

	t.Run("offset exceeds", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/out.txt", 100_000, 0)

		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("offset exceeds with limit", func(t *testing.T) {
		err := Copy("testdata/input.txt", "/tmp/out.txt", 10_000, 3_000)

		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("from /dev/urandom with offset", func(t *testing.T) {
		err := Copy("/dev/urandom", "/tmp/out.txt", 1, 0)

		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("non exists", func(t *testing.T) {
		err := Copy("testdata/input-123.txt", "/tmp/out.txt", 10_000, 3_000)

		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("from dir", func(t *testing.T) {
		err := Copy("testdata/", "/tmp/out.txt", 0, 0)

		require.Error(t, err)
	})
}
