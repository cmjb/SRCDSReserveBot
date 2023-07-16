
all:
	scripts/make.sh

install:
	scripts/install.sh

gotest:
	scripts/make.sh
	scripts/test.sh

run:
	scripts/make.sh
	scripts/run.sh