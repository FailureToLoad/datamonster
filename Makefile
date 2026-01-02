.PHONY: certs clean-certs valkey authelia postgres glossary daisyui

CERT_DIR := local/.certs
LOCAL_DIR := $(shell pwd)/local
DAISYUI_VERSION := 5.5.14
DAISYUI_CDN := https://cdn.jsdelivr.net/npm/daisyui@$(DAISYUI_VERSION)
STYLES_DIR := web/src/styles/daisyui

daisyui:
	mkdir -p $(STYLES_DIR)
	curl -sL $(DAISYUI_CDN)/base/reset.css -o $(STYLES_DIR)/reset.css
	curl -sL $(DAISYUI_CDN)/base/rootcolor.css -o $(STYLES_DIR)/rootcolor.css
	curl -sL $(DAISYUI_CDN)/base/properties.css -o $(STYLES_DIR)/properties.css
	curl -sL $(DAISYUI_CDN)/components/button.css -o $(STYLES_DIR)/button.css
	curl -sL $(DAISYUI_CDN)/components/modal.css -o $(STYLES_DIR)/modal.css
	curl -sL $(DAISYUI_CDN)/components/input.css -o $(STYLES_DIR)/input.css
	curl -sL $(DAISYUI_CDN)/components/select.css -o $(STYLES_DIR)/select.css
	curl -sL $(DAISYUI_CDN)/components/checkbox.css -o $(STYLES_DIR)/checkbox.css
	curl -sL $(DAISYUI_CDN)/components/radio.css -o $(STYLES_DIR)/radio.css
	curl -sL $(DAISYUI_CDN)/components/badge.css -o $(STYLES_DIR)/badge.css
	curl -sL $(DAISYUI_CDN)/components/card.css -o $(STYLES_DIR)/card.css
	curl -sL $(DAISYUI_CDN)/components/menu.css -o $(STYLES_DIR)/menu.css
	curl -sL $(DAISYUI_CDN)/components/tooltip.css -o $(STYLES_DIR)/tooltip.css
	curl -sL $(DAISYUI_CDN)/components/table.css -o $(STYLES_DIR)/table.css

certs: $(CERT_DIR)/authelia.crt $(CERT_DIR)/authelia.key

$(CERT_DIR):
	mkdir -p $(CERT_DIR)

$(CERT_DIR)/authelia.crt $(CERT_DIR)/authelia.key: | $(CERT_DIR)
	openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
		-keyout $(CERT_DIR)/authelia.key -out $(CERT_DIR)/authelia.crt \
		-subj "/CN=127.0.0.1" \
		-addext "subjectAltName=IP:127.0.0.1"
	sudo cp $(CERT_DIR)/authelia.crt /etc/pki/ca-trust/source/anchors/authelia-local.crt
	sudo update-ca-trust

clean-certs:
	sudo rm -f /etc/pki/ca-trust/source/anchors/authelia-local.crt
	sudo update-ca-trust
	rm -rf $(CERT_DIR)

valkey:
	-podman rm -f valkey
	podman run -d --name valkey -p 6379:6379 valkey/valkey:8.0-alpine

authelia:
	mkdir -p $(LOCAL_DIR)/data
	-podman rm -f authelia
	podman run -d --name authelia \
		-p 9091:9091 \
		-v $(LOCAL_DIR)/configuration.yml:/config/configuration.yml:ro,Z \
		-v $(LOCAL_DIR)/users_database.yml:/config/users_database.yml:ro,Z \
		-v $(LOCAL_DIR)/data:/config/data:Z \
		-v $(shell pwd)/$(CERT_DIR):/certs:ro,Z \
		-e AUTHELIA_JWT_SECRET="insecure_jwt_secret_for_local_dev" \
		-e AUTHELIA_SESSION_SECRET="insecure_session_secret_for_local_dev" \
		-e AUTHELIA_STORAGE_ENCRYPTION_KEY="insecure_encryption_key_for_local_dev" \
		-e AUTHELIA_SERVER_DISABLE_HEALTHCHECK=true \
		docker.io/authelia/authelia:latest

postgres:
	-podman rm -f postgres
	podman run -d --name postgres \
		-p 5432:5432 \
		-e POSTGRES_USER=datamonster \
		-e POSTGRES_PASSWORD=datamonster \
		-e POSTGRES_DB=datamonster \
		docker.io/postgres:18-alpine

glossary:
	podman build -t glossary -f ./local/glossary.Containerfile ./local
	-podman rm -f glossary
	podman run -d --name glossary -p 9080:80 glossary