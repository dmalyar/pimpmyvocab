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
	addVocab           = "INSERT INTO vocab(user_id) VALUES ($1) RETURNING id"
	getVocabByUserID   = "SELECT id, user_id FROM vocab WHERE user_id = $1"
	clearVocabByUserID = "DELETE FROM vocab_to_entry_link " +
		"WHERE vocab_id = (SELECT ID from vocab WHERE user_id = $1)"

	addVocabEntry = "INSERT INTO vocab_entry(text, transcription) " +
		"VALUES ($1, $2) RETURNING id"
	getVocabEntryByText = "SELECT id, text, transcription " +
		"FROM vocab_entry WHERE text = $1"
	getVocabEntryByID = "SELECT id, text, transcription " +
		"FROM vocab_entry WHERE id = $1"

	addEntryToUserVocab = "INSERT INTO vocab_to_entry_link(vocab_id, entry_id)" +
		"SELECT id, $1 FROM vocab WHERE user_id = $2"
	checkEntryInUserVocab = "SELECT l.entry_id " +
		"FROM vocab v " +
		"JOIN vocab_to_entry_link l on v.id = l.vocab_id " +
		"WHERE l.entry_id = $1 and v.user_id = $2"
	getEntryIDsByUserID = "SELECT l.entry_id " +
		"FROM vocab v " +
		"JOIN vocab_to_entry_link l on v.id = l.vocab_id " +
		"WHERE v.user_id = $1"
	getEntriesByUserID = "SELECT e.id, e.text, e.transcription, t.text " +
		"FROM vocab v " +
		"JOIN vocab_to_entry_link l on v.id = l.vocab_id " +
		"JOIN vocab_entry e on l.entry_id = e.id " +
		"JOIN translation t on e.id = t.vocab_entry_id " +
		"WHERE v.user_id = $1 AND t.position = 0"
	removeEntryFromUserVocab = "DELETE FROM vocab_to_entry_link " +
		"WHERE entry_id = $1 " +
		"AND vocab_id = (SELECT ID from vocab WHERE user_id = $2)"

	addTranslation = "INSERT INTO translation(vocab_entry_id, text, class, position) " +
		"VALUES ($1, $2, $3, $4) RETURNING id"
	getTranslationsByEntryID = "SELECT id, text, class, position " +
		"FROM translation WHERE vocab_entry_id = $1 " +
		"ORDER BY position"
)

// AddVocab inserts the given vocab to DB and returns it with inserted ID.
func (p *Postgres) AddVocab(vocab *domain.Vocab) (*domain.Vocab, error) {
	logger := p.logger.WithField("vocab", vocab)
	logger.Debug("Inserting vocab into DB")
	row := p.pool.QueryRow(context.Background(), addVocab, vocab.UserID)
	err := row.Scan(&vocab.ID)
	if err != nil {
		return nil, fmt.Errorf("inserting vocab into DB: %s", err)
	}
	logger.Debug("Vocab inserted into DB")
	return vocab, nil
}

// GetVocabByUserID returns the vocab found by the given user ID.
// Returns nil and no error if vocab was not found.
func (p *Postgres) GetVocabByUserID(userID int) (*domain.Vocab, error) {
	logger := p.logger.WithField("userID", userID)
	logger.Debug("Getting vocab by user ID from DB")
	row := p.pool.QueryRow(context.Background(), getVocabByUserID, userID)
	vocab := new(domain.Vocab)
	err := row.Scan(&vocab.ID, &vocab.UserID)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Debug("Vocab not found in DB")
			return nil, nil
		}
		return nil, fmt.Errorf("getting vocab from DB: %s", err)
	}
	logger = logger.WithField("vocab", vocab)
	logger.Debugf("Vocab found in DB")
	return vocab, nil
}

// ClearVocabByUserID removes all links to entries from the user's vocab.
func (p *Postgres) ClearVocabByUserID(userID int) error {
	logger := p.logger.WithField("userID", userID)
	logger.Debug("Removing all links to entries from the user's vocab in DB")
	_, err := p.pool.Exec(context.Background(), clearVocabByUserID, userID)
	if err != nil {
		return fmt.Errorf("removing all links to entries from user's vocab in DB: %s", err)
	}
	return nil
}

// AddVocabEntry inserts the given vocab entry to DB and returns it with inserted ID.
func (p *Postgres) AddVocabEntry(entry *domain.VocabEntry) (*domain.VocabEntry, error) {
	tx, err := p.pool.Begin(context.Background())
	if err != nil {
		return nil, fmt.Errorf("getting transaction: %s", err)
	}
	defer tx.Rollback(context.Background())

	logger := p.logger.WithField("vocabEntry", entry)
	logger.Debugf("Inserting vocab entry into DB")
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
	logger.Debug("Vocab entry inserted into DB")
	return entry, nil
}

// GetVocabEntryByText returns the vocab entry found by the given text.
// Returns nil and no error if vocab entry was not found.
func (p *Postgres) GetVocabEntryByText(text string) (*domain.VocabEntry, error) {
	logger := p.logger.WithField("text", text)
	logger.Debug("Getting vocab entry by text from DB")
	row := p.pool.QueryRow(context.Background(), getVocabEntryByText, text)
	return p.getVocabEntry(logger, row)
}

// GetVocabEntryByID returns the vocab entry found by the given ID.
// Returns nil and no error if vocab entry was not found.
func (p *Postgres) GetVocabEntryByID(id int) (*domain.VocabEntry, error) {
	logger := p.logger.WithField("id", id)
	logger.Debug("Getting vocab entry by ID from DB")
	row := p.pool.QueryRow(context.Background(), getVocabEntryByID, id)
	return p.getVocabEntry(logger, row)
}

func (p *Postgres) getVocabEntry(logger log.Logger, row pgx.Row) (*domain.VocabEntry, error) {
	entry := new(domain.VocabEntry)
	err := row.Scan(&entry.ID, &entry.Text, &entry.Transcription)
	if err != nil {
		if err == pgx.ErrNoRows {
			logger.Debug("Vocab entry not found in DB")
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
		if t.Position == 0 {
			entry.MainTranslation = t.Text
		}
	}
	logger = logger.WithField("vocabEntry", entry)
	logger.Debug("Entry found in DB")
	return entry, nil
}

// AddEntryToUserVocab links entry with the given ID to the user's vocab.
func (p *Postgres) AddEntryToUserVocab(entryID, userID int) error {
	logger := p.logger.WithFields(map[string]interface{}{
		"entryID": entryID,
		"userID":  userID,
	})
	logger.Debug("Adding the entry to the user's vocab in DB")
	_, err := p.pool.Exec(context.Background(), addEntryToUserVocab, entryID, userID)
	if err != nil {
		return fmt.Errorf("adding entry to user's vocab in DB: %s", err)
	}
	return nil
}

// CheckEntryInUserVocab returns if the entry is linked to the user's vocab.
func (p *Postgres) CheckEntryInUserVocab(entryID, userID int) (bool, error) {
	logger := p.logger.WithFields(map[string]interface{}{
		"entryID": entryID,
		"userID":  userID,
	})
	logger.Debug("Checking if the entry is added to the user's vocab in DB")
	row := p.pool.QueryRow(context.Background(), checkEntryInUserVocab, entryID, userID)
	id := new(int)
	err := row.Scan(id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("checking if entry is added to vocab DB: %s", err)
	}
	return true, nil
}

// GetEntryIDsByUserID returns IDs of all entries linked to the user's vocab.
func (p *Postgres) GetEntryIDsByUserID(userID int) ([]int, error) {
	contextLog := p.logger.WithField("userID", userID)
	contextLog.Debug("Getting entry IDs from the user's vocab from DB")
	rows, err := p.pool.Query(context.Background(), getEntryIDsByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("getting entry IDs from the user's vocab from DB: %s", err)
	}
	var ids []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, fmt.Errorf("scanning row with id: %s", err)
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// GetEntriesByUserID returns all entries linked to the user's vocab.
// Returned entries have only main translation.
func (p *Postgres) GetEntriesByUserID(userID int) ([]*domain.VocabEntry, error) {
	contextLog := p.logger.WithField("userID", userID)
	contextLog.Debug("Getting entries from the user's vocab from DB")
	rows, err := p.pool.Query(context.Background(), getEntriesByUserID, userID)
	if err != nil {
		return nil, fmt.Errorf("getting entries from the user's vocab from DB: %s", err)
	}
	var entries []*domain.VocabEntry
	for rows.Next() {
		e := new(domain.VocabEntry)
		entries = append(entries, e)
		err := rows.Scan(&e.ID, &e.Text, &e.Transcription, &e.MainTranslation)
		if err != nil {
			return nil, fmt.Errorf("scanning entry row: %s", err)
		}
	}
	return entries, nil
}

// RemoveEntryFromUserVocab removes link of the entry with given ID to the user's vocab.
func (p *Postgres) RemoveEntryFromUserVocab(entryID, userID int) error {
	logger := p.logger.WithFields(map[string]interface{}{
		"entryID": entryID,
		"userID":  userID,
	})
	logger.Debug("Removing the entry from the user's vocab in DB")
	_, err := p.pool.Exec(context.Background(), removeEntryFromUserVocab, entryID, userID)
	if err != nil {
		return fmt.Errorf("removing entry from user's vocab in DB: %s", err)
	}
	return nil
}

func (p *Postgres) ClosePool() {
	p.pool.Close()
}
