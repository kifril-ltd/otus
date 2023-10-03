package main

import (
	"errors"
	"io"
	"os"

	"github.com/kifril-ltd/otus-hw/hw07_file_copying/progressbar"
	"github.com/kifril-ltd/otus-hw/hw07_file_copying/utils"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

const BufferSize = 1024 * 1024

func Copy(fromPath, toPath string, offset, limit int64) error {
	info, err := os.Stat(fromPath)
	if err != nil {
		return err
	}
	if info.IsDir() || info.Size() == 0 {
		return ErrUnsupportedFile
	}
	if offset > info.Size() {
		return ErrOffsetExceedsFileSize
	}

	src, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer src.Close()

	if _, err = src.Seek(offset, 0); err != nil {
		return err
	}

	if limit == 0 {
		limit = info.Size()
	}

	dst, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	buffer := make([]byte, BufferSize)
	var bytesCopied int64

	total := utils.Min(limit, info.Size()-offset)
	bar := progressbar.NewBar(
		total,
		progressbar.WithSymbol("*"),
		progressbar.WithWidth(50),
	)
	for bytesCopied < limit {
		n, err := src.Read(buffer)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		if (bytesCopied + int64(n)) > limit {
			n = int(limit - bytesCopied)
		}

		if _, err := dst.Write(buffer[:n]); err != nil {
			return err
		}

		bytesCopied += int64(n)

		bar.Draw(bytesCopied)
	}

	bar.Finish()

	return nil
}
