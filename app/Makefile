# Convenient shortcuts for the main Makefile
#################################################

.PHONY: start test build api ui deploy

start:
	cd .. && $(MAKE) start

test:
	cd .. && $(MAKE) test-app

build:
	cd .. && $(MAKE) build-app

api:
	cd .. && $(MAKE) test-app-api
	cd .. && $(MAKE) build-app-api

ui:
	cd .. && $(MAKE) test-app-ui
	cd .. && $(MAKE) build-app-ui

deploy:
	cd .. && $(MAKE) deploy-app
