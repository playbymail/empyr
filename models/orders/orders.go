// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package orders

import (
	"fmt"
	"github.com/playbymail/empyr/models/units"
)

type Order interface {
	Execute() error
}

type TransferDetail struct {
	Unit     units.Unit
	Quantity int
}

func (td *TransferDetail) String() string {
	return fmt.Sprintf("{%d %s}", td.Quantity, td.Unit)
}
