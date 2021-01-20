package flywheel

import (
	"bufio"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/aluttik/go-crossplane"
)

// OverrideProvider produces nginx args for a given directive
type OverrideProvider interface {
	Override(ctx context.Context, directive string, path string) ([]string, error)
	Close() error
}

// UpdatedFile is a file with a reference to it's original location
type UpdatedFile struct {
	*os.File
	OGName string
}

// Rename renames the file to the original location
func (f UpdatedFile) Rename() error {
	return os.Rename(f.File.Name(), f.OGName)
}

// WritePayload writes a payload to a file
func WritePayload(p *crossplane.Payload, options *crossplane.BuildOptions) error {
	for _, c := range p.Config {
		err := writeConfig(c.File, c, options)
		if err != nil {
			return fmt.Errorf("failed to handle config: %w", err)
		}
	}

	return nil
}

// WritePayloadTmp writes the output to a tempdir and returns the files
//
// A caller may want to remove the files or `Rename()` them to their intended location.
func WritePayloadTmp(p *crossplane.Payload, options *crossplane.BuildOptions) ([]UpdatedFile, error) {
	tmpDir, err := ioutil.TempDir("", "nginx-flywheel-")
	if err != nil {
		return nil, fmt.Errorf("failed to create tmpdir: %w", err)
	}

	tmpFiles := make([]UpdatedFile, len(p.Config))
	for i, c := range p.Config {
		// Force flat structure; files can be given correct structure with `Rename`
		f, err := os.Create(filepath.Join(tmpDir, filepath.Base(c.File)))
		if err != nil {
			return tmpFiles, fmt.Errorf("failed to create file: %w", err)
		}
		// This defer should fail under normal operation
		defer f.Close()
		tmpFiles[i] = UpdatedFile{f, c.File}

		err = writeConfig(f.Name(), c, options)
		if err != nil {
			return tmpFiles, fmt.Errorf("failed to handle config: %w", err)
		}
	}

	return tmpFiles, nil
}

// writeConfig writes config to a file
//
// This uses a string instead of a file descriptor to match parse, which handles fd as it's what's
// searching for the NGINX conf files.
func writeConfig(f string, c crossplane.Config, options *crossplane.BuildOptions) error {
	fd, err := os.Create(f)
	if err != nil {
		return fmt.Errorf("failed to create NGINX config file: %w", err)
	}
	defer fd.Close()
	w := bufio.NewWriter(fd)
	err = crossplane.Build(w, c, options)
	if err != nil {
		return fmt.Errorf("failed to write NGINX file: %w", err)
	}
	err = w.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}
	if err = fd.Sync(); err != nil {
		return fmt.Errorf("failed to sync write: %w", err)
	}
	return nil
}

// OverridePayload overrides each config in the payload
func OverridePayload(ctx context.Context, p *crossplane.Payload, o OverrideProvider) error {
	for i := range p.Config {
		config := p.Config[i]
		err := overrideDirectives(ctx, &config.Parsed, o, config.File)
		if err != nil {
			return err
		}
	}
	return nil
}

func overrideDirectives(ctx context.Context, ds *[]crossplane.Directive, o OverrideProvider, abspath string) error {
	if ds == nil {
		return fmt.Errorf("directive list is nil for: %v", abspath)
	}
	dsValues := *ds
	for i := range dsValues {
		err := overrideDirective(ctx, &dsValues[i], o, abspath)
		if err != nil {
			return err
		}
	}
	return nil
}

// overrideDirective overrides a single directives args
func overrideDirective(ctx context.Context, d *crossplane.Directive, o OverrideProvider, abspath string) error {
	if d == nil {
		return fmt.Errorf("directive is nil for: %v", abspath)
	}
	if d.IsComment() {
		return nil
	}
	args, err := o.Override(ctx, d.Directive, abspath)
	if err != nil {
		return err
	}
	if len(args) != 0 {
		d.Args = args
	}
	if d.IsBlock() {
		if d.Block != nil {
			overrideDirectives(ctx, d.Block, o, abspath)
		}
	}

	return nil
}

type new string

var New new = "NEW"

func newDirective(ctx context.Context, p *crossplane.Payload) {

}
