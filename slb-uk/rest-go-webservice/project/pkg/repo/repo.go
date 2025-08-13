package repo

import "database/sql"

type Repo struct{ DB *sql.DB }

func (r *Repo) CheckIdempotency(tx *sql.Tx, key string) (bool, error) { return false, nil }
func (r *Repo) MarkIdempotent(tx *sql.Tx, key, traceID, status string) error { return nil }
func (r *Repo) InsertMessage(tx *sql.Tx, msg string) (int64, error) { return 0, nil }