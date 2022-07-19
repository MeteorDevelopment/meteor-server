package api

import (
	"context"
	"encoding/json"
	"github.com/plutov/paypal/v4"
	"github.com/segmentio/ksuid"
	"meteor-server/pkg/core"
	"meteor-server/pkg/db"
	"net/http"
	"strconv"
)

var client *paypal.Client
var orders = make(map[string]ksuid.KSUID)

type WebhookResponse struct {
	Resource struct {
		Id string `json:"id"`
	} `json:"resource"`
	EventType string `json:"event_type"`
}

func InitPayPal() {
	var err error
	client, err = paypal.NewClient(core.GetPrivateConfig().PayPalClientID, core.GetPrivateConfig().PayPalSecret, paypal.APIBaseLive)
	if err != nil {
		println("Failed to log in to paypal.")
		return
	}

	_, err = client.GetAccessToken(context.Background())
	if err != nil {
		println("Failed to generate PayPal access token.")
	}
}

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	amount, err := strconv.ParseFloat(r.URL.Query().Get("amount"), 32)

	if err != nil || amount < 5 {
		core.JsonError(w, core.J{"error": "Invalid order amount."})
		return
	}

	user, err := db.GetAccount(r)
	if err != nil {
		core.JsonError(w, core.J{"error": "Could not authorize account."})
		return
	}

	order, err := client.CreateOrder(
		context.Background(),
		paypal.OrderIntentCapture,
		[]paypal.PurchaseUnitRequest{
			{
				Amount: &paypal.PurchaseUnitAmount{
					Value:    strconv.FormatFloat(amount, 'f', 2, 32),
					Currency: "GBP",
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
		core.JsonError(w, core.J{"error": "Failed to create order."})
		println(err)
		return
	}

	orders[order.ID] = user.ID
	core.Json(w, order)
}

func ConfirmOrderHandler(w http.ResponseWriter, r *http.Request) {
	valid, err := client.VerifyWebhookSignature(r.Context(), r, core.GetPrivateConfig().PayPalWebhookId)
	if err != nil || valid.VerificationStatus == "FAILURE" {
		core.JsonError(w, core.J{"error": "Could not validate webhook signature"})
		println(err)
		return
	}

	var res WebhookResponse
	err = json.NewDecoder(r.Body).Decode(&res)
	if err != nil {
		println(err)
		return
	}

	acc, err := db.GetAccountId(orders[res.Resource.Id])
	if err != nil {
		println(err)
		return
	}

	acc.GiveDonator()

	core.Json(w, core.J{})
}

func CancelOrderHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if id != "" {
		delete(orders, id)
		core.Json(w, core.J{"success": "Order cancelled."})
	} else {
		core.Json(w, core.J{"error": "No active orders!"})
	}
}
