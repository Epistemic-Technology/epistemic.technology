package backend

import (
	"database/sql"
	"fmt"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func GetDB(path string) (*DB, error) {
	sqlite_vec.Auto()
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	sDB := &DB{db: db}

	err = Init(sDB)
	if err != nil {
		return nil, err
	}

	return sDB, nil
}

func Init(db *DB) error {
	_, err := db.db.Exec(`
		CREATE TABLE IF NOT EXISTS documents (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT,
			content TEXT NOT NULL,
			author TEXT,
			publication_date TEXT,
			url TEXT,
			file_path TEXT
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create documents table: %w", err)
	}

	_, err = db.db.Exec(`
		CREATE TABLE IF NOT EXISTS chunks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			document_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			hash BLOB NOT NULL,
			FOREIGN KEY (document_id) REFERENCES documents(id)
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create chunks table: %w", err)
	}

	_, err = db.db.Exec(`
		CREATE VIRTUAL TABLE IF NOT EXISTS vec_chunks
		USING vec0(
			id INTEGER PRIMARY KEY,
			embedding FLOAT[1536]
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to create vec_chunks table: %w", err)
	}

	return nil
}

func InsertDocument(db *DB, doc *Document) error {
	result, err := db.db.Exec(`
		INSERT INTO documents (title, content, author, publication_date, url, file_path)
		VALUES (?, ?, ?, ?, ?, ?)
	`, doc.Title, doc.Content, doc.Author, doc.PublicationDate, doc.URL, doc.FilePath)
	if err != nil {
		return fmt.Errorf("failed to insert document: %w", err)
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}
	doc.ID = int(lastInsertID)
	return nil
}

func GetDocumentByID(db *DB, id int) (Document, error) {
	var doc Document
	row := db.db.QueryRow(`
		SELECT id, title, content, author, publication_date, url, file_path
		FROM documents
		WHERE id = ?
	`, id)
	err := row.Scan(&doc.ID, &doc.Title, &doc.Content, &doc.Author, &doc.PublicationDate, &doc.URL, &doc.FilePath)
	if err != nil {
		return Document{}, fmt.Errorf("failed to get document: %w", err)
	}
	return doc, nil
}

func GetAllDocuments(db *DB) ([]Document, error) {
	docs := []Document{}
	rows, err := db.db.Query(`
		SELECT id, title, content, author, publication_date, url, file_path
		FROM documents
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all documents: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var doc Document
		err = rows.Scan(&doc.ID, &doc.Title, &doc.Content, &doc.Author, &doc.PublicationDate, &doc.URL, &doc.FilePath)
		if err != nil {
			return nil, fmt.Errorf("failed to scan document: %w", err)
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

func GetAllChunks(db *DB) ([]Chunk, error) {
	chunks := []Chunk{}
	rows, err := db.db.Query(`
		SELECT id, content, hash, document_id
		FROM chunks
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all chunks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var chunk Chunk
		err = rows.Scan(&chunk.ID, &chunk.Content, &chunk.Hash, &chunk.DocumentID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chunk: %w", err)
		}
		chunks = append(chunks, chunk)
	}
	return chunks, nil
}

func GetAllChunksWithDocumentID(db *DB, docID int) ([]Chunk, error) {
	chunks := []Chunk{}
	rows, err := db.db.Query(`
		SELECT id, content, hash, document_id
		FROM chunks
		WHERE document_id = ?
	`, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to get all chunks with document ID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var chunk Chunk
		err = rows.Scan(&chunk.ID, &chunk.Content, &chunk.Hash, &chunk.DocumentID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chunk: %w", err)
		}
		chunks = append(chunks, chunk)
	}
	return chunks, nil
}

func GetDocumentChunks(db *DB, docID int) ([]Chunk, error) {
	chunks := []Chunk{}
	rows, err := db.db.Query(`
		SELECT id, content, hash, document_id
		FROM chunks
		WHERE document_id = ?
	`, docID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document chunks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var chunk Chunk
		err = rows.Scan(&chunk.ID, &chunk.Content, &chunk.Hash, &chunk.DocumentID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan chunk: %w", err)
		}
		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

func InsertChunk(db *DB, chunk *Chunk) error {
	result, err := db.db.Exec(`
		INSERT INTO chunks (document_id, content, hash)
		VALUES (?, ?, ?)
	`, chunk.DocumentID, chunk.Content, chunk.Hash)
	if err != nil {
		return fmt.Errorf("failed to insert chunk: %w", err)
	}

	// Get the last inserted ID
	lastID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to get last insert ID: %w", err)
	}

	embeddingFloat := make([]float32, len(chunk.Embedding))
	for i, v := range chunk.Embedding {
		embeddingFloat[i] = float32(v)
	}
	serializedEmbedding, err := sqlite_vec.SerializeFloat32(embeddingFloat)
	if err != nil {
		return fmt.Errorf("failed to serialize embedding: %w", err)
	}

	_, err = db.db.Exec(`
		INSERT INTO vec_chunks (id, embedding)
		VALUES (?, ?)
	`, lastID, serializedEmbedding)
	if err != nil {
		return fmt.Errorf("failed to insert vec_chunk: %w", err)
	}

	// Update the chunk ID with the database ID
	chunk.ID = int(lastID)

	return nil
}

func SimilaritySearch(db *DB, embedding Embedding, limit int) ([]Chunk, error) {
	embeddingFloat := make([]float32, len(embedding))
	for i, v := range embedding {
		embeddingFloat[i] = float32(v)
	}
	serializedEmbedding, err := sqlite_vec.SerializeFloat32(embeddingFloat)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize embedding: %w", err)
	}
	results, err := db.db.Query(`
		SELECT
			chunks.id,
			chunks.content,
			chunks.hash,
			chunks.document_id
		FROM chunks
		JOIN vec_chunks ON chunks.id = vec_chunks.id
		WHERE vec_chunks.embedding MATCH ?
		AND vec_chunks.k = ?
		ORDER BY vec_chunks.distance
	`, serializedEmbedding, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to perform similarity search: %w", err)
	}
	defer results.Close()

	chunks := []Chunk{}
	for results.Next() {
		var chunk Chunk
		err = results.Scan(&chunk.ID, &chunk.Content, &chunk.Hash, &chunk.DocumentID)
		if err != nil {
			return nil, fmt.Errorf("failed to scan result: %w", err)
		}

		chunks = append(chunks, chunk)
	}

	return chunks, nil
}

func Close(db *DB) error {
	return db.db.Close()
}
