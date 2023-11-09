# FindSecrets

## Overview
The FindSecrets is a command-line tool written in Go that scans files in a given directory for specified patterns such as email addresses, IP addresses, URLs, and various tokens and keys. It supports outputting the results in HTML, JSON, and plain text formats.

## Installation

To set up the FindSecrets, you need to have Go installed on your machine. Follow these steps:

1. Clone the repository to your local machine.
2. Navigate to the cloned directory.
3. Build the program using `go build`.

## Usage

Run the program with the following command-line arguments:

- `-d`: Specify the directory path to scan.
- `-o`: Specify the output format (`html`, `json`, or `txt`).

Example:

```shell
go run main.go -d /path/to/directory -o json
```

```shell
./findsecrets -d /path/to/directory -o json
```

This will scan the directory `/path/to/directory` and output the results in JSON format.

## Output

The tool generates an output file in the chosen format with the detected patterns. The output includes the pattern's label, the file path where it was found, and the line number.

- HTML output: Provides a list of hyperlinks for each match found, allowing easy navigation to the matched line in the file.
- JSON output: Useful for further processing with other tools or importing into databases.
- TXT output: A simple text file with one match per line.

## Contributing

Contributions to the FindSecrets are welcome. Please fork the repository, make your changes, and submit a pull request.

## Contact
https://www.linkedin.com/in/huseyintintas/

### Buy Me A Coffee
https://www.buymeacoffee.com/huseyintintas
