package backend

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Document struct {
	Title           string
	Content         string
	Author          string
	PublicationDate string
	URL             string
	FilePath        string
	ID              int
	Hash            []byte
}

type Chunk struct {
	DocumentID int
	Content    string
	Embedding  Embedding
	Hash       []byte
	ID         int
}

type User struct {
	ID int
}

func ChunkDocument(doc *Document, embeddingClient *EmbeddingClient, user *User, db *DB) ([]Chunk, error) {
	if len(doc.Content) == 0 {
		return nil, fmt.Errorf("document content is empty")
	}
	
	// Calculate document hash if not already set
	if doc.Hash == nil {
		doc.Hash = MakeHash(doc.Content)
	}
	
	// Check if document has already been processed
	processed, err := DocumentHasBeenProcessed(db, doc.Hash)
	if err != nil {
		return nil, fmt.Errorf("error checking if document has been processed: %w", err)
	}
	
	// If document has been processed, return its existing chunks
	if processed {
		existingDoc, err := GetDocumentByHash(db, doc.Hash)
		if err != nil {
			return nil, fmt.Errorf("error retrieving existing document: %w", err)
		}
		
		chunks, err := GetDocumentChunks(db, existingDoc.ID)
		if err != nil {
			return nil, fmt.Errorf("error getting chunks for existing document: %w", err)
		}
		
		// Update document ID to match existing document
		doc.ID = existingDoc.ID
		return chunks, nil
	}
	
	// If document hasn't been processed, create new chunks
	chunks := []Chunk{}
	paragraphs := strings.Split(doc.Content, "\n\n")
	var chunkContents []string
	
	// Collect all paragraph chunks
	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if len(paragraph) == 0 {
			continue
		}
		
		chunks = append(chunks, Chunk{
			DocumentID: doc.ID,
			Content:    paragraph,
			Hash:       MakeHash(paragraph),
		})
		chunkContents = append(chunkContents, paragraph)
	}
	
	// Add the full document as a chunk
	chunks = append(chunks, Chunk{
		DocumentID: doc.ID,
		Content:    doc.Content,
		Hash:       MakeHash(doc.Content),
	})
	chunkContents = append(chunkContents, doc.Content)
	
	// Create embeddings for all chunks
	embeddingVectors, err := CreateEmbeddings(embeddingClient, chunkContents, user.ID)
	if err != nil {
		return nil, err
	}
	
	// Add embeddings to chunks
	for i, embeddingVector := range embeddingVectors {
		chunks[i].Embedding = embeddingVector
	}
	
	return chunks, nil
}

func HugoToDocument(filePath string) (Document, error) {
	theDocument := Document{FilePath: filePath}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return Document{}, err
	}

	lines := strings.Split(string(content), "\n")

	metadataStart := -1
	metadataEnd := -1

	for i, line := range lines {
		if strings.TrimSpace(line) == "---" {
			if metadataStart == -1 {
				metadataStart = i
			} else {
				metadataEnd = i
				break
			}
		}
	}

	if metadataStart != -1 && metadataEnd != -1 {
		for i := metadataStart + 1; i < metadataEnd; i++ {
			line := lines[i]
			if strings.HasPrefix(line, "title: ") {
				theDocument.Title = strings.TrimPrefix(line, "title: ")
			} else if strings.HasPrefix(line, "author: ") {
				theDocument.Author = strings.TrimPrefix(line, "author: ")
			} else if strings.HasPrefix(line, "date: ") {
				theDocument.PublicationDate = strings.TrimPrefix(line, "date: ")
			} else if strings.HasPrefix(line, "url: ") {
				theDocument.URL = strings.TrimPrefix(line, "url: ")
			}
		}

		contentLines := lines[metadataEnd+1:]
		theDocument.Content = strings.Join(contentLines, "\n")
	} else {
		theDocument.Content = string(content)
	}

	return theDocument, nil
}

func HugoDirectoryToDocuments(directory string, recursive bool) ([]Document, error) {
	documents := []Document{}
	files, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			if recursive {
				subDirDocuments, err := HugoDirectoryToDocuments(filepath.Join(directory, file.Name()), recursive)
				if err != nil {
					return nil, err
				}
				documents = append(documents, subDirDocuments...)
			}
			continue
		}
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		document, err := HugoToDocument(filepath.Join(directory, file.Name()))
		if len(document.Content) == 0 {
			continue
		}
		if err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}

	return documents, nil
}

func MakeHash(content string) []byte {
	hash := sha256.New()
	hash.Write([]byte(content))
	return hash.Sum(nil)
}

func DocumentsFromChunks(chunks []Chunk, database *DB) ([]Document, error) {
	documents := []Document{}
	seenDocIDs := make(map[int]bool)
	for _, chunk := range chunks {
		if seenDocIDs[chunk.DocumentID] {
			continue
		}
		newDoc, err := GetDocumentByID(database, chunk.DocumentID)
		if err != nil {
			return nil, fmt.Errorf("error getting document by ID: %w", err)
		}
		documents = append(documents, newDoc)
		seenDocIDs[chunk.DocumentID] = true
	}
	return documents, nil
}

// CalculateDocumentHash calculates the hash for a document if not already set
func CalculateDocumentHash(doc *Document) {
	if doc.Hash == nil {
		doc.Hash = MakeHash(doc.Content)
	}
}

// PrepareDocumentForProcessing calculates all necessary hashes and checks
// if the document already exists in the database
func PrepareDocumentForProcessing(doc *Document, db *DB) (bool, error) {
	CalculateDocumentHash(doc)
	
	// Check if document already exists
	exists, err := DocumentHashExists(db, doc.Hash)
	if err != nil {
		return false, fmt.Errorf("error checking document hash: %w", err)
	}
	
	if exists {
		// Document exists, get its ID
		existingDoc, err := GetDocumentByHash(db, doc.Hash)
		if err != nil {
			return true, fmt.Errorf("error retrieving existing document: %w", err)
		}
		
		// Update document ID to match existing one
		doc.ID = existingDoc.ID
		return true, nil
	}
	
	return false, nil
}

// ProcessDocumentBatch processes a batch of documents, efficiently skipping duplicates
func ProcessDocumentBatch(db *DB, docs []Document, embeddingClient *EmbeddingClient, user *User) (int, int, error) {
	totalDocuments := 0
	skippedDocuments := 0
	
	for i := range docs {
		// Calculate hash if not already set
		if docs[i].Hash == nil {
			docs[i].Hash = MakeHash(docs[i].Content)
		}
		
		// Check if document has already been processed
		processed, err := DocumentHasBeenProcessed(db, docs[i].Hash)
		if err != nil {
			return totalDocuments, skippedDocuments, fmt.Errorf("error checking if document has been processed: %w", err)
		}
		
		if processed {
			// Document has already been processed, skip it
			skippedDocuments++
			
			// Update the document ID to match the existing one
			existingDoc, err := GetDocumentByHash(db, docs[i].Hash)
			if err != nil {
				return totalDocuments, skippedDocuments, fmt.Errorf("error retrieving existing document: %w", err)
			}
			docs[i].ID = existingDoc.ID
			continue
		}
		
		// Insert the document
		err = InsertDocument(db, &docs[i])
		if err != nil {
			return totalDocuments, skippedDocuments, fmt.Errorf("error inserting document: %w", err)
		}
		
		// Process and insert chunks
		chunks, err := ChunkDocument(&docs[i], embeddingClient, user, db)
		if err != nil {
			return totalDocuments, skippedDocuments, fmt.Errorf("error chunking document: %w", err)
		}
		
		// Insert chunks
		for j := range chunks {
			err = InsertChunk(db, &chunks[j])
			if err != nil {
				return totalDocuments, skippedDocuments, fmt.Errorf("error inserting chunk: %w", err)
			}
		}
		
		totalDocuments++
	}
	
	return totalDocuments, skippedDocuments, nil
}
