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

	return o[i].insertedAt.Before(o[j].insertedAt)
}

// OrderList is a list of orders with functions wrapping it
type OrderList struct {
	orders orders
}

// Len is the length of the internal orders slice
func (ol OrderList) Len() int { return ol.orders.Len() }

// Swap allows us to swap orders
func (ol OrderList) Swap(i, j int) { ol.orders.Swap(i, j) }

// Less implements a less than check
func (ol OrderList) Less(i, j int) bool { return ol.orders.Less(i, j) }

const displayOrders = `
{{- range . -}}
{{- .Display -}}
{{- end -}}
`

// Display shows the current bids
func (ol *OrderList) Display() string {

	var buf bytes.Buffer

	t := template.Must(template.New("orders").Parse(displayOrders))
	err := t.Execute(&buf, ol.orders)
	if err != nil {
		panic(err)
	}

	return buf.String()
}

// GetBest gets the first order from the orders
func (ol *OrderList) GetBest() (order *Order) {
	order = ol.orders[0]
	return
}

// Insert puts a new order into the list
func (ol *OrderList) Insert(order *Order) {
	order.insertedAt = time.Now()
	ol.orders = append(ol.orders, order)
	sort.Sort(ol)
}

// Remove takes an order out of the list
func (ol *OrderList) Remove(orderID string) (order *Order) {
	for i := 0; i < ol.orders.Len(); i++ {
		if ol.orders[i].ID == orderID {
			order = ol.orders[i]
			ol.orders = append(ol.orders[:i], ol.orders[i+1:]...)
			return
		}
	}

	return
}
