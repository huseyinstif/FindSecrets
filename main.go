package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"regexp"
)

// Result struct defines the structure for pattern match results
type Result struct {
	Pattern  string `json:"pattern"`
	FilePath string `json:"file_path"`
	Line     int    `json:"line"`
}

var results []Result

// writeHTML function generates an HTML file with the match results
func writeHTML(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("File creation error: %v\n", err)
		return
	}
	defer file.Close()

	file.WriteString("<html><body><ul>\n")
	for _, result := range results {
		link := fmt.Sprintf("file://%s", result.FilePath)
		file.WriteString(fmt.Sprintf("<li>%s: <a href=\"%s\">%s</a>, Line: %d</li>\n", html.EscapeString(result.Pattern), link, html.EscapeString(result.FilePath), result.Line))
	}
	file.WriteString("</ul></body></html>\n")
}

// writeJSON function writes the match results to a JSON file
func writeJSON(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("File creation error: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(results)
	if err != nil {
		fmt.Printf("JSON writing error: %v\n", err)
	}
}

// writeTXT function writes the match results to a plain text file
func writeTXT(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		fmt.Printf("File creation error: %v\n", err)
		return
	}
	defer file.Close()

	for _, result := range results {
		file.WriteString(fmt.Sprintf("%s: %s, Line: %d\n", result.Pattern, result.FilePath, result.Line))
	}
}

// scanFiles function scans the given directory for files and searches for the specified patterns
func scanFiles(path string, patterns map[string]string) {
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			buf := make([]byte, 0, 1*1024)
			scanner.Buffer(buf, 500*1024*1024)
			lineNumber := 0
			for scanner.Scan() {
				lineNumber++
				for label, pattern := range patterns {
					r, err := regexp.Compile(pattern)
					if err != nil {
						fmt.Printf("Regex error: %v\n", err)
						continue
					}
					if r.MatchString(scanner.Text()) {
						fmt.Printf("%s found, FILE: %s, LINE: %d\n", label, filePath, lineNumber)
						// Get the absolute file path
						absPath, err := filepath.Abs(filePath)
						if err != nil {
							fmt.Printf("File path error: %v\n", err)
							continue
						}
						// Add the matching result to results slice
						results = append(results, Result{
							Pattern:  label,
							FilePath: absPath,
							Line:     lineNumber,
						})
					}
				}
			}

			if err := scanner.Err(); err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func main() {
	dirPath := flag.String("d", "", "Please specify the directory path to scan.")
	outputFormat := flag.String("o", "txt", "Specify the output format: html, json, or txt.")
	flag.Parse()

	if *dirPath == "" {
		fmt.Println("Please specify the directory path to scan.")
		return
	}

	patterns := map[string]string{
		"E-Mail":                        `\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,7}\b`,
		"IP Address":                    `\b\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}\b`,
		"URL:":                          "\bhttps?:\\/\\/[^\\s]*\b",
		"Json Web Token (JWT)":          "\b[eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9]+\\.[eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ]+\\.[SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c]+\b",
		"Dropbox Access Token":          `\b[a-zA-Z0-9]{64}\b`,
		"SendGrid API Key":              `\bSG\.[a-zA-Z0-9-_]{22,}\b`,
		"Slack Bot Token":               `\bxoxb-[a-zA-Z0-9]{10,}\b`,
		"Stripe Publishable Key":        `\bpk_(live|test)_[a-zA-Z0-9]{24}\b`,
		"Twilio Account SID":            `\bAC[a-zA-Z0-9]{32}\b`,
		"Firebase URL":                  ".*firebaseio\\.com",
		"Slack Token":                   "(xox[p|b|o|a]-[0-9]{12}-[0-9]{12}-[0-9]{12}-[a-z0-9]{32})",
		"RSA private key":               "-----BEGIN RSA PRIVATE KEY-----",
		"SSH (DSA) private key":         "-----BEGIN DSA PRIVATE KEY-----",
		"SSH (EC) private key":          "-----BEGIN EC PRIVATE KEY-----",
		"PGP private key block":         "-----BEGIN PGP PRIVATE KEY BLOCK-----",
		"Amazon AWS Access Key ID":      "AKIA[0-9A-Z]{16}",
		"Amazon MWS Auth Token":         "amzn\\.mws\\.[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}",
		"AWS API Key":                   "AKIA[0-9A-Z]{16}",
		"Facebook Access Token":         "EAACEdEose0cBA[0-9A-Za-z]+",
		"Facebook OAuth":                "[f|F][a|A][c|C][e|E][b|B][o|O][o|O][k|K].*['|\"][0-9a-f]{32}['|\"]",
		"GitHub":                        "[g|G][i|I][t|T][h|H][u|U][b|B].*['|\"][0-9a-zA-Z]{35,40}['|\"]",
		"Generic API Key":               "[a|A][p|P][i|I][_]?[k|K][e|E][y|Y].*['|\"][0-9a-zA-Z]{32,45}['|\"]",
		"Generic Secret":                "[s|S][e|E][c|C][r|R][e|E][t|T].*['|\"][0-9a-zA-Z]{32,45}['|\"]",
		"Google API Key":                "AIza[0-9A-Za-z\\-_]{35}",
		"Google (GCP) Service-account":  "\"type\": \"service_account\"",
		"Google OAuth Access Token":     "ya29\\.[0-9A-Za-z\\-_]+",
		"Heroku API Key":                "[h|H][e|E][r|R][o|O][k|K][u|U].*[0-9A-F]{8}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{4}-[0-9A-F]{12}",
		"MailChimp API Key":             "[0-9a-f]{32}-us[0-9]{1,2}",
		"Mailgun API Key":               "key-[0-9a-zA-Z]{32}",
		"Password in URL":               "[a-zA-Z]{3,10}://[^/\\s:@]{3,20}:[^/\\s:@]{3,20}@.{1,100}[\"'\\s]",
		"PayPal Braintree Access Token": "access_token\\$production\\$[0-9a-z]{16}\\$[0-9a-f]{32}",
		"Picatic API Key":               "sk_live_[0-9a-z]{32}",
		"Slack Webhook":                 "https://hooks.slack.com/services/T[a-zA-Z0-9_]{8}/B[a-zA-Z0-9_]{8}/[a-zA-Z0-9_]{24}",
		"Stripe API Key":                "sk_live_[0-9a-zA-Z]{24}",
		"Stripe Restricted API Key":     "rk_live_[0-9a-zA-Z]{24}",
		"Square Access Token":           "sq0atp-[0-9A-Za-z\\-_]{22}",
		"Square OAuth Secret":           "sq0csp-[0-9A-Za-z\\-_]{43}",
		"Twilio API Key":                "SK[0-9a-fA-F]{32}",
		"Twitter Access Token":          "[t|T][w|W][i|I][t|T][t|T][e|E][r|R].*[1-9][0-9]+-[0-9a-zA-Z]{40}",
	}

	dirName := filepath.Base(*dirPath)
	scanFiles(*dirPath, patterns)
	outputFilename := fmt.Sprintf("%s.%s", dirName, *outputFormat)

	switch *outputFormat {
	case "html":
		writeHTML(outputFilename)
	case "json":
		writeJSON(outputFilename)
	case "txt":
		writeTXT(outputFilename)
	default:
		fmt.Println("Unsupported output format. Please choose html, json, or txt.")
	}
}
