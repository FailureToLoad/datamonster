.PHONY: certs clean-certs valkey-new valkey-kill authelia-new authelia-kill postgres-new postgres-kill

CERT_DIR := local/.certs
LOCAL_DIR := $(shell pwd)/local

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

valkey-new:
	podman run -d --name valkey -p 6379:6379 valkey/valkey:8.0-alpine

valkey-kill:
	podman rm -f valkey

authelia-new:
	mkdir -p $(LOCAL_DIR)/data
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

authelia-kill:
	podman rm -f authelia

postgres-new:
	podman run -d --name postgres \
		-p 5432:5432 \
		-e POSTGRES_USER=datamonster \
		-e POSTGRES_PASSWORD=datamonster \
		-e POSTGRES_DB=datamonster \
		docker.io/postgres:18-alpine

postgres-kill:
	podman rm -f postgres