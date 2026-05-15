package depcheck

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
)

// DepCheck checks and optionally starts required dependencies
type DepCheck struct {
	log *zap.Logger
}

func New(log *zap.Logger) *DepCheck {
	return &DepCheck{log: log}
}

// CheckResult holds the status of a dependency check
type CheckResult struct {
	Name      string
	Installed bool
	Running   bool
	Error     string
}

// CheckAll checks all dependencies
func (d *DepCheck) CheckAll(mysqlAddr, redisAddr, zlmAddr string) []CheckResult {
	var results []CheckResult

	results = append(results, d.checkMySQL(mysqlAddr))
	results = append(results, d.checkRedis(redisAddr))
	results = append(results, d.checkZLMediaKit(zlmAddr))

	return results
}

// TryStart attempts to start a dependency if not running
func (d *DepCheck) TryStart(result CheckResult) error {
	if result.Running {
		return nil
	}
	if !result.Installed {
		return fmt.Errorf("%s is not installed, please install it first", result.Name)
	}

	switch result.Name {
	case "MySQL":
		return d.startMySQL()
	case "Redis":
		return d.startRedis()
	case "ZLMediaKit":
		return d.startZLMediaKit()
	}
	return nil
}

func (d *DepCheck) checkMySQL(addr string) CheckResult {
	result := CheckResult{Name: "MySQL"}

	// Check if mysqld or mysql is installed
	if _, err := exec.LookPath("mysqld"); err != nil {
		if _, err2 := exec.LookPath("mysql"); err2 != nil {
			if _, err3 := exec.LookPath("docker"); err3 == nil {
				result.Installed = true // Can use docker
			} else {
				result.Error = "MySQL is not installed"
				return result
			}
		}
	}
	result.Installed = true

	// Check if MySQL is running
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err == nil {
		conn.Close()
		result.Running = true
	}

	return result
}

func (d *DepCheck) checkRedis(addr string) CheckResult {
	result := CheckResult{Name: "Redis"}

	if _, err := exec.LookPath("redis-server"); err != nil {
		if _, err2 := exec.LookPath("docker"); err2 == nil {
			result.Installed = true
		} else {
			result.Error = "Redis is not installed"
			return result
		}
	}
	result.Installed = true

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err == nil {
		conn.Close()
		result.Running = true
	}

	return result
}

func (d *DepCheck) checkZLMediaKit(addr string) CheckResult {
	result := CheckResult{Name: "ZLMediaKit"}

	if _, err := exec.LookPath("MediaServer"); err != nil {
		if _, err2 := exec.LookPath("docker"); err2 == nil {
			result.Installed = true
		} else {
			result.Error = "ZLMediaKit is not installed"
			return result
		}
	}
	result.Installed = true

	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err == nil {
		conn.Close()
		result.Running = true
	}

	return result
}

func (d *DepCheck) startMySQL() error {
	d.log.Info("attempting to start MySQL...")

	// Try systemctl first
	if out, err := exec.Command("systemctl", "start", "mysql").CombinedOutput(); err == nil {
		d.log.Info("MySQL started via systemctl")
		return nil
	} else {
		d.log.Debug("systemctl start mysql failed", zap.String("output", string(out)))
	}

	// Try mysqld service
	if out, err := exec.Command("systemctl", "start", "mysqld").CombinedOutput(); err == nil {
		d.log.Info("MySQL started via systemctl (mysqld)")
		return nil
	} else {
		d.log.Debug("systemctl start mysqld failed", zap.String("output", string(out)))
	}

	// Try service command (Windows)
	if runtime.GOOS == "windows" {
		if out, err := exec.Command("net", "start", "MySQL").CombinedOutput(); err == nil {
			d.log.Info("MySQL started via net start")
			return nil
		} else {
			d.log.Debug("net start MySQL failed", zap.String("output", string(out)))
		}
	}

	// Try docker
	if _, err := exec.LookPath("docker"); err == nil {
		d.log.Info("starting MySQL via docker...")
		cmd := exec.Command("docker", "start", "mysql")
		if out, err := cmd.CombinedOutput(); err == nil {
			d.log.Info("MySQL docker container started")
			return nil
		} else {
			d.log.Debug("docker start mysql failed", zap.String("output", string(out)))
		}
	}

	return fmt.Errorf("failed to start MySQL automatically, please start it manually")
}

func (d *DepCheck) startRedis() error {
	d.log.Info("attempting to start Redis...")

	// Try systemctl
	if out, err := exec.Command("systemctl", "start", "redis").CombinedOutput(); err == nil {
		d.log.Info("Redis started via systemctl")
		return nil
	} else {
		d.log.Debug("systemctl start redis failed", zap.String("output", string(out)))
	}

	// Try redis-server directly
	if path, err := exec.LookPath("redis-server"); err == nil {
		d.log.Info("starting Redis via redis-server...")
		cmd := exec.Command(path, "--daemonize", "yes")
		if out, err := cmd.CombinedOutput(); err == nil {
			d.log.Info("Redis started via redis-server")
			return nil
		} else {
			d.log.Debug("redis-server --daemonize failed", zap.String("output", string(out)))
		}
	}

	// Try docker
	if _, err := exec.LookPath("docker"); err == nil {
		d.log.Info("starting Redis via docker...")
		cmd := exec.Command("docker", "start", "redis")
		if out, err := cmd.CombinedOutput(); err == nil {
			d.log.Info("Redis docker container started")
			return nil
		} else {
			d.log.Debug("docker start redis failed", zap.String("output", string(out)))
		}
	}

	return fmt.Errorf("failed to start Redis automatically, please start it manually")
}

func (d *DepCheck) startZLMediaKit() error {
	d.log.Info("attempting to start ZLMediaKit...")

	// Try docker
	if _, err := exec.LookPath("docker"); err == nil {
		d.log.Info("starting ZLMediaKit via docker...")
		cmd := exec.Command("docker", "start", "zlmediakit")
		if out, err := cmd.CombinedOutput(); err == nil {
			d.log.Info("ZLMediaKit docker container started")
			return nil
		} else {
			d.log.Debug("docker start zlmediakit failed", zap.String("output", string(out)))
		}
	}

	return fmt.Errorf("failed to start ZLMediaKit automatically, please start it manually")
}

// PrintInstallGuide prints installation instructions for missing dependencies
func PrintInstallGuide(results []CheckResult) {
	var missing []string
	for _, r := range results {
		if !r.Installed {
			missing = append(missing, r.Name)
		}
	}

	if len(missing) == 0 {
		return
	}

	fmt.Println()
	fmt.Println("=====================================================")
	fmt.Println("  Missing Dependencies - Installation Guide")
	fmt.Println("=====================================================")
	fmt.Println()

	for _, name := range missing {
		switch name {
		case "MySQL":
			fmt.Println("MySQL:")
			fmt.Println("  Ubuntu/Debian:  sudo apt-get install mysql-server")
			fmt.Println("  CentOS/RHEL:    sudo yum install mysql-server")
			fmt.Println("  macOS:          brew install mysql")
			fmt.Println("  Docker:         docker run -d --name mysql -e MYSQL_ROOT_PASSWORD=Root@123456 -p 3306:3306 mysql:8.0")
			fmt.Println()
		case "Redis":
			fmt.Println("Redis:")
			fmt.Println("  Ubuntu/Debian:  sudo apt-get install redis-server")
			fmt.Println("  CentOS/RHEL:    sudo yum install redis")
			fmt.Println("  macOS:          brew install redis")
			fmt.Println("  Docker:         docker run -d --name redis -p 6379:6379 redis:7")
			fmt.Println()
		case "ZLMediaKit":
			fmt.Println("ZLMediaKit:")
			fmt.Println("  Docker (recommended):")
			fmt.Println("  docker run -d --name zlmediakit -p 8080:80 -p 8443:443 -p 1935:1935 \\")
			fmt.Println("    -p 554:554 -p 554:554/udp -p 8000:8000/udp -p 9000:9000/udp \\")
			fmt.Println("    -p 10000-10500:10000-10500/udp -p 10000-10500:10000-10500/tcp \\")
			fmt.Println("    zlmediakit/zlmediakit:master")
			fmt.Println()
		}
	}

	fmt.Println("After installation, update configs/config.yaml with your settings.")
	fmt.Println("=====================================================")
	fmt.Println()
}

// PrintStartGuide prints instructions for starting non-running dependencies
func PrintStartGuide(results []CheckResult) {
	var notRunning []CheckResult
	for _, r := range results {
		if r.Installed && !r.Running {
			notRunning = append(notRunning, r)
		}
	}

	if len(notRunning) == 0 {
		return
	}

	fmt.Println()
	fmt.Println("=====================================================")
	fmt.Println("  Starting Dependencies...")
	fmt.Println("=====================================================")
	fmt.Println()

	for _, r := range notRunning {
		fmt.Printf("  %s: not running, attempting auto-start...\n", r.Name)
	}
	fmt.Println()
}

// IsLinux checks if running on Linux
func IsLinux() bool {
	return runtime.GOOS == "linux"
}

// IsMacOS checks if running on macOS
func IsMacOS() bool {
	return runtime.GOOS == "darwin"
}

// IsWindows checks if running on Windows
func IsWindows() bool {
	return runtime.GOOS == "windows"
}

// GetPackageManager returns the appropriate package manager command
func GetPackageManager() string {
	if IsLinux() {
		if _, err := exec.LookPath("apt-get"); err == nil {
			return "apt-get"
		}
		if _, err := exec.LookPath("yum"); err == nil {
			return "yum"
		}
		if _, err := exec.LookPath("dnf"); err == nil {
			return "dnf"
		}
		return "apt-get"
	}
	if IsMacOS() {
		return "brew"
	}
	return ""
}

// InstallCommand returns the install command for a package
func InstallCommand(pkgName string) string {
	pm := GetPackageManager()
	switch pm {
	case "apt-get":
		return fmt.Sprintf("sudo apt-get install -y %s", pkgName)
	case "yum":
		return fmt.Sprintf("sudo yum install -y %s", pkgName)
	case "dnf":
		return fmt.Sprintf("sudo dnf install -y %s", pkgName)
	case "brew":
		return fmt.Sprintf("brew install %s", pkgName)
	}
	return ""
}

// ServiceCommand returns the service start command
func ServiceCommand(serviceName string) string {
	if _, err := exec.LookPath("systemctl"); err == nil {
		return fmt.Sprintf("sudo systemctl start %s", serviceName)
	}
	if _, err := exec.LookPath("service"); err == nil {
		return fmt.Sprintf("sudo service %s start", serviceName)
	}
	return fmt.Sprintf("# Please start %s manually", serviceName)
}

// QuickDockerInstall prints docker-based install commands
func QuickDockerInstall() string {
	var cmds []string
	cmds = append(cmds, "# Quick setup with Docker:")
	cmds = append(cmds, "")
	cmds = append(cmds, "# 1. MySQL")
	cmds = append(cmds, `docker run -d --name mysql -e MYSQL_ROOT_PASSWORD=Root@123456 -p 3306:3306 mysql:8.0 --default-authentication-plugin=mysql_native_password`)
	cmds = append(cmds, "")
	cmds = append(cmds, "# 2. Redis")
	cmds = append(cmds, `docker run -d --name redis -p 6379:6379 redis:7`)
	cmds = append(cmds, "")
	cmds = append(cmds, "# 3. ZLMediaKit")
	cmds = append(cmds, `docker run -d --name zlmediakit -p 8080:80 -p 8443:443 -p 1935:1935 \`)
	cmds = append(cmds, `  -p 554:554 -p 554:554/udp -p 8000:8000/udp -p 9000:9000/udp \`)
	cmds = append(cmds, `  -p 10000-10500:10000-10500/udp -p 10000-10500:10000-10500/tcp \`)
	cmds = append(cmds, `  zlmediakit/zlmediakit:master`)
	cmds = append(cmds, "")
	cmds = append(cmds, "# 4. Create MySQL database")
	cmds = append(cmds, `docker exec -i mysql mysql -uroot -pRoot@123456 -e "CREATE DATABASE IF NOT EXISTS wvp DEFAULT CHARACTER SET utf8mb4;"`)
	return strings.Join(cmds, "\n")
}
