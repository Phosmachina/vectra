package generator

import (
	"bytes"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

type WatcherConfig struct {
	PugConfig  `yaml:"pug_config"`
	SassConfig `yaml:"sass_config"`
	JsConfig   `yaml:"js_config"`
	I18nConfig `yaml:"i18n_config"`
}

type Watcher struct {
	IsEnabled bool `yaml:"is_enabled"`
}

type PugConfig struct {
	Watcher            `yaml:",inline"`
	LayoutFilePattern  string `yaml:"layout_file_pattern"`
	IgnoredFilePattern string `yaml:"ignored_file_pattern"`
}

type SassConfig struct {
	Watcher `yaml:",inline"`
}

type JsConfig struct {
	Watcher `yaml:",inline"`
}
type I18nConfig struct {
	Watcher `yaml:",inline"`
}

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
		fmt.Println("Container does not exist, creating new one...")
	} else if err != nil {
		return err
	} else {
		return nil
	}
	// Create a new Docker container
	fullPathOfProject, _ := filepath.Abs(projectPath)
	command = fmt.Sprintf("docker create --name=%s -v '%s:/vectra' %s", containerName,
		fullPathOfProject, imageName)
	err = ExecuteCommand(command, false, true)
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

// WatchFiles watches recursively a root folder for file changes
// and triggers a task when a file that matches the include patterns
// and does not match the exclude patterns is written.
//
// Parameters:
//
// - rootFolder: The root folder to watch for file changes.
//
// - includePatterns: The patterns for files to include in the watch.
// Only files that match any of the include patterns are considered.
//
// - excludePatterns: The patterns for files to exclude from the watch.
// Files that match any of the exclude patterns are ignored.
//
// - delay: The delay in milliseconds before triggering the task after a file is written.
//
// - task: The task to be executed when a file that matches the include patterns
// and does not match the exclude patterns is written.
//
// Returns:
//
// - error: An error if any occurred, otherwise nil.
func WatchFiles(rootFolder string, includePatterns, excludePatterns []string, delay int, task func(string)) error {

	// TODO handle remove event (for pug for example)

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

				// If a new directory is created, watch it.
				if info, err := os.Stat(filePath); err == nil && info.IsDir() {
					watcher.Add(filePath)
				}

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

func watchPug(v *Vectra) error {

	pugConfig := v.WatcherConfig.PugConfig
	layoutPattern, _ := regexp.Compile(pugConfig.LayoutFilePattern)
	ignoredPattern, _ := regexp.Compile(pugConfig.IgnoredFilePattern)
	root := filepath.Join(v.ProjectPath, "src", "view", "pug")

	pugFilesToBeCompiled := func() []string {
		var files []string
		if err := filepath.Walk(root, visit(&files, ".pug")); err != nil {
			fmt.Printf("error walking the path %v: %v\n", root, err)
			return nil
		}

		var filteredFiles []string
		for _, file := range files {
			// Exclude files that match layoutPattern or ignoredPattern
			if !layoutPattern.MatchString(file) && !ignoredPattern.MatchString(file) {
				filteredFiles = append(filteredFiles, file)
			}
		}

		return filteredFiles
	}

	compileInDocker := func(file string) {
		c := fmt.Sprintf(
			"docker exec %s jade -writer -pkg view -d /vectra/src/view/go /vectra/%s",
			v.ProjectName+"_Pug",
			file,
		)
		_ = ExecuteCommand(c, false, true)
	}

	return WatchFiles(
		root,
		[]string{".*\\.pug$"},
		[]string{".*completion_variable.*"},
		200,
		func(pth string) {
			relPth, _ := filepath.Rel(v.ProjectPath, pth)
			if layoutPattern.MatchString(relPth) {
				// Get all files excluding layout and ignored ones.
				files := pugFilesToBeCompiled()
				for _, file := range files {
					rel, _ := filepath.Rel(v.ProjectPath, file)
					compileInDocker(rel)
				}
			} else if !ignoredPattern.MatchString(relPth) {
				compileInDocker(relPth)
			} else {
				return // Avoid log.
			}
			log.Print("PUG ", relPth, " | Transpile DONE.")
		},
	)
}

func watchJS(v *Vectra) error {
	return WatchFiles(filepath.Join(v.ProjectPath, "static", "js"),
		[]string{"main.js$"},
		[]string{"prod"},
		200, func(pth string) {
			_ = ExecuteCommand(
				fmt.Sprintf("docker start %s_MinifyJS", v.ProjectName), false, true)

			relPth, _ := filepath.Rel(v.ProjectPath, pth)
			log.Print("JS ", relPth, " | Minify DONE.")
		},
	)
}

func watchI18n(v *Vectra) error {
	return WatchFiles(filepath.Join(v.ProjectPath, "data", "i18n"),
		[]string{".*en.*\\.ini$"},
		[]string{},
		200, func(pth string) {
			v.Generate("i18n")
			relPth, _ := filepath.Rel(v.ProjectPath, pth)
			log.Print("I18N helpers ", relPth, " | Generation DONE.")
		},
	)
}

func watchSass(v *Vectra) error {
	return WatchFiles(filepath.Join(v.ProjectPath, "static", "css"),
		[]string{".*\\.sass$", ".*\\.scss$"},
		[]string{},
		200, func(pth string) {
			_ = ExecuteCommand(
				fmt.Sprintf("docker start %s_Sass", v.ProjectName), false, true)
			_ = ExecuteCommand(
				fmt.Sprintf("docker start %s_Autoprefixer", v.ProjectName), false, true)
			time.Sleep(400 * time.Millisecond)
			_ = ExecuteCommand(
				fmt.Sprintf("docker start %s_MinifyCSS", v.ProjectName), false, true)

			relPth, _ := filepath.Rel(v.ProjectPath, pth)
			log.Print("CSS ", relPth, " | Sass, Autoprefixer, Minify DONE.")
		},
	)
}
