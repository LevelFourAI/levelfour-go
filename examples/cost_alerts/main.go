package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/LevelFourAI/levelfour-go/levelfour"
)

// Creates a monthly budget alert, raises its budget, pauses and resumes it, then deletes it.
//
// LEVELFOUR_API_KEY must be a read-write key. ALERT_SOURCE_ID is the account or cost view the
// alert watches, an id from client.CostAlerts.ListSources, and ALERT_RECIPIENT is the address
// of a member of your organization.
func main() {
	ctx := context.Background()
	client, err := levelfour.NewClient("")
	if err != nil {
		log.Fatal(err)
	}

	created, err := client.CostAlerts.CreateAlert(ctx, &levelfour.CreateAlertRequest{
		Name:     "Production budget",
		SourceID: os.Getenv("ALERT_SOURCE_ID"),
		Params:   &levelfour.CreateAlertRequestParams{Budget: monthlyBudget(40000)},
		Recipients: []*levelfour.AlertRecipient{
			{Name: "Ana", Email: os.Getenv("ALERT_RECIPIENT"), Member: true},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	id := created.Data.ID
	fmt.Println("created", id)

	// params replaces the whole rule, so a new budget sends the rule back with that one figure changed.
	if _, err := client.CostAlerts.UpdateAlert(ctx, id, &levelfour.UpdateAlertRequest{
		Params: &levelfour.UpdateAlertRequestParams{Budget: monthlyBudget(50000)},
	}); err != nil {
		log.Fatal(err)
	}
	fmt.Println("raised the budget to 50000")

	for _, enabled := range []bool{false, true} {
		if _, err := client.CostAlerts.UpdateAlert(ctx, id, &levelfour.UpdateAlertRequest{Enabled: &enabled}); err != nil {
			log.Fatal(err)
		}
		fmt.Println("enabled:", enabled)
	}

	// Deleting takes the alert's history with it. Pausing keeps it.
	if _, err := client.CostAlerts.DeleteAlert(ctx, id); err != nil {
		log.Fatal(err)
	}
	fmt.Println("deleted", id)
}

// Fires when month-to-date spend passes 80 percent of the budget.
func monthlyBudget(amount float64) *levelfour.BudgetParamsInput {
	return &levelfour.BudgetParamsInput{
		Budget: &levelfour.BudgetParamsInputBudget{Double: amount},
		Reset:  "monthly",
		Rules: []*levelfour.BudgetRuleInput{
			{Amount: &levelfour.BudgetRuleInputAmount{Double: 80}, Unit: "percent"},
		},
	}
}
