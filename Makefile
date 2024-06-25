ARCH=amd64
LDFLAGS=-v -X main.version=@dev
OUTPUT_DIR=./bin
DEBSRC_DIR=$(OUTPUT_DIR)/deb

.PHONY: build clean templates package

all:
	$(MAKE) templates
	$(MAKE) build

build:
	mkdir -p ./bin
	$(MAKE) ./bin/tgbot_$(ARCH)

clean:
	$(RM) -rf ./bin/*

templates:
	$(shell qtc -dir ./internal/controller/telegram/templates)

package: $(OUTPUT_DIR)/tgbot_$(ARCH)
	$(MAKE) $(OUTPUT_DIR)/tgbot_$(ARCH).deb

# #################### Binary files ###################

$(OUTPUT_DIR)/bin:
	mkdir -p './bin'

bin/tgbot_amd64: $(OUTPUT_DIR)/bin
	GOARCH=amd64 CGO_ENABLED=0 go build --ldflags '$(LDFLAGS)' -v -o ./bin/tgbot_amd64 cmd/bot/main.go

# ################## Debian package ##################

$(OUTPUT_DIR)/tgbot_$(ARCH).deb:
	mkdir -p $(DEBSRC_DIR)
	cp -r ./build/package/debian/* $(DEBSRC_DIR)
	mv $(DEBSRC_DIR)/tgbot/DEBIAN/control_$(ARCH) $(DEBSRC_DIR)/tgbot/DEBIAN/control

	mkdir -p $(DEBSRC_DIR)/tgbot/usr/bin $(DEBSRC_DIR)/tgbot/etc/tgbot

	cp $(OUTPUT_DIR)/tgbot_$(ARCH) $(DEBSRC_DIR)/tgbot/usr/bin/tgbot
	cp ./configs/config.yml $(DEBSRC_DIR)/tgbot/etc/tgbot/config.dist.yml
	cd $(DEBSRC_DIR)/tgbot; md5deep -l -o f -r usr -r lib > DEBIAN/md5sums

	export DEB_CONTROL_SIZE=`du -s $(DEBSRC_DIR) | cut -f1` && \
		envsubst < ./build/package/debian/tgbot/DEBIAN/control_$(ARCH) > $(DEBSRC_DIR)/tgbot/DEBIAN/control

	cd $(DEBSRC_DIR); fakeroot dpkg-deb --build tgbot
	mv $(DEBSRC_DIR)/tgbot.deb $(OUTPUT_DIR)/tgbot_$(ARCH).deb
	rm -rf $(DEBSRC_DIR)
