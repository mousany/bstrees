package ordered

import (
	"bstrees/pkg/type/number"
)

type Ordered interface {
	number.Number | ~string
}
