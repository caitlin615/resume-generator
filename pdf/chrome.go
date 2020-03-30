package pdf

import (
	"fmt"
	"log"
	"os"

	"github.com/fedesog/webdriver"
)

var chromePort = getEnvDefault("CHROME_PORT", "9222")

func startChrome() (chromeDriver *webdriver.ChromeDriver, err error) {
	log.Printf("##### starting chromeDriver")

	chromeDriver = webdriver.NewChromeDriver(getEnvDefault("CHROME_DRIVER_BINARY_PATH", "/usr/local/bin/chromedriver"))
	err = chromeDriver.Start()
	if err != nil {
		return
	}

	desired := webdriver.Capabilities{
		"platform":    "Linux",
		"browserName": "chrome",
		"chromeOptions": webdriver.Capabilities{
			// "debuggerAddress": "127.0.0.1:9222",
			"binary": getEnvDefault("CHROME_BINARY_PATH", "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"),
			// args: http://peter.sh/experiments/chromium-command-line-switches/
			"args": []string{
				// "window-size=800,800", // TODO: Make this customizable by the server
				"--no-sandbox",
				// "ignore-gpu-blacklist",
				"--disable-dev-shm-usage",
				"--disable-dev-shm-using",
				fmt.Sprintf("--remote-debugging-port=%s", chromePort),
				"--headless",
				"--disable-gpu",
			},
		},
		// https://github.com/SeleniumHQ/selenium/wiki/DesiredCapabilities
		// It appears that we still get *some* logging even though we said OFF. We
		// tried SEVERE and WARNING but they seemed to have no effect.
		"loggingPrefs": map[string]string{
			"driver": "OFF",
		},
	}
	required := webdriver.Capabilities{}
	_, err = chromeDriver.NewSession(desired, required)
	if err != nil {
		log.Printf("Failed to start chromeSession: %s", err)
		chromeDriver.Stop()
		return
	}
	return
}

func getEnvDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
