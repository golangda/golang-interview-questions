package message

import "context"

type Repository interface {
    Create(ctx context.Context, m Message) (Message, error)
    GetByID(ctx context.Context, id int) (Message, error)
    Update(ctx context.Context, m Message) (Message, error)
    Delete(ctx context.Context, id int) error
}