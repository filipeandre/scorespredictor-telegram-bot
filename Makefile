SHELL := /bin/bash
.ONESHELL:

APP_NAME = scorespredictor
DATA_DIR = /var/lib/$(APP_NAME)
CRON_FILE = /etc/cron.d/$(APP_NAME)
OUTPUT_BIN = /usr/local/$(APP_NAME)/bin
ENV_FILE = /etc/profile.d/$(APP_NAME).sh

define cron
SHELL=/bin/bash
MAILTO=$(USER)
CRON_TZ=Europe/Rome

0 10 * * * $(USER) source /etc/profile && $(APP_NAME) >/dev/null
endef
export cron

build:
	env CGO_ENABLED=0
	source /etc/profile && go build -i -o compiled

install: build
	sudo mkdir -p $(DATA_DIR)
	sudo mkdir -p $(OUTPUT_BIN)
	sudo chown $(USER):$(USER) $(DATA_DIR)
	cp .env.yaml $(DATA_DIR)
	sudo cp compiled $(OUTPUT_BIN)/$(APP_NAME)
	sudo sh -c 'echo "$(cron)" > $(CRON_FILE)'
	sudo sh -c 'echo "export SCORESPREDICTOR_HOME=$(DATA_DIR)" > $(ENV_FILE)'
	sudo sh -c 'echo "export PATH=$(OUTPUT_BIN):\$$PATH" >> $(ENV_FILE)'
	sudo systemctl restart crond.service

clean:
	rm compiled

uninstall: clean
	sudo rm -Rf $(OUTPUT_BIN)
	sudo rm -Rf $(DATA_DIR)
	sudo rm $(CRON_FILE)
	sudo rm $(ENV_FILE)

reset:
	git fetch && git reset --hard origin/master

update: reset build
	sudo cp compiled $(OUTPUT_BIN)/$(APP_NAME)