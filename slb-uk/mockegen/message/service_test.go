// crud/message/service_test.go
package message

import (
	"context"
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestService_Create(t *testing.T) {
    t.Parallel()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := NewMockRepository(ctrl)
    svc := NewService(mockRepo)
    ctx := context.Background()

    t.Run("success", func(t *testing.T) {
        in := " hello "
        want := Message{ID: 1, Content: "hello"}

        // Expect repo.Create with trimmed content
        mockRepo.
            EXPECT().
            Create(gomock.Any(), Message{Content: "hello"}).
            Return(want, nil).
            Times(1)

        got, err := svc.Create(ctx, in)
        require.NoError(t, err)
        require.Equal(t, want, got)
    })

    t.Run("validation: empty content", func(t *testing.T) {
        _, err := svc.Create(ctx, "   ")
        require.ErrorIs(t, err, ErrEmptyContent)
    })

    t.Run("repo error", func(t *testing.T) {
        mockRepo.
            EXPECT().
            Create(gomock.Any(), Message{Content: "x"}).
            Return(Message{}, errors.New("db down")).
            Times(1)

        _, err := svc.Create(ctx, "x")
        require.EqualError(t, err, "db down")
    })
}

func TestService_Get(t *testing.T) {
    t.Parallel()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := NewMockRepository(ctrl)
    svc := NewService(mockRepo)
    ctx := context.Background()

    t.Run("invalid id", func(t *testing.T) {
        _, err := svc.Get(ctx, 0)
        require.ErrorIs(t, err, ErrInvalidID)
    })

    t.Run("success", func(t *testing.T) {
        want := Message{ID: 42, Content: "hi"}

        mockRepo.
            EXPECT().
            GetByID(gomock.Any(), 42).
            Return(want, nil).
            Times(1)

        got, err := svc.Get(ctx, 42)
        require.NoError(t, err)
        require.Equal(t, want, got)
    })

    t.Run("repo error", func(t *testing.T) {
        mockRepo.
            EXPECT().
            GetByID(gomock.Any(), 7).
            Return(Message{}, errors.New("not found")).
            Times(1)

        _, err := svc.Get(ctx, 7)
        require.EqualError(t, err, "not found")
    })
}

func TestService_Update(t *testing.T) {
    t.Parallel()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := NewMockRepository(ctrl)
    svc := NewService(mockRepo)
    ctx := context.Background()

    t.Run("invalid id", func(t *testing.T) {
        _, err := svc.Update(ctx, -1, "x")
        require.ErrorIs(t, err, ErrInvalidID)
    })

    t.Run("empty content", func(t *testing.T) {
        _, err := svc.Update(ctx, 1, "   ")
        require.ErrorIs(t, err, ErrEmptyContent)
    })

    t.Run("success", func(t *testing.T) {
        want := Message{ID: 3, Content: "updated"}

        mockRepo.
            EXPECT().
            Update(gomock.Any(), Message{ID: 3, Content: "updated"}).
            Return(want, nil).
            Times(1)

        got, err := svc.Update(ctx, 3, "updated")
        require.NoError(t, err)
        require.Equal(t, want, got)
    })

    t.Run("repo error", func(t *testing.T) {
        mockRepo.
            EXPECT().
            Update(gomock.Any(), Message{ID: 5, Content: "xxx"}).
            Return(Message{}, errors.New("conflict")).
            Times(1)

        _, err := svc.Update(ctx, 5, "xxx")
        require.EqualError(t, err, "conflict")
    })
}

func TestService_Delete(t *testing.T) {
    t.Parallel()
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := NewMockRepository(ctrl)
    svc := NewService(mockRepo)
    ctx := context.Background()

    t.Run("invalid id", func(t *testing.T) {
        err := svc.Delete(ctx, 0)
        require.ErrorIs(t, err, ErrInvalidID)
    })

    t.Run("success", func(t *testing.T) {
        mockRepo.
            EXPECT().
            Delete(gomock.Any(), 9).
            Return(nil).
            Times(1)

        err := svc.Delete(ctx, 9)
        require.NoError(t, err)
    })

    t.Run("repo error", func(t *testing.T) {
        mockRepo.
            EXPECT().
            Delete(gomock.Any(), 10).
            Return(errors.New("foreign key")).
            Times(1)

        err := svc.Delete(ctx, 10)
        require.EqualError(t, err, "foreign key")
    })
}
