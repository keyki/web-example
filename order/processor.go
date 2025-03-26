package order

import (
	"context"
	"errors"
	"fmt"
	"web-example/database"
	"web-example/log"
	"web-example/product"
	"web-example/util"
)

type CreateMessage struct {
	Order       *Order
	ErrResponse chan error
	IdResponse  chan int
	Context     context.Context
}

func ProcessOrder(queue chan *CreateMessage, orderStore Repository,
	productStore product.Repository, txService database.Transactional) {
	log.BaseLogger().Info("Started order processing...")
processLoop:
	for msg := range queue {
		order := msg.Order
		ctx := msg.Context
		log.Logger(ctx).Info("Processing order: ", order)

		products, err := productStore.FindAllByIds(ctx, order.AllProductIds())
		if err != nil {
			log.Logger(ctx).Infof("Error finding products: %v", err)
			msg.ErrResponse <- util.NewInternalError()
			continue
		}

		if err := validateOrderEligibility(products, order); err != nil {
			msg.ErrResponse <- err
			continue
		}

		tx := txService.BeginTransaction()
		if tx.Error != nil {
			log.Logger(ctx).Infof("Error starting transaction: %v", tx.Error)
			msg.ErrResponse <- util.NewInternalError()
			continue
		}

		id, err := orderStore.Create(ctx, order, tx)
		if err != nil {
			log.Logger(ctx).Infof("Error creating order: %v", err)
			rollbackTransaction(ctx, tx)
			msg.ErrResponse <- util.NewInternalError()
			continue
		}
		log.Logger(ctx).Infof("Created order: %v", id)

		for _, prod := range products {
			prod.Quantity -= order.FindOrderProductByName(prod.Name).RequestedQuantity
			err := productStore.UpdateQuantity(ctx, prod, tx)
			if err != nil {
				log.Logger(ctx).Infof("Error updating product: %v", err)
				rollbackTransaction(ctx, tx)
				var transactionError database.TransactionError
				if errors.As(err, &transactionError) {
					msg.ErrResponse <- errors.New("Product quantity has changed, please try again: " + prod.Name)
				}
				continue processLoop
			}
		}

		result := tx.Commit()
		if result.Error != nil {
			log.Logger(ctx).Infof("Error committing transaction: %v", result.Error)
			rollbackTransaction(ctx, tx)
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
