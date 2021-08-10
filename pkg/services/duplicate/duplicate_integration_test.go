// +build integration

package duplicate

import (
	"testing"

	"best-practice-action/pkg/helper"

	"best-practice-action/pkg/models"
	fa "best-practice-action/pkg/services/fileaction"

	"github.com/stretchr/testify/require"
)

// проверка интеграции файловой системы
func TestDuplicate_ReadDirectory(t *testing.T) {
	path := "/home/d/projects/gb/best-practice-action/test/test_integration"

	i := &ioutilStruct{}
	f := &fa.FileActions{FS: i}

	expected := models.FilesInfo{}
	helper.AddFileInfo(path+"/copy/aaaa", "aaaa", &expected)

	// вынуждена тестить тот же метод, что и для мок, но входные данные такие,
	// что буду просто тестировать открытие одной папки и одного файла
	res, err := f.GetAllFiles(path)

	require.NoError(t, err)
	require.Len(t, res.List, len(expected.List))
	require.Equal(t, res.List[0], expected.List[0])
}
