package shifts

import (
	"github.com/nemesisesq/flexable/company"
	"github.com/nemesisesq/flexable/manager"
)

type Location interface{}
type ShiftSite struct {
	Name     string            `json:"name"`
	Location Location          `json:"location"`
	Company  company.Company   `json:"company"`
	Manager  []manager.Manager `json:"manager"`
}
