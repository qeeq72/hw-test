package interfaces

import (
	"context"
	"fmt"
)

type ICollector interface {
	Run(context.Context, chan fmt.Stringer)
}
