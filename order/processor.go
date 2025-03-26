package order

import (
	"errors"
	"fmt"
	"log"
	"web-example/database"
	"web-example/product"
	"web-example/util"
)

type CreateMessage struct {
	Order       *Order
	ErrResponse chan error
	IdResponse  chan int
}

func ProcessOrder(queue chan *CreateMessage, orderStore Repository, productStore product.Repository, txService database.Transactional) {
	log.Println("Started order processing...")
processLoop:
	for msg := range queue {
		order := msg.Order
		log.Println("Processing order: ", order)

		products, err := productStore.FindAllByIds(order.AllProductIds())
		if err != nil {
			log.Printf("Error finding products: %v", err)
			msg.ErrResponse <- util.NewInternalError()
			continue
		}

		if err := validateOrderEligibility(products, order); err != nil {
			msg.ErrResponse <- err
			continue
		}

		tx := txService.BeginTransaction()
		if tx.Error != nil {
			log.Printf("Error starting transaction: %v", tx.Error)
			msg.ErrResponse <- util.NewInternalError()
			continue
		}

		id, err := orderStore.Create(order, tx)
		if err != nil {
			log.Printf("Error creating order: %v", err)
			rollbackTransaction(tx)
			msg.ErrResponse <- util.NewInternalError()
			continue
		}
		log.Printf("Created order: %v", id)

		for _, prod := range products {
			prod.Quantity -= order.FindOrderProductByName(prod.Name).RequestedQuantity
			err := productStore.UpdateQuantity(prod, tx)
			if err != nil {
				log.Printf("Error updating product: %v", err)
				rollbackTransaction(tx)
				var transactionError database.TransactionError
				if errors.As(err, &transactionError) {
					msg.ErrResponse <- errors.New("Product quantity has changed, please try again: " + prod.Name)
				}
				continue processLoop
			}
		}

		result := tx.Commit()
		if result.Error != nil {
			log.Printf("Error committing transaction: %v", result.Error)
			rollbackTransaction(tx)
			msg.ErrResponse <- util.NewInternalError()
			continue
		}

		msg.IdResponse <- order.ID
	}
}

func validateOrderEligibility(products []*product.Product, order *Order) error {
	for _, p := range products {
		if !checkProductQuantity(p, order.FindOrderProductByName(p.Name)) {
			return fmt.Errorf("There are not enough product '%s' to place the order, max item(s): %d", p.Name, p.Quantity)
		}
	}
	return nil
}

func checkProductQuantity(product *product.Product, orderProduct *OrderProduct) bool {
	if product.Quantity >= orderProduct.RequestedQuantity {
		return true
	}
	return false
}
