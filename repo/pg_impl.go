package repo

import (
	"context"
	"fmt"
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Postgres implements repo.Vocab interface for working with PostgreSQL DB.
type Postgres struct {
	pool *pgxpool.Pool
}

func NewPostgresRepo(p *pgxpool.Pool) *Postgres {
	return &Postgres{pool: p}
}

const (
	addVocab         = "INSERT INTO vocab(user_id) VALUES ($1) RETURNING id"
	getVocabByUserID = "SELECT id, user_id FROM vocab WHERE user_id = $1"
)

// AddVocab inserts the given vocab to DB and returns it with inserted ID.
func (p *Postgres) AddVocab(vocab *domain.Vocab) (*domain.Vocab, error) {
	row := p.pool.QueryRow(context.Background(), addVocab, vocab.UserID)
	err := row.Scan(&vocab.ID)
	if err != nil {
		return nil, fmt.Errorf("error adding vocab: %s", err)
	}
	return vocab, nil
}

// GetVocabByUserID returns the vocab found by the given user ID.
// Returns nil and no error if vocab is not found.
func (p *Postgres) GetVocabByUserID(userID int) (*domain.Vocab, error) {
	row := p.pool.QueryRow(context.Background(), getVocabByUserID, userID)
	vocab := new(domain.Vocab)
	err := row.Scan(&vocab.ID, &vocab.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting vocab by user ID: %s", err)
	}
	return vocab, nil
}

func (p *Postgres) ClosePool() {
	p.pool.Close()
}
