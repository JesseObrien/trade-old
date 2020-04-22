package orders

import (
	"fmt"
	"sort"
	"strings"
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

type OrderList struct {
	orders orders
}

func (ol OrderList) Len() int           { return ol.orders.Len() }
func (ol OrderList) Swap(i, j int)      { ol.orders.Swap(i, j) }
func (ol OrderList) Less(i, j int) bool { return ol.orders.Less(i, j) }

func (ol *OrderList) Display() string {
	bids := strings.Builder{}
	for _, bid := range ol.orders {
		bids.WriteString(fmt.Sprintf("Price: $%s | Quantity: %s", bid.Price.StringFixed(2), bid.OpenQuantity().String()))
	}

	return bids.String()
}
func (ol *OrderList) Iter(fn func(*Order)) {
	for _, order := range ol.orders {
		fn(order)
	}
}

func (ol *OrderList) GetBest() (order *Order) {
	order = ol.orders[0]
	return
}

func (ol *OrderList) Insert(order *Order) {
	order.insertedAt = time.Now()
	ol.orders = append(ol.orders, order)
	sort.Sort(ol)
}

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
