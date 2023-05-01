package cli_cmd

import (
	"fmt"
	flag "github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"log"
	"strings"
)

var (
	cliRepository    string
	cliMarketName    string
	cliMarketSection string
	cliSymbolCode    string
	cliPath          string
	cliIsFile        bool
)

func InitRootCmd() *flag.Command {
	rootCmd := flag.Command{
		Use:     "market-trades-loader",
		Version: "1.0",
		Run: func(cmd *flag.Command, args []string) {

		},
	}

	rootCmd.Flags().SetNormalizeFunc(wordSepNormalizeFunc)

	rootCmd.Flags().StringVarP(&cliMarketName, "market", "m", "", "the name of a market, that has produced trades data")
	rootCmd.Flags().StringVarP(&cliMarketSection, "section", "s", "", "the market section")
	rootCmd.Flags().StringVarP(&cliSymbolCode, "symbol", "", "", "the symbol code, which trades should be imported in database")
	rootCmd.Flags().StringVarP(&cliPath, "path", "p", "", "the path to a csv file with trades data or to the directory with csv files")
	rootCmd.Flags().BoolVar(&cliIsFile, "isFile", true, "if true passed path is full path to the scv file")
	rootCmd.Flags().StringVarP(&cliRepository, "repository", "r", "mysql", "used to chose repository to collect trades")

	err := rootCmd.MarkFlagRequired("market")
	if err != nil {
		panic(fmt.Sprintf("can't make field \"market\" required | %v", err))
	}

	err = rootCmd.MarkFlagRequired("section")
	if err != nil {
		panic(fmt.Sprintf("can't make field \"section\" required | %v", err))
	}

	err = rootCmd.MarkFlagRequired("symbol")
	if err != nil {
		panic(fmt.Sprintf("can't make field \"symbol\" required | %v", err))
	}

	err = rootCmd.MarkFlagRequired("path")
	if err != nil {
		panic(fmt.Sprintf("can't make field \"path\" required | %v", err))
	}

	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(fmt.Sprintf("recuired flags are missed | %v", err))
	}

	return &rootCmd
}

func wordSepNormalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	from := []string{"-", "_"}
	to := "."
	for _, sep := range from {
		name = strings.Replace(name, sep, to, -1)
	}
	return pflag.NormalizedName(name)
}
