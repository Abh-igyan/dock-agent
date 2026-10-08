package docker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type Builder struct {
	dir string
}

func NewBuilder(dir string) *Builder {
	return &Builder{
		dir: dir,
	}
}

// Build attempts to build the Dockerfile and returns the combined stdout/stderr output if it fails.
func (b *Builder) Build(dockerfileContent string) error {
	// Write the Dockerfile
	err := os.WriteFile(fmt.Sprintf("%s/Dockerfile", b.dir), []byte(dockerfileContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write Dockerfile: %v", err)
	}

	// Execute `docker build .`
	cmd := exec.Command("docker", "build", ".")
	cmd.Dir = b.dir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("Docker build failed:\n%s", out.String())
	}

	return nil
}

