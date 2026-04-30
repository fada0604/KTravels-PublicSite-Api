package domain

type ProviderService struct {
	ID                       string                    `json:"provider_service_id" bson:"_id"`
	ProviderID               string                    `json:"provider_id" bson:"provider_id"`
	Title                    string                    `json:"title" bson:"title"`
	ProviderName             string                    `json:"provider_name" bson:"provider_name"`
	Category                 Category                  `json:"category" bson:"category"`
	SubCategory              SubCategory               `json:"sub_category" bson:"sub_category"`
	Attributes               []Attribute               `json:"attributes" bson:"attributes"`
	Destination              Destination               `json:"destination" bson:"destination"`
	SubcategoriesDestination []SubcategoryDestination  `json:"subcategories_destination" bson:"subcategories_destination"`
	LabelsDestination        []LabelDestination        `json:"labels_destination" bson:"labels_destination"`
	GeoLocation              GeoLocation               `json:"geo_location" bson:"geo_location"`
	Location                 string                    `json:"location" bson:"location"`
	DistributionType         int                       `json:"distribution_type" bson:"distribution_type"`
	Description              Description               `json:"description" bson:"description"`
	Units                    []Unit                    `json:"units" bson:"units"`
}

type Category struct {
	ID   int    `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type SubCategory struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type Attribute struct {
	ID    string `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

type Destination struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type SubcategoryDestination struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type LabelDestination struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type GeoLocation struct {
	Type        string    `json:"type" bson:"type"`
	Coordinates []float64 `json:"coordinates" bson:"coordinates"`
}

type Description struct {
	Description       string `json:"description" bson:"description"`
	HourOperationFrom string `json:"hour_operation_from" bson:"hour_operation_from"`
	HourOperationTo   string `json:"hour_operation_to" bson:"hour_operation_to"`
	ServicePolicies   string `json:"service_policies" bson:"service_policies"`
}

type Unit struct {
	ID         string        `json:"id" bson:"id"`
	Name       string        `json:"name" bson:"name"`
	Capacity   int           `json:"capacity" bson:"capacity"`
	Quantity   int           `json:"quantity" bson:"quantity"`
	Amenities  []interface{} `json:"amenities" bson:"amenities"`
	Images     []Image       `json:"images" bson:"images"`
	Rates      Rate          `json:"rates" bson:"rates"`
	Attributes []Attribute   `json:"attibutes" bson:"attibutes"`
}

type Image struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

type Rate struct {
	ID                    string                `json:"id" bson:"id"`
	RateType              int                   `json:"rate_type" bson:"rate_type"`
	Calendar              Calendar              `json:"calendar" bson:"calendar"`
	CalendarRate          *interface{}          `json:"calendar_rate" bson:"calendar_rate"`
	CalendarPerPersonRate CalendarPerPersonRate `json:"calendar_per_person_rate" bson:"calendar_per_person_rate"`
}

type Calendar struct {
	ID         string              `json:"id" bson:"id"`
	Name       string              `json:"name" bson:"name"`
	StartDate  string              `json:"start_date" bson:"start_date"`
	EndDate    string              `json:"end_date" bson:"end_date"`
	RangeType  int                 `json:"range_type" bson:"range_type"`
	ActiveDays []int               `json:"active_days" bson:"active_days"`
	Exceptions []CalendarException `json:"exceptions" bson:"exceptions"`
}

type CalendarException struct {
	ID          string `json:"id" bson:"id"`
	Name        string `json:"name" bson:"name"`
	StartDate   string `json:"start_date" bson:"start_date"`
	EndDate     string `json:"end_date" bson:"end_date"`
	ActiveDays  []int  `json:"active_days" bson:"active_days"`
	IsRecurring bool   `json:"is_recurring" bson:"is_recurring"`
	Color       string `json:"color" bson:"color"`
}

type CalendarPerPersonRate struct {
	AgeRangeRates  []AgeRangeRate  `json:"age_range_rates" bson:"age_range_rates"`
	ExceptionRates []ExceptionRate `json:"exception_rates" bson:"exception_rates"`
}

type AgeRangeRate struct {
	ID            string  `json:"id" bson:"id"`
	GroupAgeRange int     `json:"group_age_range" bson:"group_age_range"`
	Name          string  `json:"name" bson:"name"`
	Status        bool    `json:"status" bson:"status"`
	MinAge        int     `json:"min_age" bson:"min_age"`
	MaxAge        int     `json:"max_age" bson:"max_age"`
	Rate          float64 `json:"rate" bson:"rate"`
}

type ExceptionRate struct {
	ID                  string         `json:"id" bson:"id"`
	CalendarExceptionID string         `json:"calendar_exception_id" bson:"calendar_exception_id"`
	ExceptionName       string         `json:"exception_name" bson:"exception_name"`
	AgeRangeRates       []AgeRangeRate `json:"age_range_rates" bson:"age_range_rates"`
}
