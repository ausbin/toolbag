# honestly, i made this number up by trial-and-error. it seems to work,
# and i'm too lazy to do the math because i don't want some exact dpi
DPI = 256
OUT = ../notfound-img
PROBLEMS = $(patsubst %.tex,$(OUT)/%,$(wildcard *.tex))

.PHONY: all clean

all: $(PROBLEMS)

$(OUT):
	mkdir $@

$(OUT)/%: %.base64 $(OUT)
	{ tail -1 $(notdir $@).tex | cut -b 3-; cat $<; } >$@

%.base64: %.png
	base64 -w 0 <$< >$@

%.png: %.dvi
	dvipng -T tight -D $(DPI) -o $@ $<

%.dvi: %.tex
	latex $<

clean:
	rm -rvf $(OUT) *.base64 *.png *.dvi *.log *.aux
