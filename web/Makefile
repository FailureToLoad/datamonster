.PHONY: daisyui
DAISYUI_VERSION := 5.5.14
DAISYUI_CDN := https://cdn.jsdelivr.net/npm/daisyui@$(DAISYUI_VERSION)
STYLES_DIR := src/styles/daisyui

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
	curl -sL $(DAISYUI_CDN)/components/navbar.css -o $(STYLES_DIR)/navbar.css
	curl -sL $(DAISYUI_CDN)/components/dropdown.css -o $(STYLES_DIR)/dropdown.css
	curl -sL $(DAISYUI_CDN)/components/tab.css -o $(STYLES_DIR)/tab.css
	curl -sL $(DAISYUI_CDN)/components/collapse.css -o $(STYLES_DIR)/collapse.css
