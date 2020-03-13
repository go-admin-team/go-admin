// package rotatelogs is a port of File-RotateLogs from Perl
// (https://metacpan.org/release/File-RotateLogs), and it allows
// you to automatically rotate output files when you write to them
// according to the filename pattern that you can specify.
package rotatelogs

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	strftime "github.com/lestrrat/go-strftime"
	"github.com/pkg/errors"
)

func (c clockFn) Now() time.Time {
	return c()
}

func (o OptionFn) Configure(rl *RotateLogs) error {
	return o(rl)
}

// WithClock creates a new Option that sets a clock
// that the RotateLogs object will use to determine
// the current time.
//
// By default rotatelogs.Local, which returns the
// current time in the local time zone, is used. If you
// would rather use UTC, use rotatelogs.UTC as the argument
// to this option, and pass it to the constructor.
func WithClock(c Clock) Option {
	return OptionFn(func(rl *RotateLogs) error {
		rl.clock = c
		return nil
	})
}

// WithLocation creates a new Option that sets up a
// "Clock" interface that the RotateLogs object will use
// to determine the current time.
//
// This optin works by always returning the in the given
// location.
func WithLocation(loc *time.Location) Option {
	return WithClock(clockFn(func() time.Time {
		return time.Now().In(loc)
	}))
}

// WithLinkName creates a new Option that sets the
// symbolic link name that gets linked to the current
// file name being used.
func WithLinkName(s string) Option {
	return OptionFn(func(rl *RotateLogs) error {
		rl.linkName = s
		return nil
	})
}

// WithMaxAge creates a new Option that sets the
// max age of a log file before it gets purged from
// the file system.
func WithMaxAge(d time.Duration) Option {
	return OptionFn(func(rl *RotateLogs) error {
		if rl.rotationCount > 0 && d > 0 {
			return errors.New("attempt to set MaxAge when RotationCount is also given")
		}
		rl.maxAge = d
		return nil
	})
}

// WithRotationTime creates a new Option that sets the
// time between rotation.
func WithRotationTime(d time.Duration) Option {
	return OptionFn(func(rl *RotateLogs) error {
		rl.rotationTime = d
		return nil
	})
}

// WithRotationCount creates a new Option that sets the
// number of files should be kept before it gets
// purged from the file system.
func WithRotationCount(n int) Option {
	return OptionFn(func(rl *RotateLogs) error {
		if rl.maxAge > 0 && n > 0 {
			return errors.New("attempt to set RotationCount when MaxAge is also given")
		}
		rl.rotationCount = n
		return nil
	})
}

// New creates a new RotateLogs object. A log filename pattern
// must be passed. Optional `Option` parameters may be passed
func New(pattern string, options ...Option) (*RotateLogs, error) {
	globPattern := pattern
	for _, re := range patternConversionRegexps {
		globPattern = re.ReplaceAllString(globPattern, "*")
	}

	strfobj, err := strftime.New(pattern)
	if err != nil {
		return nil, errors.Wrap(err, `invalid strftime pattern`)
	}

	var rl RotateLogs
	rl.clock = Local
	rl.globPattern = globPattern
	rl.pattern = strfobj
	rl.rotationTime = 24 * time.Hour
	// Keeping forward compatibility, maxAge is prior to rotationCount.
	rl.maxAge = 7 * 24 * time.Hour
	rl.rotationCount = -1
	for _, opt := range options {
		opt.Configure(&rl)
	}

	return &rl, nil
}

func (rl *RotateLogs) genFilename() string {
	now := rl.clock.Now()
	diff := time.Duration(now.UnixNano()) % rl.rotationTime
	t := now.Add(time.Duration(-1 * diff))
	return rl.pattern.FormatString(t)
}

// Write satisfies the io.Writer interface. It writes to the
// appropriate file handle that is currently being used.
// If we have reached rotation time, the target file gets
// automatically rotated, and also purged if necessary.
func (rl *RotateLogs) Write(p []byte) (n int, err error) {
	// Guard against concurrent writes
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	out, err := rl.getTargetWriter()
	if err != nil {
		return 0, errors.Wrap(err, `failed to acquite target io.Writer`)
	}

	return out.Write(p)
}

// must be locked during this operation
func (rl *RotateLogs) getTargetWriter() (io.Writer, error) {
	// This filename contains the name of the "NEW" filename
	// to log to, which may be newer than rl.currentFilename
	filename := rl.genFilename()
	if rl.curFn == filename {
		// nothing to do
		return rl.outFh, nil
	}

	// if we got here, then we need to create a file
	fh, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return nil, errors.Errorf("failed to open file %s: %s", rl.pattern, err)
	}

	if err := rl.rotate(filename); err != nil {
		// Failure to rotate is a problem, but it's really not a great
		// idea to stop your application just because you couldn't rename
		// your log. For now, we're just going to punt it and write to
		// os.Stderr
		fmt.Fprintf(os.Stderr, "failed to rotate: %s\n", err)
	}

	rl.outFh.Close()
	rl.outFh = fh
	rl.curFn = filename

	return fh, nil
}

// CurrentFileName returns the current file name that
// the RotateLogs object is writing to
func (rl *RotateLogs) CurrentFileName() string {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()
	return rl.curFn
}

var patternConversionRegexps = []*regexp.Regexp{
	regexp.MustCompile(`%[%+A-Za-z]`),
	regexp.MustCompile(`\*+`),
}

type cleanupGuard struct {
	enable bool
	fn     func()
	mutex  sync.Mutex
}

func (g *cleanupGuard) Enable() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.enable = true
}
func (g *cleanupGuard) Run() {
	g.fn()
}

func (rl *RotateLogs) rotate(filename string) error {
	lockfn := filename + `_lock`
	fh, err := os.OpenFile(lockfn, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		// Can't lock, just return
		return err
	}

	var guard cleanupGuard
	guard.fn = func() {
		fh.Close()
		os.Remove(lockfn)
	}
	defer guard.Run()

	if rl.linkName != "" {
		tmpLinkName := filename + `_symlink`
		if err := os.Symlink(filename, tmpLinkName); err != nil {
			return errors.Wrap(err, `failed to create new symlink`)
		}

		if err := os.Rename(tmpLinkName, rl.linkName); err != nil {
			return errors.Wrap(err, `failed to rename new symlink`)
		}
	}

	if rl.maxAge <= 0 && rl.rotationCount <= 0 {
		return errors.New("panic: maxAge and rotationCount are both set")
	}

	matches, err := filepath.Glob(rl.globPattern)
	if err != nil {
		return err
	}

	cutoff := rl.clock.Now().Add(-1 * rl.maxAge)
	var toUnlink []string
	for _, path := range matches {
		// Ignore lock files
		if strings.HasSuffix(path, "_lock") || strings.HasSuffix(path, "_symlink") {
			continue
		}

		fi, err := os.Stat(path)
		if err != nil {
			continue
		}

		fl, err := os.Lstat(path)
		if err != nil {
			continue
		}

		if rl.maxAge > 0 && fi.ModTime().After(cutoff) {
			continue
		}

		if rl.rotationCount > 0 && fl.Mode()&os.ModeSymlink == os.ModeSymlink {
			continue
		}
		toUnlink = append(toUnlink, path)
	}

	if rl.rotationCount > 0 {
		// Only delete if we have more than rotationCount
		if rl.rotationCount >= len(toUnlink) {
			return nil
		}

		toUnlink = toUnlink[:len(toUnlink)-rl.rotationCount]
	}

	if len(toUnlink) <= 0 {
		return nil
	}

	guard.Enable()
	go func() {
		// unlink files on a separate goroutine
		for _, path := range toUnlink {
			os.Remove(path)
		}
	}()

	return nil
}

// Close satisfies the io.Closer interface. You must
// call this method if you performed any writes to
// the object.
func (rl *RotateLogs) Close() error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if rl.outFh == nil {
		return nil
	}

	rl.outFh.Close()
	rl.outFh = nil
	return nil
}
