"use client";

import axios from "axios";

const API_CART = process.env.NEXT_PUBLIC_CART_API || "http://localhost:8004";
const API_SEARCH =
  process.env.NEXT_PUBLIC_SEARCH_API || "http://localhost:8003";

interface Product {
  id: string;
  name: string;
  description: string;
  price: number;
  category: string;
  stock: number;
  image_url: string;
}

interface ProductCardProps {
  product: Product;
}

export default function ProductCard({ product }: ProductCardProps) {
  const addToCart = async () => {
    try {
      await axios.post(`${API_CART}/api/cart/admin/items`, {
        product_id: product.id,
        name: product.name,
        price: product.price,
        quantity: 1,
      });

      // Also index in search
      await axios
        .post(`${API_SEARCH}/api/search/index`, {
          id: product.id,
          name: product.name,
          description: product.description,
          price: product.price,
          category: product.category,
          stock: product.stock,
          image_url: product.image_url,
        })
        .catch((err) => console.log("Search indexing error:", err));

      alert("Added to cart!");
    } catch (error) {
      console.error("Error adding to cart:", error);
      alert("Failed to add to cart");
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden hover:shadow-xl transition-shadow">
      <div className="h-48 bg-gray-200 flex items-center justify-center">
        {product.image_url ? (
          <img
            src={`https://placehold.co/300`}
            alt={product.name}
            className="h-full w-full object-cover"
          />
        ) : (
          <svg
            className="w-20 h-20 text-gray-400"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
            />
          </svg>
        )}
      </div>
      <div className="p-4">
        <h3 className="text-lg font-semibold text-gray-800 mb-2 truncate">
          {product.name}
        </h3>
        <p className="text-gray-600 text-sm mb-3 line-clamp-2">
          {product.description}
        </p>
        <div className="flex items-center justify-between mb-3">
          <span className="text-sm text-gray-500 bg-gray-100 px-2 py-1 rounded">
            {product.category}
          </span>
          <span className="text-sm text-gray-500">Stock: {product.stock}</span>
        </div>
        <div className="flex items-center justify-between">
          <span className="text-2xl font-bold text-blue-600">
            ${product.price.toFixed(2)}
          </span>
          <button
            onClick={addToCart}
            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
          >
            Add to Cart
          </button>
        </div>
      </div>
    </div>
  );
}
