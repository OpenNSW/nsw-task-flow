package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/OpenNSW/nsw-task-flow/orchestrator"
)

// loadTemplates scans all *.json files recursively in templatesDir and registers them in the registry.
func loadTemplates(registry *orchestrator.TaskTemplateRegistry, templatesDir string) error {
	err := filepath.WalkDir(templatesDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".json" {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", path, err)
		}

		var entry orchestrator.TaskTemplateEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			// Skip files that aren't valid JSON or are other structures
			return nil
		}
		if entry.TemplateID == "" || entry.PluginName == "" {
			// Skip non-template JSONs (like workflow graphs, UI schemas, or JSONForms files)
			return nil
		}

		registry.Register(entry)
		log.Printf("[Registry] Loaded template: %s (task_type=%s, plugin=%s)", entry.TemplateID, entry.TaskType, entry.PluginName)
		return nil
	})

	if err != nil {
		return fmt.Errorf("recursive template search failed: %w", err)
	}
	return nil
}
