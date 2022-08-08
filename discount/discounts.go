package discount

import (
	config "final_project/config"
	logger "final_project/logger"
	model "final_project/model"
	"strconv"
	"time"
)

/*
	If the customer made purchase which is more than given amount in a month
	then all subsequent purchases should have %10 off.

	How it works? When a customer makes a purchase that has a total price more than given amount,
	we set it's consecutive discount to today() + 30 days, so when we're discounting the next purchase,
	we check if the consecutive discount date is greater than today(). If it is, we discount the purchase.
*/

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
   respectively.

	How this works? We check if the total price of the cart is more than given amount.
	If it is, we check if the VAT of the each cart item.
	Depending on the VAT values we apply different discounts.

*/
func CalculateGivenAmountDiscount(cartid string) float64 {
	logger.AppLogger.Info().Println("Function hit : CalculateGivenAmountDiscount")
	cart := model.GetCart(cartid)                            //get the cart
	if cart.TotalPrice < config.ConfigInstance.GivenAmount { //if the total price of the cart is less than given amount, return 0
		logger.AppLogger.Info().Printf("CalculateGivenAmountDiscount: Cart total price: %f is less than given amount not applying any discount: %f \n", cart.TotalPrice, config.ConfigInstance.GivenAmount)
		return 0
	}
	logger.AppLogger.Info().Printf("CalculateGivenAmountDiscount: Cart total price: %f is greater than given amount: %f \n", cart.TotalPrice, config.ConfigInstance.GivenAmount)
	owner := model.GetCustomer(model.GetCartOwnerId(cart))
	if owner.Consecutive_discount == 0 { //if the customer has no consecutive purchase discount, return 0
		logger.AppLogger.Info().Printf("CalculateGivenAmountDiscount: Customer id: %d has no consecutive purchase discount \n", owner.Id)
		return 0
	}
	if (owner.Consecutive_discount%4 + 1) != 4 { //if customer had 3 previous orders that satisfies the given amount
		logger.AppLogger.Info().Printf("CalculateGivenAmountDiscount: Customer consecutive discount: %d is not divisible by 4 not applying any discount. \n", owner.Consecutive_discount)
		return 0
	}
	totalDiscount := 0.0

	logger.AppLogger.Info().Printf("CalculateGivenAmountDiscount: Customer consecutive discount: %d is divisible by 4 \n", owner.Consecutive_discount)

	cart_items := model.GetCartItems(cartid) //get the cart items

	for _, cartItem := range cart_items { //iterate every cart item
		product := model.GetProductById(strconv.Itoa(cartItem.ProductId))
		switch product.Vat { //depending on the VAT value, apply different discounts
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
	return totalDiscount //return the total discount
}

/*If there are more than 3 items of the same product, then fourth and
subsequent ones would have %8 off.

How this works? We check if the cart has more than 3 items of the same product.
If there are more than 3 items of the same product, we apply discount the next item
Then we keep iterating until we reach the end of the cart.
*/
func CalculateThreeSubsequentPurchaseDiscount(cartid string) float64 {
	logger.AppLogger.Info().Println("Function hit : CalculateThreeSubsequentPurchaseDiscount")
	cart_items := model.GetCartItems(cartid) //get the cart items
	cartItemCounts := map[string]int{}       //map to store the count of the cart items
	totalDiscount := 0.0

	for _, cartItem := range cart_items { //iterate every cart item
		cartItemCounts[strconv.Itoa(cartItem.ProductId)]++ //increment the count of the cart item
	}

	keysSlice := make([]string, 0, len(cartItemCounts)) //create a slice of keys
	for k := range cartItemCounts {
		keysSlice = append(keysSlice, k) //append the keys to the slice
	}

	for index, key := range keysSlice {
		for cartItemCounts[key] >= 3 { //for each product that has been purchased 3 or more times, iterate every 4th other product, I've used a for here because using an if will only catch the discount for once, if there are 6 items, it'll only apply discount once.
			totalDiscount += model.GetProductById(keysSlice[index]).Price * config.ConfigInstance.ThreeSubsequentPurchaseDiscount
			cartItemCounts[key] -= 3 //decrement the count of the cart item since we want to apply the discount NOT only once
			logger.AppLogger.Info().Printf("CalculateThreeSubsequentPurchaseDiscount: Item %s has been discounted for 3 consecutive purchases \n", key)
			logger.AppLogger.Info().Printf("CalculateThreeSubsequentPurchaseDiscount: Item %s has %d remaining \n", key, cartItemCounts[key])
		}
	}
	return totalDiscount //return the total discount

}
