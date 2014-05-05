package omtr

import (
	"encoding/json"
	"time"
)

type Locale string

const (
	English    Locale = "en_US"
	German     Locale = "de_DE"
	Spanish    Locale = "es_ES"
	French     Locale = "fr_FR"
	Japanese   Locale = "jp_JP"
	Portuguese Locale = "pt_BR"
	Korean     Locale = "ko_KR"
	Chinese    Locale = "zh_CN"
	Chinese_CN Locale = "zh_CN"
	Chinese_TW Locale = "zh_TW"

	en_US Locale = "en_US"
	de_DE Locale = "de_DE"
	es_ES Locale = "es_ES"
	fr_FR Locale = "fr_FR"
	jp_JP Locale = "jp_JP"
	pt_BR Locale = "pt_BR"
	ko_KR Locale = "ko_KR"
	zh_CN Locale = "zh_CN"
	zh_TW Locale = "zh_TW"
)

type SearchType string

const (
	SearchTypeAND SearchType = "AND"
	SearchTypeOR  SearchType = "OR"
	SearchTypeNOT SearchType = "NOT"
)

type ReportQuery struct {
	ReportDescription *Description `json:"reportDescription"`
}

type Description struct {
	ReportSuiteID    string     `json:"reportSuiteID"`
	Date             string     `json:"date,omitempty"`
	DateFrom         string     `json:"dateFrom,omitempty"`
	DateTo           string     `json:"dateTo,omitempty"`
	DateGranularity  string     `json:"dateGranularity,omitempty"`
	Metrics          []*Metric  `json:"metrics,omitempty"`
	Elements         []*Element `json:"elements,omitempty"`
	Locale           Locale     `json:"locale,omitempty"`
	SortBy           string     `json:"sortBy,omitempty"`
	Segments         []*Segment `json:"segments,omitempty"`
	SegmentId        string     `json:"segment_id,omitempty"`
	AnomalyDetection bool       `json:"anomalyDetection,omitempty"`
	CurrentData      bool       `json:"currentData,omitempty"`
	Expedite         bool       `json:"expedite,omitempty"`
}

type Metric struct {
	Id string `json:"id"`
}

type Element struct {
	Id             string     `json:"id"`
	Classification string     `json:"classification,omitempty"`
	Top            int        `json:"top,omitempty"`
	StartingWith   int        `json:"startingWith,omitempty"`
	Search         *Search    `json:"search,omitempty"`
	Selected       []string   `json:"selected,omitempty"`
	ParentID       string     `json:"parentID,omitempty"`
	Checkpoints    []string   `json:"checkpoints,omitempty"`
	Pattern        [][]string `json:"pattern,omitempty"`
}

type Search struct {
	Type     SearchType `json:"type"`
	Keywords []string   `json:"keywords"`
	Searches []*Search  `json:"searches,omitempty"`
}

type Segment struct {
	Id             string  `json:"id"`
	Element        string  `json:"element,omitempty"`
	Search         *Search `json:"search,omitempty"`
	Classification string  `json:"classification,omitempty"`
}

type ReportResponse struct {
	WaitSeconds   OmtrFloat `json:"waitSeconds"`
	RunSeconds    OmtrFloat `json:"runSeconds"`
	Report        *Report   `json:"report"`
	TimeRetrieved time.Time
}

type Report struct {
	Type        string       `json:"type"`
	ReportSuite *ReportSuite `json:"reportSuite"`
	Period      string       `json:"period"`
	Elements    []*Element   `json:"elements"`
	Metrics     []*Metric    `json:"metrics"`
	Segments    []*Segment   `json:"segments"`
	Data        []*Data      `json:"data"`
	Totals      []OmtrFloat  `json:"totals"`
}

type ReportSuite struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Data struct {
	Name           string      `json:"name"`
	Url            string      `json:"url"`
	Path           *DataPath   `json:"path"`
	ParentID       string      `json:"parentID"`
	Year           int         `json:"year"`
	Month          int         `json:"month"`
	Day            int         `json:"day"`
	Hour           int         `json:"hour"`
	Counts         []OmtrFloat `json:"counts"`
	UpperBounds    []OmtrFloat `json:"upperBounds"`
	LowerBounds    []OmtrFloat `json:"lowerBounds"`
	Forecasts      []OmtrFloat `json:"forecasts"`
	BreakdownTotal []OmtrFloat `json:"breakdownTotal"`
	Breakdown      []*Data     `json:"breakdown"`
}

type DataPath struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (d *Data) DayOfWeek() int {
	date := time.Date(d.Year, time.Month(d.Month), d.Day, d.Hour, 0, 0, 0, time.UTC)
	return int(date.Weekday())
}

type OmtrInt int64
type OmtrFloat float64

func (n *OmtrInt) UnmarshalJSON(data []byte) error {
	var num json.Number
	err := json.Unmarshal(data, &num)
	if err != nil {
		return err
	}

	*(*int64)(n), err = num.Int64()
	return err
}

func (n *OmtrFloat) UnmarshalJSON(data []byte) error {
	var num json.Number
	err := json.Unmarshal(data, &num)
	if err != nil {
		return err
	}

	*(*float64)(n), err = num.Float64()
	return err
}

type queueReport_response struct {
	ReportID OmtrInt `json:"reportID"`
}

type getError struct {
	ErrorName        string `json:"error"`
	ErrorDescription string `json:"error_description"`
	ErrorUri         string `json:"error_uri"`
}

// bind an Error() method to getError type makes it fulfill the error interface
func (e getError) Error() string {
	return e.ErrorDescription
}
