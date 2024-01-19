package traverser

type TraversalRequest struct {
	TraversalDirectory string   `yaml:"traversal_directory"`
	ExcludedPaths      []string `yaml:"excluded_paths"`
}
