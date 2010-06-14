include $(GOROOT)/src/Make.$(GOARCH)

TARG=cachemap
GOFILES=\
	cachemap.go

include $(GOROOT)/src/Make.pkg
