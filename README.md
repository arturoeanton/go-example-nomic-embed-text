# go-example-nomic-embed-text

This project demonstrates how to integrate text embeddings using **nomic-embed-text** and **granite-embedding** models with PostgreSQL and pgvector. You can perform similarity searches, text analysis, and more.

---

## Getting Started

### Install Models for Play

To get started, pull the required models using **Ollama**:

```bash
ollama pull nomic-embed-text
ollama run granite-embedding:30m
ollama run granite-embedding:278m
```

---

## Database Setup

### Table Creation

Create a table to store your texts and their embeddings:

```sql
CREATE TABLE IF NOT EXISTS documents (
     id SERIAL PRIMARY KEY,
     text TEXT NOT NULL,
     embedding VECTOR(768) -- 768 for nomic-embed-text or granite-embedding:30m
                           -- 1024 for granite-embedding:278m
);
```

### Index Creation

Speed up similarity searches by creating an index on the embedding column:

```sql
CREATE INDEX documents_embedding_idx ON documents USING ivfflat (embedding) WITH (lists = 100);
```

---

## Models Overview

### **nomic-embed-text**

A lightweight model for generating 768-dimensional embeddings. Ideal for text similarity and retrieval tasks.

```sql
CREATE TABLE IF NOT EXISTS documents (
     id SERIAL PRIMARY KEY,
     text TEXT NOT NULL,
     embedding VECTOR(768) -- 768 dimensions for nomic-embed-text
);
```

### **granite-embedding:30m**

A smaller, faster model for efficient embedding generation with 768 dimensions.

```sql
-- Add a column for granite-embedding:30m
ALTER TABLE documents ADD COLUMN embedding VECTOR(768);
```

### **granite-embedding:278m**

A multilingual model with 1024-dimensional embeddings for richer representation.

```sql
-- Add a column for granite-embedding:278m
ALTER TABLE documents ADD COLUMN embedding VECTOR(1024);
```

---

## How to Use

1. **Insert Texts with Embeddings**
   Use the provided Go script to insert texts and their embeddings into the database.

2. **Query Similar Texts**
   Perform similarity searches using SQL queries like:
   ```sql
   SELECT id, text, 1 - (embedding <=> '[QUERY_VECTOR]') AS similarity
   FROM documents
   ORDER BY similarity DESC
   LIMIT 5;
   ```

3. **Experiment with Different Models**
   Compare the performance and results of nomic-embed-text and granite-embedding models.

---

Feel free to contribute or explore further! 🎉
