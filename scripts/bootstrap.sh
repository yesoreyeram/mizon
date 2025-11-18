#!/bin/bash

echo "Bootstrapping Mizon catalog with 100 products..."

CATALOG_API="http://localhost:8002"
SEARCH_API="http://localhost:8003"

# Categories
categories=("Electronics" "Clothing" "Books" "Home & Kitchen" "Sports" "Toys" "Beauty" "Automotive" "Garden" "Health")

# Product names by category
electronics=("Laptop" "Smartphone" "Tablet" "Headphones" "Smartwatch" "Camera" "TV" "Gaming Console" "Keyboard" "Mouse")
clothing=("T-Shirt" "Jeans" "Dress" "Jacket" "Sneakers" "Boots" "Hat" "Scarf" "Gloves" "Sweater")
books=("Fiction Novel" "Science Book" "History Book" "Cookbook" "Biography" "Fantasy Novel" "Mystery Book" "Self-Help Book" "Comic Book" "Poetry Book")
home=("Coffee Maker" "Blender" "Vacuum Cleaner" "Toaster" "Microwave" "Lamp" "Pillow" "Blanket" "Curtains" "Rug")
sports=("Basketball" "Football" "Tennis Racket" "Yoga Mat" "Dumbbells" "Running Shoes" "Bicycle" "Swimming Goggles" "Golf Clubs" "Skateboard")
toys=("Action Figure" "Doll" "Board Game" "Puzzle" "RC Car" "Building Blocks" "Stuffed Animal" "Play-Doh" "Video Game" "Coloring Book")
beauty=("Shampoo" "Conditioner" "Face Cream" "Lipstick" "Mascara" "Perfume" "Nail Polish" "Hair Dryer" "Makeup Brush Set" "Skin Serum")
automotive=("Car Vacuum" "Dash Cam" "Car Charger" "Phone Mount" "Seat Covers" "Floor Mats" "Air Freshener" "Tool Kit" "Jumper Cables" "Tire Pressure Gauge")
garden=("Garden Hose" "Lawn Mower" "Plant Pot" "Seeds" "Fertilizer" "Garden Gloves" "Pruning Shears" "Watering Can" "Outdoor Lights" "Bird Feeder")
health=("Vitamins" "First Aid Kit" "Blood Pressure Monitor" "Thermometer" "Hand Sanitizer" "Face Mask" "Pain Relief Gel" "Protein Powder" "Fitness Tracker" "Massage Gun")

counter=0

# Function to create a product
create_product() {
    local name="$1"
    local category="$2"
    local price="$3"
    local stock="$4"
    local description="$5"
    
    curl -X POST "$CATALOG_API/api/catalog/products" \
        -H "Content-Type: application/json" \
        -d "{
            \"name\": \"$name\",
            \"description\": \"$description\",
            \"price\": $price,
            \"category\": \"$category\",
            \"stock\": $stock,
            \"image_url\": \"https://via.placeholder.com/300x300?text=$name\"
        }" \
        -s -o /dev/null
    
    echo "Created: $name ($category) - \$$price"
    ((counter++))
}

# Generate 10 products per category
for i in {0..9}; do
    # Electronics
    price=$(awk -v min=50 -v max=2000 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    stock=$(awk -v min=5 -v max=100 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    create_product "${electronics[$i]}" "Electronics" "$price" "$stock" "High-quality ${electronics[$i],,} with advanced features"
    
    # Clothing
    price=$(awk -v min=15 -v max=200 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    stock=$(awk -v min=10 -v max=200 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    create_product "${clothing[$i]}" "Clothing" "$price" "$stock" "Comfortable and stylish ${clothing[$i],,}"
    
    # Books
    price=$(awk -v min=10 -v max=50 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    stock=$(awk -v min=20 -v max=150 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    create_product "${books[$i]}" "Books" "$price" "$stock" "Fascinating ${books[$i],,} that will captivate you"
    
    # Home & Kitchen
    price=$(awk -v min=20 -v max=300 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    stock=$(awk -v min=10 -v max=100 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    create_product "${home[$i]}" "Home & Kitchen" "$price" "$stock" "Essential ${home[$i],,} for your home"
    
    # Sports
    price=$(awk -v min=15 -v max=500 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    stock=$(awk -v min=10 -v max=80 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    create_product "${sports[$i]}" "Sports" "$price" "$stock" "Professional-grade ${sports[$i],,}"
    
    # Toys
    price=$(awk -v min=10 -v max=100 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    stock=$(awk -v min=15 -v max=200 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    create_product "${toys[$i]}" "Toys" "$price" "$stock" "Fun and educational ${toys[$i],,}"
    
    # Beauty
    price=$(awk -v min=10 -v max=150 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    stock=$(awk -v min=20 -v max=150 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    create_product "${beauty[$i]}" "Beauty" "$price" "$stock" "Premium ${beauty[$i],,} for your beauty routine"
    
    # Automotive
    price=$(awk -v min=15 -v max=200 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    stock=$(awk -v min=10 -v max=100 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    create_product "${automotive[$i]}" "Automotive" "$price" "$stock" "Quality ${automotive[$i],,} for your vehicle"
    
    # Garden
    price=$(awk -v min=10 -v max=300 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    stock=$(awk -v min=15 -v max=120 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    create_product "${garden[$i]}" "Garden" "$price" "$stock" "Perfect ${garden[$i],,} for your garden"
    
    # Health
    price=$(awk -v min=10 -v max=200 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    stock=$(awk -v min=20 -v max=150 'BEGIN{srand(); print int(min+rand()*(max-min+1))}')
    create_product "${health[$i]}" "Health" "$price" "$stock" "Essential ${health[$i],,} for wellness"
done

echo ""
echo "Bootstrap complete! Created $counter products."
echo "Products are now available in the catalog."
echo "The search service will be automatically indexed when products are added to cart."
