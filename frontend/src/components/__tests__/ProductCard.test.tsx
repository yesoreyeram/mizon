import { render, screen } from '@testing-library/react'
import '@testing-library/jest-dom'
import ProductCard from '../ProductCard'

// Mock axios
jest.mock('axios', () => ({
  post: jest.fn(() => Promise.resolve({ data: {} })),
}))

// Mock window.alert
global.alert = jest.fn()

describe('ProductCard', () => {
  const mockProduct = {
    id: '123',
    name: 'Test Product',
    description: 'Test Description',
    price: 29.99,
    category: 'Electronics',
    stock: 10,
    image_url: 'https://example.com/image.jpg',
  }

  it('renders product information correctly', () => {
    render(<ProductCard product={mockProduct} />)

    expect(screen.getByText('Test Product')).toBeInTheDocument()
    expect(screen.getByText('Test Description')).toBeInTheDocument()
    expect(screen.getByText('$29.99')).toBeInTheDocument()
    expect(screen.getByText('Electronics')).toBeInTheDocument()
    expect(screen.getByText('Stock: 10')).toBeInTheDocument()
  })

  it('renders Add to Cart button', () => {
    render(<ProductCard product={mockProduct} />)

    const button = screen.getByRole('button', { name: /add to cart/i })
    expect(button).toBeInTheDocument()
  })

  it('displays price with two decimal places', () => {
    const productWithPrice = { ...mockProduct, price: 10 }
    render(<ProductCard product={productWithPrice} />)

    expect(screen.getByText('$10.00')).toBeInTheDocument()
  })
})
