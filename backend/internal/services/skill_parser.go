package services

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"fnchatbot/internal/models"

	"gopkg.in/yaml.v3"
)

// SkillFrontmatter defines the structure for YAML frontmatter in skill files
type SkillFrontmatter struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}

// ParseSkillFile parses a skill from an uploaded file (.md or .zip)
func ParseSkillFile(fileHeader *multipart.FileHeader) (*models.Skill, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	var content []byte
	var skillName string

	if ext == ".zip" {
		// Read entire zip file into memory
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, file); err != nil {
			return nil, fmt.Errorf("failed to read zip file: %w", err)
		}

		readerAt := bytes.NewReader(buf.Bytes())
		zipReader, err := zip.NewReader(readerAt, int64(buf.Len()))
		if err != nil {
			return nil, fmt.Errorf("failed to parse zip: %w", err)
		}

		// Find SKILL.md or any .md file
		var skillFile *zip.File
		for _, f := range zipReader.File {
			if strings.EqualFold(filepath.Base(f.Name), "SKILL.md") {
				skillFile = f
				break
			}
		}

		if skillFile == nil {
			// Fallback to any .md file if SKILL.md not found
			for _, f := range zipReader.File {
				if strings.ToLower(filepath.Ext(f.Name)) == ".md" {
					skillFile = f
					break
				}
			}
		}

		if skillFile == nil {
			return nil, errors.New("no markdown file found in zip archive")
		}

		rc, err := skillFile.Open()
		if err != nil {
			return nil, fmt.Errorf("failed to open file in zip: %w", err)
		}
		defer rc.Close()

		content, err = io.ReadAll(rc)
		if err != nil {
			return nil, fmt.Errorf("failed to read file content: %w", err)
		}

		// Use filename without extension as default name
		skillName = strings.TrimSuffix(filepath.Base(skillFile.Name), filepath.Ext(skillFile.Name))

	} else if ext == ".md" {
		content, err = io.ReadAll(file)
		if err != nil {
			return nil, fmt.Errorf("failed to read file content: %w", err)
		}
		skillName = strings.TrimSuffix(fileHeader.Filename, ext)
	} else {
		return nil, fmt.Errorf("unsupported file type: %s", ext)
	}

	return parseMarkdownContent(content, skillName)
}

func parseMarkdownContent(content []byte, defaultName string) (*models.Skill, error) {
	skill := &models.Skill{
		Name:    defaultName,
		Enabled: true,
	}

	sContent := string(content)

	// Check for YAML frontmatter
	if strings.HasPrefix(sContent, "---") {
		parts := strings.SplitN(sContent, "---", 3)
		if len(parts) >= 3 {
			var meta SkillFrontmatter
			if err := yaml.Unmarshal([]byte(parts[1]), &meta); err == nil {
				if meta.Name != "" {
					skill.Name = meta.Name
				}
				if meta.Description != "" {
					skill.Description = meta.Description
				}
			}
			// Use the content after frontmatter as description if not provided in YAML
			if skill.Description == "" {
				skill.Description = strings.TrimSpace(parts[2])
			}
			return skill, nil
		}
	}

	// Fallback: simple parsing if no frontmatter or parsing failed
	lines := strings.Split(sContent, "\n")
	var descBuilder strings.Builder
	capturingDesc := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") && skill.Name == defaultName {
			skill.Name = strings.TrimSpace(strings.TrimPrefix(line, "# "))
		} else if strings.HasPrefix(line, "## Description") {
			capturingDesc = true
			continue
		} else if strings.HasPrefix(line, "##") && capturingDesc {
			capturingDesc = false
		} else if capturingDesc {
			descBuilder.WriteString(line + "\n")
		}
	}

	if descBuilder.Len() > 0 {
		skill.Description = strings.TrimSpace(descBuilder.String())
	} else if skill.Description == "" {
		// If no description found, use the whole content
		skill.Description = sContent
	}

	return skill, nil
}
