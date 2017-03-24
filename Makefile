CH = ch
CHD = chd
CODEGEN = codegen

BUILD_DIR = ./build
CMD_DIR = ./cmd
CFG_DIR = ./etc
SCRIPTS_DIR = ./scripts
COVERAGE_DIR = ./test_coverage

SRC_CH = ${CMD_DIR}/${CH}/${CH}.go
TARGET_CH = ${BUILD_DIR}/${CH}
SRC_CHD = ${CMD_DIR}/${CHD}/${CHD}.go
TARGET_CHD = ${BUILD_DIR}/${CHD}
SRC_CODEGEN = ${CMD_DIR}/${CODEGEN}/${CODEGEN}.go
TARGET_CODEGEN = ${BUILD_DIR}/${CODEGEN}

all: ${TARGET_CODEGEN} ${CODEGEN} ${TARGET_CHD} ${TARGET_CH}

${TARGET_CODEGEN}:
	[[ -d ${BUILD_DIR} ]] || mkdir -p ${BUILD_DIR}
	CGO_ENABLED=0 go build -v -o ${TARGET_CODEGEN} ${SRC_CODEGEN}

${CODEGEN}: ${TARGET_CODEGEN}
	${TARGET_CODEGEN}

${TARGET_CHD}:
	[[ -d ${BUILD_DIR} ]] || mkdir -p ${BUILD_DIR}
	CGO_ENABLED=0 go build -v -o ${TARGET_CHD} ${SRC_CHD}

${TARGET_CH}:
	[[ -d ${BUILD_DIR} ]] || mkdir -p ${BUILD_DIR}
	CGO_ENABLED=0 go build -v -o ${TARGET_CH} ${SRC_CH}

test:
	${SCRIPTS_DIR}/run_all_tests.sh

install:
	install -d -o root -g root -m 0750 /etc/ch
	install -d -o root -g root -m 0750 /etc/ch/commands.d
	install -o root -g root -m 0640 ${CFG_DIR}/server.yml \
		/etc/ch/server.yml
	install -o root -g root -m 0640 ${CFG_DIR}/client.yml \
		/etc/ch/client.yml
	install -o root -g root -m 0640 ${CFG_DIR}/commands.d/uname.yml \
		/etc/ch/commands.d/uname.yml

clean:
	rm -rf ${BUILD_DIR} ${COVERAGE_DIR} || true
