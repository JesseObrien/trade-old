package orders

import (
	"bytes"
	"html/template"
	"sort"
	"time"
)

type orders []*Order

func (o orders) Len() int      { return len(o) }
func (o orders) Swap(i, j int) { o[i], o[j] = o[j], o[i] }
func (o orders) Less(i, j int) bool {
	switch o[i].Price.Cmp(o[j].Price) {
	case 1:
		return true
	case -1:
		return false
	}

	// FIFO requires that we reverse sort by time, we want the first in to be the first out
	return o[i].insertedAt.After(o[j].insertedAt)
}

// OrderList is a list of orders with functions wrapping it
type OrderList struct {
	Orders orders `json:"Orders"`
}

// Len is the length of the internal orders slice
func (ol OrderList) Len() int { return ol.Orders.Len() }

// Swap allows us to swap orders
func (ol OrderList) Swap(i, j int) { ol.Orders.Swap(i, j) }

// Less implements a less than check
func (ol OrderList) Less(i, j int) bool { return ol.Orders.Less(i, j) }

const displayOrders = `
{{- range . -}}
{{- .Display -}}
{{- end -}}
`

// Display shows the current bids
func (ol *OrderList) Display() string {

	var buf bytes.Buffer

	t := template.Must(template.New("orders").Parse(displayOrders))
	err := t.Execute(&buf, ol.Orders)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

// GetBest gets the first order from the orders
func (ol *OrderList) GetMax() (order *Order) {
	order = ol.Orders[0]
	return
}

func (ol *OrderList) GetMin() (order *Order) {
	order = ol.Orders[len(ol.Orders)-1]
	return
}

// Insert puts a new order into the list
func (ol *OrderList) Insert(order *Order) {
	order.insertedAt = time.Now()
	ol.Orders = append(ol.Orders, order)
	sort.Sort(ol)
}

// Remove takes an order out of the list
func (ol *OrderList) Remove(orderID string) (order *Order) {
	for i := 0; i < ol.Orders.Len(); i++ {
		if ol.Orders[i].ID == orderID {
			order = ol.Orders[i]
			ol.Orders = append(ol.Orders[:i], ol.Orders[i+1:]...)
			return
		}
	}

	return
}
