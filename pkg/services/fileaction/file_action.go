package fileaction

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"sync"

	"best-practice-action/pkg/models"

	log "github.com/sirupsen/logrus"
)

type FS interface {
	ReadDir(string) ([]fs.FileInfo, error)
	ReadFile(string) ([]byte, error)
}

type FileActions struct {
	FS FS
	m  sync.Mutex
	wg sync.WaitGroup

	fi models.FilesInfo
}

func (f *FileActions) GetAllFiles(path string) (models.FilesInfo, error) {
	l := log.WithField("FuncName", "getAllFiles").WithField("path", path)
	l.Debugf("run get all files")

	f.fi = models.FilesInfo{}

	files, err := f.FS.ReadDir(path)

	if err != nil {
		return f.fi, fmt.Errorf("directory %s  does not open. error: %v", path, err)
	}

	// Линтер errcheck  - не было проверки на ошибку
	err = f.ReadDirectory(path, files)
	if err != nil {
		return f.fi, err
	}

	return f.fi, nil
}

func (f *FileActions) ReadDirectory(path string, files []fs.FileInfo) error {
	l := log.WithField("FuncName", "readDirectory").WithField("path", path)

	for _, file := range files {
		newPath := path + "/" + file.Name()
		if !file.IsDir() {
			l.Debug("read file:", newPath)

			fileByte, err := f.FS.ReadFile(newPath)
			if err != nil {
				return fmt.Errorf("file: %s error: %s", path, err)
			}
			// Линтер wsl - скзаал не прижимать выражение к блокам
			f.wg.Add(1)

			go f.addFileInfo(fileByte, newPath, file.Name())
		} else {
			l.Debug("read directory:", newPath)
			dir, err := f.FS.ReadDir(newPath)

			if err != nil {
				return fmt.Errorf("directory %s  does not open. err: %v", newPath, err)
			}

			// Линтер errcheck  - не было проверки на ошибку
			// Линтер wsl - названия переменной
			if errRead := f.ReadDirectory(newPath, dir); errRead != nil {
				return err
			}
		}
	}

	f.wg.Wait()

	return nil
}

// Линтер gocriric - параметры типа стринг, можно объединить
// Линтер revive - исправила название переменных
func (f *FileActions) addFileInfo(file []byte, path, fileName string) {
	defer f.wg.Done()

	h1 := md5.New()
	h1.Write(file)
	hashMd5 := hex.EncodeToString(h1.Sum(nil))

	h2 := sha256.New()
	h2.Write(file)
	hashSha256 := hex.EncodeToString(h2.Sum(nil))

	f.m.Lock()
	f.fi.AddItem(fileName, path, hashMd5, hashSha256)
	f.m.Unlock()
}
