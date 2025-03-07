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

func ChunkDocument(doc *Document, embeddingClient *EmbeddingClient, user *User) ([]Chunk, error) {
	if len(doc.Content) == 0 {
		return nil, fmt.Errorf("document content is empty")
	}
	chunks := []Chunk{}
	paragraphs := strings.Split(doc.Content, "\n\n")
	var chunkContents []string

	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if len(paragraph) == 0 {
			continue
		}
		chunks = append(chunks, Chunk{
			DocumentID: doc.ID,
			Content:    paragraph,
			Hash:       MakeHash(paragraph),
			ID:         len(chunks),
		})
		chunkContents = append(chunkContents, paragraph)
	}

	chunks = append(chunks, Chunk{
		DocumentID: doc.ID,
		Content:    doc.Content,
		Hash:       MakeHash(doc.Content),
		ID:         len(chunks),
	})
	chunkContents = append(chunkContents, doc.Content)

	embeddingVectors, err := CreateEmbeddings(embeddingClient, chunkContents, user.ID)
	if err != nil {
		return nil, err
	}

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
	for _, chunk := range chunks {
		newDoc, err := GetDocumentByID(database, chunk.DocumentID)
		if err != nil {
			return nil, fmt.Errorf("error getting document by ID: %w", err)
		}
		documents = append(documents, newDoc)
	}
	return documents, nil
}
