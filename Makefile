include $(GOROOT)/src/Make.inc

TARG=fcmp
GOFILES=\
	fcmp.go
# for binary programs --> copied into $GOROOT/bin
include $(GOROOT)/src/Make.cmd
