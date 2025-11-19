import { render, screen } from '@testing-library/react'
import '@testing-library/jest-dom'
import Cart from '../Cart'
import axios from 'axios'

// Mock axios
jest.mock('axios')
const mockedAxios = axios as jest.Mocked<typeof axios>

// Mock window.alert
global.alert = jest.fn()

describe('Cart', () => {
  const mockProps = {
    onClose: jest.fn(),
  }

  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('renders Shopping Cart title', () => {
    mockedAxios.get.mockResolvedValue({ data: { items: [] } })
    render(<Cart {...mockProps} />)
    expect(screen.getByText('Shopping Cart')).toBeInTheDocument()
  })

  it('shows loading state initially', () => {
    mockedAxios.get.mockResolvedValue({ data: { items: [] } })
    render(<Cart {...mockProps} />)
    // Loading spinner should be visible initially
    const spinner = document.querySelector('.animate-spin')
    expect(spinner).toBeInTheDocument()
  })

  it('displays empty cart message when no items', async () => {
    mockedAxios.get.mockResolvedValue({ data: { items: [] } })
    render(<Cart {...mockProps} />)

    // Wait for loading to complete
    const emptyMessage = await screen.findByText('Your cart is empty')
    expect(emptyMessage).toBeInTheDocument()
  })

  it('renders close button', () => {
    mockedAxios.get.mockResolvedValue({ data: { items: [] } })
    render(<Cart {...mockProps} />)

    const closeButton = screen.getByRole('button', { name: '' })
    expect(closeButton).toBeInTheDocument()
  })
})
