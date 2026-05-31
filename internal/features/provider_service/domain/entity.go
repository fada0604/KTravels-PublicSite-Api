package domain

type EntityAssociation string

const (
	EntityAssociationProviderIdentity EntityAssociation = "ProviderDigitalIdentity"
	EntityAssociationServiceUnit      EntityAssociation = "ServiceDistributionUnit"
)

type ProviderService struct {
	ID                       string                   `json:"provider_service_id"`
	ProviderID               string                   `json:"provider_id"`
	Status                   int                      `json:"status"`
	Title                    string                   `json:"title"`
	ProviderName             string                   `json:"provider_name"`
	ProviderLogo             *Image                   `json:"provider_logo"`
	Category                 Category                 `json:"category"`
	SubCategory              SubCategory              `json:"sub_category"`
	Attributes               []Attribute              `json:"attributes"`
	Destination              Destination              `json:"destination"`
	SubcategoriesDestination []SubcategoryDestination `json:"subcategories_destination"`
	LabelsDestination        []LabelDestination       `json:"labels_destination"`
	GeoLocation              GeoLocation              `json:"geo_location"`
	Location                 string                   `json:"location"`
	DistributionType         int                      `json:"distribution_type"`
	Description              Description              `json:"description"`
	Units                    []Unit                   `json:"units"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SubCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Attribute struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Destination struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type SubcategoryDestination struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type LabelDestination struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GeoLocation struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Description struct {
	Description       string `json:"description"`
	HourOperationFrom string `json:"hour_operation_from"`
	HourOperationTo   string `json:"hour_operation_to"`
	ServicePolicies   string `json:"service_policies"`
}

type Unit struct {
	ID         string        `json:"id"`
	Name       string        `json:"name"`
	Capacity   int           `json:"capacity"`
	Quantity   int           `json:"quantity"`
	Amenities  []any         `json:"amenities"`
	Images     []Image       `json:"images"`
	Rates      Rate          `json:"rates"`
	Attributes []Attribute   `json:"attibutes"`
}

type Image struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Rate struct {
	ID                    string               `json:"id"`
	RateType              int                  `json:"rate_type"`
	Calendar              Calendar             `json:"calendar"`
	CalendarRate          *any                 `json:"calendar_rate"`
	CalendarPerPersonRate CalendarPerPersonRate `json:"calendar_per_person_rate"`
}

type Calendar struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	StartDate  string              `json:"start_date"`
	EndDate    string              `json:"end_date"`
	RangeType  int                 `json:"range_type"`
	ActiveDays []int               `json:"active_days"`
	Exceptions []CalendarException `json:"exceptions"`
}

type CalendarException struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
	ActiveDays  []int  `json:"active_days"`
	IsRecurring bool   `json:"is_recurring"`
	Color       string `json:"color"`
}

type CalendarPerPersonRate struct {
	AgeRangeRates  []AgeRangeRate  `json:"age_range_rates"`
	ExceptionRates []ExceptionRate `json:"exception_rates"`
}

type AgeRangeRate struct {
	ID            string  `json:"id"`
	GroupAgeRange int     `json:"group_age_range"`
	Name          string  `json:"name"`
	Status        bool    `json:"status"`
	MinAge        int     `json:"min_age"`
	MaxAge        int     `json:"max_age"`
	Rate          float64 `json:"rate"`
}

type ExceptionRate struct {
	ID                  string         `json:"id"`
	CalendarExceptionID string         `json:"calendar_exception_id"`
	ExceptionName       string         `json:"exception_name"`
	AgeRangeRates       []AgeRangeRate `json:"age_range_rates"`
}
