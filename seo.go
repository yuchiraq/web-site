package main

const (
	siteURL         = "https://avayusstroi.by"
	siteName        = "АВАЮССТРОЙ"
	defaultRobots   = "index,follow,max-image-preview:large"
	defaultOGType   = "website"
	defaultSEOImage = siteURL + "/static/images/hero_new.webp"
)

type SEOData struct {
	Title       string
	Description string
	Canonical   string
	Image       string
	ImageAlt    string
	Robots      string
	Type        string
}

var seoDataByPath = map[string]SEOData{
	"/": {
		Title:       "Автономная канализация и водопонижение в Бресте | АВАЮССТРОЙ",
		Description: "АВАЮССТРОЙ выполняет монтаж автономной канализации, водопонижение, дренаж и аренду спецтехники в Бресте и Брестской области.",
		Image:       siteURL + "/static/images/hero_new.webp",
		ImageAlt:    "Автономная канализация и водопонижение от АВАЮССТРОЙ",
	},
	"/services": {
		Title:       "Инженерные услуги в Бресте | АВАЮССТРОЙ",
		Description: "Водопонижение, дренаж, автономная канализация, наружные сети и другие инженерные работы в Бресте и Брестской области.",
		Image:       siteURL + "/static/images/topas/cover.jpg",
		ImageAlt:    "Инженерные услуги АВАЮССТРОЙ в Бресте",
	},
	"/rent": {
		Title:       "Аренда спецтехники в Бресте | АВАЮССТРОЙ",
		Description: "Аренда спецтехники и оборудования в Бресте и области: техника для земляных, строительных и монтажных работ с выездом на объект.",
		Image:       siteURL + "/static/images/cars/lgce_summer.png",
		ImageAlt:    "Аренда спецтехники АВАЮССТРОЙ",
	},
	"/contacts": {
		Title:       "Контакты АВАЮССТРОЙ | Брест",
		Description: "Телефоны, адрес, реквизиты, быстрые способы связи и вакансии АВАЮССТРОЙ в Бресте. Можно позвонить, написать или оставить заявку.",
		Image:       siteURL + "/static/images/visitka.jpg",
		ImageAlt:    "Контакты АВАЮССТРОЙ",
	},
	"/services/drainage": {
		Title:       "Дренажные системы в Бресте | АВАЮССТРОЙ",
		Description: "Проектирование и устройство дренажных систем в Бресте и Брестской области для отвода воды и защиты участка от переувлажнения.",
		Image:       siteURL + "/static/images/drainage/cover.jpg",
		ImageAlt:    "Дренажные системы от АВАЮССТРОЙ",
	},
	"/services/plumbing": {
		Title:       "Водоснабжение в Бресте | АВАЮССТРОЙ",
		Description: "Проектирование и монтаж наружных сетей водоснабжения под ключ в Бресте и области для частных и коммерческих объектов.",
		Image:       siteURL + "/static/images/plumbing/cover.jpg",
		ImageAlt:    "Наружные сети водоснабжения от АВАЮССТРОЙ",
	},
	"/services/sewerage": {
		Title:       "Монтаж канализации в Бресте | АВАЮССТРОЙ",
		Description: "Монтаж канализационных систем любой сложности в Бресте и Брестской области для домов, участков и коммерческих объектов.",
		Image:       siteURL + "/static/images/sewerage/cover.jpg",
		ImageAlt:    "Монтаж канализации от АВАЮССТРОЙ",
	},
	"/services/storm_sewer": {
		Title:       "Ливневая канализация в Бресте | АВАЮССТРОЙ",
		Description: "Устройство ливневой канализации в Бресте и области для защиты участков, площадок и дорог от подтопления.",
		Image:       siteURL + "/static/images/storm_sewer/cover.jpg",
		ImageAlt:    "Ливневая канализация от АВАЮССТРОЙ",
	},
	"/services/topas": {
		Title:       "Установка автономной канализации ТОПАС в Бресте | Под ключ",
		Description: "Монтаж автономной канализации ТОПАС в Бресте и области. Без запаха, с понятной сметой, выездом инженера и установкой под ключ.",
		Image:       siteURL + "/static/images/topas/cover.jpg",
		ImageAlt:    "Установка автономной канализации ТОПАС",
	},
	"/services/water_lowering": {
		Title:       "Водопонижение в Бресте и Брестской области | Иглофильтры и дренаж — АВАЮССТРОЙ",
		Description: "Услуги водопонижения в Бресте и области: иглофильтры, насосы, дренаж для котлованов и строительных площадок. Точный расчет, безопасные работы, гарантия.",
		Image:       siteURL + "/static/images/water_lowering/cover.png",
		ImageAlt:    "Работы по водопонижению от АВАЮССТРОЙ",
	},
	"/404": {
		Title:       "404 - Страница не найдена | АВАЮССТРОЙ",
		Description: "Страница не найдена. Вернитесь на сайт АВАЮССТРОЙ, чтобы перейти к услугам, аренде спецтехники или контактам.",
		Robots:      "noindex,follow,max-image-preview:large",
		Image:       siteURL + "/static/images/hero_new.webp",
		ImageAlt:    "Страница не найдена",
	},
}

func seoDataForPath(path string) SEOData {
	seo, ok := seoDataByPath[path]
	if !ok {
		seo = SEOData{
			Title:       siteName,
			Description: "Инженерные и строительные услуги в Бресте и Брестской области.",
		}
	}

	if seo.Canonical == "" && path != "/404" {
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

	return seo
}

func canonicalURL(path string) string {
	if path == "" || path == "/" {
		return siteURL + "/"
	}
	if path[0] != '/' {
		return siteURL + "/" + path
	}
	return siteURL + path
}
