package crawler

import (
	"database/sql"
	"fmt"
	"mma-events-api/internal/models"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Consulta URL de eventos no banco de dados das organizações
func GetOrganizationEventURL(db *sql.DB, orgName string) (string, error) {
	var eventURL string
	query := `SELECT eventurl FROM organizations WHERE name = ?`
	err := db.QueryRow(query, orgName).Scan(&eventURL)
	if err != nil {
		return "", fmt.Errorf("erro ao obter URL de eventos para %s: %w", orgName, err)
	}
	return eventURL, nil
}

// ScrapeUFCEvents coleta dados dos próximos eventos da UFC e os retorna em uma lista
func ScrapeUFCEvents(db *sql.DB) ([]models.Event, error) {
	// Obtém URL do banco de dados
	url, err := GetOrganizationEventURL(db, "UFC (Ultimate Fighting Championship)")
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar a página de eventos UFC: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro: status %d ao acessar %s", resp.StatusCode, url)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler HTML da página: %w", err)
	}

	var events []models.Event

	// Percorre cada evento contido em um <article class="c-card-event--result">
	doc.Find("article.c-card-event--result").Each(func(i int, event *goquery.Selection) {
		// Extrai o nome do evento, removendo espaços extras
		eventOrg := "UFC"
		eventName := strings.TrimSpace(event.Find("h3.c-card-event--result__headline a").Text())

		// Extrai a data e hora do evento principal e do card preliminar (usando atributos de datas)
		eventDateMain := strings.TrimSpace(event.Find(".c-card-event--result__date").AttrOr("data-main-card", "Data não disponível"))
		eventDatePrelims := strings.TrimSpace(event.Find(".c-card-event--result__date").AttrOr("data-prelims-card", "Data não disponível"))

		// Extrai a localização e limpa espaços
		eventLocation := strings.TrimSpace(event.Find(".field--name-taxonomy-term-title h5").Text())
		eventCity := strings.TrimSpace(event.Find(".address .locality").Text())
		eventState := strings.TrimSpace(event.Find(".address .administrative-area").Text())
		eventCountry := strings.TrimSpace(event.Find(".address .country").Text())

		// Filtra dados inválidos ou SVG capturados por engano
		if eventName != "" && !strings.Contains(eventName, "svg") {
			events = append(events, models.Event{
				Org:             eventOrg,
				Name:            eventName,
				DateMainCard:    eventDateMain,
				DatePrelimsCard: eventDatePrelims,
				Location:        eventLocation,
				City:            eventCity,
				State:           eventState,
				Country:         eventCountry,
			})
		}
	})

	return events, nil
}

// Scrape outra organização
