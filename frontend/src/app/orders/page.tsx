"use client";

import { useEffect, useState } from "react";
import axios from "axios";
import Navbar from "@/components/Navbar";
import Cart from "@/components/Cart";

const API_ORDER = process.env.NEXT_PUBLIC_ORDER_API || "http://localhost:8005";

interface OrderItem {
  product_id: string;
  name: string;
  price: number;
  quantity: number;
}

interface Order {
  id: string;
  user_id: string;
  items: OrderItem[];
  total: number;
  status: string;
  created_at: string;
}

export default function OrdersPage() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
  const [searchQuery, setSearchQuery] = useState("");
  const [showCart, setShowCart] = useState(false);

  useEffect(() => {
    loadOrders();
  }, []);

  const loadOrders = async () => {
    try {
      setIsLoading(true);
      const response = await axios.get(`${API_ORDER}/api/orders/admin`);
      setOrders(response.data || []);
    } catch (error) {
      console.error("Error loading orders:", error);
      setOrders([]);
    } finally {
      setIsLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString("en-US", {
      year: "numeric",
      month: "short",
      day: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case "pending":
        return "bg-yellow-100 text-yellow-800";
      case "processing":
        return "bg-blue-100 text-blue-800";
      case "completed":
        return "bg-green-100 text-green-800";
      case "cancelled":
        return "bg-red-100 text-red-800";
      default:
        return "bg-gray-100 text-gray-800";
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar
        onSearchChange={setSearchQuery}
        searchQuery={searchQuery}
        onCartClick={() => setShowCart(!showCart)}
        useSearchPage={true}
      />

      <div className="container mx-auto px-4 py-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">My Orders</h1>
          <p className="text-gray-600 mt-2">
            Track and view your order history
          </p>
        </div>

        {isLoading ? (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-600">Loading orders...</p>
          </div>
        ) : orders.length === 0 ? (
          <div className="bg-white rounded-lg shadow-md p-12 text-center">
            <svg
              className="w-24 h-24 mx-auto text-gray-400 mb-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
            <h3 className="text-xl font-semibold text-gray-700 mb-2">
              No orders yet
            </h3>
            <p className="text-gray-600 mb-6">
              Start shopping to place your first order!
            </p>
            <a
              href="/"
              className="inline-block bg-blue-600 hover:bg-blue-700 text-white px-6 py-3 rounded-lg font-medium"
            >
              Browse Products
            </a>
          </div>
        ) : (
          <div className="space-y-4">
            {orders.map((order) => (
              <div
                key={order.id}
                className="bg-white rounded-lg shadow-md overflow-hidden"
              >
                <div className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <h3 className="text-lg font-semibold text-gray-900">
                        Order #{order.id.substring(0, 8)}
                      </h3>
                      <p className="text-sm text-gray-600">
                        {formatDate(order.created_at)}
                      </p>
                    </div>
                    <div className="text-right">
                      <span
                        className={`inline-block px-3 py-1 rounded-full text-sm font-medium ${getStatusColor(
                          order.status
                        )}`}
                      >
                        {order.status.charAt(0).toUpperCase() +
                          order.status.slice(1)}
                      </span>
                      <p className="text-lg font-bold text-gray-900 mt-2">
                        ${order.total.toFixed(2)}
                      </p>
                    </div>
                  </div>

                  <div className="border-t pt-4">
                    <button
                      onClick={() =>
                        setSelectedOrder(
                          selectedOrder?.id === order.id ? null : order
                        )
                      }
                      className="text-blue-600 hover:text-blue-800 font-medium text-sm flex items-center"
                    >
                      {selectedOrder?.id === order.id ? (
                        <>
                          <svg
                            className="w-4 h-4 mr-1"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M5 15l7-7 7 7"
                            />
                          </svg>
                          Hide Details
                        </>
                      ) : (
                        <>
                          <svg
                            className="w-4 h-4 mr-1"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M19 9l-7 7-7-7"
                            />
                          </svg>
                          View Details ({order.items.length} items)
                        </>
                      )}
                    </button>

                    {selectedOrder?.id === order.id && (
                      <div className="mt-4 space-y-3">
                        {order.items.map((item, index) => (
                          <div
                            key={index}
                            className="flex items-center justify-between bg-gray-50 p-3 rounded-lg"
                          >
                            <div className="flex-1">
                              <h4 className="font-medium text-gray-900">
                                {item.name}
                              </h4>
                              <p className="text-sm text-gray-600">
                                Quantity: {item.quantity}
                              </p>
                            </div>
                            <div className="text-right">
                              <p className="font-semibold text-gray-900">
                                ${(item.price * item.quantity).toFixed(2)}
                              </p>
                              <p className="text-sm text-gray-600">
                                ${item.price.toFixed(2)} each
                              </p>
                            </div>
                          </div>
                        ))}
                        <div className="border-t pt-3 flex justify-between items-center font-bold">
                          <span className="text-gray-900">Total</span>
                          <span className="text-xl text-blue-600">
                            ${order.total.toFixed(2)}
                          </span>
                        </div>
                      </div>
                    )}
                  </div>

                  {/* Cart Sidebar */}
                  {showCart && <Cart onClose={() => setShowCart(false)} />}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
