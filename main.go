package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/IsaTippens/kucoin-account/config"
	"github.com/IsaTippens/kucoin-account/kucoin"
)

func main() {
	if !config.FileExists("config.yml") {
		config.CreateConfigFile()
		fmt.Println("Config file created. Please edit config.yml and restart the program.")
		os.Exit(0)
	}

	cfg, err := config.LoadConfig("config")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	client := kucoin.NewClient(cfg.ApiKey, cfg.ApiSecret, cfg.ApiPassphrase)

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Value:    "config.yml",
				Usage:    "config file to use",
				Required: false,
				Action: func(cCtx *cli.Context, value string) error {
					cfg, err := config.LoadConfig(value)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					client = kucoin.NewClient(cfg.ApiKey, cfg.ApiSecret, cfg.ApiPassphrase)
					return nil

				},
			},
		},

		Commands: []*cli.Command{
			{
				Name:    "accounts",
				Aliases: []string{"a"},
				Usage:   "List all accounts",
				Action: func(cCtx *cli.Context) error {
					acc, err := client.GetAccounts()
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					accModel := kucoin.AccountsModel{}
					if err := acc.Unmarshal(&accModel); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					for _, a := range accModel {
						fmt.Printf("Type: %s, Currency: %s, Balance: %s, Available: %s\n", a.Type, a.Currency, a.Balance, a.Available)
					}
					return nil
				},
			},
			{
				Name:    "orders",
				Aliases: []string{"o"},
				Usage:   "Fetch Orders",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "active",
						Usage: "fetch active orders",
					},
					&cli.StringFlag{
						Name:     "symbol",
						Value:    "",
						Usage:    "coin pair to fetch orders for",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "startDate",
						Value:    "",
						Usage:    "start date for orders",
						Required: false,
					},
					&cli.StringFlag{
						Name:     "endDate",
						Value:    "",
						Usage:    "end date for orders",
						Required: false,
					},
				},
				Action: func(cCtx *cli.Context) error {
					active := cCtx.Bool("active")
					symbol := cCtx.String("symbol")
					startDate := cCtx.String("startDate")
					endDate := cCtx.String("endDate")

					or := kucoin.OrdersRequest{
						Symbol: symbol,
					}

					// parse dates
					date_parse := func(s string) time.Time {
						// yyyy-mm-dd hh:mm:ss.000
						var t time.Time
						var err error
						if len(s) == 10 {
							t, err = time.Parse("2006-01-02", s)
						} else if len(s) == 19 {
							t, err = time.Parse("2006-01-02 15:04:05", s)
						} else {
							t, err = time.Parse("2006-01-02 15:04:05.000", s)
						}
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						return t
					}
					if startDate != "" {
						or.StartDate = date_parse(startDate)
					}
					if endDate != "" {
						or.EndDate = date_parse(endDate)
					}

					if active {
						orders, err := client.GetActiveHFOrders(or)
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						ordersModel := kucoin.ActiveHFOrders{}
						if err := orders.Unmarshal(&ordersModel); err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						if len(ordersModel) == 0 {
							fmt.Println("No active orders")
							os.Exit(0)
						}
						for _, o := range ordersModel {
							tyd := o.CreatedAt
							timestamp := time.Unix(tyd/1000, tyd%1000*1000000)
							fmt.Printf("ID: %s Symbol: %s Side: %s Price: %s Size: %s Type: %s Time: %s\n", o.Id, o.Symbol, o.Side, o.Price, o.Size, o.Type, timestamp.Format("15:04:05.999 02/01/2006"))
						}
					} else {
						orders, err := client.GetFilledHFOrders(or)
						if err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						ordersModel := kucoin.FilledHFOrder{}

						if err := orders.Unmarshal(&ordersModel); err != nil {
							fmt.Println(err)
							os.Exit(1)
						}
						if len(ordersModel.Items) == 0 {
							fmt.Println("No filled orders")
							os.Exit(0)
						}
						for _, o := range ordersModel.Items {
							tyd := o.CreatedAt
							// convert unix time to hh:mm:ss.000 dd/mm/yyyy
							timestamp := time.Unix(tyd/1000, tyd%1000*1000000)
							dealFunds := o.DealFunds
							dealSize := o.DealSize
							// convert to float
							dealFundsFloat, _ := strconv.ParseFloat(dealFunds, 64)
							dealSizeFloat, _ := strconv.ParseFloat(dealSize, 64)
							price := dealFundsFloat / dealSizeFloat
							fmt.Printf("ID: %s Symbol: %s Side: %s Price: %.8f Size: %s Type: %s Time: %s\n", o.Id, o.Symbol, o.Side, price, o.Size, o.Type, timestamp.Format("15:04:05.999 02/01/2006"))
						}
					}
					return nil
				},
			},
			{
				Name: "fills",
				Action: func(cCtx *cli.Context) error {
					fills, err := client.GetFills()
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					// unmarschal to string
					j, err := json.Marshal(fills.Data)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					fmt.Println(string(j))
					return nil
				},
			},
			{
				Name:     "transfer",
				Aliases:  []string{"t"},
				Usage:    "transfer funds between accounts",
				HelpName: "help",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "from",
						Value:    "",
						Usage:    "accounts to transfer from - main, trade, trade_hf",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "to",
						Value:    "",
						Usage:    "accounts to transfer to  - main, trade, trade_hf",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "currency",
						Value:    "",
						Usage:    "currency to transfer",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "amount",
						Value:    "",
						Usage:    "amount to transfer",
						Required: true,
					},
				},
				Action: func(cCtx *cli.Context) error {
					tr := kucoin.TransferRequest{
						Currency: cCtx.String("currency"),
						Amount:   cCtx.Float64("amount"),
						From:     cCtx.String("from"),
						To:       cCtx.String("to"),
					}

					resp, err := client.Transfer(tr)
					if err != nil {
						fmt.Println(err)
						os.Exit(1)
					}
					trModel := kucoin.TransferModel{}
					if err := resp.Unmarshal(&trModel); err != nil {
						fmt.Println(err)
						os.Exit(1)
					}

					fmt.Printf("OrderID %s. Transfered %f %s from %s to %s\n", trModel.OrderId, tr.Amount, tr.Currency, tr.From, tr.To)
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
