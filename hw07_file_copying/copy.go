package main

import (
	"errors"
	"io"
	"os"
	"time"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrInvalidOffset         = errors.New("invalid offset")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrInvalidLimit          = errors.New("invalid limit")
	ErrInvalidFile           = errors.New("invalid file")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrFileCreationFail      = errors.New("file creation fail")
	ErrFileOpenningFail      = errors.New("file openning fail")
	ErrFileSeekingFail       = errors.New("file seeking fail")
	ErrCopyingFail           = errors.New("copying fail")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	// 1. Проверяем смещение и лимит на валидность
	if offset < 0 {
		return ErrInvalidOffset
	}
	if limit < 0 {
		return ErrInvalidLimit
	}

	// 2. Получаем мета-данные по файлу-источнику
	srcFileInfo, err := os.Stat(fromPath)
	if err != nil {
		return ErrInvalidFile
	}

	// 3. Проверяем файл-источник на валидность
	if srcFileInfo.IsDir() {
		return ErrInvalidFile
	}
	if srcFileInfo.Size() < 1 {
		return ErrUnsupportedFile
	}

	// 4. Проверяем смещение и лимит на соответствие файлу-истоничку
	if offset > srcFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}
	if limit == 0 {
		limit = srcFileInfo.Size()
	}

	// 3. Создаем новый файл для копирования
	dstFile, err := os.Create(toPath)
	if err != nil {
		return ErrFileCreationFail
	}
	defer dstFile.Close()

	// 4. Открываем файл-источник и смещаемся на Offset
	srcFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0o777)
	if err != nil {
		return ErrFileOpenningFail
	}
	defer srcFile.Close()
	_, err = srcFile.Seek(offset, 0)
	if err != nil {
		return ErrFileSeekingFail
	}

	// 5. Создаем прогресс-бар
	bar := pb.New(int(limit))
	bar.SetRefreshRate(time.Millisecond * 100)
	bar.Start()
	defer bar.Finish()
	proxy := bar.NewProxyReader(srcFile)

	// 6. Копируем
	_, err = io.CopyN(dstFile, proxy, limit)
	if errors.Is(err, io.EOF) {
		return nil
	}
	if errors.Is(err, nil) {
		return nil
	}
	return ErrCopyingFail
}
