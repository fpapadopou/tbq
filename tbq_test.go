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

func TestTBQ_Consume_ContextCancelled(t *testing.T) {

	ctrl := gomock.NewController(t)
	source := NewMockSource(ctrl)

	ctx, cancelContext := context.WithCancel(context.TODO())
	q := &TBQ{
		s: source,
	}

	cancelContext()
	err := q.Consume(ctx)

	if err != nil {
		t.Errorf("Consume() error = %v", err)
	}
}

func TestTBQ_Consume_SourceReceiveError(t *testing.T) {

	ctrl := gomock.NewController(t)
	source := NewMockSource(ctrl)
	source.EXPECT().Receive(gomock.Any()).Return(nil, errors.New("source receiver error")).AnyTimes()

	q := &TBQ{
		s: source,
	}

	err := q.Consume(context.Background())

	if err == nil {
		t.Error("Consume() expected error but did not get any")
	}
}

func TestTBQ_Consume_ProcessorFuncError(t *testing.T) {

	ctrl := gomock.NewController(t)
	source := NewMockSource(ctrl)
	source.EXPECT().Receive(gomock.Any()).Return("received item", nil).AnyTimes()

	processor := func(ctx context.Context, i interface{}) error {
		return errors.New("processing of item failed")
	}
	q := &TBQ{
		s:    source,
		proc: processor,
	}

	err := q.Consume(context.Background())

	if err == nil {
		t.Error("Consume() expected error but did not get any")
	}
}
