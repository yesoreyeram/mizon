"use client";

import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import axios from "axios";
import Navbar from "@/components/Navbar";
import ProductCard from "@/components/ProductCard";
import Cart from "@/components/Cart";

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

export default function SearchPage() {
  const searchParams = useSearchParams();
  const [products, setProducts] = useState<Product[]>([]);
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [isLoading, setIsLoading] = useState(false);
  const [showCart, setShowCart] = useState(false);
  const [totalResults, setTotalResults] = useState(0);

  useEffect(() => {
    const query = searchParams.get("q");
    if (query) {
      setSearchQuery(query);
      performSearch(query);
    }
  }, [searchParams]);

  const performSearch = async (query: string) => {
    if (!query.trim()) {
      setProducts([]);
      setTotalResults(0);
      return;
    }

    try {
      setIsLoading(true);
      const response = await axios.get(
        `${API_SEARCH}/api/search?q=${encodeURIComponent(query)}`
      );
      setProducts(response.data.products || []);
      setTotalResults(response.data.total || 0);
    } catch (error) {
      console.error("Error searching products:", error);
      setProducts([]);
      setTotalResults(0);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearch = () => {
    if (searchQuery.trim()) {
      window.history.pushState(
        {},
        "",
        `/search?q=${encodeURIComponent(searchQuery)}`
      );
      performSearch(searchQuery);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Navbar
        onSearchChange={setSearchQuery}
        onSearch={handleSearch}
        searchQuery={searchQuery}
        onCartClick={() => setShowCart(!showCart)}
      />

      <div className="container mx-auto px-4 py-8">
        {/* Search Results Header */}
        {searchParams.get("q") && (
          <div className="mb-6">
            <h1 className="text-2xl font-bold text-gray-800">
              Search Results for &quot;{searchParams.get("q")}&quot;
            </h1>
            <p className="text-gray-600 mt-2">
              {isLoading
                ? "Searching..."
                : `Found ${totalResults} ${
                    totalResults === 1 ? "result" : "results"
                  }`}
            </p>
          </div>
        )}

        {/* Loading State */}
        {isLoading ? (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-600">Searching products...</p>
          </div>
        ) : !searchParams.get("q") ? (
          /* No Search Query */
          <div className="text-center py-12">
            <svg
              className="mx-auto h-24 w-24 text-gray-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
              />
            </svg>
            <h2 className="mt-4 text-xl font-semibold text-gray-700">
              Start Your Search
            </h2>
            <p className="mt-2 text-gray-600">
              Enter keywords in the search bar above to find products
            </p>
            <a
              href="/"
              className="mt-6 inline-block bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700"
            >
              Browse All Products
            </a>
          </div>
        ) : products.length === 0 ? (
          /* No Results */
          <div className="text-center py-12">
            <svg
              className="mx-auto h-24 w-24 text-gray-400"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <h2 className="mt-4 text-xl font-semibold text-gray-700">
              No Results Found
            </h2>
            <p className="mt-2 text-gray-600">
              We couldn&apos;t find any products matching &quot;{searchParams.get("q")}&quot;
            </p>
            <p className="mt-2 text-gray-500 text-sm">
              Try different keywords or browse all products
            </p>
            <a
              href="/"
              className="mt-6 inline-block bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700"
            >
              Browse All Products
            </a>
          </div>
        ) : (
          /* Products Grid */
          <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
            {products.map((product) => (
              <ProductCard key={product.id} product={product} />
            ))}
          </div>
        )}
      </div>

      {/* Cart Sidebar */}
      {showCart && <Cart onClose={() => setShowCart(false)} />}
    </div>
  );
}
