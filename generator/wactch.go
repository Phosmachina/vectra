package generator

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func IsDockerInstalled() bool {
	err := ExecuteCommand("docker version", false, true)
	return err == nil
}

func CreateDockerImage(dockerfileFileName string, imageName string) error {

	// Check if docker image already exists
	cmd := exec.Command("docker", "images", "-q", imageName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if strings.TrimSpace(string(output)) != "" {
		fmt.Printf("Image %s exists\n", imageName)
		return nil
	}

	dockerfileContent, _ := EmbedFS.ReadFile("template/.pipe/" + dockerfileFileName)

	// Create a temporary directory
	dir, err := os.MkdirTemp("", "docker")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir) // clean up

	// Write Dockerfile content
	err = os.WriteFile(fmt.Sprintf("%s/Dockerfile", dir), dockerfileContent, 0666)
	if err != nil {
		return err
	}

	// Build docker image
	err = ExecuteCommand("docker build -t "+imageName+" "+dir, true, true)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully built image %s\n", imageName)
	return nil
}

func CreateDockerContainer(containerName, projectPath, imageName string) error {

	// Check if a container with given name already exists
	command := fmt.Sprintf("--filter=name=%s", containerName)

	cmd := exec.Command("docker", "ps", "-a", "-q", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}

	if string(output) == "" {
		fmt.Println("Container does not exist, creating new one.")
	} else if err != nil {
		return err
	} else {
		fmt.Printf("Container %s exists.\n", containerName)
		return nil
	}
	// Create a new Docker container
	fullPathOfProject, _ := filepath.Abs(projectPath)
	command = fmt.Sprintf("docker create --name=%s -v '%s:/vectra' %s", containerName,
		fullPathOfProject, imageName)
	err = ExecuteCommand(command, true, true)
	if err != nil {
		return err
	}

	fmt.Printf("Successfully created and start Docker container %s\n", containerName)
	return nil
}

func StartDockerContainer(containerName string) error {
	return ExecuteCommand("docker start "+containerName, false, true)
}

func ExecuteCommand(command string, printStandardOutput bool, printErrorOutput bool) error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()

	if err != nil {
		if printErrorOutput {
			fmt.Printf("Error Output: %s\n", stderr.String())
		}
		return fmt.Errorf("failed to execute command: %w", err)
	}

	if printStandardOutput {
		fmt.Printf("Standard Output: %s\n", out.String())
	}

	return nil
}

func WatchFiles(rootFolder string, includePatterns, excludePatterns []string, delay int, task func(string)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create watcher: %w", err)
	}
	defer watcher.Close()

	done := make(chan bool)

	err = filepath.Walk(rootFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path: %w", err)
		}
		if info.IsDir() {
			err = watcher.Add(path)
			if err != nil {
				return fmt.Errorf("error adding path to watcher: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error setting watcher paths: %w", err)
	}

	go func() {
		timer := time.NewTimer(time.Duration(delay) * time.Millisecond)

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				filePath := event.Name

				include := false
				exclude := false

				for _, pattern := range includePatterns {
					matched, _ := regexp.MatchString(pattern, filePath)
					if matched {
						include = true
						break
					}
				}

				for _, pattern := range excludePatterns {
					matched, _ := regexp.MatchString(pattern, filePath)
					if matched {
						exclude = true
						break
					}
				}

				if include && !exclude && event.Op&fsnotify.Write == fsnotify.Write {
					timer.Stop()
					timer = time.AfterFunc(time.Duration(delay)*time.Millisecond, func() {
						task(filePath)
					})
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()

	<-done

	return nil
}
