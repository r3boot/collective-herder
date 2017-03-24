package utils

import (
	"regexp"
	"testing"
)

const (
	UUID_LEN      int    = 36
	RE_VALID_UUID string = "^[a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12}$"
)

var (
	reValidUuid = regexp.MustCompile(RE_VALID_UUID)
)

func TestUuidgen(t *testing.T) {
	var (
		uuid   string
		uuid_l int
	)

	uuid = Uuidgen()

	uuid_l = len(uuid)
	if uuid_l != UUID_LEN {
		t.Error("Length of uuid is", uuid_l, ", expected", UUID_LEN)
	}

	if !reValidUuid.MatchString(uuid) {
		t.Error("Uuid does not match regexp", RE_VALID_UUID)
	}
}
