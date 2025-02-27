import * as api from '@opentelemetry/api';
import { randomUUID } from "crypto";
import { CartItem, CheckoutEvent, UserCheckout } from "../model";
import getCustomerCartItems from "../cart/cartService";
import Payment from "../payment/paymentService";
import { logger } from "../logger";
import checkoutPublisher from "../messagging/publisher";

export default async function Checkout(checkout: UserCheckout) {
    const span = api.trace.getTracer('checkout-api').startSpan('checkout-api.Checkout');
    await api.context.with(api.trace.setSpan(api.ROOT_CONTEXT, span), async () => {
        try {
            const checkoutID = randomUUID().toString();
            const cart = await getCustomerCartItems(checkout.customer_id)
            const pay = await Payment(checkoutID, checkout, cart!!.items!!)
            logger.info(pay?.transactionId)

            const cartItems = cart?.items?.map<CartItem>((i) => {
                return {
                    item_id: i.itemId?.toString(),
                    price: i.price,
                    quantity: Number(i.quantity)
                } as CartItem
            })

            const checkoutEvent: CheckoutEvent = {
                transaction_id: pay.transactionId,
                user_checkout: checkout,
                customer_cart: {
                    customer_id: cart?.customerId!!,
                    items: cartItems,
                }
            }
            logger.info("checkout-event: ", checkoutEvent)
            await checkoutPublisher.Publish(checkoutEvent)
        } catch (error) {
            span.setStatus({
                code: api.SpanStatusCode.ERROR,
                message: (error as Error)?.message,
            })
            throw error
        } finally {
            span.end();
        }
    });
}