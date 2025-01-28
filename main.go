package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq" // Driver para PostgreSQL
)

// Configuración de la conexión a la base de datos
const dbConnString = "postgres://user:password@localhost:5432/embeddings?sslmode=disable"

// Estructura para la respuesta de Ollama
type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

// Genera un embedding usando la API de Ollama
func generateEmbedding(prompt string) ([]float64, error) {
	url := "http://localhost:11434/api/embeddings"
	payload := map[string]interface{}{
		"model":  "nomic-embed-text",
		"prompt": prompt,
	}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error al llamar a la API de Ollama: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("error en la respuesta de Ollama: %s", string(body))
	}

	var embeddingResponse EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResponse); err != nil {
		return nil, fmt.Errorf("error al decodificar respuesta de Ollama: %v", err)
	}
	return embeddingResponse.Embedding, nil
}

// Convierte un slice de float64 a una cadena compatible con pgvector
func floatsToPgVector(floats []float64) string {
	var str strings.Builder
	str.WriteString("[")
	for i, v := range floats {
		str.WriteString(fmt.Sprintf("%f", v))
		if i < len(floats)-1 {
			str.WriteString(",")
		}
	}
	str.WriteString("]")
	return str.String()
}

// Convierte una lista de floats a una cadena separada por comas
func joinFloats(floats []float64) string {
	var s []string
	for _, v := range floats {
		s = append(s, fmt.Sprintf("%f", v))
	}
	return join(s, ",")
}

// Función auxiliar para unir cadenas con un separador
func join(elems []string, sep string) string {
	if len(elems) == 0 {
		return ""
	}
	return elems[0] + sep + join(elems[1:], sep)
}

// Inserta un texto y su embedding en la base de datos
func insertText(db *sql.DB, text string) error {
	embedding, err := generateEmbedding(text)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO documents (text, embedding) VALUES ($1, $2)
	`
	_, err = db.Exec(query, text, floatsToPgVector(embedding))
	return err
}

// Busca textos similares en la base de datos
func querySimilarTexts(db *sql.DB, query string) error {
	queryEmbedding, err := generateEmbedding(query)
	if err != nil {
		return err
	}

	sqlQuery := `
		SELECT id, text, 1 - (embedding <=> $1) AS similarity, embedding
		FROM documents
		WHERE (embedding <=> $1) <= 0.5
		ORDER BY similarity DESC
		--LIMIT 5
	`
	rows, err := db.Query(sqlQuery, floatsToPgVector(queryEmbedding))
	if err != nil {
		return err
	}
	defer rows.Close()

	fmt.Println("Resultados similares:")
	i := 0
	for rows.Next() {
		i++
		var id int
		var text string
		var similarity float64
		var embeddingRaw []byte // Extraer el embedding como bytes
		if err := rows.Scan(&id, &text, &similarity, &embeddingRaw); err != nil {
			return err
		}
		// Convertir el campo raw a []float64
		embedding, err := parseVector(embeddingRaw)
		if err != nil {
			return fmt.Errorf("error al convertir embedding: %v", err)
		}

		// Calcular la similitud coseno
		calculateInGoSimilarity := cosineSimilarity(queryEmbedding, embedding)

		fmt.Printf("%.2d - ID: %.4d | Similaridad: %.4f | Go Code Similaridad: %.4f  | Texto: %s\n", i, id, similarity, calculateInGoSimilarity, text)
	}

	return nil
}

func parseVector(raw []byte) ([]float64, error) {
	// Elimina los corchetes iniciales y finales
	data := strings.Trim(string(raw), "[]")

	// Divide la cadena en los valores individuales
	values := strings.Split(data, ",")

	// Convierte los valores a float64
	vector := make([]float64, len(values))
	for i, v := range values {
		val, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return nil, fmt.Errorf("error al parsear valor %q: %v", v, err)
		}
		vector[i] = val
	}

	return vector, nil
}

// Calcula la distancia coseno entre dos vectores
func cosineSimilarity(vec1, vec2 []float64) float64 {
	var dotProduct, normA, normB float64
	for i := 0; i < len(vec1); i++ {
		dotProduct += vec1[i] * vec2[i]
		normA += vec1[i] * vec1[i]
		normB += vec2[i] * vec2[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
}

func main() {
	insertFlag := flag.String("insert", "", "Texto a insertar en la base de datos")
	queryFlag := flag.String("query", "", "Texto para buscar por similitud")
	flag.Parse()

	// Conectar a la base de datos PostgreSQL
	db, err := sql.Open("postgres", dbConnString)
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	defer db.Close()

	// Verifica las banderas
	if *insertFlag != "" {
		err = insertText(db, *insertFlag)
		if err != nil {
			log.Fatalf("Error al insertar texto: %v", err)
		}
		fmt.Println("Texto insertado correctamente.")
	} else if *queryFlag != "" {
		err = querySimilarTexts(db, *queryFlag)
		if err != nil {
			log.Fatalf("Error al buscar textos: %v", err)
		}
	} else {
		fmt.Println("Usa -insert para insertar texto o -query para buscar por similitud.")
		os.Exit(1)
	}
}
