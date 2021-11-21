package main

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const ChunkSize = 1024

func TestCopy(t *testing.T) {
	from := "./testdata/input.txt"
	to := "./testdata/out.txt"

	t.Run("copying whole file without offset", func(t *testing.T) {
		err := Copy(from, to, 0, 0)
		defer os.Remove(to)
		require.NoError(t, err)

		// для начала просто сравним размер файлов
		srcFileInfo, err := os.Stat(from)
		require.NoError(t, err)
		dstFileInfo, err := os.Stat(to)
		require.NoError(t, err)
		require.Equal(t, srcFileInfo.Size(), dstFileInfo.Size())

		// размер совпал, значит имеет смысл сравнить данные в копии и "оригинале"
		templateFile, err := os.ReadFile(from)
		require.NoError(t, err)
		copiedFile, err := os.ReadFile(to)
		require.NoError(t, err)
		require.Equal(t, templateFile, copiedFile)
	})

	t.Run("copying 10 bytes without offset", func(t *testing.T) {
		err := Copy(from, to, 0, 10)
		defer os.Remove(to)
		require.NoError(t, err)

		templateFileInfo, err := os.Stat("./testdata/out_offset0_limit10.txt")
		require.NoError(t, err)
		copiedFileInfo, err := os.Stat(to)
		require.NoError(t, err)
		require.Equal(t, templateFileInfo.Size(), copiedFileInfo.Size())

		templateFile, err := os.ReadFile("./testdata/out_offset0_limit10.txt")
		require.NoError(t, err)
		copiedFile, err := os.ReadFile(to)
		require.NoError(t, err)
		require.Equal(t, templateFile, copiedFile)
	})

	t.Run("copying 1000 bytes without offset", func(t *testing.T) {
		err := Copy(from, to, 0, 1000)
		defer os.Remove(to)
		require.NoError(t, err)

		templateFileInfo, err := os.Stat("./testdata/out_offset0_limit1000.txt")
		require.NoError(t, err)
		copiedFileInfo, err := os.Stat(to)
		require.NoError(t, err)
		require.Equal(t, templateFileInfo.Size(), copiedFileInfo.Size())

		templateFile, err := os.ReadFile("./testdata/out_offset0_limit1000.txt")
		require.NoError(t, err)
		copiedFile, err := os.ReadFile(to)
		require.NoError(t, err)
		require.Equal(t, templateFile, copiedFile)
	})

	t.Run("copying more bytes than file size without offset", func(t *testing.T) {
		err := Copy(from, to, 0, 10000)
		defer os.Remove(to)
		require.NoError(t, err)

		templateFileInfo, err := os.Stat("./testdata/out_offset0_limit10000.txt")
		require.NoError(t, err)
		copiedFileInfo, err := os.Stat(to)
		require.NoError(t, err)
		require.Equal(t, templateFileInfo.Size(), copiedFileInfo.Size())

		templateFile, err := os.ReadFile("./testdata/out_offset0_limit10000.txt")
		require.NoError(t, err)
		copiedFile, err := os.ReadFile(to)
		require.NoError(t, err)
		require.Equal(t, templateFile, copiedFile)
	})

	t.Run("copying 1000 bytes with 100 bytes offset", func(t *testing.T) {
		err := Copy(from, to, 100, 1000)
		defer os.Remove(to)
		require.NoError(t, err)

		// для начала просто сравним размер файлов
		templateFileInfo, err := os.Stat("./testdata/out_offset100_limit1000.txt")
		require.NoError(t, err)
		copiedFileInfo, err := os.Stat(to)
		require.NoError(t, err)
		require.Equal(t, templateFileInfo.Size(), copiedFileInfo.Size())

		// размер совпал, значит имеет смысл сравнить данные в копии и "оригинале"
		// *здесь сравним чанками, потому что задуман типа большой файл
		templateFile, err := os.Open("./testdata/out_offset100_limit1000.txt")
		require.NoError(t, err)
		defer templateFile.Close()
		copiedFile, err := os.Open(to)
		require.NoError(t, err)
		defer copiedFile.Close()
		var isTemplateFileEOF, isCopiedFileEOF bool
		for {
			templateFileChunk := make([]byte, ChunkSize)
			_, err := templateFile.Read(templateFileChunk)
			if err == io.EOF {
				isTemplateFileEOF = true
			} else {
				require.NoError(t, err)
			}

			copiedFileChunk := make([]byte, ChunkSize)
			_, err = copiedFile.Read(copiedFileChunk)
			if err == io.EOF {
				isCopiedFileEOF = true
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, isTemplateFileEOF, isCopiedFileEOF)
			require.Equal(t, templateFileChunk, copiedFileChunk)

			if isTemplateFileEOF && isCopiedFileEOF {
				break
			}
		}
	})

	t.Run("copying 1000 bytes with 6000 bytes offset", func(t *testing.T) {
		err := Copy(from, to, 6000, 1000)
		defer os.Remove(to)
		require.NoError(t, err)

		templateFileInfo, err := os.Stat("./testdata/out_offset6000_limit1000.txt")
		require.NoError(t, err)
		copiedFileInfo, err := os.Stat(to)
		require.NoError(t, err)
		require.Equal(t, templateFileInfo.Size(), copiedFileInfo.Size())

		templateFile, err := os.ReadFile("./testdata/out_offset6000_limit1000.txt")
		require.NoError(t, err)
		copiedFile, err := os.ReadFile(to)
		require.NoError(t, err)
		require.Equal(t, templateFile, copiedFile)
	})

	t.Run("copying with extra large offset", func(t *testing.T) {
		err := Copy(from, to, 1000000, 0)
		defer os.Remove(to)
		require.True(t, err != nil)
	})
}
