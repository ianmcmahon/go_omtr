package omtr

import (
	"encoding/json"
)

type Locale string

const (
	English		Locale = "en_US"
	German		Locale = "de_DE"
	Spanish		Locale = "es_ES"
	French		Locale = "fr_FR"
	Japanese	Locale = "jp_JP"
	Portuguese	Locale = "pt_BR"
	Korean		Locale = "ko_KR"
	Chinese		Locale = "zh_CN"
	Chinese_CN	Locale = "zh_CN"
	Chinese_TW	Locale = "zh_TW"	

	en_US		Locale = "en_US"
	de_DE		Locale = "de_DE"
	es_ES		Locale = "es_ES"
	fr_FR		Locale = "fr_FR"
	jp_JP		Locale = "jp_JP"
	pt_BR		Locale = "pt_BR"
	ko_KR		Locale = "ko_KR"
	zh_CN		Locale = "zh_CN"
	zh_TW		Locale = "zh_TW"	
)

type SearchType string 

const (
	SearchTypeAND SearchType = "AND"
	SearchTypeOR  SearchType = "OR"
	SearchTypeNOT SearchType = "NOT"
)

type ReportQuery struct {
	ReportDesc *ReportDescription `json:"reportDescription"`
}

type ReportDescription struct {
	ReportSuiteID		string 		`json:"reportSuiteID"`
	Date				string 		`json:"date,omitempty"`
	DateFrom			string 		`json:"dateFrom,omitempty"`
	DateTo				string 		`json:"dateTo,omitempty"`
	DateGranularity		string 		`json:"dateGranularity,omitempty"`
	Metrics				[]*Metric 	`json:"metrics,omitempty"`
	Elements			[]*Element 	`json:"elements,omitempty"`
	Locale				Locale		`json:"locale,omitempty"`
	SortBy				string		`json:"sortBy,omitempty"`
	Segments			[]*Segment 	`json:"segments,omitempty"`
	SegmentId			string	 	`json:"segment_id,omitempty"`
	AnomalyDetection	bool		`json:"anomalyDetection,omitempty"`
	CurrentData			bool		`json:"currentData,omitempty"`
	Expedite			bool		`json:"expedite,omitempty"`
}

type Metric struct {
	Id 		string 		`json:"id"`
}

type Element struct {
	Id 					string 		`json:"id"`
	Classification		string		`json:"classification,omitempty"`
	Top					int 		`json:"top,omitempty"`
	StartingWith 		int 		`json:"startingWith,omitempty"`
	Search 				*Search 	`json:"search,omitempty"`
	Selected 			[]string 	`json:"selected,omitempty"`
	ParentID			string 		`json:"parentID,omitempty"`
	Checkpoints 		[]string 	`json:"checkpoints,omitempty"`
	Pattern 			[][]string 	`json:"pattern,omitempty"`
}

type Search struct {
	Type 		SearchType 	`json:"type"`
	Keywords	[]string 	`json:"keywords"`
	Searches	[]*Search 	`json:"searches,omitempty"`
}

type Segment struct {
	Id 				string 		`json:"id"`
	Element 		string 		`json:"element,omitempty"`
	Search 			*Search 	`json:"search,omitempty"`
	Classification 	string 		`json:"classification,omitempty"`
}


type Number int64

var _ = json.Unmarshaler(new(Number))

func (n *Number) UnmarshalJSON(data []byte) error {
	var num json.Number
	err := json.Unmarshal(data, &num)
	if err != nil { return err }

	*(*int64)(n), err = num.Int64()
	return err
}


type queueReport_response struct {
	ReportID 		Number		`json:"reportID"`		
}

type getError struct {
	ErrorName			string 		`json:"error"`
	ErrorDescription	string 		`json:"error_description"`
	ErrorUri 			string 		`json:"error_uri"`
}

// bind an Error() method to getError type makes it fulfill the error interface
func (e getError) Error() string {
	return e.ErrorDescription
}
