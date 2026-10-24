package main

import (
	"net/url"
	"sort"
	"strings"
)

const (
	siteURL         = "https://avayusstroi.by"
	siteName        = "ЧСУП АВАЮССТРОЙ"
	defaultRobots   = "index,follow,max-image-preview:large"
	defaultOGType   = "website"
	defaultSEOImage = siteURL + "/static/images/hero_new.webp"
)

type SEOData struct {
	Title           string
	Description     string
	Canonical       string
	Image           string
	ImageAlt        string
	Robots          string
	Type            string
	PageType        string
	HeroImage       string
	HeroImageMobile string
}

var seoDataByPath = map[string]SEOData{
	"/": {
		Title:           "Канализация и водопонижение в Бресте | ЧСУП АВАЮССТРОЙ",
		Description:     "Автономная канализация ТОПАС, монтаж водопровода, дренаж и водопонижение в Бресте и области. Собственная спецтехника. Обсудим участок и рассчитаем работы.",
		Image:           siteURL + "/static/images/hero_new.webp",
		ImageAlt:        "Автономная канализация и водопонижение от АВАЮССТРОЙ",
		HeroImage:       "/static/images/hero_new.webp",
		HeroImageMobile: "/static/images/hero_mobile.png",
	},
	"/services": {
		PageType:    "CollectionPage",
		Title:       "Инженерные услуги в Бресте | ЧСУП АВАЮССТРОЙ",
		Description: "Водопонижение, дренаж, автономная канализация, наружные сети и другие инженерные работы в Бресте и Брестской области.",
		Image:       siteURL + "/static/images/topas/cover.jpg",
		ImageAlt:    "Инженерные услуги АВАЮССТРОЙ в Бресте",
	},
	"/rent": {
		PageType:    "CollectionPage",
		Title:       "Аренда спецтехники в Бресте: экскаватор, манипулятор | АВАЮССТРОЙ",
		Description: "Экскаватор-погрузчик, миниэкскаватор, манипулятор DAF и оборудование в аренду в Бресте и области. Подберём технику под задачу и согласуем подачу на объект.",
		Image:       siteURL + "/static/images/cars/lgce_summer.png",
		ImageAlt:    "Аренда спецтехники АВАЮССТРОЙ",
	},
	"/contacts": {
		Title:       "Контакты ЧСУП АВАЮССТРОЙ | Брест",
		Description: "Телефоны, адрес, реквизиты, быстрые способы связи и вакансии АВАЮССТРОЙ в Бресте. Можно позвонить, написать или оставить заявку.",
		Image:       siteURL + "/static/images/visitka.jpg",
		ImageAlt:    "Контакты АВАЮССТРОЙ",
		PageType:    "ContactPage",
	},
	"/services/drainage": {
		Title:       "Дренаж участка и фундамента в Бресте | АВАЮССТРОЙ",
		Description: "Устройство дренажа участка и фундамента в Бресте и области. Подбор схемы отвода воды, земляные работы, трубы и колодцы. Расчёт после оценки условий участка.",
		Image:       siteURL + "/static/images/drainage/cover.jpg",
		ImageAlt:    "Дренажные системы от АВАЮССТРОЙ",
	},
	"/services/plumbing": {
		Title:       "Монтаж водопровода в Бресте — наружные сети | АВАЮССТРОЙ",
		Description: "Прокладка наружного водопровода в Бресте и области: подбор трассы, земляные работы, монтаж и замена труб. Согласуем подключение и составим смету под объект.",
		Image:       siteURL + "/static/images/plumbing/cover.jpg",
		ImageAlt:    "Наружные сети водоснабжения от АВАЮССТРОЙ",
	},
	"/services/sewerage": {
		Title:       "Канализация под ключ в Бресте — монтаж наружных сетей | АВАЮССТРОЙ",
		Description: "Монтаж канализации для частного дома и строительных объектов в Бресте и области. Наружные трубы, колодцы, земляные работы и подключение. Расчёт по участку.",
		Image:       siteURL + "/static/images/sewerage/cover.jpg",
		ImageAlt:    "Монтаж канализации от АВАЮССТРОЙ",
	},
	"/services/storm_sewer": {
		Title:       "Ливневая канализация в Бресте — монтаж под ключ | АВАЮССТРОЙ",
		Description: "Монтаж ливневой канализации в Бресте и области. Дождеприёмники, лотки, трубы и колодцы для участка и площадки. Подбор схемы водоотвода и расчёт стоимости.",
		Image:       siteURL + "/static/images/storm_sewer/cover.jpg",
		ImageAlt:    "Ливневая канализация от АВАЮССТРОЙ",
	},
	"/services/topas": {
		Title:       "Автономная канализация ТОПАС в Бресте | ЧСУП АВАЮССТРОЙ",
		Description: "Подбор и монтаж автономной канализации ТОПАС для дома и дачи в Бресте и области. Учтём нагрузку, грунт и отвод воды. Рассчитаем оборудование и установку.",
		Image:       siteURL + "/static/images/topas/cover.jpg",
		ImageAlt:    "Установка автономной канализации ТОПАС",
	},
	"/services/water_lowering": {
		Title:       "Водопонижение в Бресте и области | ЧСУП АВАЮССТРОЙ",
		Description: "Строительное водопонижение в Бресте и области: котлованы, траншеи, инженерные сети. Подбор метода по геологии и уровню воды. Расчёт по исходным данным объекта.",
		Image:       siteURL + "/static/images/water_lowering/cover.png",
		ImageAlt:    "Работы по водопонижению от АВАЮССТРОЙ",
	},
	"/404": {
		Title:       "404 - Страница не найдена | ЧСУП АВАЮССТРОЙ",
		Description: "Страница не найдена. Вернитесь на сайт АВАЮССТРОЙ, чтобы перейти к услугам, аренде спецтехники или контактам.",
		Robots:      "noindex,follow,max-image-preview:large",
		Image:       siteURL + "/static/images/hero_new.webp",
		ImageAlt:    "Страница не найдена",
	},
}

// This registry is shared by routes, canonical tags, sitemap and RSS.
// A template alone must never accidentally become an indexable public page.
func publicPagePaths() []string {
	paths := make([]string, 0, len(seoDataByPath)+len(rentalPages)+len(topicPages))
	for path, seo := range seoDataByPath {
		if path != "/404" && !strings.Contains(seo.Robots, "noindex") {
			paths = append(paths, path)
		}
	}
	for _, rental := range rentalPages {
		paths = append(paths, rental.Path)
	}
	for _, topic := range topicPages {
		paths = append(paths, topic.Path)
	}
	sort.Strings(paths)
	return paths
}

func templateForPath(path string) string {
	if topicPageForPath(path) != nil {
		return "templates/static/topic_page.html"
	}
	if rentalPageForPath(path) != nil {
		return "templates/static/rental_detail.html"
	}
	if path == "/" {
		return "templates/index.html"
	}
	return "templates" + path + ".html"
}

func seoDataForPath(path string) SEOData {
	seo, ok := seoDataByPath[path]
	if topic := topicPageForPath(path); topic != nil {
		seo, ok = topic.SEO, true
		seo.Image, seo.ImageAlt = siteURL+topic.Image, topic.ImageAlt
	}
	if rental := rentalPageForPath(path); rental != nil {
		seo, ok = rental.SEO, true
	}
	if !ok {
		seo = SEOData{
			Title:       "Страница не найдена | " + siteName,
			Description: "Перейдите к инженерным услугам, аренде спецтехники или контактам АВАЮССТРОЙ в Бресте.",
			Robots:      "noindex,follow",
		}
	}

	if seo.Canonical == "" && ok && path != "/404" {
		seo.Canonical = canonicalURL(path)
	}
	if seo.Image == "" {
		seo.Image = defaultSEOImage
	}
	if seo.ImageAlt == "" {
		seo.ImageAlt = seo.Title
	}
	if seo.Robots == "" {
		seo.Robots = defaultRobots
	}
	if seo.Type == "" {
		seo.Type = defaultOGType
	}
	if seo.PageType == "" {
		seo.PageType = "WebPage"
	}

	return seo
}

func canonicalURL(path string) string {
	parsed, err := url.Parse(path)
	if err != nil {
		return siteURL + "/"
	}
	cleanPath := "/" + strings.Trim(parsed.Path, "/")
	if cleanPath == "/index.html" || cleanPath == "/index" {
		cleanPath = "/"
	}
	return siteURL + cleanPath
}
