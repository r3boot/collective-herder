package facts

var (
	FACT_DIRS = []string{"/etc/ansible/facts.d"}
)

type Facts struct {
	facts map[string]interface{}
}
