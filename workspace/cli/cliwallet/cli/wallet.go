package cli

import (
	"cliwallet/internal/logic"
	"cliwallet/internal/model"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	Cli.AddCommand(cliWallet)

	cliWallet.AddCommand(cliWalletCommandShow)
	cliWallet.AddCommand(cliWalletCommandDeposit)
	cliWallet.AddCommand(cliWalletCommandWithdraw)
}

var cliWallet = &cobra.Command{
	Use:   "wallet",
	Short: "Wallet cli management commands",
}

var cliWalletCommandShow = &cobra.Command{
	Use:   "show",
	Short: "Show wallet information",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		println("Wallet information: <show>")
		return nil
	},
}

var cliWalletCommandDeposit = &cobra.Command{
	Use:   "deposit <amount>",
	Short: "Deposit funds into wallet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// TODO: Add in a busine logic

		amount := args[0]
		money, err := logic.ParseMoney(amount, model.EUR)
		if err != nil {
			return err
		}

		if money.Amount <= 0 {
			return logic.ErOperationWithZeroAmount
		}

		fmt.Printf("Deposited %+v\n", money)
		return nil
	},
}

var cliWalletCommandWithdraw = &cobra.Command{
	Use:   "withdraw <amount>",
	Short: "Withdraw funds from wallet",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// TODO: Add in a busine logic

		amount := args[0]
		money, err := logic.ParseMoney(amount, model.EUR)
		if err != nil {
			return err
		}
		if money.Amount <= 0 {
			return logic.ErOperationWithZeroAmount
		}

		fmt.Printf("Withdrawed %+v\n", money)
		return nil
	},
}
