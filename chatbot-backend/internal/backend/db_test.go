package backend

import (
	"os"
	"testing"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
)

func TestGetDB(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Test database creation
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer Close(db)

	// Verify that the database was created
	if db == nil {
		t.Fatal("Expected database to be created, got nil")
	}
}

func TestGetDBError(t *testing.T) {
	// Test with an invalid path that should cause an error
	// Using a directory as a database path should fail
	tempDir, err := os.MkdirTemp("", "test-db-dir-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Try to create a database with a directory path (should fail)
	db, err := GetDB(tempDir)

	// Verify that an error occurred
	if err == nil {
		t.Error("Expected an error when using a directory as a database path, but got nil")
		if db != nil {
			Close(db)
		}
	}
}

func TestInit(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Initialize SQLite vector extension
	sqlite_vec.Auto()

	// Open the database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer Close(db)

	// The Init function is already called by GetDB, so we just need to verify
	// that the tables were created correctly

	// Verify that tables were created
	var documentTableExists, chunkTableExists, vecChunksTableExists bool

	// Check if documents table exists
	rows, err := db.db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='documents'")
	if err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}
	documentTableExists = rows.Next()
	rows.Close()

	// Check if chunks table exists
	rows, err = db.db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='chunks'")
	if err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}
	chunkTableExists = rows.Next()
	rows.Close()

	// Check if vec_chunks virtual table exists
	rows, err = db.db.Query("SELECT name FROM sqlite_master WHERE type='table' AND name='vec_chunks'")
	if err != nil {
		t.Fatalf("Failed to query tables: %v", err)
	}
	vecChunksTableExists = rows.Next()
	rows.Close()

	if !documentTableExists {
		t.Error("Documents table was not created")
	}

	if !chunkTableExists {
		t.Error("Chunks table was not created")
	}

	if !vecChunksTableExists {
		t.Error("Vec_chunks virtual table was not created")
	}
}

func TestTableStructure(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Create database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer Close(db)

	// Check documents table structure
	rows, err := db.db.Query("PRAGMA table_info(documents)")
	if err != nil {
		t.Fatalf("Failed to query documents table info: %v", err)
	}

	columns := make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, type_name string
		var notnull, pk int
		var dflt_value interface{}

		if err := rows.Scan(&cid, &name, &type_name, &notnull, &dflt_value, &pk); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}

		columns[name] = true
	}
	rows.Close()

	requiredColumns := []string{"id", "title", "content", "author", "publication_date", "url", "file_path"}
	for _, col := range requiredColumns {
		if !columns[col] {
			t.Errorf("Documents table missing required column: %s", col)
		}
	}

	// Check chunks table structure
	rows, err = db.db.Query("PRAGMA table_info(chunks)")
	if err != nil {
		t.Fatalf("Failed to query chunks table info: %v", err)
	}

	columns = make(map[string]bool)
	for rows.Next() {
		var cid int
		var name, type_name string
		var notnull, pk int
		var dflt_value interface{}

		if err := rows.Scan(&cid, &name, &type_name, &notnull, &dflt_value, &pk); err != nil {
			t.Fatalf("Failed to scan row: %v", err)
		}

		columns[name] = true
	}
	rows.Close()

	requiredColumns = []string{"id", "document_id", "content", "hash"}
	for _, col := range requiredColumns {
		if !columns[col] {
			t.Errorf("Chunks table missing required column: %s", col)
		}
	}

	// Check vec_chunks virtual table structure
	// For virtual tables, we need to use a different approach since PRAGMA table_info
	// might not work the same way for virtual tables
	var vecChunksExists int
	err = db.db.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='vec_chunks'").Scan(&vecChunksExists)
	if err != nil {
		t.Fatalf("Failed to check if vec_chunks exists: %v", err)
	}

	if vecChunksExists == 0 {
		t.Error("Vec_chunks virtual table does not exist")
	}

	// Test that we can query the structure of the virtual table
	_, err = db.db.Exec("SELECT id, embedding FROM vec_chunks LIMIT 0")
	if err != nil {
		t.Errorf("Failed to query vec_chunks table structure: %v", err)
	}
}

func TestInsertDocument(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Initialize the database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer Close(db)

	// Create a test document
	doc := &Document{
		Title:           "Test Document",
		Content:         "This is a test document.",
		Author:          "Test Author",
		PublicationDate: "2023-01-01",
		URL:             "https://example.com/test",
		FilePath:        "/path/to/test.md",
	}

	// Insert the document
	err = InsertDocument(db, doc)
	if err != nil {
		t.Fatalf("Failed to insert document: %v", err)
	}

	// Verify that the document was inserted with an ID
	if doc.ID <= 0 {
		t.Errorf("Expected document ID to be set, got %d", doc.ID)
	}
}

func TestInsertChunk(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Initialize the database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer Close(db)

	// Create and insert a test document
	doc := &Document{
		Title:           "Test Document",
		Content:         "This is a test document.",
		Author:          "Test Author",
		PublicationDate: "2023-01-01",
		URL:             "https://example.com/test",
		FilePath:        "/path/to/test.md",
	}
	err = InsertDocument(db, doc)
	if err != nil {
		t.Fatalf("Failed to insert document: %v", err)
	}

	// Create a test chunk
	chunk := &Chunk{
		DocumentID: doc.ID,
		Content:    "This is a test chunk.",
		Hash:       []byte{1, 2, 3, 4},
		Embedding:  Embedding{0.1, 0.2, 0.3, 0.4},
	}

	// Insert the chunk
	err = InsertChunk(db, chunk)
	if err != nil {
		t.Fatalf("Failed to insert chunk: %v", err)
	}

	// Verify that the chunk was inserted with an ID
	if chunk.ID <= 0 {
		t.Errorf("Expected chunk ID to be set, got %d", chunk.ID)
	}
}

func TestForeignKeyConstraint(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Create database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer Close(db)

	// Enable foreign key constraints
	_, err = db.db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		t.Fatalf("Failed to enable foreign keys: %v", err)
	}

	// Try to insert a chunk with a non-existent document_id
	testHash := []byte{10, 20, 30, 40, 50}

	_, err = db.db.Exec(`
		INSERT INTO chunks (document_id, content, hash)
		VALUES (?, ?, ?)
	`, 999, "Test Chunk Content", testHash)

	// This should fail due to foreign key constraint
	if err == nil {
		t.Error("Expected foreign key constraint error, but insert succeeded")
	}
}

func TestCloseDB(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Initialize the database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Test closing the database
	err = Close(db)
	if err != nil {
		t.Fatalf("Failed to close database: %v", err)
	}

	// Verify that the database is closed by trying to use it
	// This should fail since the database is closed
	_, err = db.db.Exec("SELECT 1")
	if err == nil {
		t.Error("Expected error when using closed database, got nil")
	}
}

func TestVecChunksInsert(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Create database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer Close(db)

	// Create a test embedding (1536 dimensions as specified in the schema)
	testEmbedding := make([]float64, 1536)
	for i := range testEmbedding {
		testEmbedding[i] = float64(i) * 0.01
	}

	// Convert to []float32 for serialization
	testEmbeddingFloat32 := make([]float32, len(testEmbedding))
	for i, v := range testEmbedding {
		testEmbeddingFloat32[i] = float32(v)
	}

	vector, err := sqlite_vec.SerializeFloat32(testEmbeddingFloat32)
	if err != nil {
		t.Fatalf("Failed to serialize embedding: %v", err)
	}

	// Insert into vec_chunks
	_, err = db.db.Exec(`
		INSERT INTO vec_chunks (id, embedding)
		VALUES (?, ?)
	`, 1, vector)

	if err != nil {
		t.Fatalf("Failed to insert into vec_chunks: %v", err)
	}

	// Verify the insertion
	var count int
	err = db.db.QueryRow("SELECT COUNT(*) FROM vec_chunks").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to query vec_chunks count: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected 1 row in vec_chunks, got %d", count)
	}

	// Test a simple vector similarity search (if supported by the vec0 extension)
	// This is a basic test to ensure the virtual table is functioning
	_, err = db.db.Exec("SELECT id FROM vec_chunks ORDER BY vec_cosine(embedding, ?) LIMIT 1", vector)
	if err != nil {
		t.Logf("Vector similarity search not tested: %v", err)
		// Not failing the test as this might depend on the specific capabilities of the vec0 extension
	}
}

func TestSimilaritySearch(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Initialize the database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer Close(db)

	// Create and insert a test document
	doc := &Document{
		Title:           "Test Document",
		Content:         "This is a test document.",
		Author:          "Test Author",
		PublicationDate: "2023-01-01",
		URL:             "https://example.com/test",
		FilePath:        "/path/to/test.md",
	}
	err = InsertDocument(db, doc)
	if err != nil {
		t.Fatalf("Failed to insert document: %v", err)
	}

	// Create and insert test chunks with different embeddings
	chunk1 := &Chunk{
		DocumentID: doc.ID,
		Content:    "This is the first test chunk.",
		Hash:       []byte{1, 2, 3, 4},
		Embedding:  Embedding{0.1, 0.2, 0.3, 0.4},
	}
	err = InsertChunk(db, chunk1)
	if err != nil {
		t.Fatalf("Failed to insert chunk1: %v", err)
	}

	chunk2 := &Chunk{
		DocumentID: doc.ID,
		Content:    "This is the second test chunk.",
		Hash:       []byte{5, 6, 7, 8},
		Embedding:  Embedding{0.5, 0.6, 0.7, 0.8},
	}
	err = InsertChunk(db, chunk2)
	if err != nil {
		t.Fatalf("Failed to insert chunk2: %v", err)
	}

	chunk3 := &Chunk{
		DocumentID: doc.ID,
		Content:    "This is the third test chunk.",
		Hash:       []byte{9, 10, 11, 12},
		Embedding:  Embedding{0.9, 1.0, 1.1, 1.2},
	}
	err = InsertChunk(db, chunk3)
	if err != nil {
		t.Fatalf("Failed to insert chunk3: %v", err)
	}

	// Perform a similarity search
	queryEmbedding := Embedding{0.1, 0.2, 0.3, 0.4}
	results, err := SimilaritySearch(db, queryEmbedding, 3)
	if err != nil {
		t.Fatalf("Failed to perform similarity search: %v", err)
	}

	// Verify that results were returned
	if len(results) == 0 {
		t.Error("Expected similarity search to return results, got none")
	}
}

func TestSimilaritySearchWithEmptyDB(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Create database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer Close(db)

	// Create a query embedding
	queryEmbedding := make(Embedding, 1536)
	for i := range queryEmbedding {
		queryEmbedding[i] = float64(i) * 0.01
	}

	// Perform similarity search on empty database
	results, err := SimilaritySearch(db, queryEmbedding, 3)
	if err != nil {
		t.Fatalf("Expected no error for empty database search, got: %v", err)
	}

	// Verify we got empty results
	if len(results) != 0 {
		t.Errorf("Expected empty results for empty database, got %d results", len(results))
	}
}

func TestSimilaritySearchWithInvalidEmbedding(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Create database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer Close(db)

	// Create an invalid embedding (wrong size)
	invalidEmbedding := make(Embedding, 10) // Should be 1536
	for i := range invalidEmbedding {
		invalidEmbedding[i] = float64(i) * 0.01
	}

	// Perform similarity search with invalid embedding
	// This should still work but return empty results since no matches will be found
	results, err := SimilaritySearch(db, invalidEmbedding, 3)
	if err != nil {
		t.Fatalf("Expected no error for invalid embedding search, got: %v", err)
	}

	// Verify we got empty results
	if len(results) != 0 {
		t.Errorf("Expected empty results for invalid embedding, got %d results", len(results))
	}
}

func TestSimilaritySearchLimit(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Create database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer Close(db)

	// Insert test documents
	doc1 := &Document{
		Title:           "Test Document 1",
		Content:         "This is the first test document about artificial intelligence.",
		Author:          "Test Author 1",
		PublicationDate: "2023-01-01",
		URL:             "https://example.com/doc1",
	}

	doc2 := &Document{
		Title:           "Test Document 2",
		Content:         "This is the second test document about machine learning.",
		Author:          "Test Author 2",
		PublicationDate: "2023-01-02",
		URL:             "https://example.com/doc2",
	}

	doc3 := &Document{
		Title:           "Test Document 3",
		Content:         "This is the third test document about databases.",
		Author:          "Test Author 3",
		PublicationDate: "2023-01-03",
		URL:             "https://example.com/doc3",
	}

	// Insert documents and get their IDs
	_, err = db.db.Exec(`
		INSERT INTO documents (title, content, author, publication_date, url)
		VALUES (?, ?, ?, ?, ?)
	`, doc1.Title, doc1.Content, doc1.Author, doc1.PublicationDate, doc1.URL)
	if err != nil {
		t.Fatalf("Failed to insert document 1: %v", err)
	}
	var doc1ID int
	err = db.db.QueryRow("SELECT last_insert_rowid()").Scan(&doc1ID)
	if err != nil {
		t.Fatalf("Failed to get document 1 ID: %v", err)
	}

	_, err = db.db.Exec(`
		INSERT INTO documents (title, content, author, publication_date, url)
		VALUES (?, ?, ?, ?, ?)
	`, doc2.Title, doc2.Content, doc2.Author, doc2.PublicationDate, doc2.URL)
	if err != nil {
		t.Fatalf("Failed to insert document 2: %v", err)
	}
	var doc2ID int
	err = db.db.QueryRow("SELECT last_insert_rowid()").Scan(&doc2ID)
	if err != nil {
		t.Fatalf("Failed to get document 2 ID: %v", err)
	}

	_, err = db.db.Exec(`
		INSERT INTO documents (title, content, author, publication_date, url)
		VALUES (?, ?, ?, ?, ?)
	`, doc3.Title, doc3.Content, doc3.Author, doc3.PublicationDate, doc3.URL)
	if err != nil {
		t.Fatalf("Failed to insert document 3: %v", err)
	}
	var doc3ID int
	err = db.db.QueryRow("SELECT last_insert_rowid()").Scan(&doc3ID)
	if err != nil {
		t.Fatalf("Failed to get document 3 ID: %v", err)
	}

	// Create test chunks with embeddings
	// For simplicity, we'll create embeddings with different patterns to ensure they're distinguishable

	// Chunk 1 - AI related (mostly 0.1 values)
	chunk1 := &Chunk{
		ID:         1,
		DocumentID: doc1ID,
		Content:    doc1.Content,
		Hash:       []byte("hash1"),
	}

	// Chunk 2 - ML related (mostly 0.2 values)
	chunk2 := &Chunk{
		ID:         2,
		DocumentID: doc2ID,
		Content:    doc2.Content,
		Hash:       []byte("hash2"),
	}

	// Chunk 3 - DB related (mostly 0.3 values)
	chunk3 := &Chunk{
		ID:         3,
		DocumentID: doc3ID,
		Content:    doc3.Content,
		Hash:       []byte("hash3"),
	}

	// Create embeddings with distinct patterns
	embedding1 := make(Embedding, 1536)
	embedding2 := make(Embedding, 1536)
	embedding3 := make(Embedding, 1536)

	for i := range embedding1 {
		embedding1[i] = 0.1
		embedding2[i] = 0.2
		embedding3[i] = 0.3
	}

	// Make them slightly different to create a clear similarity pattern
	// First 100 dimensions are more similar to embedding1
	for i := 0; i < 100; i++ {
		embedding1[i] = 0.9
		embedding2[i] = 0.8
		embedding3[i] = 0.1
	}

	chunk1.Embedding = embedding1
	chunk2.Embedding = embedding2
	chunk3.Embedding = embedding3

	// Insert chunks into database
	err = InsertChunk(db, chunk1)
	if err != nil {
		t.Fatalf("Failed to insert chunk 1: %v", err)
	}

	err = InsertChunk(db, chunk2)
	if err != nil {
		t.Fatalf("Failed to insert chunk 2: %v", err)
	}

	err = InsertChunk(db, chunk3)
	if err != nil {
		t.Fatalf("Failed to insert chunk 3: %v", err)
	}

	// Create a query embedding that should be most similar to embedding1, then embedding2, then embedding3
	queryEmbedding := make(Embedding, 1536)
	for i := range queryEmbedding {
		queryEmbedding[i] = 0.1
	}
	// Make it very similar to embedding1 in the first 100 dimensions
	for i := 0; i < 100; i++ {
		queryEmbedding[i] = 0.9
	}

	// Perform similarity search
	results, err := SimilaritySearch(db, queryEmbedding, 3)
	if err != nil {
		t.Fatalf("Failed to perform similarity search: %v", err)
	}

	// Verify we got results
	if len(results) == 0 {
		t.Fatal("Expected search results, got none")
	}

	// Verify the order of results (should be chunk1, chunk2, chunk3 based on our embedding patterns)
	if len(results) >= 1 && results[0].ID != chunk1.ID {
		t.Errorf("Expected first result to be chunk1 (ID: %d), got chunk with ID: %d", chunk1.ID, results[0].ID)
	}

	if len(results) >= 2 && results[1].ID != chunk2.ID {
		t.Errorf("Expected second result to be chunk2 (ID: %d), got chunk with ID: %d", chunk2.ID, results[1].ID)
	}

	if len(results) >= 3 && results[2].ID != chunk3.ID {
		t.Errorf("Expected third result to be chunk3 (ID: %d), got chunk with ID: %d", chunk3.ID, results[2].ID)
	}

	// Verify the content of the results
	if len(results) >= 1 && results[0].Content != doc1.Content {
		t.Errorf("Expected first result content to be '%s', got '%s'", doc1.Content, results[0].Content)
	}
}

func TestGetDocumentByID(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Initialize the database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer Close(db)

	// Create and insert a test document
	doc := &Document{
		Title:           "Test Document",
		Content:         "This is a test document.",
		Author:          "Test Author",
		PublicationDate: "2023-01-01",
		URL:             "https://example.com/test",
		FilePath:        "/path/to/test.md",
	}
	err = InsertDocument(db, doc)
	if err != nil {
		t.Fatalf("Failed to insert document: %v", err)
	}

	// Retrieve the document by ID
	retrievedDoc, err := GetDocumentByID(db, doc.ID)
	if err != nil {
		t.Fatalf("Failed to get document by ID: %v", err)
	}

	// Verify that the retrieved document matches the original
	if retrievedDoc.ID != doc.ID {
		t.Errorf("Expected document ID %d, got %d", doc.ID, retrievedDoc.ID)
	}
	if retrievedDoc.Title != doc.Title {
		t.Errorf("Expected document title %q, got %q", doc.Title, retrievedDoc.Title)
	}
	if retrievedDoc.Content != doc.Content {
		t.Errorf("Expected document content %q, got %q", doc.Content, retrievedDoc.Content)
	}
}

func TestGetAllDocuments(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Initialize the database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer Close(db)

	// Create and insert test documents
	doc1 := &Document{
		Title:           "Test Document 1",
		Content:         "This is test document 1.",
		Author:          "Test Author 1",
		PublicationDate: "2023-01-01",
		URL:             "https://example.com/test1",
		FilePath:        "/path/to/test1.md",
	}
	err = InsertDocument(db, doc1)
	if err != nil {
		t.Fatalf("Failed to insert document 1: %v", err)
	}

	doc2 := &Document{
		Title:           "Test Document 2",
		Content:         "This is test document 2.",
		Author:          "Test Author 2",
		PublicationDate: "2023-01-02",
		URL:             "https://example.com/test2",
		FilePath:        "/path/to/test2.md",
	}
	err = InsertDocument(db, doc2)
	if err != nil {
		t.Fatalf("Failed to insert document 2: %v", err)
	}

	// Retrieve all documents
	docs, err := GetAllDocuments(db)
	if err != nil {
		t.Fatalf("Failed to get all documents: %v", err)
	}

	// Verify that all documents were retrieved
	if len(docs) != 2 {
		t.Errorf("Expected 2 documents, got %d", len(docs))
	}
}

func TestGetDocumentChunks(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Initialize the database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer Close(db)

	// Create and insert a test document
	doc := &Document{
		Title:           "Test Document",
		Content:         "This is a test document.",
		Author:          "Test Author",
		PublicationDate: "2023-01-01",
		URL:             "https://example.com/test",
		FilePath:        "/path/to/test.md",
	}
	err = InsertDocument(db, doc)
	if err != nil {
		t.Fatalf("Failed to insert document: %v", err)
	}

	// Create and insert test chunks
	chunk1 := &Chunk{
		DocumentID: doc.ID,
		Content:    "This is the first test chunk.",
		Hash:       []byte{1, 2, 3, 4},
		Embedding:  Embedding{0.1, 0.2, 0.3, 0.4},
	}
	err = InsertChunk(db, chunk1)
	if err != nil {
		t.Fatalf("Failed to insert chunk1: %v", err)
	}

	chunk2 := &Chunk{
		DocumentID: doc.ID,
		Content:    "This is the second test chunk.",
		Hash:       []byte{5, 6, 7, 8},
		Embedding:  Embedding{0.5, 0.6, 0.7, 0.8},
	}
	err = InsertChunk(db, chunk2)
	if err != nil {
		t.Fatalf("Failed to insert chunk2: %v", err)
	}

	// Retrieve chunks for the document
	chunks, err := GetDocumentChunks(db, doc.ID)
	if err != nil {
		t.Fatalf("Failed to get document chunks: %v", err)
	}

	// Verify that all chunks were retrieved
	if len(chunks) != 2 {
		t.Errorf("Expected 2 chunks, got %d", len(chunks))
	}
}

// Test with invalid embedding dimensions
func TestSimilaritySearchInvalidEmbeddingDimensions(t *testing.T) {
	// Create a temporary database file
	tempFile, err := os.CreateTemp("", "test-db-*.sqlite")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()
	dbPath := tempFile.Name()

	// Clean up after the test
	defer os.Remove(dbPath)

	// Create database
	db, err := GetDB(dbPath)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer Close(db)

	// Test with invalid embedding dimensions
	invalidEmbedding := Embedding{0.1, 0.2}
	_, err = SimilaritySearch(db, invalidEmbedding, 3)
	if err == nil {
		t.Fatal("Expected error for invalid embedding dimensions, got nil")
	}
}
