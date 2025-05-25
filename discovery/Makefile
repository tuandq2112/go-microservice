.PHONY: start-consul all install clean

CONSUL_VERSION := 1.15.3
CONSUL_ZIP := consul_${CONSUL_VERSION}_linux_amd64.zip
DOWNLOAD_URL := https://releases.hashicorp.com/consul/${CONSUL_VERSION}/${CONSUL_ZIP}

install:
	@echo "Downloading Consul ${CONSUL_VERSION}..."
	curl -LO ${DOWNLOAD_URL}
	@echo "Unzipping..."
	unzip -o ${CONSUL_ZIP}
	@echo "Moving consul binary to /usr/local/bin (requires sudo)..."
	sudo mv consul /usr/local/bin/
	@echo "Setting executable permission..."
	sudo chmod +x /usr/local/bin/consul
	@echo "Cleaning up zip file..."
	rm -f ${CONSUL_ZIP}
	@echo "Consul installed successfully!"

run:
	@echo "Starting Consul in dev mode..."
	consul agent -dev

clean:
	rm -f consul ${CONSUL_ZIP}

start-consul:
	docker compose -f docker-compose.yaml up -d

