package main

import (
	"github.com/sujit-baniya/phonenumbers"
	"github.com/sujit-baniya/utils/countries"
	"github.com/sujit-baniya/utils/pool"
	"strings"
)

var PhoneNumberType = map[int]string{
	0: "FIXED_LINE",
	1: "MOBILE",
	// In some regions (e.g. the USA), it is impossible to distinguish between fixed-line and
	// mobile numbers by looking at the phone number itself.
	2: "FIXED_LINE_OR_MOBILE",
	// Freephone lines
	3: "TOLL_FREE",
	4: "PREMIUM_RATE",
	// The cost of this call is shared between the caller and the recipient, and is hence typically
	// less than PREMIUM_RATE calls. See // http://en.wikipedia.org/wiki/Shared_Cost_Service for
	// more information.
	5: "SHARED_COST",
	// Voice over IP numbers. This includes TSoIP (Telephony Service over IP).
	6: "VOIP",
	// A personal number is associated with a particular person, and may be routed to either a
	// MOBILE or FIXED_LINE number. Some more information can be found here:
	// http://en.wikipedia.org/wiki/Personal_Numbers
	7: "PERSONAL_NUMBER",
	8: "PAGER",
	// Used for "Universal Access Numbers" or "Company Numbers". They may be further routed to
	// specific offices, but allow one number to be used for a company.
	9: "UAN",
	// A phone number is of type UNKNOWN when it does not fit any of the known patterns for a
	// specific region.
	10: "UNKNOWN",
	// A phone number is of type UNKNOWN when it does not fit any of the known patterns for a
	// specific region.
	11: "UNKNOWN",

	// Emergency
	27: "EMERGENCY",

	// Voicemail
	28: "VOICEMAIL",

	// Short Code
	29: "SHORT_CODE",
}

type Number struct {
	Phone          string `query:"phone" json:"phone" csv:"phone"`
	DefaultPrefix  string `query:"default_prefix" json:"default_prefix,omitempty" csv:"default_prefix,omitempty"`
	PhoneTypeHuman string `json:"phone_type_human,omitempty" csv:"phone_type_human,omitempty"`
	CarrierName    string `json:"carrier_name,omitempty" csv:"carrier_name,omitempty"`
	CarrierMcc     string `json:"carrier_mcc,omitempty" csv:"carrier_mcc,omitempty"`
	CarrierMnc     string `json:"carrier_mnc,omitempty" csv:"carrier_mnc,omitempty"`
	CarrierNnc     string `json:"carrier_nnc,omitempty" csv:"carrier_nnc,omitempty"`
	CountryName    string `json:"country_name,omitempty" csv:"country_name,omitempty"`
	CountryCode    string `json:"country_code,omitempty" csv:"country_code,omitempty"`
	Location       string `json:"location,omitempty" csv:"location,omitempty"`
	Currency       string `json:"currency,omitempty" csv:"currency,omitempty"`
	CurrencySymbol string `json:"currency_symbol,omitempty" csv:"currency_symbol,omitempty"`
	Timezone       string `json:"timezone,omitempty" csv:"timezone,omitempty"`
	Invalid        bool   `json:"invalid,omitempty" csv:"invalid,omitempty"`
	DialCode       int32  `json:"dial_code,omitempty" csv:"dial_code,omitempty"`
	PhoneType      int    `json:"phone_type,omitempty" csv:"phone_type,omitempty"`
}

type Numbers struct {
	Phones []Number `json:"phones"`
}

type Phones struct {
	Phones []string `json:"phones"`
}

type Ops struct {
	Raw      string `json:"raw"`
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
	Alias    string `json:"alias"`
}

type Query struct {
	Raw        string   `json:"raw"`
	Fields     []string `json:"fields"`
	Operations []Ops    `json:"operations,omitempty"`
	GroupBy    []string `json:"group_by"`
}

type NumberList struct {
	Phone         []string `json:"phones"`
	DefaultPrefix string   `query:"default_prefix" json:"default_prefix,omitempty"`
	PhoneTypes    []string `json:"phone_types"`
	Query         Query    `json:"query"`
	PhoneOnly     bool     `json:"phone_only"`
}

type AnalyzeResult struct {
	CarrierName string `json:"carrier_name,omitempty"`
	CarrierCode string `json:"carrier_code,omitempty"`
	CountryCode string `json:"country,omitempty"`
	PhoneType   string `json:"phone_type,omitempty"`
	PhoneCount  int    `json:"phone_count"`
}

type Stats struct {
	Result       map[string]*AnalyzeResult `json:"-"`
	RS           []*AnalyzeResult          `json:"result"`
	InvalidCount int                       `json:"invalid_count"`
	TotalCount   int                       `json:"total_count"`
}

func (p *Number) Verify() {
	var num *phonenumbers.PhoneNumber
	var e error
	p.Phone = strings.Trim(p.Phone, " ")
	if strings.HasPrefix(p.Phone, "00") {
		p.Phone = strings.Replace(p.Phone, "00", "+", 1)
	}
	hasPlusSymbol := strings.HasPrefix(p.Phone, "+")
	phoneWithoutPlus := strings.Replace(p.Phone, "+", "", 1)
	phoneWithPlus := "+" + phoneWithoutPlus
	if hasPlusSymbol {
		num, e = phonenumbers.Parse(phoneWithPlus, "")
	} else {
		num, e = phonenumbers.Parse(phoneWithoutPlus, strings.ToUpper(p.DefaultPrefix))
	}
	if e != nil {
		num, e = phonenumbers.Parse(phoneWithPlus, "")
	}
	if e != nil {
		num = nil
		p.Invalid = true
		return
	}
	p.PhoneType = int(phonenumbers.GetNumberType(num))
	p.PhoneTypeHuman = PhoneNumberType[p.PhoneType]
	if p.PhoneType == 10 || p.PhoneType == 11 {
		num = nil
		p.Invalid = true
		return
	}

	region := phonenumbers.GetRegionCodeForNumber(num)
	p.Phone = phonenumbers.Format(num, phonenumbers.E164)
	carrier, _ := phonenumbers.GetCarrierForNumber(num, "EN")
	geo, _ := phonenumbers.GetGeocodingForNumber(num, "EN")
	timezones, e := phonenumbers.GetTimezonesForNumber(num)

	if e != nil {
		num = nil
		p.Invalid = true
		return
	}
	if len(timezones) > 0 {
		p.Timezone = timezones[0]
	}
	p.DialCode = num.GetCountryCode()
	p.CountryCode = region

	country, _ := countries.Countries[region]
	p.CountryName = country.Name
	p.Currency = country.Currency
	p.CurrencySymbol = country.CurrencySymbol
	p.Location = geo
	p.CarrierName = carrier
	networks := CountryNetwork[region]
	for _, net := range networks {
		if carrier != "" && (strings.Contains(strings.ToLower(net.Brand), strings.ToLower(carrier)) ||
			strings.Contains(strings.ToLower(net.Operator), strings.ToLower(carrier))) {
			p.CarrierMnc = net.Mnc
			p.CarrierMcc = net.Mcc
			p.CarrierNnc = net.Plmn
			break
		}
	}
}

func (p *Numbers) Verify() {
	for _, phone := range p.Phones {
		phone.Verify()
	}
}

var allowedPhoneTypes = map[string]int{
	"landline":           4,
	"mobile":             1,
	"mobile_or_landline": 2,
	"toll_free":          3,
	"voip":               6,
}

func (p *NumberList) Verify() Numbers {
	numbers := Numbers{}
	workerPool := pool.New()
	defer workerPool.Close()

	batch := workerPool.Batch()

	go func() {
		for _, phone := range p.Phone {
			num := Number{Phone: phone, DefaultPrefix: p.DefaultPrefix}
			batch.Queue(verify(num))
		}
		batch.QueueComplete()
	}()
	for phone := range batch.Results() {
		numbers.Phones = append(numbers.Phones, phone.Value().(Number))
	}
	return numbers
}

func (p *NumberList) Clean() (Numbers, Phones) {
	phones := Phones{}
	numbers := Numbers{}
	workerPool := pool.New()
	defer workerPool.Close()

	batch := workerPool.Batch()

	go func() {
		for _, phone := range p.Phone {
			num := Number{Phone: phone, DefaultPrefix: p.DefaultPrefix}
			batch.Queue(verify(num))
		}
		batch.QueueComplete()
	}()
	for phone := range batch.Results() {
		num := phone.Value().(Number)
		if !num.Invalid {
			if len(p.PhoneTypes) > 0 {
				for _, phoneType := range p.PhoneTypes {
					allowed := allowedPhoneTypes[phoneType]
					if allowed == 4 && num.PhoneType == 0 {
						if p.PhoneOnly {
							phones.Phones = append(phones.Phones, num.Phone)
						} else {
							numbers.Phones = append(numbers.Phones, num)
						}

					} else if allowed > 0 && num.PhoneType == allowed {
						if p.PhoneOnly {
							phones.Phones = append(phones.Phones, num.Phone)
						} else {
							numbers.Phones = append(numbers.Phones, num)
						}
					}
				}
			} else {
				if p.PhoneOnly {
					phones.Phones = append(phones.Phones, num.Phone)
				} else {
					numbers.Phones = append(numbers.Phones, num)
				}
			}

		}
	}
	return numbers, phones
}

func (p *NumberList) Stats() Stats {
	stats := Stats{}
	rs := make(map[string]*AnalyzeResult)
	workerPool := pool.New()
	defer workerPool.Close()

	batch := workerPool.Batch()
	stats.TotalCount = len(p.Phone)
	go func() {
		for _, phone := range p.Phone {
			num := Number{Phone: phone, DefaultPrefix: p.DefaultPrefix}
			batch.Queue(verify(num))
		}
		batch.QueueComplete()
	}()
	for phone := range batch.Results() {
		num := phone.Value().(Number)
		if num.Invalid {
			stats.InvalidCount += 1
			continue
		}
		hash := num.CountryCode
		if num.CarrierNnc != "" {
			hash = hash + ":" + num.CarrierNnc
		}
		hash = hash + ":" + num.PhoneTypeHuman
		analyzeResult := &AnalyzeResult{
			CarrierName: num.CarrierName,
			CarrierCode: num.CarrierNnc,
			CountryCode: num.CountryCode,
			PhoneType:   num.PhoneTypeHuman,
		}

		if rs[hash] == nil {
			analyzeResult.PhoneCount = 1
			rs[hash] = analyzeResult
		} else {
			rs[hash].PhoneCount += 1
		}
	}

	for _, val := range rs {
		stats.RS = append(stats.RS, val)
	}
	rs = nil
	return stats
}

func verify(phone Number) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		if wu.IsCancelled() {
			return nil, nil
		}
		phone.Verify()
		return phone, nil
	}
}
