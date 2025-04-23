package changelog

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq" // Postgres driver
)

var ErrChangeNotFound = errors.New("change not found")

// ChangeRepository defines the interface for changelog operations
type ChangeRepository interface {
	List(limit, offset int) ([]Change, error)
	Get(id string) (Change, error)
	Create(change Change) (Change, error)
	Delete(id string) error
}

type changeRepository struct {
	db *sql.DB
}

func Init(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS changes (
        id BIGSERIAL PRIMARY KEY,
        timestamp TIMESTAMP NOT NULL,
        actor TEXT NOT NULL,
        service TEXT NOT NULL,
        details TEXT NOT NULL
    );`

	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// NewChangeRepository creates a new repository instance
func NewChangeRepository(db *sql.DB) ChangeRepository {
	return &changeRepository{db: db}
}

// List retrieves changes within last 1 day
func (r *changeRepository) List(limit, offset int) ([]Change, error) {
	// TODO: Add time range filtering
	query := `
        SELECT id, timestamp, actor, service, details
        FROM changes
        ORDER BY timestamp ASC
        LIMIT $1 OFFSET $2`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var changes []Change
	for rows.Next() {
		var c Change
		if err := rows.Scan(&c.ID, &c.Timestamp, &c.Actor, &c.Service, &c.Details); err != nil {
			return nil, err
		}
		changes = append(changes, c)
	}
	return changes, rows.Err()
}

// Get fetches a single change by ID
func (r *changeRepository) Get(id string) (Change, error) {
	query := `
        SELECT id, timestamp, actor, service, details
        FROM changes
        WHERE id = $1`
	var c Change
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Timestamp, &c.Actor, &c.Service, &c.Details)
	if err == sql.ErrNoRows {
		return Change{}, ErrChangeNotFound
	}
	return c, err
}

// Create adds a new change to the database
func (r *changeRepository) Create(change Change) (Change, error) {
	query := `
        INSERT INTO changes (timestamp, actor, service, details)
        VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(query, change.Timestamp, change.Actor, change.Service, change.Details)
	return change, err
}

// Delete removes a change by ID
func (r *changeRepository) Delete(id string) error {
	query := `DELETE FROM changes WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrChangeNotFound
	}
	return nil
}
