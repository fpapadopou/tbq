package tbqueue

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"testing"
	"time"
)

func TestTBQ_Publish(t *testing.T) {

	currentTime := time.Now()

	tests := []struct {
		name       string
		mockSource func() Source
		wantErr    bool
	}{
		{
			"success",
			func() Source {
				ctrl := gomock.NewController(t)
				source := NewMockSource(ctrl)
				source.EXPECT().Send(context.TODO(), "foo-item", currentTime.UTC().Unix()).Return(nil).Times(1)

				return source
			},
			false,
		},
		{
			"source error",
			func() Source {
				ctrl := gomock.NewController(t)
				source := NewMockSource(ctrl)
				source.EXPECT().Send(context.TODO(), "foo-item", currentTime.UTC().Unix()).Return(errors.New("source error")).Times(1)

				return source
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := tt.mockSource()
			q := &TBQ{
				s: s,
			}
			if err := q.Publish(context.TODO(), "foo-item", currentTime); (err != nil) != tt.wantErr {
				t.Errorf("Publish() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
