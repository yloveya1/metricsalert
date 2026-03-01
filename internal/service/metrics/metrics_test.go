package metrics

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/yloveya1/metricsalert/internal/mocks"
)

func Test_New(t *testing.T) {
	ctrl := gomock.NewController(t)
	store := mocks.NewMockIStorage(ctrl)

	s := NewService(store)

	assert.NotEmpty(t, s)
}

func TestService_UpdateCounterMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		prepare func(store *mocks.MockIStorage)
		wantErr bool
	}{
		{
			name: "success updating counter metric",
			prepare: func(store *mocks.MockIStorage) {
				store.EXPECT().UpdateCounterMetric(nil).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failed updating counter metric",
			prepare: func(store *mocks.MockIStorage) {
				store.EXPECT().UpdateCounterMetric(nil).Return(errors.New("error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewMockIStorage(ctrl)
			s := &Service{
				storage: store,
			}

			if tt.prepare != nil {
				tt.prepare(store)
			}
			err := s.UpdateCounterMetric(nil)
			assert.Equal(t, err == nil, tt.wantErr)
		})
	}
}

func TestService_UpdateGaugerMetric(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name    string
		prepare func(store *mocks.MockIStorage)
		wantErr bool
	}{
		{
			name: "success updating gauge metric",
			prepare: func(store *mocks.MockIStorage) {
				store.EXPECT().UpdateGaugeMetric(nil).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failed updating gauge metric",
			prepare: func(store *mocks.MockIStorage) {
				store.EXPECT().UpdateGaugeMetric(nil).Return(errors.New("error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := mocks.NewMockIStorage(ctrl)
			s := &Service{
				storage: store,
			}

			if tt.prepare != nil {
				tt.prepare(store)
			}
			err := s.UpdateGaugeMetric(nil)

			assert.Contains(t, err == nil, tt.wantErr)
		})
	}
}
