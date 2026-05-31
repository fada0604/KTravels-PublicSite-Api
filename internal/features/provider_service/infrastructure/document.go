package infrastructure

import "ktravels-publicsite-api/internal/features/provider_service/domain"

type ProviderServiceDocument struct {
	ID                       string                           `bson:"_id"`
	ProviderID               string                           `bson:"provider_id"`
	Status                   int                              `bson:"status"`
	Title                    string                           `bson:"title"`
	ProviderName             string                           `bson:"provider_name"`
	ProviderLogo             *ImageDocument                   `bson:"provider_logo,omitempty"`
	Category                 CategoryDocument                 `bson:"category"`
	SubCategory              SubCategoryDocument              `bson:"sub_category"`
	Attributes               []AttributeDocument              `bson:"attributes"`
	Destination              DestinationDocument              `bson:"destination"`
	SubcategoriesDestination []SubcategoryDestinationDocument `bson:"subcategories_destination"`
	LabelsDestination        []LabelDestinationDocument       `bson:"labels_destination"`
	GeoLocation              GeoLocationDocument              `bson:"geo_location"`
	Location                 string                           `bson:"location"`
	DistributionType         int                              `bson:"distribution_type"`
	Description              DescriptionDocument              `bson:"description"`
	Units                    []UnitDocument                   `bson:"units"`
}

type CategoryDocument struct {
	ID   int    `bson:"id"`
	Name string `bson:"name"`
}

type SubCategoryDocument struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}

type AttributeDocument struct {
	ID    string `bson:"id"`
	Name  string `bson:"name"`
	Value string `bson:"value"`
}

type DestinationDocument struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}

type SubcategoryDestinationDocument struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}

type LabelDestinationDocument struct {
	ID   string `bson:"id"`
	Name string `bson:"name"`
}

type GeoLocationDocument struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

type DescriptionDocument struct {
	Description       string `bson:"description"`
	HourOperationFrom string `bson:"hour_operation_from"`
	HourOperationTo   string `bson:"hour_operation_to"`
	ServicePolicies   string `bson:"service_policies"`
}

type UnitDocument struct {
	ID         string              `bson:"id"`
	Name       string              `bson:"name"`
	Capacity   int                 `bson:"capacity"`
	Quantity   int                 `bson:"quantity"`
	Amenities  []any               `bson:"amenities"`
	Images     []ImageDocument     `bson:"images"`
	Rates      RateDocument        `bson:"rates"`
	Attributes []AttributeDocument `bson:"attibutes"` // preserves original field name in MongoDB
}

type ImageDocument struct {
	Name string `bson:"name"`
	URL  string `bson:"url"`
}

type RateDocument struct {
	ID                    string                       `bson:"id"`
	RateType              int                          `bson:"rate_type"`
	Calendar              CalendarDocument             `bson:"calendar"`
	CalendarRate          *any                         `bson:"calendar_rate"`
	CalendarPerPersonRate CalendarPerPersonRateDocument `bson:"calendar_per_person_rate"`
}

type CalendarDocument struct {
	ID         string                      `bson:"id"`
	Name       string                      `bson:"name"`
	StartDate  string                      `bson:"start_date"`
	EndDate    string                      `bson:"end_date"`
	RangeType  int                         `bson:"range_type"`
	ActiveDays []int                       `bson:"active_days"`
	Exceptions []CalendarExceptionDocument `bson:"exceptions"`
}

type CalendarExceptionDocument struct {
	ID          string `bson:"id"`
	Name        string `bson:"name"`
	StartDate   string `bson:"start_date"`
	EndDate     string `bson:"end_date"`
	ActiveDays  []int  `bson:"active_days"`
	IsRecurring bool   `bson:"is_recurring"`
	Color       string `bson:"color"`
}

type CalendarPerPersonRateDocument struct {
	AgeRangeRates  []AgeRangeRateDocument  `bson:"age_range_rates"`
	ExceptionRates []ExceptionRateDocument `bson:"exception_rates"`
}

type AgeRangeRateDocument struct {
	ID            string  `bson:"id"`
	GroupAgeRange int     `bson:"group_age_range"`
	Name          string  `bson:"name"`
	Status        bool    `bson:"status"`
	MinAge        int     `bson:"min_age"`
	MaxAge        int     `bson:"max_age"`
	Rate          float64 `bson:"rate"`
}

type ExceptionRateDocument struct {
	ID                  string                `bson:"id"`
	CalendarExceptionID string                `bson:"calendar_exception_id"`
	ExceptionName       string                `bson:"exception_name"`
	AgeRangeRates       []AgeRangeRateDocument `bson:"age_range_rates"`
}

func toDocument(s *domain.ProviderService) *ProviderServiceDocument {
	doc := &ProviderServiceDocument{
		ID:               s.ID,
		ProviderID:       s.ProviderID,
		Status:           s.Status,
		Title:            s.Title,
		ProviderName:     s.ProviderName,
		Category:         CategoryDocument{ID: s.Category.ID, Name: s.Category.Name},
		SubCategory:      SubCategoryDocument{ID: s.SubCategory.ID, Name: s.SubCategory.Name},
		Destination:      DestinationDocument{ID: s.Destination.ID, Name: s.Destination.Name},
		GeoLocation:      GeoLocationDocument{Type: s.GeoLocation.Type, Coordinates: s.GeoLocation.Coordinates},
		Location:         s.Location,
		DistributionType: s.DistributionType,
		Description: DescriptionDocument{
			Description:       s.Description.Description,
			HourOperationFrom: s.Description.HourOperationFrom,
			HourOperationTo:   s.Description.HourOperationTo,
			ServicePolicies:   s.Description.ServicePolicies,
		},
	}

	if s.ProviderLogo != nil {
		doc.ProviderLogo = &ImageDocument{Name: s.ProviderLogo.Name, URL: s.ProviderLogo.URL}
	}

	doc.Attributes = make([]AttributeDocument, len(s.Attributes))
	for i, a := range s.Attributes {
		doc.Attributes[i] = AttributeDocument{ID: a.ID, Name: a.Name, Value: a.Value}
	}

	doc.SubcategoriesDestination = make([]SubcategoryDestinationDocument, len(s.SubcategoriesDestination))
	for i, sd := range s.SubcategoriesDestination {
		doc.SubcategoriesDestination[i] = SubcategoryDestinationDocument{ID: sd.ID, Name: sd.Name}
	}

	doc.LabelsDestination = make([]LabelDestinationDocument, len(s.LabelsDestination))
	for i, ld := range s.LabelsDestination {
		doc.LabelsDestination[i] = LabelDestinationDocument{ID: ld.ID, Name: ld.Name}
	}

	doc.Units = make([]UnitDocument, len(s.Units))
	for i, u := range s.Units {
		doc.Units[i] = toUnitDocument(u)
	}

	return doc
}

func toUnitDocument(u domain.Unit) UnitDocument {
	ud := UnitDocument{
		ID:        u.ID,
		Name:      u.Name,
		Capacity:  u.Capacity,
		Quantity:  u.Quantity,
		Amenities: u.Amenities,
		Rates:     toRateDocument(u.Rates),
	}

	ud.Images = make([]ImageDocument, len(u.Images))
	for i, img := range u.Images {
		ud.Images[i] = ImageDocument{Name: img.Name, URL: img.URL}
	}

	ud.Attributes = make([]AttributeDocument, len(u.Attributes))
	for i, a := range u.Attributes {
		ud.Attributes[i] = AttributeDocument{ID: a.ID, Name: a.Name, Value: a.Value}
	}

	return ud
}

func toRateDocument(r domain.Rate) RateDocument {
	cal := r.Calendar
	calDoc := CalendarDocument{
		ID:         cal.ID,
		Name:       cal.Name,
		StartDate:  cal.StartDate,
		EndDate:    cal.EndDate,
		RangeType:  cal.RangeType,
		ActiveDays: cal.ActiveDays,
	}
	calDoc.Exceptions = make([]CalendarExceptionDocument, len(cal.Exceptions))
	for i, e := range cal.Exceptions {
		calDoc.Exceptions[i] = CalendarExceptionDocument{
			ID:          e.ID,
			Name:        e.Name,
			StartDate:   e.StartDate,
			EndDate:     e.EndDate,
			ActiveDays:  e.ActiveDays,
			IsRecurring: e.IsRecurring,
			Color:       e.Color,
		}
	}

	cppr := r.CalendarPerPersonRate
	cppDoc := CalendarPerPersonRateDocument{
		AgeRangeRates:  make([]AgeRangeRateDocument, len(cppr.AgeRangeRates)),
		ExceptionRates: make([]ExceptionRateDocument, len(cppr.ExceptionRates)),
	}
	for i, ar := range cppr.AgeRangeRates {
		cppDoc.AgeRangeRates[i] = toAgeRangeRateDocument(ar)
	}
	for i, er := range cppr.ExceptionRates {
		erd := ExceptionRateDocument{
			ID:                  er.ID,
			CalendarExceptionID: er.CalendarExceptionID,
			ExceptionName:       er.ExceptionName,
			AgeRangeRates:       make([]AgeRangeRateDocument, len(er.AgeRangeRates)),
		}
		for j, ar := range er.AgeRangeRates {
			erd.AgeRangeRates[j] = toAgeRangeRateDocument(ar)
		}
		cppDoc.ExceptionRates[i] = erd
	}

	return RateDocument{
		ID:                    r.ID,
		RateType:              r.RateType,
		Calendar:              calDoc,
		CalendarRate:          r.CalendarRate,
		CalendarPerPersonRate: cppDoc,
	}
}

func toAgeRangeRateDocument(ar domain.AgeRangeRate) AgeRangeRateDocument {
	return AgeRangeRateDocument{
		ID:            ar.ID,
		GroupAgeRange: ar.GroupAgeRange,
		Name:          ar.Name,
		Status:        ar.Status,
		MinAge:        ar.MinAge,
		MaxAge:        ar.MaxAge,
		Rate:          ar.Rate,
	}
}

func fromDocument(doc *ProviderServiceDocument) *domain.ProviderService {
	s := &domain.ProviderService{
		ID:               doc.ID,
		ProviderID:       doc.ProviderID,
		Status:           doc.Status,
		Title:            doc.Title,
		ProviderName:     doc.ProviderName,
		Category:         domain.Category{ID: doc.Category.ID, Name: doc.Category.Name},
		SubCategory:      domain.SubCategory{ID: doc.SubCategory.ID, Name: doc.SubCategory.Name},
		Destination:      domain.Destination{ID: doc.Destination.ID, Name: doc.Destination.Name},
		GeoLocation:      domain.GeoLocation{Type: doc.GeoLocation.Type, Coordinates: doc.GeoLocation.Coordinates},
		Location:         doc.Location,
		DistributionType: doc.DistributionType,
		Description: domain.Description{
			Description:       doc.Description.Description,
			HourOperationFrom: doc.Description.HourOperationFrom,
			HourOperationTo:   doc.Description.HourOperationTo,
			ServicePolicies:   doc.Description.ServicePolicies,
		},
	}

	if doc.ProviderLogo != nil {
		s.ProviderLogo = &domain.Image{Name: doc.ProviderLogo.Name, URL: doc.ProviderLogo.URL}
	}

	s.Attributes = make([]domain.Attribute, len(doc.Attributes))
	for i, a := range doc.Attributes {
		s.Attributes[i] = domain.Attribute{ID: a.ID, Name: a.Name, Value: a.Value}
	}

	s.SubcategoriesDestination = make([]domain.SubcategoryDestination, len(doc.SubcategoriesDestination))
	for i, sd := range doc.SubcategoriesDestination {
		s.SubcategoriesDestination[i] = domain.SubcategoryDestination{ID: sd.ID, Name: sd.Name}
	}

	s.LabelsDestination = make([]domain.LabelDestination, len(doc.LabelsDestination))
	for i, ld := range doc.LabelsDestination {
		s.LabelsDestination[i] = domain.LabelDestination{ID: ld.ID, Name: ld.Name}
	}

	s.Units = make([]domain.Unit, len(doc.Units))
	for i, u := range doc.Units {
		s.Units[i] = fromUnitDocument(u)
	}

	return s
}

func fromUnitDocument(u UnitDocument) domain.Unit {
	unit := domain.Unit{
		ID:        u.ID,
		Name:      u.Name,
		Capacity:  u.Capacity,
		Quantity:  u.Quantity,
		Amenities: u.Amenities,
		Rates:     fromRateDocument(u.Rates),
	}

	unit.Images = make([]domain.Image, len(u.Images))
	for i, img := range u.Images {
		unit.Images[i] = domain.Image{Name: img.Name, URL: img.URL}
	}

	unit.Attributes = make([]domain.Attribute, len(u.Attributes))
	for i, a := range u.Attributes {
		unit.Attributes[i] = domain.Attribute{ID: a.ID, Name: a.Name, Value: a.Value}
	}

	return unit
}

func fromRateDocument(r RateDocument) domain.Rate {
	cal := r.Calendar
	calDomain := domain.Calendar{
		ID:         cal.ID,
		Name:       cal.Name,
		StartDate:  cal.StartDate,
		EndDate:    cal.EndDate,
		RangeType:  cal.RangeType,
		ActiveDays: cal.ActiveDays,
	}
	calDomain.Exceptions = make([]domain.CalendarException, len(cal.Exceptions))
	for i, e := range cal.Exceptions {
		calDomain.Exceptions[i] = domain.CalendarException{
			ID:          e.ID,
			Name:        e.Name,
			StartDate:   e.StartDate,
			EndDate:     e.EndDate,
			ActiveDays:  e.ActiveDays,
			IsRecurring: e.IsRecurring,
			Color:       e.Color,
		}
	}

	cppr := r.CalendarPerPersonRate
	cppDomain := domain.CalendarPerPersonRate{
		AgeRangeRates:  make([]domain.AgeRangeRate, len(cppr.AgeRangeRates)),
		ExceptionRates: make([]domain.ExceptionRate, len(cppr.ExceptionRates)),
	}
	for i, ar := range cppr.AgeRangeRates {
		cppDomain.AgeRangeRates[i] = fromAgeRangeRateDocument(ar)
	}
	for i, er := range cppr.ExceptionRates {
		erd := domain.ExceptionRate{
			ID:                  er.ID,
			CalendarExceptionID: er.CalendarExceptionID,
			ExceptionName:       er.ExceptionName,
			AgeRangeRates:       make([]domain.AgeRangeRate, len(er.AgeRangeRates)),
		}
		for j, ar := range er.AgeRangeRates {
			erd.AgeRangeRates[j] = fromAgeRangeRateDocument(ar)
		}
		cppDomain.ExceptionRates[i] = erd
	}

	return domain.Rate{
		ID:                    r.ID,
		RateType:              r.RateType,
		Calendar:              calDomain,
		CalendarRate:          r.CalendarRate,
		CalendarPerPersonRate: cppDomain,
	}
}

func fromAgeRangeRateDocument(ar AgeRangeRateDocument) domain.AgeRangeRate {
	return domain.AgeRangeRate{
		ID:            ar.ID,
		GroupAgeRange: ar.GroupAgeRange,
		Name:          ar.Name,
		Status:        ar.Status,
		MinAge:        ar.MinAge,
		MaxAge:        ar.MaxAge,
		Rate:          ar.Rate,
	}
}

func toImageDocument(img domain.Image) ImageDocument {
	return ImageDocument{Name: img.Name, URL: img.URL}
}

func toImageDocuments(imgs []domain.Image) []ImageDocument {
	docs := make([]ImageDocument, len(imgs))
	for i, img := range imgs {
		docs[i] = toImageDocument(img)
	}
	return docs
}
