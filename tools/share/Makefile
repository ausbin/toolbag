NOTFOUND = notfound
IMAGES = $(NOTFOUND)-img
TARBALL = share.tar.xz
DIRS = $(IMAGES) templates

.PHONY: clean

$(TARBALL): $(DIRS)
	tar cvJf $@ $^

$(IMAGES): $(NOTFOUND)
	$(MAKE) -C $<

clean:
	rm -vf $(TARBALL)
	$(MAKE) -C $(NOTFOUND) $@
