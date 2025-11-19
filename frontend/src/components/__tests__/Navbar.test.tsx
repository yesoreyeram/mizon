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

  it('renders Cart button', () => {
    render(<Navbar {...mockProps} />)
    const cartButton = screen.getByText('Cart')
    expect(cartButton).toBeInTheDocument()
  })

  it('calls onCartClick when cart button is clicked', () => {
    render(<Navbar {...mockProps} />)
    const cartButton = screen.getByText('Cart')

    fireEvent.click(cartButton)
    expect(mockProps.onCartClick).toHaveBeenCalled()
  })

  it('renders Orders link', () => {
    render(<Navbar {...mockProps} />)
    const ordersLink = screen.getByText('Orders')
    expect(ordersLink).toBeInTheDocument()
  })

  it('displays admin user', () => {
    render(<Navbar {...mockProps} />)
    expect(screen.getByText('admin')).toBeInTheDocument()
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
