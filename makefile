build:
	@echo "Building..."
	@go build -o bin/$(NAME) -v
	
dev:
	@go run main.go	

tailwind:
	@echo "🎨 Building tailwind"
	@cd templates && npx tailwindcss -i ./tailwind.css -o ./public/tailwind.css --watch