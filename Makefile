
BINARY = main
SOURCE = ./cmd

.PHONY: all build run clean

all: build run


build:
	@echo "Building..."
	GOOS=linux GOARCH=amd64 go build -o $(BINARY) $(SOURCE)


run: build
	@echo "Starting application..."
	sudo ./$(BINARY) &
	sleep 5  # Даем приложению время на инициализацию

	@echo "Configuring network..."
	sudo ip addr add 10.0.0.1/24 dev tun0 || true
	sudo ip link set tun0 up || true

	@echo "Testing connection..."
	curl --interface tun0 http://example.com || true

# Очистка
clean:
	@echo "Cleaning..."
	rm -f $(BINARY)
	-sudo pkill $(BINARY) || true
	-sudo ip link set tun0 down || true
	-sudo ip addr del 10.0.0.1/24 dev tun0 || true