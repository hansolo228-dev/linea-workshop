package main

import (
	"linea/workshop/blockchain"
	"log"
	"net/http"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

func (h *Header) Render() app.UI {
	return app.Div().Body(
		app.H1().Body(
			app.Text("Connect Blockchain"),
		).Style("width", "10em").Style("margin-left", "auto").Style("margin-right", "auto"),
		app.P().Body(
			app.Input().Style("width", "500px").
				Type("text").
				Value(h.Address).
				Placeholder("Address blockchain").
				AutoFocus(true).
				OnChange(h.ValueTo(&h.Address)),
		).Style("width", "10em").Style("margin-right", "auto"),

		// Block Number
		app.H1().Body(
			app.Text("Request BlockNumber"),
		).Style("width", "10em").Style("margin-left", "auto").Style("margin-right", "auto"),
		app.Div().Body(
			app.Div().Body(
				app.Button().Text("BlockNumber").OnClick(h.GetBlockNumber).Style("float", "right").Style("height", "200px").Style("width", "500px"),
				app.Div().Body(
					app.H2().Text("Block Number:"),
					app.If(h.BlockNumber != "",
						app.H2().Text(h.BlockNumber),
					).Else(
						app.H2().Text("0"),
					)).Style("overflow", "hidden").Style("height", "200px").Style("width", "300px"),
			),
		).Style("background-color", "deepskyblue").Style("display", "table").Style("clear", "both").Style("display", "block").Style("height", "200px"),

		// Transaction
		app.H1().Body(
			app.Text("Block Transaction"),
		).Style("width", "10em").Style("margin-left", "auto").Style("margin-right", "auto"),
		app.Div().Body(
			app.Div().Body(
				app.Button().Text("Transaction").OnClick(h.GetTransaction).Style("float", "right").Style("height", "200px").Style("width", "500px"),
				app.Div().Body(
					app.H2().Text("Transaction: \n"),
					app.If(h.Transaction.Hash != "",
						app.H2().Text("Hash transaction:"),
						app.H2().Text(h.Transaction.Hash),
					).Else(
						app.H2().Text("None"),
					)).Style("overflow", "hidden").Style("height", "200px"),
			),
		).Style("background-color", "deepskyblue").Style("display", "table").Style("clear", "both").Style("display", "block").Style("height", "200px"),
	)
}

type Header struct {
	app.Compo
	Address     string
	BlockNumber string
	Name        string
	Transaction Transaction
	TxStatus    int
	ChainId     string
}
type Transaction struct {
	Hash     string
	Value    string
	Gas      uint64
	GasPrice uint64
	Nonce    uint64
	Data     []byte
	To       string
}

func (header *Header) GetBlockNumber(ctx app.Context, e app.Event) {
	header.BlockNumber = blockchain.BlockNumber(header.Address)
}
func (header *Header) GetTransaction(ctx app.Context, e app.Event) {
	if header.BlockNumber != "" {
		tx, status := blockchain.Transaction(header.BlockNumber, header.Address)
		header.TxStatus = status
		header.Transaction.Data = tx.Data()
		header.Transaction.To = tx.To().Hex()
		header.Transaction.Nonce = tx.Nonce()
		header.Transaction.GasPrice = tx.GasPrice().Uint64()
		header.Transaction.Value = tx.Value().String()
		header.Transaction.Hash = tx.Hash().Hex()
		header.Transaction.Gas = tx.Gas()
	}
}
func (header *Header) ChainID(ctx app.Context, e app.Event) {
	if header.ChainId != blockchain.ChainID(header.Address) {
		header.ChainId = blockchain.ChainID(header.Address)
	} else {
		header.ChainId = "None"
	}
}

func main() {
	// Components routing:
	app.Route("/", &Header{})
	app.RunWhenOnBrowser()

	// HTTP routing:
	http.Handle("/", &app.Handler{
		Name:        "Hello",
		Description: "An Hello World! example",
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal(err)
	}
}
