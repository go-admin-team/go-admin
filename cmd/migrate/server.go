package migrate

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/go-admin-team/go-admin-core/v2/sdk"
	"github.com/go-admin-team/go-admin-core/v2/sdk/pkg"
	"gorm.io/gorm"

	"github.com/go-admin-team/go-admin-core/v2/config/source/file"
	"github.com/go-admin-team/go-admin-core/v2/sdk/contract/app"
	"github.com/spf13/cobra"

	"github.com/go-admin-team/go-admin-core/v2/sdk/config"
	"go-admin/cmd/migrate/migration"
	_ "go-admin/cmd/migrate/migration/version"
	_ "go-admin/cmd/migrate/migration/version-local"
	"go-admin/common/database"
	"go-admin/common/models"
)

var (
	configYml string
	generate  bool
	goAdmin   bool
	host      string
	appCode   string
	dryRun    bool
	StartCmd  = &cobra.Command{
		Use:     "migrate",
		Short:   "Initialize the database",
		Example: "go-admin migrate -c config/settings.yml",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}
	statusCmd = &cobra.Command{
		Use:     "status",
		Short:   "List applied and pending migrations, grouped by app",
		Example: "go-admin migrate status -c config/settings.yml",
		Run: func(cmd *cobra.Command, args []string) {
			runStatus()
		},
	}
	// Under migrate rather than under the existing `app` command, which
	// already means "generate the skeleton of a new app" - a directory that
	// does not exist yet, not an application already compiled into this
	// binary. Installing an application is running its migrations, which is
	// what this command is; --app, --domain and resolveDB are all already
	// here, including the guard that refuses a mistyped code instead of
	// reporting a successful no-op.
	installCmd = &cobra.Command{
		Use:     "install <code>",
		Short:   "Install one application: run its migrations and record it in sys_app",
		Example: "go-admin migrate install order -c config/settings.yml",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runInstall(args[0])
		},
	}
	uninstallCmd = &cobra.Command{
		Use:     "uninstall <code>",
		Short:   "Remove one application's menus, apis and permission grants; its own tables are left alone",
		Example: "go-admin migrate uninstall order -c config/settings.yml",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runUninstall(args[0])
		},
	}
)

// fixme 在您看不见代码的时候运行迁移，我觉得是不安全的，所以编译后最好不要去执行迁移
func init() {
	StartCmd.PersistentFlags().StringVarP(&configYml, "config", "c", "config/settings.yml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().BoolVarP(&generate, "generate", "g", false, "generate migration file")
	StartCmd.PersistentFlags().BoolVarP(&goAdmin, "goAdmin", "a", false, "with -g, write the generated file to version/ instead of version-local/ (does not affect which migrations run)")
	StartCmd.PersistentFlags().StringVarP(&host, "domain", "d", "*", "select tenant host")

	// --app is deliberately long-only. -a already means "generate into
	// version/ rather than version-local/", which is about writing a template
	// file, not about which migrations run; giving the two the same letter
	// would be a trap.
	StartCmd.PersistentFlags().StringVar(&appCode, "app", "", "limit to the migrations of one app (\""+migration.FrameworkAppCode+"\" for the framework's own)")
	StartCmd.Flags().BoolVar(&dryRun, "dry-run", false, "list what would be applied, in order, and write nothing")

	StartCmd.AddCommand(statusCmd)
	StartCmd.AddCommand(installCmd)
	StartCmd.AddCommand(uninstallCmd)
}

func run() {

	if !generate {
		fmt.Println(`start init`)
		//1. 读取配置
		config.Setup(
			file.NewSource(file.WithPath(configYml)),
			initDB,
		)
	} else {
		fmt.Println(`generate migration file`)
		_ = genFile()
	}
}

// resolveDB picks the tenant database and hands it to the registry.
//
// It creates and alters nothing, which is what lets status and --dry-run share
// it: those two must be able to run against a production database without
// leaving a trace.
func resolveDB() (*gorm.DB, error) {
	if host == "" {
		host = "*"
	}
	db := sdk.Runtime.GetDbByTenant(host)
	if db == nil {
		if len(sdk.Runtime.GetAllDb()) == 1 && host == "*" {
			for k, v := range sdk.Runtime.GetAllDb() {
				db = v
				host = k
				break
			}
		}
	}
	if db == nil {
		return nil, fmt.Errorf("未找到数据库配置")
	}
	if config.DatabasesConfig[host].Driver == "mysql" {
		//初始化数据库时候用
		db.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4")
	}
	return db, nil
}

// exitUnlessAppRegistered ends the command when --app names something no
// migration was registered under.
//
// Every path took a typo as "nothing matched" and reported success: `migrate`
// printed that the app was unknown and still exited 0, while `--dry-run` and
// `status` said "nothing to apply" and "none recorded" - which is what an
// up-to-date database says too, so the output does not even hint at the typo.
// An operator running `go-admin migrate --app crmm && deploy` gets the deploy.
//
// Checked against the registry, which init() has already filled, so this runs
// before any database work and costs nothing. It lives in the command layer
// because the exit code does: the migration package stays callable from a test
// without taking the process down with it.
func exitUnlessAppRegistered() {
	if err := appRegistrationError(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// appRegistrationError carries the decision on its own so it can be tested;
// exitUnlessAppRegistered is only the os.Exit around it. Nil means --app was
// either empty or names a registered app.
func appRegistrationError() error {
	if appCode == "" {
		return nil
	}
	want := migration.DisplayAppCode(migration.AppFilter(appCode))
	registered := migration.Migrate.AppCodes()
	for _, c := range registered {
		if c == want {
			return nil
		}
	}
	return fmt.Errorf("no migrations are registered for app %q; registered: %s",
		want, strings.Join(registered, ", "))
}

func migrateModel() error {
	db, err := resolveDB()
	if err != nil {
		return err
	}
	// sys_migration is the one table that never goes through a versioned
	// migration - it is the table that records them. AutoMigrate realigns it
	// on every run, which is how the app_code column reaches an existing
	// database without anyone writing a migration for it.
	if err = db.Debug().AutoMigrate(&models.Migration{}); err != nil {
		return err
	}
	migration.Migrate.SetDb(db.Debug())
	if appCode != "" {
		return migration.Migrate.MigrateApp(appCode)
	}
	return migration.Migrate.Migrate()
}

func initDB() {
	// Before the database is touched, so a typo cannot get as far as looking
	// like a successful no-op on either path below.
	exitUnlessAppRegistered()

	//3. 初始化数据库链接
	database.Setup()

	if dryRun {
		db, err := resolveDB()
		if err != nil {
			fmt.Println(err)
			return
		}
		migration.Migrate.SetDb(db)
		entries, err := migration.Migrate.Status()
		if err != nil {
			fmt.Println(err)
			return
		}
		if err = printPending(os.Stdout, entries, appCode); err != nil {
			fmt.Println(err)
		}
		return
	}

	//4. 数据库迁移
	fmt.Println("数据库迁移开始")
	exitOnError(os.Stderr, migrateModel())
	fmt.Println(`数据库基础数据初始化成功`)
}

// exitOnError ends the command non-zero when the migration did not go through.
//
// A caller that migrates before starting a server decides whether to go ahead
// on the exit code alone - the deploy workflow does exactly that. Every path
// out of migrateModel used to return without one: an unreachable tenant
// database or a failed AutoMigrate printed a line and exited 0, so a
// deployment carried on onto a schema that had not been brought forward. A
// failing migration function was the only one reported, and only because it
// ended the process from inside the migration engine - which is the call this
// batch moved out here, so without this the last reported failure would have
// stopped being reported too.
//
// Split from the exit itself, the way appRegistrationError is split from
// exitUnlessAppRegistered, so what it decides can be tested without a
// subprocess. osExit is a variable for the same reason.
func exitOnError(w io.Writer, err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(w, err)
	osExit(1)
}

// osExit is a variable so a test can watch the decision without ending the
// test binary; origExit is what it is put back to.
var (
	osExit   = os.Exit
	origExit = os.Exit
)

func runStatus() {
	config.Setup(
		file.NewSource(file.WithPath(configYml)),
		func() {
			exitUnlessAppRegistered()

			database.Setup()
			db, err := resolveDB()
			if err != nil {
				fmt.Println(err)
				return
			}
			migration.Migrate.SetDb(db)
			entries, err := migration.Migrate.Status()
			if err != nil {
				fmt.Println(err)
				return
			}
			// Which applications exist is a different question from which
			// migrations ran, and an install that stopped partway is only
			// visible in the answer to the first.
			apps, err := loadApps(db)
			if err != nil {
				fmt.Println(err)
				return
			}
			if err = printStatus(os.Stdout, entries, apps, appCode); err != nil {
				fmt.Println(err)
			}
		},
	)
}

func runInstall(code string) {
	config.Setup(
		file.NewSource(file.WithPath(configYml)),
		func() {
			database.Setup()
			db, err := resolveDB()
			if err != nil {
				exitOnError(os.Stderr, err)
				return
			}
			registered := app.Snapshot()
			m, err := manifestFor(registered, code)
			if err != nil {
				exitOnError(os.Stderr, err)
				return
			}
			// Over every registered manifest, not just this one's closure: a
			// cycle between two other applications is still an authoring
			// mistake, and the day somebody installs into it is the worse
			// time to find out.
			if err := refuseOnDependencyCycle(registered); err != nil {
				exitOnError(os.Stderr, err)
				return
			}
			rep, err := install(db, migration.Migrate, m)
			if err != nil {
				exitOnError(os.Stderr, err)
				return
			}
			reportInstall(os.Stdout, rep)
		},
	)
}

func runUninstall(code string) {
	config.Setup(
		file.NewSource(file.WithPath(configYml)),
		func() {
			database.Setup()
			db, err := resolveDB()
			if err != nil {
				exitOnError(os.Stderr, err)
				return
			}
			// No manifest lookup. An application whose code has already been
			// taken out of the binary registers nothing, and that is exactly
			// when somebody needs to clear its rows out of the database.
			rep, err := uninstall(db, code)
			if err != nil {
				exitOnError(os.Stderr, err)
				return
			}
			reportUninstall(os.Stdout, rep)
		},
	)
}

func genFile() error {
	t1, err := template.ParseFiles("template/migrate.template")
	if err != nil {
		return err
	}
	m := map[string]string{}
	m["GenerateTime"] = strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	m["Package"] = "version_local"
	if goAdmin {
		m["Package"] = "version"
	}
	var b1 bytes.Buffer
	err = t1.Execute(&b1, m)
	if goAdmin {
		pkg.FileCreate(b1, "./cmd/migrate/migration/version/"+m["GenerateTime"]+"_migrate.go")
	} else {
		pkg.FileCreate(b1, "./cmd/migrate/migration/version-local/"+m["GenerateTime"]+"_migrate.go")
	}
	return nil
}
