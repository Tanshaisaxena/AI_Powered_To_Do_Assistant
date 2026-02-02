# 🧠 RAG-powered To-Do Assistant (Go + Gin + Ollama)

A Retrieval-Augmented Generation (RAG) based To-Do Assistant built in Go using Gin and Ollama.
This project evolves step-by-step from a basic TODO API into a fully functional semantic AI assistant.

---

## 🚀 Features

- CRUD operations for tasks
- Keyword-based search (baseline)
- Semantic search using embeddings
- Local embeddings with Ollama (no API keys required)
- Cosine similarity for relevance ranking
- Top-K and similarity threshold filtering
- LLM-powered natural language answers
- Hallucination-safe prompt design
- Fully local RAG pipeline

---

## 🏗️ Architecture Overview

User Question
  → Question Embedding
  → Similarity Search (Task Embeddings)
  → Top-K + Threshold Filtering
  → Prompt Construction
  → LLM Generation (Ollama)
  → Final Answer

---

## 🧩 Tech Stack

- Language: Go
- Framework: Gin
- Embeddings: Ollama (nomic-embed-text)
- LLM: Ollama (mistral / llama3)
- Similarity: Cosine Similarity
- Storage: In-memory (learning phase)

---

## 📁 Project Structure

handlers/    → API handlers (CRUD, AskRAG)
helpers/     → Embeddings, similarity, prompt & LLM helpers
models/      → Data models
constants/   → Log constants
main.go      → App entry point

---

## 🔧 Setup Instructions

### Prerequisites
- Go 1.21+
- Ollama installed

### Install models
ollama pull nomic-embed-text
ollama pull mistral

### Run
go run main.go

Server runs on http://localhost:8080

---

## 📌 API Endpoints

POST /task        → Create task (embedding generated automatically)
GET /task         → List tasks
DELETE /task/:id  → Delete task
POST /ask         → Ask questions using token count
POST /ask-rag     → Ask questions using RAG

---

## 🧠 RAG Flow

1. Task creation → embedding stored
2. User question → embedding generated
3. Similarity computed against all tasks
4. Results sorted by relevance
5. Threshold + Top-K applied
6. Prompt built with grounded context
7. LLM generates final answer

---

## 🔐 Hallucination Safety

The LLM is restricted to answering only from retrieved tasks.
If information is missing, it responds with "I don't know".

---

## 🛠️ Future Enhancements

- pgvector + Postgres
- Streaming responses
- Citations per task
- Conversation memory
- UI frontend

---

## 🏆 Key Takeaway

This project demonstrates how real-world RAG systems are built:
retrieval first, generation second.

---

## 📜 License

MIT
