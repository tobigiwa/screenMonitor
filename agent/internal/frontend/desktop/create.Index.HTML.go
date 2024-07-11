package desktop

import (
	frontend "agent/internal/frontend/components"
	"context"
	"fmt"
	"os"
	"path/filepath"
)

func CreateIndexHTML() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	parentPath := filepath.Dir(cwd)

	p := fmt.Sprintf("%s/agent/internal/frontend/desktop/index.html", parentPath)
	indexHTML, err := os.Create(p)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}

	if err = frontend.IndexPage().Render(context.Background(), indexHTML); err != nil {
		return fmt.Errorf("could not generate html: %v", err)
	}

	return nil
}
