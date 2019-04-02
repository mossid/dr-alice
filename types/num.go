package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// I'm just using cosmos-sdk dec implementation
// might be changed
type Dec = sdk.Dec

func NewDec(i int64) Dec {
	return sdk.NewDec(i)
}
