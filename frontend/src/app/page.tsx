"use client";

import { useEffect, useState } from "react";
import axios from "axios";
import Navbar from "@/components/Navbar";
import ProductCard from "@/components/ProductCard";
import Cart from "@/components/Cart";

const API_CATALOG =
  process.env.NEXT_PUBLIC_CATALOG_API || "http://localhost:8002";
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

export default function Home() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<string[]>([]);
  const [selectedCategory, setSelectedCategory] = useState<string>("");
  const [searchQuery, setSearchQuery] = useState<string>("");
  const [isLoading, setIsLoading] = useState(true);
  const [showCart, setShowCart] = useState(false);

  useEffect(() => {
    loadProducts();
    loadCategories();
  }, []);

  const loadProducts = async () => {
    try {
      setIsLoading(true);
      const response = await axios.get(`${API_CATALOG}/api/catalog/products`);
      setProducts(response.data || []);
    } catch (error) {
      console.error("Error loading products:", error);
      setProducts([]);
    } finally {
      setIsLoading(false);
    }
  };

  const loadCategories = async () => {
    try {
      const response = await axios.get(`${API_CATALOG}/api/catalog/categories`);
      setCategories(response.data || []);
    } catch (error) {
      console.error("Error loading categories:", error);
    }
  };

  const handleSearch = async () => {
    if (!searchQuery.trim()) {
      loadProducts();
      return;
    }

    try {
      setIsLoading(true);
      const response = await axios.get(
        `${API_SEARCH}/api/search?q=${encodeURIComponent(searchQuery)}`
      );
      setProducts(response.data.products || []);
    } catch (error) {
      console.error("Error searching products:", error);
      setProducts([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCategoryFilter = async (category: string) => {
    setSelectedCategory(category);
    setSearchQuery("");

    if (!category) {
      loadProducts();
      return;
    }

    try {
      setIsLoading(true);
      const response = await axios.get(
        `${API_CATALOG}/api/catalog/categories/${category}/products`
      );
      setProducts(response.data || []);
    } catch (error) {
      console.error("Error filtering by category:", error);
      setProducts([]);
    } finally {
      setIsLoading(false);
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
        {/* Categories */}
        <div className="mb-8">
          <h2 className="text-xl font-semibold mb-4">Categories</h2>
          <div className="flex flex-wrap gap-2">
            <button
              onClick={() => handleCategoryFilter("")}
              className={`px-4 py-2 rounded-lg ${
                !selectedCategory
                  ? "bg-blue-600 text-white"
                  : "bg-white text-gray-700 hover:bg-gray-100"
              }`}
            >
              All Products
            </button>
            {categories.map((category) => (
              <button
                key={category}
                onClick={() => handleCategoryFilter(category)}
                className={`px-4 py-2 rounded-lg ${
                  selectedCategory === category
                    ? "bg-blue-600 text-white"
                    : "bg-white text-gray-700 hover:bg-gray-100"
                }`}
              >
                {category}
              </button>
            ))}
          </div>
        </div>

        {/* Products Grid */}
        {isLoading ? (
          <div className="text-center py-12">
            <div className="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
            <p className="mt-4 text-gray-600">Loading products...</p>
          </div>
        ) : products.length === 0 ? (
          <div className="text-center py-12">
            <p className="text-gray-600 text-lg">No products found</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
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
