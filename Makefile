.PHONY: certs clean-certs  deps-up deps-down glossary-build

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

deps-up: glossary-build
	podman compose -f $(LOCAL_DIR)/compose.yaml up -d

deps-down:
	-podman compose -f $(LOCAL_DIR)/compose.yaml down

glossary-build:
	podman build -t glossary -f ./local/glossary.Containerfile ./local