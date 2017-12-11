package shifts

import "time"
import (
	"github.com/nemesisesq/flexable/employees"
	"github.com/nemesisesq/flexable/manager"
)

type Shift struct {
	Employee     employees.Employee   `json:"employee"`
	Manager      manager.Manager      `json:"manager"`
	ShiftSite    ShiftSite            `json:"shift_site"`
	DateTime     time.Time            `json:"date_time"`
	Replacements []employees.Employee `json:"replacements"`
	Chosen       employees.Employee   `json:"chosen"`
}
