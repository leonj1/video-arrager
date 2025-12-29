package app

import (
	"encoding/json"
	"os"
)

type Project struct {
	Videos []string `json:"videos"`
}

func SaveProject(videos []*Video, path string) error {
	project := Project{
		Videos: make([]string, len(videos)),
	}

	for i, v := range videos {
		project.Videos[i] = v.Path
	}

	data, err := json.MarshalIndent(project, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func LoadProject(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var project Project
	if err := json.Unmarshal(data, &project); err != nil {
		return nil, err
	}

	return project.Videos, nil
}
