package main

import (
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	cases := []struct {
		name, from, to string
		limit, offset  int64
		err            error
	}{
		{
			name: "Unsupported Input",
			from: "/dev/urandom",
			to:   "./output.txt",
			err:  ErrUnsupportedFile,
		},
		{
			name:   "Invalid Offset",
			from:   "./testdata/input.txt",
			offset: 1_000_000,
			to:     "./output.txt",
			err:    ErrOffsetExceedsFileSize,
		},
		{
			name:   "Input File Not Found",
			from:   "./testdata/404.txt",
			offset: 1_000_000,
			to:     "./output.txt",
			err:    os.ErrNotExist,
		},
		{
			name: "Success Case",
			from: "./testdata/input.txt",
			to:   "./output.txt",
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			result := Copy(testCase.from, testCase.to, testCase.offset, testCase.limit)

			if testCase.err != nil {
				require.ErrorIs(t, result, testCase.err)
			} else {
				require.NoError(t, result)
			}

			_ = os.Remove(testCase.to)
		})
	}
}

func TestCopyLargeFile(t *testing.T) {
	srcSize := int64(1 * 1024 * 1024 * 1024)
	srcPath := "./testdata/large.txt"
	dstPath := "./output.txt"

	defer func() {
		_ = os.Remove(srcPath)
		_ = os.Remove(dstPath)
	}()

	generateFile := func(path string, size int64) {
		f, err := os.Create(path)
		require.NoError(t, err)
		defer f.Close()

		buffer := make([]byte, 10*1024*1024)
		rand.Seed(42)

		var writtenBytes int64
		for writtenBytes < size {
			n, err := rand.Read(buffer)
			require.NoError(t, err)

			_, err = f.Write(buffer[:n])
			require.NoError(t, err)

			writtenBytes += int64(n)
		}
	}

	generateFile(srcPath, srcSize)

	result := Copy(srcPath, dstPath, 0, 0)
	require.NoError(t, result)
}
