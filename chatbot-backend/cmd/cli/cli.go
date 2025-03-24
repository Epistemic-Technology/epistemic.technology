package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Epistemic-Technology/epistemic.technology/chatbot-backend/internal/backend"
	"github.com/joho/godotenv"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	err := godotenv.Load("../../.env")
	if err != nil {
		fmt.Printf("Error loading .env file: %v\n", err)
	}

	switch os.Args[1] {
	case "embed-hugo-directory":
		embedHugoDirectory(os.Args[2:])
	case "embed-hugo-file":
		embedHugoFile(os.Args[2:])
	case "list-documents":
		listDocuments(os.Args[2:])
	case "get-document-details":
		getDocumentDetails(os.Args[2:])
	case "get-db-stats":
		getDBStats(os.Args[2:])
	case "list-chunks":
		listChunks(os.Args[2:])
	case "list-chunks-for-document":
		listChunksForDocument(os.Args[2:])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  cli embed-hugo-directory <directory> [--recursive] [--db=<path>]")
	fmt.Println("  cli embed-hugo-file <file> [--db=<path>]")
	fmt.Println("  cli list-documents [--db=<path>]")
	fmt.Println("  cli get-document-details <document_id> [--db=<path>]")
	fmt.Println("  cli get-db-stats [--db=<path>]")
	fmt.Println("  cli list-chunks [--db=<path>]")
	fmt.Println("  cli list-chunks-for-document <document_id> [--db=<path>]")
}

func parseArgs(args []string) ([]string, map[string]string) {
	positionalArgs := []string{}
	namedArgs := make(map[string]string)

	for _, arg := range args {
		if strings.HasPrefix(arg, "--") {
			parts := strings.SplitN(arg[2:], "=", 2)
			if len(parts) == 2 {
				namedArgs[parts[0]] = parts[1]
			} else {
				namedArgs[parts[0]] = "true"
			}
		} else {
			positionalArgs = append(positionalArgs, arg)
		}
	}

	return positionalArgs, namedArgs
}

func getDBPath(args []string) string {
	_, namedArgs := parseArgs(args)
	dbPath := namedArgs["db"]
	if dbPath == "" {
		dbPath = os.Getenv("DATABASE_PATH")
		if dbPath == "" {
			log.Fatal("Error: No database path provided. Use --db flag or set DATABASE_PATH environment variable")
		}
	}
	return dbPath
}

func embedHugoDirectory(args []string) {
	positionalArgs, namedArgs := parseArgs(args)
	if len(positionalArgs) < 1 {
		fmt.Println("Error: Missing directory path")
		printUsage()
		os.Exit(1)
	}

	directory := positionalArgs[0]
	recursive := namedArgs["recursive"] == "true"
	dbPath := getDBPath(args)

	if !fileExists(directory) {
		log.Fatalf("Error: Directory %s does not exist", directory)
	}

	fmt.Printf("Embedding Hugo directory: %s (recursive: %v)\n", directory, recursive)
	fmt.Printf("Using database: %s\n", dbPath)

	// Connect to database
	database, err := backend.GetDB(dbPath)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer backend.Close(database)

	embeddingClient, err := backend.NewEmbeddingClient()
	if err != nil {
		log.Fatalf("Error creating embeddings client: %v", err)
	}

	// Process directory
	docs, err := backend.HugoDirectoryToDocuments(directory, recursive)
	if err != nil {
		log.Fatalf("Error processing directory: %v", err)
	}

	fmt.Printf("Found %d documents\n", len(docs))

	// Create a user for embedding generation
	user := &backend.User{ID: 1}

	for i, doc := range docs {
		fmt.Printf("Processing document %d/%d: %s\n", i+1, len(docs), doc.Title)

		if err := backend.InsertDocument(database, &doc); err != nil {
			log.Fatalf("Error inserting document: %v", err)
		}

		chunks, err := backend.ChunkDocument(&doc, embeddingClient, user, database)
		if err != nil {
			log.Fatalf("Error creating chunks: %v", err)
		}

		fmt.Printf("Created %d chunks\n", len(chunks))

		for j, chunk := range chunks {
			chunk.DocumentID = doc.ID
			if err := backend.InsertChunk(database, &chunk); err != nil {
				log.Fatalf("Error inserting chunk %d: %v", j, err)
			}
		}
	}

	fmt.Println("Done!")
}

func embedHugoFile(args []string) {
	positionalArgs, _ := parseArgs(args)
	if len(positionalArgs) < 1 {
		fmt.Println("Error: Missing file path")
		printUsage()
		os.Exit(1)
	}

	filePath := positionalArgs[0]
	dbPath := getDBPath(args)

	if !fileExists(filePath) {
		log.Fatalf("Error: File %s does not exist", filePath)
	}

	fmt.Printf("Embedding Hugo file: %s\n", filePath)
	fmt.Printf("Using database: %s\n", dbPath)

	// Connect to database
	database, err := backend.GetDB(dbPath)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer backend.Close(database)

	embeddingClient, err := backend.NewEmbeddingClient()
	if err != nil {
		log.Fatalf("Error creating embeddings client: %v", err)
	}

	// Process Hugo file
	doc, err := backend.HugoToDocument(filePath)
	if err != nil {
		log.Fatalf("Error processing Hugo file: %v", err)
	}

	if err := backend.InsertDocument(database, &doc); err != nil {
		log.Fatalf("Error inserting document: %v", err)
	}

	// Create a user for embedding generation
	user := &backend.User{ID: 1}

	chunks, err := backend.ChunkDocument(&doc, embeddingClient, user, database)
	if err != nil {
		log.Fatalf("Error creating chunks: %v", err)
	}

	fmt.Printf("Created %d chunks\n", len(chunks))

	for j, chunk := range chunks {
		chunk.DocumentID = doc.ID
		if err := backend.InsertChunk(database, &chunk); err != nil {
			log.Fatalf("Error inserting chunk %d: %v", j, err)
		}
	}

	fmt.Println("Done!")
}

func listDocuments(args []string) {
	dbPath := getDBPath(args)

	// Connect to database
	database, err := backend.GetDB(dbPath)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer backend.Close(database)

	documents, err := backend.GetAllDocuments(database)
	if err != nil {
		log.Fatalf("Error getting documents: %v", err)
	}

	fmt.Printf("Found %d documents:\n", len(documents))
	for _, doc := range documents {
		fmt.Printf("ID: %d, Title: %s\n", doc.ID, doc.Title)
		if doc.Author != "" {
			fmt.Printf("  Author: %s\n", doc.Author)
		}
		if doc.PublicationDate != "" {
			fmt.Printf("  Date: %s\n", doc.PublicationDate)
		}
		if doc.URL != "" {
			fmt.Printf("  URL: %s\n", doc.URL)
		}
		if doc.FilePath != "" {
			fmt.Printf("  File: %s\n", doc.FilePath)
		}
		fmt.Println()
	}
}

func listChunks(args []string) {
	dbPath := getDBPath(args)

	// Connect to database
	database, err := backend.GetDB(dbPath)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer backend.Close(database)

	chunks, err := backend.GetAllChunks(database)
	if err != nil {
		log.Fatalf("Error getting chunks: %v", err)
	}

	fmt.Printf("Found %d chunks:\n", len(chunks))
	for i, chunk := range chunks {
		// Limit the number of chunks displayed to avoid overwhelming output
		if i >= 10 {
			fmt.Printf("... and %d more chunks\n", len(chunks)-10)
			break
		}

		fmt.Printf("ID: %d, Document ID: %d\n", chunk.ID, chunk.DocumentID)
		// Truncate content if it's too long
		content := chunk.Content
		if len(content) > 100 {
			content = content[:97] + "..."
		}
		fmt.Printf("  Content: %s\n", content)
		fmt.Printf("  Hash: %x\n", chunk.Hash)
		fmt.Println()
	}
}

func listChunksForDocument(args []string) {
	positionalArgs, _ := parseArgs(args)
	if len(positionalArgs) < 1 {
		fmt.Println("Error: Missing document ID")
		printUsage()
		os.Exit(1)
	}

	documentID, err := strconv.Atoi(positionalArgs[0])
	if err != nil {
		log.Fatalf("Error: Invalid document ID: %v", err)
	}

	dbPath := getDBPath(args)

	// Connect to database
	database, err := backend.GetDB(dbPath)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer backend.Close(database)

	chunks, err := backend.GetAllChunksWithDocumentID(database, documentID)
	if err != nil {
		log.Fatalf("Error getting chunks: %v", err)
	}

	fmt.Printf("Found %d chunks for document ID %d:\n", len(chunks), documentID)
	for i, chunk := range chunks {
		// Limit the number of chunks displayed to avoid overwhelming output
		if i >= 10 {
			fmt.Printf("... and %d more chunks\n", len(chunks)-10)
			break
		}

		fmt.Printf("ID: %d\n", chunk.ID)
		// Truncate content if it's too long
		content := chunk.Content
		if len(content) > 100 {
			content = content[:97] + "..."
		}
		fmt.Printf("  Content: %s\n", content)
		fmt.Printf("  Hash: %x\n", chunk.Hash)
		fmt.Println()
	}
}

func getDBStats(args []string) {
	dbPath := getDBPath(args)

	// Connect to database
	database, err := backend.GetDB(dbPath)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer backend.Close(database)

	documents, err := backend.GetAllDocuments(database)
	if err != nil {
		log.Fatalf("Error getting documents: %v", err)
	}
	fmt.Printf("Found %d documents\n", len(documents))

	chunks, err := backend.GetAllChunks(database)
	if err != nil {
		log.Fatalf("Error getting chunks: %v", err)
	}
	fmt.Printf("Found %d chunks\n", len(chunks))
}

func getDocumentDetails(args []string) {
	positionalArgs, _ := parseArgs(args)
	if len(positionalArgs) < 1 {
		fmt.Println("Error: Missing document ID")
		printUsage()
		os.Exit(1)
	}

	documentID, err := strconv.Atoi(positionalArgs[0])
	if err != nil {
		log.Fatalf("Error: Invalid document ID: %v", err)
	}

	dbPath := getDBPath(args)

	// Connect to database
	database, err := backend.GetDB(dbPath)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer backend.Close(database)

	doc, err := backend.GetDocumentByID(database, documentID)
	if err != nil {
		log.Fatalf("Error getting document: %v", err)
	}

	chunks, err := backend.GetDocumentChunks(database, documentID)
	if err != nil {
		log.Fatalf("Error getting document chunks: %v", err)
	}

	fmt.Printf("Document ID: %d\n", doc.ID)
	fmt.Printf("Title: %s\n", doc.Title)
	if doc.Author != "" {
		fmt.Printf("Author: %s\n", doc.Author)
	}
	if doc.PublicationDate != "" {
		fmt.Printf("Date: %s\n", doc.PublicationDate)
	}
	if doc.URL != "" {
		fmt.Printf("URL: %s\n", doc.URL)
	}
	if doc.FilePath != "" {
		fmt.Printf("File: %s\n", doc.FilePath)
	}
	fmt.Printf("Content length: %d characters\n", len(doc.Content))
	fmt.Printf("Number of chunks: %d\n", len(chunks))

	// Print a preview of the content
	contentPreview := doc.Content
	if len(contentPreview) > 200 {
		contentPreview = contentPreview[:197] + "..."
	}
	fmt.Printf("\nContent preview:\n%s\n", contentPreview)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
