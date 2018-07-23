package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"strconv"

	"github.com/insighted4/insighted-go/eventsourcing"
	"github.com/sirupsen/logrus"
)

// Order is an example of state generated from left fold of events
type Order struct {
	ID        string
	Version   int
	CreatedAt time.Time
	UpdatedAt time.Time
	State     string
}

// OrderCreated event used a marker of order created
type OrderCreated struct {
	eventsourcing.Model
}

// OrderShipped event used a marker of order shipped
type OrderShipped struct {
	eventsourcing.Model
}

// On implements Aggregate interface
func (item *Order) On(event eventsourcing.Event) error {
	switch v := event.(type) {
	case *OrderCreated:
		item.State = "created"

	case *OrderShipped:
		item.State = "shipped"

	default:
		return fmt.Errorf("unable to handle event, %v", v)
	}

	item.Version = event.EventVersion()
	item.ID = event.AggregateID()
	item.UpdatedAt = event.EventAt()

	return nil
}

// CreateOrder command
type CreateOrder struct {
	eventsourcing.CommandModel
}

// ShipOrder command
type ShipOrder struct {
	eventsourcing.CommandModel
}

// Apply implements the CommandHandler interface
func (item *Order) Apply(ctx context.Context, command eventsourcing.Command) ([]eventsourcing.Event, error) {
	switch v := command.(type) {
	case *CreateOrder:
		orderCreated := &OrderCreated{
			Model: eventsourcing.Model{ID: command.AggregateID(), Version: item.Version + 1, At: time.Now()},
		}
		return []eventsourcing.Event{orderCreated}, nil

	case *ShipOrder:
		if item.State != "created" {
			return nil, fmt.Errorf("order, %v, has already shipped", command.AggregateID())
		}
		orderShipped := &OrderShipped{
			Model: eventsourcing.Model{ID: command.AggregateID(), Version: item.Version + 1, At: time.Now()},
		}
		return []eventsourcing.Event{orderShipped}, nil

	default:
		return nil, fmt.Errorf("unhandled command, %v", v)
	}
}

func check(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	store := eventsourcing.NewMemStore()
	serializer := eventsourcing.NewJSONSerializer(
		OrderCreated{},
		OrderShipped{},
	)
	repo := eventsourcing.NewRepository(&Order{}, store, serializer, logrus.New())

	id := strconv.FormatInt(time.Now().UnixNano(), 36)
	ctx := context.Background()

	// Create Order
	create := &CreateOrder{
		CommandModel: eventsourcing.CommandModel{ID: id},
	}
	_, err := repo.Apply(ctx, create)
	check(err)

	// Ship Order
	ship := &ShipOrder{
		CommandModel: eventsourcing.CommandModel{ID: id},
	}
	_, err = repo.Apply(ctx, ship)
	check(err)

	// Aggregate
	aggregate, err := repo.Load(ctx, id)
	check(err)

	found := aggregate.(*Order)
	fmt.Printf("Order %v [version %v] %v %v\n", found.ID, found.Version, found.State, found.UpdatedAt.Format(time.RFC822))
}
