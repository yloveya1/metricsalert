package filestore

import (
	"encoding/json"
	"os"

	models "github.com/yloveya1/metricsalert/internal/model"
	"github.com/yloveya1/metricsalert/internal/repository"
)

type FileStorage struct {
	filePath string
}

func NewFileStorage(filename string) repository.IFile {
	return &FileStorage{filePath: filename}
}

func (f *FileStorage) WriteMetrics(metrics []*models.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}

	data, err := json.Marshal(metrics)
	if err != nil {
		return err
	}

	return os.WriteFile(f.filePath, data, 0644)
}

func (f *FileStorage) UploadMetrics() ([]*models.Metrics, error) {
	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return nil, err
	}

	var metrics []*models.Metrics
	if len(data) == 0 {
		return metrics, nil
	}

	if err = json.Unmarshal(data, &metrics); err != nil {
		return nil, err
	}

	return metrics, nil
}
