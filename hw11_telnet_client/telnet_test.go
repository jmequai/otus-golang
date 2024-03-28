package main

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeout, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeout, io.NopCloser(in), out, nil)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("no connection", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("127.0.0.1:4242", 10*time.Second, io.NopCloser(in), out, nil)

		require.ErrorContains(t, client.Connect(), "connection refused")
	})

	t.Run("timeout", func(t *testing.T) {
		_, err := net.Listen("tcp", "127.0.0.1:4242")
		require.NoError(t, err)

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("127.0.0.1:4242", 10*time.Nanosecond, io.NopCloser(in), out, nil)

		require.Error(t, client.Connect())
		require.ErrorContains(t, client.Connect(), "i/o timeout")

		client.Close()
	})

	t.Run("closed connection", func(t *testing.T) {
		_, err := net.Listen("tcp", "127.0.0.1:4243")
		require.NoError(t, err)

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("127.0.0.1:4243", 10*time.Second, io.NopCloser(in), out, nil)

		require.NoError(t, client.Connect())
		require.NoError(t, client.Close())

		in.WriteString("hello\n")
		require.ErrorIs(t, client.Send(), ErrNoConnection)
	})

	t.Run("closed connection", func(t *testing.T) {
		_, err := net.Listen("tcp", "127.0.0.1:4245")
		require.NoError(t, err)

		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		client := NewTelnetClient("127.0.0.1:4245", 10*time.Second, io.NopCloser(in), out, nil)

		require.NoError(t, client.Connect())
		require.NoError(t, client.Close())

		require.ErrorIs(t, client.Receive(), ErrNoConnection)
	})
}
