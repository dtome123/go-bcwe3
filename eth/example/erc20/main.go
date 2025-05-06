package main

import (
	"fmt"
	"math/big"

	"github.com/dtome123/go-bcwe3/eth"
)

func main() {

	eth := eth.NewEth("wss://holesky.infura.io/ws/v3/da05d3dc31244bd483a28d746233d32f")

	defer eth.Provider.Close()

	// info, err := eth.ERC20.GetInfo(context.Background(), "0xA73aAE60B392d2E46d9693851bFcA872a9c54635")

	// if err != nil {
	// 	panic(err)
	// }

	// isMatch, _ := eth.ERC20.IsPossiblyERC20(context.Background(), "0xA73aAE60B392d2E46d9693851bFcA872a9c54635")

	// fmt.Println("name: ", info.Name)
	// fmt.Println("symbol: ", info.Symbol)
	// fmt.Println("decimals: ", info.Decimals)
	// fmt.Println("totalSupply: ", info.TotalSupply)

	// fmt.Println("is erc20: ", isMatch)

	cmd, err := eth.ERC20.NewCmd("0x55d2EC94ffc9f7A2042317022Af4B758D5A1Dc36")

	if err != nil {
		panic(err)
	}
	wei := new(big.Int)
	wei.SetString("10000000000000000000", 10)
	fmt.Println(cmd.BalanceOf("0x7556989c2A60E60F0c66A2b9D77079BC9F189037"))
	tx, err := cmd.Transfer("0xa84e540D1eb5458DFC2bC25760bD64fbECb8e345", wei, "private")

	if err != nil {
		panic(err)
	}

	fmt.Println(tx)
}
