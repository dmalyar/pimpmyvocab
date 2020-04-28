package repo

import (
	"context"
	"fmt"
	"github.com/dmalyar/pimpmyvocab/domain"
	"github.com/dmalyar/pimpmyvocab/log"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Postgres implements repo.Vocab interface for working with PostgreSQL DB.
type Postgres struct {
	logger log.Logger
	pool   *pgxpool.Pool
}

func NewPostgresRepo(logger log.Logger, pool *pgxpool.Pool) *Postgres {
	return &Postgres{logger: logger, pool: pool}
}

const (
	addVocab         = "INSERT INTO vocab(user_id) VALUES ($1) RETURNING id"
	getVocabByUserID = "SELECT id, user_id FROM vocab WHERE user_id = $1"

	addVocabEntry = "INSERT INTO vocab_entry(text, transcription) " +
		"VALUES ($1, $2) RETURNING id"
	getVocabEntryByText = "SELECT id, text, transcription " +
		"FROM vocab_entry WHERE text = $1"
	getVocabEntryByID = "SELECT id, text, transcription " +
		"FROM vocab_entry WHERE id = $1"

	addTranslation = "INSERT INTO translation(vocab_entry_id, text, class, position) " +
		"VALUES ($1, $2, $3, $4) RETURNING id"
	getTranslationsByEntryID = "SELECT id, text, class, position " +
		"FROM translation WHERE vocab_entry_id = $1"
)

// AddVocab inserts the given vocab to DB and returns it with inserted ID.
func (p *Postgres) AddVocab(vocab *domain.Vocab) (*domain.Vocab, error) {
	contextLogger := p.logger.WithField("vocab", vocab)
	contextLogger.Debug("Inserting vocab into DB")
	row := p.pool.QueryRow(context.Background(), addVocab, vocab.UserID)
	err := row.Scan(&vocab.ID)
	if err != nil {
		return nil, fmt.Errorf("inserting vocab into DB: %s", err)
	}
	contextLogger.Debug("Vocab inserted into DB")
	return vocab, nil
}

// GetVocabByUserID returns the vocab found by the given user ID.
// Returns nil and no error if vocab is not found.
func (p *Postgres) GetVocabByUserID(userID int) (*domain.Vocab, error) {
	contextLogger := p.logger.WithField("userID", userID)
	contextLogger.Debug("Getting vocab by user ID from DB")
	row := p.pool.QueryRow(context.Background(), getVocabByUserID, userID)
	vocab := new(domain.Vocab)
	err := row.Scan(&vocab.ID, &vocab.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			contextLogger.Debug("Vocab not found in DB")
			return nil, nil
		}
		return nil, fmt.Errorf("getting vocab from DB: %s", err)
	}
	contextLogger = contextLogger.WithField("vocab", vocab)
	contextLogger.Debugf("Vocab found in DB")
	return vocab, nil
}

// AddVocabEntry inserts the given vocab entry to DB and returns it with inserted ID.
func (p *Postgres) AddVocabEntry(entry *domain.VocabEntry) (*domain.VocabEntry, error) {
	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("getting transaction: %s", err)
	}
	defer tx.Rollback(context.Background())

	contextLogger := p.logger.WithField("vocabEntry", entry)
	contextLogger.Debugf("Inserting vocab entry into DB")
	row := p.pool.QueryRow(context.Background(), addVocabEntry,
		entry.Text, entry.Transcription)
	err = row.Scan(&entry.ID)
	if err != nil {
		return nil, fmt.Errorf("inserting vocab entry into DB: %s", err)
	}
	for _, t := range entry.Translations {
		row := p.pool.QueryRow(context.Background(), addTranslation,
			entry.ID, t.Text, t.Class, t.Position)
		err = row.Scan(&t.ID)
		if err != nil {
			return nil, fmt.Errorf("inserting translation into DB: %s", err)
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, fmt.Errorf("commiting transaction: %s", err)
	}
	contextLogger.Debug("Vocab entry inserted into DB")
	return entry, nil
}

// GetVocabEntryByText returns the vocab entry found by the given text.
// Returns nil and no error if vocab entry is not found.
func (p *Postgres) GetVocabEntryByText(text string) (*domain.VocabEntry, error) {
	contextLogger := p.logger.WithField("text", text)
	contextLogger.Debug("Getting vocab entry by text from DB")
	row := p.pool.QueryRow(context.Background(), getVocabEntryByText, text)
	return p.getVocabEntry(contextLogger, row)
}

// GetVocabEntryByID returns the vocab entry found by the given ID.
// Returns nil and no error if vocab entry is not found.
func (p *Postgres) GetVocabEntryByID(id int) (*domain.VocabEntry, error) {
	contextLogger := p.logger.WithField("ID", id)
	contextLogger.Debug("Getting vocab entry by ID from DB")
	row := p.pool.QueryRow(context.Background(), getVocabEntryByID, id)
	return p.getVocabEntry(contextLogger, row)
}

func (p *Postgres) getVocabEntry(contextLogger log.Logger, row pgx.Row) (*domain.VocabEntry, error) {
	entry := new(domain.VocabEntry)
	err := row.Scan(&entry.ID, &entry.Text, &entry.Transcription)
	if err != nil {
		if err == pgx.ErrNoRows {
			contextLogger.Debug("Vocab entry not found in DB")
			return nil, nil
		}
		return nil, fmt.Errorf("getting vocab entry from DB: %s", err)
	}
	rows, err := p.pool.Query(context.Background(), getTranslationsByEntryID, entry.ID)
	if err != nil {
		return nil, fmt.Errorf("getting translations by entry ID: %s", err)
	}
	for rows.Next() {
		t := new(domain.Translation)
		entry.Translations = append(entry.Translations, t)
		err := rows.Scan(&t.ID, &t.Text, &t.Class, &t.Position)
		if err != nil {
			return nil, fmt.Errorf("scanning translation row: %s", err)
		}
	}
	contextLogger = contextLogger.WithField("vocabEntry", entry)
	contextLogger.Debug("Entry found in DB")
	return entry, nil
}

func (p *Postgres) ClosePool() {
	p.pool.Close()
}
