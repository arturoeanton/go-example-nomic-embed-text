# go-example-nomic-embed-text


### install modeles for play :)

```bash
ollama pull nomic-embed-text
ollama run granite-embedding:30m
ollama run granite-embedding:278m
```

## nomic-embed-text

```sql
CREATE TABLE IF NOT EXISTS documents (
     id SERIAL PRIMARY KEY,
     text TEXT NOT NULL,
     embedding VECTOR(768) -- 768 para nomic-embed-text o  granite-embedding:278m 
                           -- 1024 para granite-embedding:30m
);
```

```sql
CREATE INDEX documents_embedding_idx ON documents USING ivfflat (embedding) WITH (lists = 100);
```

## granite-embedding:30m

```sql
-- Para 30M
ALTER TABLE documents ADD COLUMN embedding VECTOR(768);
```
## granite-embedding:278m

```sql
-- Para 278M
ALTER TABLE documents ADD COLUMN embedding VECTOR(1024);
```

 