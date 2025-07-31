package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

// Page представляет данные для одной страницы в RSS-фиде
type Page struct {
	Title        string
	URL          string
	Description  string
	PubDate      string
	Author       string
	TurboContent string
}

func generateRSSFeed(c *gin.Context) {
	// Путь к директории с HTML-шаблонами
	templateDir := "templates"

	// Извлекаем номера телефонов из contacts.html
	contactFilePath := "templates/contacts.html"
	contactDoc, err := goquery.NewDocumentFromReader(strings.NewReader(readFile(contactFilePath)))
	if err != nil {
		c.String(500, fmt.Sprintf("Ошибка парсинга contacts.html: %v", err))
		return
	}
	var phoneNumbers []string
	contactDoc.Find("a[href^='tel:']").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			phoneNumber := strings.TrimPrefix(href, "tel:")
			phoneNumbers = append(phoneNumbers, phoneNumber)
		}
	})
	if len(phoneNumbers) == 0 {
		phoneNumbers = append(phoneNumbers, "+375123456789") // Запасной номер
	}

	// Список страниц для RSS
	var pages []Page

	// Рекурсивно обходим директорию templates
	err = filepath.Walk(templateDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Обрабатываем только HTML-файлы в поддиректории services/
		if !info.IsDir() && strings.HasSuffix(filePath, ".html") {
			relativePath := strings.TrimPrefix(filePath, "templates/")
			pagePath := strings.TrimSuffix(relativePath, ".html")

			// Пропускаем корневые страницы
			if pagePath == "index" || pagePath == "services" || pagePath == "rent" || pagePath == "contacts" {
				return nil
			}

			// Проверяем, что файл находится в services/
			if !strings.HasPrefix(pagePath, "services/") {
				return nil
			}

			// Определяем URL страницы
			pageName := strings.TrimPrefix(pagePath, "services/")
			pageURL := fmt.Sprintf("https://avayusstroi.by/services/%s", pageName)

			// Парсим HTML-файл
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(readFile(filePath)))
			if err != nil {
				return fmt.Errorf("ошибка парсинга файла %s: %v", filePath, err)
			}

			// Извлекаем заголовок
			title := doc.Find("title").Text()

			// Извлекаем описание
			description, _ := doc.Find(`meta[name="description"]`).Attr("content")
			if description == "" {
				description = "Услуги по строительству и дренажу от ЧСУП «АВАЮССТРОЙ»"
			}

			// Извлекаем контент для <turbo:content>
			turboContent := `<![CDATA[`
			turboContent += `<header>`

			// Заголовок для <h1> (берём первый <h2> или <h1> из <main>)
			h1 := doc.Find("main h1").First().Text()
			if h1 == "" {
				h1 = doc.Find("main h2").First().Text()
			}
			if h1 != "" {
				turboContent += fmt.Sprintf(`<h1>%s</h1>`, h1)
			}

			// Изображение (берём из .service-image)
			imgSrc, _ := doc.Find(".service-image img").Attr("src")
			if imgSrc != "" {
				turboContent += `<figure>`
				turboContent += fmt.Sprintf(`<img src="https://avayusstroi.by%s">`, imgSrc)
				turboContent += `</figure>`
			}

			// Меню
			turboContent += `<menu>`
			turboContent += `<a href="https://avayusstroi.by/">Главная</a>`
			turboContent += `<a href="https://avayusstroi.by/services">Услуги</a>`
			turboContent += `</menu>`
			turboContent += `</header>`

			// Извлекаем весь контент из <main>
			mainContent := doc.Find("main")
			mainContent.Children().Each(func(i int, section *goquery.Selection) {
				section.Children().Each(func(j int, container *goquery.Selection) {
					if container.HasClass("container") {
						// Заголовки <h2>
						container.Find("h2").Each(func(k int, h2 *goquery.Selection) {
							turboContent += fmt.Sprintf(`<h2>%s</h2>`, h2.Text())
						})

						// Текст из <p>
						container.Find("p").Each(func(k int, p *goquery.Selection) {
							text := strings.TrimSpace(p.Text())
							if text != "" {
								turboContent += fmt.Sprintf(`<p>%s</p>`, text)
							}
						})

						// Списки <ul>
						container.Find("ul").Each(func(k int, ul *goquery.Selection) {
							turboContent += `<ul>`
							ul.Find("li").Each(func(m int, li *goquery.Selection) {
								text := strings.TrimSpace(li.Text())
								if text != "" {
									turboContent += fmt.Sprintf(`<li>%s</li>`, text)
								}
							})
							turboContent += `</ul>`
						})

						// Секция .service-content
						if container.HasClass("service-content") {
							container.Find(".service-description").Children().Each(func(k int, descChild *goquery.Selection) {
								if descChild.Is("h3") {
									turboContent += fmt.Sprintf(`<h2>%s</h2>`, descChild.Text())
								} else if descChild.Is("p") {
									text := strings.TrimSpace(descChild.Text())
									if text != "" {
										turboContent += fmt.Sprintf(`<p>%s</p>`, text)
									}
								} else if descChild.Is("ul") {
									turboContent += `<ul>`
									descChild.Find("li").Each(func(m int, li *goquery.Selection) {
										text := strings.TrimSpace(li.Text())
										if text != "" {
											turboContent += fmt.Sprintf(`<li>%s</li>`, text)
										}
									})
									turboContent += `</ul>`
								}
							})
						}

						// Галерея (например, для topas.html)
						if container.HasClass("service-gallery") {
							container.Find("h3").Each(func(k int, h3 *goquery.Selection) {
								turboContent += fmt.Sprintf(`<h2>%s</h2>`, h3.Text())
							})
							container.Find(".gallery img").Each(func(k int, img *goquery.Selection) {
								if src, exists := img.Attr("src"); exists {
									turboContent += `<figure>`
									turboContent += fmt.Sprintf(`<img src="https://avayusstroi.by%s">`, src)
									turboContent += `</figure>`
								}
							})
						}
					}
				})
			})

			// Добавляем кнопки "Заказать" с номерами телефонов
			for _, phoneNumber := range phoneNumbers {
				turboContent += fmt.Sprintf(`<button formaction="tel:%s" data-background-color="#5B97B0" data-color="white" data-primary="true">Заказать</button>`, phoneNumber)
			}

			turboContent += `]]>`

			// Добавляем страницу в список
			pages = append(pages, Page{
				Title:        title,
				URL:          pageURL,
				Description:  description,
				PubDate:      time.Now().Format(time.RFC1123Z),
				Author:       "АВАЮССТРОЙ",
				TurboContent: turboContent,
			})
		}
		return nil
	})

	if err != nil {
		c.String(500, fmt.Sprintf("Ошибка обхода директории шаблонов: %v", err))
		return
	}

	// Генерируем XML
	xmlContent := `<?xml version="1.0" encoding="UTF-8"?>
    <rss xmlns:yandex="http://news.yandex.ru" xmlns:media="http://search.yahoo.com/mrss/" xmlns:turbo="http://turbo.yandex.ru" xmlns:atom="http://www.w3.org/2005/Atom" version="2.0">
        <channel>
            <title>АВАЮССТРОЙ | Услуги</title>
            <link>https://avayusstroi.by</link>
            <description>Услуги по строительству и дренажу в Бресте и области</description>
            <language>ru</language>
    `

	for _, page := range pages {
		xmlContent += `
            <item turbo="true">
                <title>` + page.Title + `</title>
                <link>` + page.URL + `</link>
                <description>` + page.Description + `</description>
                <pubDate>` + page.PubDate + `</pubDate>
                <author>` + page.Author + `</author>
                <turbo:content>` + page.TurboContent + `</turbo:content>
            </item>`
	}

	xmlContent += `
        </channel>
    </rss>`

	// Отдаём XML
	c.Data(200, "application/xml", []byte(xmlContent))
}

// Вспомогательная функция для чтения файла
func readFile(filePath string) string {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка чтения файла %s: %v\n", filePath, err)
		os.Exit(1) // <- аварийный выход, чтобы явно увидеть ошибку
	}
	return string(content)
}
