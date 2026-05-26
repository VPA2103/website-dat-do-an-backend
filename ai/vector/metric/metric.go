package metric

import "fmt"

// Spec returns the pgvector operator class and distance operator for the given metric.
// Supported metrics: "cosine", "l2", "ip".
func Spec(metric string) (opClass string, distOp string, err error) {
	switch metric {
	case "cosine":
		return "vector_cosine_ops", "<=>", nil
	case "l2":
		return "vector_l2_ops", "<->", nil
	case "ip":
		// pgvector uses <#> for (negative) inner product distance.
		return "vector_ip_ops", "<#>", nil
	default:
		return "", "", fmt.Errorf("unsupported VECTOR_METRIC: %q", metric)
	}
}
