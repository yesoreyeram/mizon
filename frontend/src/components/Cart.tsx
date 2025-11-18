"use client";

import { useEffect, useState } from "react";
import axios from "axios";

const API_CART = process.env.NEXT_PUBLIC_CART_API || "http://localhost:8004";
const API_ORDER = process.env.NEXT_PUBLIC_ORDER_API || "http://localhost:8005";

interface CartItem {
  product_id: string;
  name: string;
  price: number;
  quantity: number;
}

interface CartProps {
  onClose: () => void;
}

export default function Cart({ onClose }: CartProps) {
  const [items, setItems] = useState<CartItem[]>([]);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    loadCart();
  }, []);

  const loadCart = async () => {
    try {
      setIsLoading(true);
      const response = await axios.get(`${API_CART}/api/cart/admin`);
      setItems(response.data.items || []);
    } catch (error) {
      console.error("Error loading cart:", error);
      setItems([]);
    } finally {
      setIsLoading(false);
    }
  };

  const updateQuantity = async (productId: string, quantity: number) => {
    try {
      await axios.put(`${API_CART}/api/cart/admin/items/${productId}`, {
        quantity,
      });
      loadCart();
    } catch (error) {
      console.error("Error updating quantity:", error);
    }
  };

  const removeItem = async (productId: string) => {
    try {
      await axios.delete(`${API_CART}/api/cart/admin/items/${productId}`);
      loadCart();
    } catch (error) {
      console.error("Error removing item:", error);
    }
  };

  const placeOrder = async () => {
    if (items.length === 0) {
      alert("Cart is empty");
      return;
    }

    const total = items.reduce(
      (sum, item) => sum + item.price * item.quantity,
      0
    );

    try {
      await axios.post(`${API_ORDER}/api/orders`, {
        user_id: "admin",
        items: items,
        total: total,
      });

      await axios.delete(`${API_CART}/api/cart/admin`);
      alert("Order placed successfully! Redirecting to orders page...");
      loadCart();
      onClose();
      // Redirect to orders page
      window.location.href = "/orders";
    } catch (error) {
      console.error("Error placing order:", error);
      alert("Failed to place order");
    }
  };

  const total = items.reduce(
    (sum, item) => sum + item.price * item.quantity,
    0
  );

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex justify-end z-50">
      <div className="bg-white w-full max-w-md h-full overflow-y-auto shadow-xl">
        <div className="p-6">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-bold">Shopping Cart</h2>
            <button
              onClick={onClose}
              className="text-gray-500 hover:text-gray-700"
            >
              <svg
                className="w-6 h-6"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M6 18L18 6M6 6l12 12"
                />
              </svg>
            </button>
          </div>

          {isLoading ? (
            <div className="text-center py-8">
              <div className="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
          ) : items.length === 0 ? (
            <div className="text-center py-8">
              <p className="text-gray-600">Your cart is empty</p>
            </div>
          ) : (
            <>
              <div className="space-y-4 mb-6">
                {items.map((item) => (
                  <div
                    key={item.product_id}
                    className="flex items-center space-x-4 border-b pb-4"
                  >
                    <div className="flex-1">
                      <h3 className="font-semibold text-gray-800">
                        {item.name}
                      </h3>
                      <p className="text-gray-600">${item.price.toFixed(2)}</p>
                    </div>
                    <div className="flex items-center space-x-2">
                      <button
                        onClick={() =>
                          updateQuantity(item.product_id, item.quantity - 1)
                        }
                        className="bg-gray-200 hover:bg-gray-300 w-8 h-8 rounded flex items-center justify-center"
                      >
                        -
                      </button>
                      <span className="w-8 text-center">{item.quantity}</span>
                      <button
                        onClick={() =>
                          updateQuantity(item.product_id, item.quantity + 1)
                        }
                        className="bg-gray-200 hover:bg-gray-300 w-8 h-8 rounded flex items-center justify-center"
                      >
                        +
                      </button>
                    </div>
                    <button
                      onClick={() => removeItem(item.product_id)}
                      className="text-red-600 hover:text-red-800"
                    >
                      <svg
                        className="w-5 h-5"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                        />
                      </svg>
                    </button>
                  </div>
                ))}
              </div>

              <div className="border-t pt-4">
                <div className="flex items-center justify-between mb-4">
                  <span className="text-lg font-semibold">Total:</span>
                  <span className="text-2xl font-bold text-blue-600">
                    ${total.toFixed(2)}
                  </span>
                </div>
                <button
                  onClick={placeOrder}
                  className="w-full bg-blue-600 hover:bg-blue-700 text-white py-3 rounded-lg font-semibold"
                >
                  Place Order
                </button>
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
