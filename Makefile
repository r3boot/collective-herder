CH = ch
CHD = chd

BUILD_DIR = ./build
CMD_DIR = ./cmd

SRC_CH = ${CMD_DIR}/${CH}/${CH}.go
TARGET_CH = ${BUILD_DIR}/${CH}
SRC_CHD = ${CMD_DIR}/${CHD}/${CHD}.go
TARGET_CHD = ${BUILD_DIR}/${CHD}

all: ${TARGET_CHD} ${TARGET_CH}

${TARGET_CHD}:
	go build -v -o ${TARGET_CHD} ${SRC_CHD}

${TARGET_CH}:
	go build -v -o ${TARGET_CH} ${SRC_CH}

clean:
	rm -rf ${BUILD_DIR} || true
