package wave

// Address represents an address associated with a Wave Customer or Business.
type Address struct {
	Address1   *string   `json:"address1,omitempty"`
	Address2   *string   `json:"address2,omitempty"`
	City       *string   `json:"city,omitempty"`
	Province   *Province `json:"province,omitempty"`
	Country    *Country  `json:"country,omitempty"`
	PostalCode *string   `json:"postal_code,omitempty"`
}
