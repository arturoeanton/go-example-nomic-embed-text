CREATE EXTENSION IF NOT EXISTS vector;
CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    text TEXT NOT NULL,
    embedding VECTOR(768) -- Ajusta el tamaño del vector según el modelo.
    -- El tamaño del vector en los embeddings (por ejemplo, VECTOR(768)) depende del modelo que genera dichos embeddings. Cada modelo tiene una arquitectura específica que determina el número de dimensiones en su espacio vectorial.
    -- Los modelos están diseñados para representar texto en un espacio de dimensiones fijas.
    -- Por ejemplo,
    -- - El modelo BERT base tiene 768 dimensiones.
    -- - GPT (nomic-embed-text): 768 dimensiones.
    -- - OpenAI text-embedding-ada-002: 1536 dimensiones.
    -- - LLaMA: puede variar, pero los embeddings de sus variantes suelen tener más de 1024 dimensiones.
    -- - Otros modelos pueden tener más o menos dimensiones.
    -- Para obtener el tamaño del vector de un modelo específico, consulta la documentación del modelo.
);

CREATE INDEX documents_embedding_idx ON documents USING ivfflat (embedding) WITH (lists = 100);

