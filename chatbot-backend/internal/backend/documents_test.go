package backend

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func TestChunkWithEmbeddings(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		t.Fatalf("Error loading .env file: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test: OPENAI_API_KEY not set in environment")
	}

	embeddingClient, err := NewEmbeddingClient()
	if err != nil {
		t.Fatalf("Failed to create embedding client: %v", err)
	}

	doc := &Document{
		ID:       1,
		Content:  "This is a test chunk for embedding generation.\n\nThis is another test chunk with different content.",
		Title:    "Test Document",
		FilePath: "/path/to/test/document.md",
	}

	user := &User{ID: 123}

	chunks, err := ChunkDocument(doc, embeddingClient, user)
	if err != nil {
		t.Fatalf("Failed to chunk document: %v", err)
	}

	// Verify that chunks were created
	if len(chunks) == 0 {
		t.Fatal("Expected chunks to be created, got none")
	}

	// Verify that the chunks have the correct document ID
	for i, chunk := range chunks {
		if chunk.DocumentID != doc.ID {
			t.Errorf("Chunk %d has incorrect document ID: expected %d, got %d", i, doc.ID, chunk.DocumentID)
		}
	}

	// Verify that the chunks have content
	for i, chunk := range chunks {
		if chunk.Content == "" {
			t.Errorf("Chunk %d has empty content", i)
		}
	}

	// Verify that the chunks have embeddings
	for i, chunk := range chunks {
		if len(chunk.Embedding) == 0 {
			t.Errorf("Chunk %d has no embedding", i)
		}
	}

	// Verify that the chunks have hashes
	for i, chunk := range chunks {
		if len(chunk.Hash) == 0 {
			t.Errorf("Chunk %d has no hash", i)
		}
	}
}

func TestHugoDirectoryToDocuments(t *testing.T) {
	// Test with the existing test-docs directory
	docs, err := HugoDirectoryToDocuments("test-docs", false)
	if err != nil {
		t.Fatalf("Failed to load documents: %v", err)
	}

	// Verify we have the expected number of documents
	// The test-docs directory contains markdown files
	if len(docs) == 0 {
		t.Errorf("Expected documents to be loaded from test-docs, but got none")
	}

	// Check that each document has content
	for i, doc := range docs {
		if doc.Content == "" {
			t.Errorf("Document %d has empty content", i)
		}

		// Note: We don't check title or publication date as some test files might not have them
	}

	// Note: We can't test with a non-existent directory because the current implementation
	// calls log.Fatal which would terminate the test. This would need a code change to
	// make it testable.
}

func TestHugoDirectoryToDocumentsWithMixedContent(t *testing.T) {
	// Create a temporary directory with mixed content
	tempDir, err := os.MkdirTemp("", "hugo-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a markdown file with frontmatter
	mdContent := `---
title: Test Document
date: 2024-03-20
author: Test Author
url: https://example.com/test
---

# Test Content

This is a test document.`

	err = os.WriteFile(filepath.Join(tempDir, "test.md"), []byte(mdContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test markdown file: %v", err)
	}

	// Create a non-markdown file
	err = os.WriteFile(filepath.Join(tempDir, "test.txt"), []byte("This is not a markdown file"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test text file: %v", err)
	}

	// Create a subdirectory
	err = os.Mkdir(filepath.Join(tempDir, "subdir"), 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Test the function
	docs, err := HugoDirectoryToDocuments(tempDir, false)
	if err != nil {
		t.Fatalf("Failed to load documents: %v", err)
	}

	// Should only find one document (the markdown file)
	if len(docs) != 1 {
		t.Errorf("Expected 1 document, got %d", len(docs))
	}

	if len(docs) > 0 {
		doc := docs[0]
		if doc.Title != "Test Document" {
			t.Errorf("Expected title 'Test Document', got '%s'", doc.Title)
		}
		if doc.Author != "Test Author" {
			t.Errorf("Expected author 'Test Author', got '%s'", doc.Author)
		}
		if doc.PublicationDate != "2024-03-20" {
			t.Errorf("Expected date '2024-03-20', got '%s'", doc.PublicationDate)
		}
		if doc.URL != "https://example.com/test" {
			t.Errorf("Expected URL 'https://example.com/test', got '%s'", doc.URL)
		}
		if !strings.Contains(doc.Content, "# Test Content") {
			t.Errorf("Expected content to contain '# Test Content', got '%s'", doc.Content)
		}
		expectedFilePath := filepath.Join(tempDir, "test.md")
		if doc.FilePath != expectedFilePath {
			t.Errorf("Expected FilePath to be '%s', got '%s'", expectedFilePath, doc.FilePath)
		}
	}
}

func TestHugoToDocument(t *testing.T) {
	// Create a temporary markdown file with frontmatter
	tempFile, err := os.CreateTemp("", "hugo-test-*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	mdContent := `---
title: Test Hugo Document
author: Test Author
date: 2024-03-20
url: https://example.com/test-hugo
---

# Test Hugo Content

This is a test document with Hugo frontmatter.`

	err = os.WriteFile(tempFile.Name(), []byte(mdContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test markdown file: %v", err)
	}

	// Test the function
	doc, err := HugoToDocument(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to load document: %v", err)
	}

	// Verify the document fields
	if doc.Title != "Test Hugo Document" {
		t.Errorf("Expected title 'Test Hugo Document', got '%s'", doc.Title)
	}
	if doc.Author != "Test Author" {
		t.Errorf("Expected author 'Test Author', got '%s'", doc.Author)
	}
	if doc.PublicationDate != "2024-03-20" {
		t.Errorf("Expected date '2024-03-20', got '%s'", doc.PublicationDate)
	}
	if doc.URL != "https://example.com/test-hugo" {
		t.Errorf("Expected URL 'https://example.com/test-hugo', got '%s'", doc.URL)
	}
	if !strings.Contains(doc.Content, "# Test Hugo Content") {
		t.Errorf("Expected content to contain '# Test Hugo Content', got '%s'", doc.Content)
	}
	if doc.FilePath != tempFile.Name() {
		t.Errorf("Expected FilePath to be '%s', got '%s'", tempFile.Name(), doc.FilePath)
	}

	// Test with a file that has no frontmatter
	noFrontmatterFile, err := os.CreateTemp("", "no-frontmatter-*.md")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(noFrontmatterFile.Name())

	noFrontmatterContent := `# No Frontmatter Document

This document has no frontmatter section.`

	err = os.WriteFile(noFrontmatterFile.Name(), []byte(noFrontmatterContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test markdown file: %v", err)
	}

	// Test the function with no frontmatter
	noFrontmatterDoc, err := HugoToDocument(noFrontmatterFile.Name())
	if err != nil {
		t.Fatalf("Failed to load document: %v", err)
	}

	// Verify the document fields
	if noFrontmatterDoc.Title != "" {
		t.Errorf("Expected empty title for no frontmatter, got '%s'", noFrontmatterDoc.Title)
	}
	if noFrontmatterDoc.Author != "" {
		t.Errorf("Expected empty author for no frontmatter, got '%s'", noFrontmatterDoc.Author)
	}
	if noFrontmatterDoc.PublicationDate != "" {
		t.Errorf("Expected empty date for no frontmatter, got '%s'", noFrontmatterDoc.PublicationDate)
	}
	if noFrontmatterDoc.URL != "" {
		t.Errorf("Expected empty URL for no frontmatter, got '%s'", noFrontmatterDoc.URL)
	}
	if noFrontmatterDoc.Content != noFrontmatterContent {
		t.Errorf("Expected content to be the entire file for no frontmatter, got '%s'", noFrontmatterDoc.Content)
	}
	if noFrontmatterDoc.FilePath != noFrontmatterFile.Name() {
		t.Errorf("Expected FilePath to be '%s', got '%s'", noFrontmatterFile.Name(), noFrontmatterDoc.FilePath)
	}
}

func TestHugoDirectoryToDocumentsWithRealFiles(t *testing.T) {
	// Test with the actual files in test-docs directory
	docs, err := HugoDirectoryToDocuments("test-docs", false)
	if err != nil {
		t.Fatalf("Failed to load documents from test-docs: %v", err)
	}

	// Verify we have documents
	if len(docs) == 0 {
		t.Fatalf("Expected documents to be loaded from test-docs, but got none")
	}

	// Count the number of markdown files in the directory to verify all are loaded
	files, err := os.ReadDir("test-docs")
	if err != nil {
		t.Fatalf("Failed to read test-docs directory: %v", err)
	}

	mdFileCount := 0
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".md") {
			mdFileCount++
		}
	}

	if len(docs) != mdFileCount {
		t.Errorf("Expected %d documents to be loaded, but got %d", mdFileCount, len(docs))
	}

	// Create a map of documents by title for easier lookup
	docsByTitle := make(map[string]Document)
	for _, doc := range docs {
		docsByTitle[doc.Title] = doc
	}

	// Test specific document: Introduction to Transformers
	transformersDoc, exists := docsByTitle["Understanding Transformer Architecture: The Building Blocks of Modern AI"]
	if !exists {
		t.Fatalf("Expected to find document about transformers")
	}

	// Verify metadata extraction
	if transformersDoc.PublicationDate != "2024-03-15" {
		t.Errorf("Expected publication date '2024-03-15', got '%s'", transformersDoc.PublicationDate)
	}

	// Verify content extraction
	if !strings.Contains(transformersDoc.Content, "The transformer architecture has revolutionized") {
		t.Errorf("Expected content to contain introduction text about transformers")
	}

	if !strings.Contains(transformersDoc.Content, "Self-attention mechanisms") {
		t.Errorf("Expected content to contain information about self-attention mechanisms")
	}

	// Test another specific document: Ethical AI Development
	ethicalAIDoc, exists := docsByTitle["Ethical Considerations in AI Development"]
	if exists {
		// Verify content extraction for this document
		if !strings.Contains(ethicalAIDoc.Content, "ethical") && !strings.Contains(ethicalAIDoc.Content, "Ethical") {
			t.Errorf("Expected content to contain text about ethics")
		}
	}

	// Test document with a specific filename
	var foundDoc001 bool
	for _, doc := range docs {
		// Check if this document corresponds to 001-introduction-to-transformers.md
		if strings.Contains(doc.Content, "The transformer architecture has revolutionized") {
			foundDoc001 = true

			// Verify the document has the expected sections
			if !strings.Contains(doc.Content, "## The Basic Building Blocks") {
				t.Errorf("Expected document to have 'The Basic Building Blocks' section")
			}

			if !strings.Contains(doc.Content, "## Why Transformers Matter") {
				t.Errorf("Expected document to have 'Why Transformers Matter' section")
			}

			if !strings.Contains(doc.Content, "## Looking Forward") {
				t.Errorf("Expected document to have 'Looking Forward' section")
			}

			// Verify the document has the expected content structure
			if !strings.Contains(doc.Content, "1. Self-attention mechanisms") {
				t.Errorf("Expected document to list self-attention mechanisms")
			}

			break
		}
	}

	if !foundDoc001 {
		t.Errorf("Could not find document 001-introduction-to-transformers.md by content")
	}

	// Test document structure - verify paragraphs are present
	if len(docs) > 0 {
		// Use the first document for content structure test
		testDoc := docs[0]

		// Check that the document has paragraphs (separated by double newlines)
		paragraphs := strings.Split(testDoc.Content, "\n\n")
		if len(paragraphs) <= 1 {
			t.Errorf("Expected document to have multiple paragraphs, but found only %d", len(paragraphs))
		}

		// Check that the document has headings (lines starting with #)
		hasHeadings := false
		lines := strings.Split(testDoc.Content, "\n")
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "#") {
				hasHeadings = true
				break
			}
		}
		if !hasHeadings {
			t.Errorf("Expected document to have markdown headings")
		}
	}

	// Test that frontmatter is properly removed from content
	for _, doc := range docs {
		// Frontmatter should be removed from content
		if strings.Contains(doc.Content, "---\ntitle:") {
			t.Errorf("Document content should not contain frontmatter")
		}

		// Content should not start with --- markers
		if strings.TrimSpace(doc.Content)[:3] == "---" {
			t.Errorf("Document content should not start with frontmatter markers")
		}
	}
}
