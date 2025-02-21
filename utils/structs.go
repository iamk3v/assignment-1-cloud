package utils

type PopulationObject struct {
	Year  int `json:"year"`
	Value int `json:"value"`
}

type PopulationInfo struct {
	Country          string             `json:"country"`
	Code             string             `json:"code"`
	Iso3             string             `json:"iso3"`
	PopulationCounts []PopulationObject `json:"populationCounts"`
}

type PopulationData struct {
	Error bool           `json:"error"`
	Msg   string         `json:"msg"`
	Data  PopulationInfo `json:"data"`
}

type CitiesJson struct {
	Error  bool     `json:"error"`
	Msg    string   `json:"msg"`
	Cities []string `json:"data"`
}

type Name struct {
	Common   string `json:"common"`
	Official string `json:"official"`
}

type CountryName struct {
	Name Name `json:"name"`
}

type Flags struct {
	Png string `json:"png"`
}

type RestCountriesJson struct {
	Name       Name              `json:"name"`
	Continents []string          `json:"continents"`
	Population int               `json:"population"`
	Languages  map[string]string `json:"languages"`
	Borders    []string          `json:"borders"`
	Flags      Flags             `json:"flags"`
	Capital    []string          `json:"capital"`
}

type PopulationResponseJson struct {
	Mean   int
	Values []PopulationObject
}

type InfoResponseJson struct {
	Name       string
	Continents []string
	Population int
	Languages  map[string]string
	Borders    []string
	Flag       string
	Capital    []string
	Cities     []string
}
