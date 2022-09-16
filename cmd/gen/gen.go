package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"
	"time"

	sprig "github.com/Masterminds/sprig/v3"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/moby/buildkit/frontend/dockerfile/instructions"
	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

const (
	topComment = "# Generated by docker-gen version %s, build %s\n# At: %s\n\n"
)

var (
	// Version is the version of the binary.
	Version = "dev"
	// GitCommit is the git commit that was compiled
	GitCommit string
)

type DockerfileValues struct {
	// AppName is the name of the project/application, will be used as the image name
	AppName string `yaml:"appName"`
	// AppVersion is the version of the project/application, will be used as the image tag
	AppVersion string `yaml:"appVersion"`
	// MainPackage is the build target. Should default to "." which will build the current
	// directory's main package.
	MainPackage string `yaml:"mainPackage"`
	// BuilderImage is the image used to build the application.
	BuilderImage string `yaml:"builderImage"`
	// RuntTimeImage is the docker image and tag of the final runtime container.
	RuntimeImage string `yaml:"runtimeImage"`
}

func main() {

	rootCmd := &cobra.Command{
		Use:   "docker-gen",
		Short: "Generate Dockerfiles from templates",
		Long:  `docker-gen generates validated Dockerfiles from go templated Dockerfiles.`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version of docker-gen",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("docker-gen version %s, build %s\n", Version, GitCommit)
		},
	}

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Renders Dockerfiles from templates with provided values",

		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("requires at least one argument")
			}

			return nil
		},
		RunE: renderTemplate,
	}
	rootCmd.AddCommand(versionCmd)

	runCmd.PersistentFlags().StringP("template-dir", "t", "", "Directory containing templates")
	runCmd.PersistentFlags().StringP("data-file", "d", "", "File containing data to use in template")
	runCmd.PersistentFlags().StringP("output-dir", "o", "", "Directory to write rendered Dockerfiles to. Defaults to data directory if not piped to something else.")
	runCmd.PersistentFlags().BoolP("stdout", "s", false, "Write rendered Dockerfiles to stdout")
	runCmd.PersistentFlags().BoolP("force", "f", false, "Print rendered Dockerfiles even if they don't pass validation")

	runCmd.MarkFlagRequired("template-dir")
	runCmd.MarkFlagRequired("data-file")

	rootCmd.AddCommand(runCmd)
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func renderTemplate(cmd *cobra.Command, args []string) error {
	templateDir, err := cmd.Flags().GetString("template-dir")
	if err != nil {
		return err
	}
	if templateDir == "" {
		return errors.New("template-dir is required")
	}

	dataFile, err := cmd.Flags().GetString("data-file")
	if err != nil {
		return err
	}

	// Grab stdin to check if we are piping in the values
	fi, _ := os.Stdin.Stat()

	if (fi.Mode() & os.ModeCharDevice) == 0 {
		dataFile = os.Stdin.Name()
	} else if dataFile == "" {
		return errors.New("data-file is required, or data must be piped through stdin")
	}

	stdout, err := cmd.Flags().GetBool("stdout")
	if err != nil {
		return err
	}

	// Grab stdout to check if are piping out to something
	fo, _ := os.Stdout.Stat()
	var pipedOut bool
	if (fo.Mode() & os.ModeCharDevice) == 0 {
		pipedOut = true
	}

	outputDir, err := cmd.Flags().GetString("output-dir")
	if err != nil {
		return err
	}

	// If an output dir wasn't provided and the data isn't coming from stdin, AND we're not piping it out, use the data file's directory.
	if outputDir == "" && dataFile != os.Stdin.Name() && !pipedOut {
		outputDir = filepath.Dir(dataFile)
	}

	// If an output dir wasn't provided, and there's no datafile (because we got it from stdin) AND we're not piping to stdout
	// then we return an error.
	if outputDir == "" && (dataFile == os.Stdin.Name() || dataFile == "") && (!stdout && !pipedOut) {
		return errors.New("output-dir is required if data is piped in and it's not piped out")
	}

	force, err := cmd.Flags().GetBool("force")
	if err != nil {
		return err
	}

	cleanDir := filepath.Clean(templateDir)
	tpl, err := template.New("base").Funcs(sprig.FuncMap()).ParseGlob(fmt.Sprintf("%s/*.Dockerfile", cleanDir))
	if err != nil {
		return fmt.Errorf("failed to read template(s), %w", err)
	}

	f, err := os.Open(dataFile)
	if err != nil {
		return fmt.Errorf("failed to open data file, %w", err)
	}
	data, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read data file, %w", err)
	}

	var values DockerfileValues
	err = yaml.Unmarshal(data, &values)
	if err != nil {
		return fmt.Errorf("failed to unmarshal yaml from data file, %w", err)
	}

	var b bytes.Buffer
	// Iterate through the args and render the specified templates
	for _, arg := range args {
		finalName := filepath.Clean(arg)
		// Check for ".Dockerfile" suffix
		if filepath.Ext(finalName) != ".Dockerfile" {
			finalName = fmt.Sprintf("%s.Dockerfile", finalName)
		}
		// Look up the template by name and emit a nicer error if it doesn't exist
		t := tpl.Lookup(finalName)
		if t == nil {
			// We exit even if one of multiple templates are missing
			// because we don't want unexpected output.
			return fmt.Errorf("template %s not found", arg)
		}
		// Render the template

		// Add comment to the top of the file
		b.WriteString(fmt.Sprintf(topComment, Version, GitCommit, time.Now().Format(time.RFC3339)))

		err = t.Execute(&b, values)
		if err != nil {
			return fmt.Errorf("failed to render template, %w", err)
		}
		b.WriteString("\n")

		// Validate the Dockerfile
		dockerFile := b.Bytes()
		err = parseDockerfile(dockerFile)
		if err != nil {
			errMsg := fmt.Errorf("failed to validate Dockerfile, %w", err)
			if !force {
				return errMsg
			}
			fmt.Println(errMsg)
		}
	}

	if outputDir != "" {
		err := os.WriteFile(filepath.Join(outputDir, "Dockerfile"), b.Bytes(), 0644)
		if err != nil {
			return fmt.Errorf("failed to write Dockerfile, %w", err)
		}
	}
	if stdout || pipedOut {
		fmt.Println(b.String())
	}
	return nil
}

func parseDockerfile(dockerFile []byte) error {
	results, err := parser.Parse(bytes.NewReader(dockerFile))
	if err != nil {
		return fmt.Errorf("failed to parse Dockerfile, %w", err)
	}

	_, _, err = instructions.Parse(results.AST)
	if err != nil {
		return fmt.Errorf("failed to parse Dockerfile instructions, %w", err)
	}
	return nil
}
