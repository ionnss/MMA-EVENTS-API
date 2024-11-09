// cmd/main.go
package main

import (
	"fmt"
	"log"
	"mma-events-api/internal/crawler"
	"mma-events-api/internal/storage"
	"os"
)

func main() {
	fmt.Println("Iniciando o MMA Events API")

	// Caminho do arquivo do banco de dados
	dbPath := "internal/data/mma_events.db"

	// Exclui o banco de dados se ele já existe
	if _, err := os.Stat(dbPath); err == nil {
		if err := os.Remove(dbPath); err != nil {
			log.Fatalf("Erro ao excluir o banco de dados existente: %v\n", err)
		}
		fmt.Println("Banco de dados SQLITE reiniciado com sucesso!")
	}

	// Inicializa o banco de dados
	db, err := storage.InitDB()
	if err != nil {
		fmt.Printf("Erro ao inicializar o banco de dados: %v\n", err)
		return
	}
	defer db.Close()

	// Carrega as organizações do arquivo mma-orgs.txt
	organizations, err := crawler.LoadOrganizationsFromFile("internal/data/mma-orgs.txt")
	if err != nil {
		fmt.Printf("Erro ao carregar organizações: %v\n", err)
		return
	}

	// Insere apenas novas organizações no banco de dados
	if err := crawler.InsertNewOrganizations(db, organizations); err != nil {
		fmt.Printf("Erro ao inserir novas organizações: %v\n", err)
		return
	}

	fmt.Println("Organizações carregadas e inseridas no banco de dados com sucesso!")

	// Executa o crawler de eventos UFC e insere os eventos no banco de dados
	events, err := crawler.ScrapeUFCEvents(db)
	if err != nil {
		fmt.Printf("Erro no scraper da UFC: %v\n", err)
		return
	}

	// Salva os eventos no banco de dados
	for _, event := range events {
		err := storage.InsertEvent(db, event.Org, event.Name, event.DateMainCard, event.DatePrelimsCard, event.Location, event.City, event.State, event.Country)
		if err != nil {
			fmt.Printf("Erro ao salvar evento %s: %v\n", event.Name, err)
		}
	}

	fmt.Println("Eventos carregados e inseridos no banco de dados com sucesso!")
}
