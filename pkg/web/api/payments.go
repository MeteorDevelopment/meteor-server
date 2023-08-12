package api

import (
	"context"
	"encoding/json"
	"github.com/plutov/paypal/v4"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/ksuid"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
	"net/http"
	"strconv"
)

type WebhookResponse struct {
	Resource struct {
		Id string `json:"id"`
	} `json:"resource"`
	EventType string `json:"event_type"`
}

var client *paypal.Client
var orders = make(map[string]ksuid.KSUID)
var amounts = make(map[string]float64)

func InitPayPal() {
	var err error
	client, err = paypal.NewClient(core.GetPrivateConfig().PayPalClientID, core.GetPrivateConfig().PayPalSecret, paypal.APIBaseLive)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to log in to PayPal")
		return
	}

	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to generate PayPal access token")
	}
}

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 32)

	if err != nil || amount < 5 {
		core.JsonError(w, "Invalid order amount.")
		return
	}

	user, err := db.GetAccount(r)
	if err != nil {
		core.JsonError(w, "Could not authorize account.")
		return
	}

	order, err := client.CreateOrder(
		context.Background(),
		paypal.OrderIntentCapture,
		[]paypal.PurchaseUnitRequest{
			{
				Amount: &paypal.PurchaseUnitAmount{
					Value:    strconv.FormatFloat(amount, 'f', 2, 32),
					Currency: "EUR",
				},
			},
		},
		&paypal.CreateOrderPayer{
			Name: &paypal.CreateOrderPayerName{
				GivenName: user.Username,
			},
			EmailAddress: user.Email,
		},
		nil,
	)

	if err != nil {
		msg := "Failed to create order"
		core.JsonError(w, msg)
		log.Err(err).Msg(msg)
		return
	}

	orders[order.ID] = user.ID
	amounts[order.ID] = amount

	core.Json(w, order)
}

func ConfirmOrderHandler(w http.ResponseWriter, r *http.Request) {
	valid, err := client.VerifyWebhookSignature(r.Context(), r, core.GetPrivateConfig().PayPalWebhookId)
	if err != nil || valid.VerificationStatus == "FAILURE" {
		msg := "Could not validate webhook signature"
		core.JsonError(w, msg)
		log.Err(err).Msg(msg)
		return
	}

	var res WebhookResponse
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		log.Err(err).Msg("Failed to decode request")
		return
	}

	acc, err := db.GetAccountId(orders[res.Resource.Id])
	if err != nil {
		log.Err(err).Msg("Failed to get account by ID")
		return
	}

	acc.GiveDonator(amounts[res.Resource.Id])

	core.Json(w, core.J{})
}

func CancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id != "" {
		delete(orders, id)
		delete(amounts, id)
		core.Json(w, core.J{})
	} else {
		core.JsonError(w, "No active orders.")
	}
}
