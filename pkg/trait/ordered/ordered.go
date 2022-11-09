package ordered

import (
	"bstrees/pkg/trait/number"
)

type Ordered interface {
	number.Number | ~string
}
