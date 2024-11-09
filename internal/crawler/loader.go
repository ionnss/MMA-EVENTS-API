// internal/crawler/loader.go
package crawler

import (
	"bufio"
	"database/sql"
	"fmt"
	"mma-events-api/internal/models"
	"os"
	"strings"
	//"mma-events-api/internal/storage"
)

// LoadOrganizationsFromFile lê o arquivo e retorna uma lista de organizações
func LoadOrganizationsFromFile(filePath string) ([]models.Organization, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir o arquivo: %w", err)
	}
	defer file.Close()

	var organizations []models.Organization
	var orgName, orgURL, orgEventURL string

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "URL:") {
			orgURL = strings.TrimSpace(line[4:])
		} else if strings.HasPrefix(line, "EventURL:") {
			orgEventURL = strings.TrimSpace(line[9:])
			organizations = append(organizations, models.Organization{
				Name:     orgName,
				URL:      orgURL,
				EventURL: orgEventURL,
			})
		} else if line != "" {
			orgName = line
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo: %w", err)
	}

	return organizations, nil
}

// InsertNewOrganizations insere somente as organizações que ainda não estão no banco
func InsertNewOrganizations(db *sql.DB, organizations []models.Organization) error {
	for _, org := range organizations {
		// Verifica se a organização já existe
		var exists bool
		query := `SELECT EXISTS(SELECT 1 FROM organizations WHERE name = ?);`
		if err := db.QueryRow(query, org.Name).Scan(&exists); err != nil {
			return fmt.Errorf("erro ao verificar organização %s: %w", org.Name, err)
		}

		// Insere apenas se a organização não existir
		if !exists {
			insertQuery := `INSERT INTO organizations (name, url, eventurl) VALUES (?, ?, ?);`
			if _, err := db.Exec(insertQuery, org.Name, org.URL, org.EventURL); err != nil {
				return fmt.Errorf("erro ao inserir organização %s: %w", org.Name, err)
			}
			fmt.Printf("Organização inserida: %s\n", org.Name)
		} else {
			fmt.Printf("Organização já existe: %s\n", org.Name)
		}
	}
	return nil
}
