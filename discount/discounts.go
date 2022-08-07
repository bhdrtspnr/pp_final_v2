package discount

import (
	config "final_project/config"
	logger "final_project/logger"
	model "final_project/model"
	"strconv"
	"time"
)

/*If the customer made purchase which is more than given amount in a month
then all subsequent purchases should have %10 off.*/
func CalculateConsecutivePurchaseDiscount(cartid string) float64 {
	logger.AppLogger.Info().Println("Function hit : CalculateConsecutivePurchaseDiscount")
	cart := model.GetCart(cartid)
	owner_id := model.GetCartOwnerId(cart)
	customer := model.GetCustomer(owner_id)
	model.UpdateCartPrice(cartid)
	if customer.Has_subsequent_discount_until.After(time.Now()) {
		discountAmount := cart.TotalPrice * config.ConfigInstance.SubsequentPurchaseDiscount
		logger.AppLogger.Info().Printf("CalculateConsecutivePurchaseDiscount: Customer id: %d has consecutive purchase discount, discount amount: %f \n", customer.Id, discountAmount)
		return discountAmount
	}
	logger.AppLogger.Info().Printf("CalculateConsecutivePurchaseDiscount: Customer id: %d has no consecutive purchase discount \n", customer.Id)
	return 0
}

/* Every fourth order whose total is more than given amount may have discount
   depending on products. Products whose VAT is %1 don’t have any discount
   but products whose VAT is %8 and %18 have discount of %10 and %15
   respectively. */
func CalculateGivenAmountDiscount(cartid string) float64 {
	logger.AppLogger.Info().Println("Function hit : CalculateGivenAmountDiscount")
	cart := model.GetCart(cartid)
	if cart.TotalPrice < config.ConfigInstance.GivenAmount {
		logger.AppLogger.Info().Printf("CalculateGivenAmountDiscount: Cart total price: %f is less than given amount not applying any discount: %f \n", cart.TotalPrice, config.ConfigInstance.GivenAmount)
		return 0
	}
	logger.AppLogger.Info().Printf("CalculateGivenAmountDiscount: Cart total price: %f is greater than given amount: %f \n", cart.TotalPrice, config.ConfigInstance.GivenAmount)
	owner := model.GetCustomer(model.GetCartOwnerId(cart))
	if owner.Consecutive_discount == 0 {
		logger.AppLogger.Info().Printf("CalculateGivenAmountDiscount: Customer id: %d has no consecutive purchase discount \n", owner.Id)
		return 0
	}
	if (owner.Consecutive_discount%4 + 1) != 4 { //if customer had 3 previous orders that satisfies the given amount
		logger.AppLogger.Info().Printf("CalculateGivenAmountDiscount: Customer consecutive discount: %d is not divisible by 4 not applying any discount. \n", owner.Consecutive_discount)
		return 0
	}
	totalDiscount := 0.0

	logger.AppLogger.Info().Printf("CalculateGivenAmountDiscount: Customer consecutive discount: %d is divisible by 4 \n", owner.Consecutive_discount)

	cart_items := model.GetCartItems(cartid)

	for _, cartItem := range cart_items {
		product := model.GetProductById(strconv.Itoa(cartItem.ProductId))
		switch product.Vat {
		case 0.18:
			totalDiscount += product.Price * config.ConfigInstance.Point18VatDiscount
			logger.AppLogger.Info().Printf("CalculateGivenAmountDiscount: Product %s has been discounted for 18%% vat total discount applied: %f \n", product.Name, totalDiscount)
		case 0.08:
			totalDiscount += product.Price * config.ConfigInstance.Point8VatDiscount
			logger.AppLogger.Info().Printf("CalculateGivenAmountDiscount: Product %s has been discounted for 8%% vat total discount applied: %f \n", product.Name, totalDiscount)
		case 0.01:
			totalDiscount += product.Price * config.ConfigInstance.Point1VatDiscount
			logger.AppLogger.Info().Printf("CalculateGivenAmountDiscount: Product %s has been discounted for 1%% vat total discount applied: %f \n", product.Name, totalDiscount)
		default:
			logger.AppLogger.Error().Printf("CalculateGivenAmountDiscount: Product vat: %f is not supported \n", product.Vat)

		}

	}
	return totalDiscount
}

/*If there are more than 3 items of the same product, then fourth and
subsequent ones would have %8 off.*/
func CalculateThreeSubsequentPurchaseDiscount(cartid string) float64 {
	logger.AppLogger.Info().Println("Function hit : CalculateThreeSubsequentPurchaseDiscount")
	cart_items := model.GetCartItems(cartid)
	cartItemCounts := map[string]int{}
	totalDiscount := 0.0

	for _, cartItem := range cart_items {
		cartItemCounts[strconv.Itoa(cartItem.ProductId)]++
	}

	keysSlice := make([]string, 0, len(cartItemCounts))
	for k := range cartItemCounts {
		keysSlice = append(keysSlice, k)
	}

	for index, key := range keysSlice {
		for cartItemCounts[key] >= 3 { //for each product that has been purchased 3 or more times, iterate every 4th other product
			totalDiscount += model.GetProductById(keysSlice[index]).Price * config.ConfigInstance.ThreeSubsequentPurchaseDiscount
			cartItemCounts[key] -= 3
			logger.AppLogger.Info().Printf("CalculateThreeSubsequentPurchaseDiscount: Item %s has been discounted for 3 consecutive purchases \n", key)
			logger.AppLogger.Info().Printf("CalculateThreeSubsequentPurchaseDiscount: Item %s has %d remaining \n", key, cartItemCounts[key])
		}
	}
	return totalDiscount

}
