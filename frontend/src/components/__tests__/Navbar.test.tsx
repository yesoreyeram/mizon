import { render, screen, fireEvent } from '@testing-library/react'
import '@testing-library/jest-dom'
import Navbar from '../Navbar'

describe('Navbar', () => {
  const mockProps = {
    onSearchChange: jest.fn(),
    onSearch: jest.fn(),
    searchQuery: '',
    onCartClick: jest.fn(),
    useSearchPage: false,
  }

  beforeEach(() => {
    jest.clearAllMocks()
    // Clear localStorage before each test
    Storage.prototype.getItem = jest.fn(() => null)
  })

  it('renders the Mizon logo', () => {
    render(<Navbar {...mockProps} />)
    expect(screen.getByText('Mizon')).toBeInTheDocument()
  })

  it('renders search input with placeholder', () => {
    render(<Navbar {...mockProps} />)
    const searchInput = screen.getByPlaceholderText('Search products...')
    expect(searchInput).toBeInTheDocument()
  })

  it('calls onSearchChange when typing in search input', () => {
    render(<Navbar {...mockProps} />)
    const searchInput = screen.getByPlaceholderText('Search products...')

    fireEvent.change(searchInput, { target: { value: 'laptop' } })
    expect(mockProps.onSearchChange).toHaveBeenCalledWith('laptop')
  })

  it('renders Cart button when logged in', () => {
    // Mock localStorage with user data
    Storage.prototype.getItem = jest.fn((key) => {
      if (key === 'username') return 'testuser'
      return null
    })
    
    render(<Navbar {...mockProps} />)
    const cartButton = screen.getByText('Cart')
    expect(cartButton).toBeInTheDocument()
  })

  it('calls onCartClick when cart button is clicked', () => {
    // Mock localStorage with user data
    Storage.prototype.getItem = jest.fn((key) => {
      if (key === 'username') return 'testuser'
      return null
    })
    
    render(<Navbar {...mockProps} />)
    const cartButton = screen.getByText('Cart')

    fireEvent.click(cartButton)
    expect(mockProps.onCartClick).toHaveBeenCalled()
  })

  it('renders Orders link when logged in', () => {
    // Mock localStorage with user data
    Storage.prototype.getItem = jest.fn((key) => {
      if (key === 'username') return 'testuser'
      return null
    })
    
    render(<Navbar {...mockProps} />)
    const ordersLink = screen.getByText('Orders')
    expect(ordersLink).toBeInTheDocument()
  })

  it('displays Sign In button when not logged in', () => {
    // Clear localStorage
    Storage.prototype.getItem = jest.fn(() => null)
    
    render(<Navbar {...mockProps} />)
    expect(screen.getByText('Sign In')).toBeInTheDocument()
  })

  it('displays username when logged in', () => {
    // Mock localStorage with user data
    Storage.prototype.getItem = jest.fn((key) => {
      if (key === 'username') return 'testuser'
      return null
    })
    
    render(<Navbar {...mockProps} />)
    expect(screen.getByText('testuser')).toBeInTheDocument()
  })

  it('calls onSearch when search form is submitted', () => {
    render(<Navbar {...mockProps} />)
    const searchForm = screen.getByPlaceholderText('Search products...').closest('form')

    if (searchForm) {
      fireEvent.submit(searchForm)
      expect(mockProps.onSearch).toHaveBeenCalled()
    }
  })
})
