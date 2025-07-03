package cmd

import (
    "encoding/json"
    "errors"
    "fmt"
    "io/fs"
    "os"
    "sort"
	"github.com/spf13/cobra"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/plot/plotutil" // Add this import
	"path/filepath"
	"strings"
    "time"
)

type Report struct {
    Matches []Match `json:"matches"`
}

type Match struct {
    Vulnerability Vulnerability `json:"vulnerability"`
}

type Vulnerability struct {
    Severity string  `json:"severity"`
    Risk     float64 `json:"risk"`
    EPSS     []EPSS  `json:"epss"`
}

type EPSS struct {
    EPSS float64 `json:"epss"`
}

// visualizeCmd represents the visualize command
var visualizeCmd = &cobra.Command{
	Use:   "visualize [fiware-ge]",
	Short: "Visualize the analysis of a FIWARE Generic Enabler",
	Long: `Operation to generate the output data of the security analysis of a FIWARE Generic Enabler..`,

	Run: func(cmd *cobra.Command, args []string) {
        enablerName := args[0]
        reportFile, err := findLatestReportFile(enablerName)
        if err != nil {
            fmt.Printf("❌ Error: %v\n", err)
            return
        }

        fmt.Printf("📄 Using report: %s\n", reportFile)
        analyzeAndVisualize(reportFile)

        //dataAnalysis(filePath)
	},
}


func analyzeAndVisualize(filePath string) {
    file, err := os.ReadFile(filePath)
    if err != nil {
        panic(err)
    }

    var report Report
    if err := json.Unmarshal(file, &report); err != nil {
        panic(err)
    }

    total := 0
    severityCount := make(map[string]int)
    var epssAbove90, riskAbove90 int
    combinedRiskProduct := 1.0
    combinedEPSSProduct := 1.0

    for _, match := range report.Matches {
        vuln := match.Vulnerability
        total++
        severityCount[vuln.Severity]++

        // EPSS: accumulate product of (1 - P) — use the first if available
        if len(vuln.EPSS) > 0 {
            epss := vuln.EPSS[0].EPSS
            combinedEPSSProduct *= (1 - epss)
            if epss >= 0.9 {
                epssAbove90++
            }
        }

        // Risk: accumulate product of (1 - P)
        risk := vuln.Risk / 100.0
        combinedRiskProduct *= (1 - risk)
        if vuln.Risk >= 90 {
            riskAbove90++
        }
    }

    combinedRisk := 1 - combinedRiskProduct
    combinedEPSS := 1 - combinedEPSSProduct

    // Output
    fmt.Printf("📌 Number of vulnerabilities: %d\n", total)
    fmt.Printf("📊 Severity Histogram:\n")
    keys := make([]string, 0, len(severityCount))
    for k := range severityCount {
        keys = append(keys, k)
    }
    sort.Strings(keys)
    for _, k := range keys {
        fmt.Printf("  %-10s: %s (%d)\n", k, string(repeat('#', severityCount[k])), severityCount[k])
    }

    if total > 0 {
        fmt.Printf("📈 Combined EPSS Score: %.5f\n", combinedEPSS)
        fmt.Printf("📈 Combined Risk Score: %.5f\n", combinedRisk)
    }
    fmt.Printf("🔥 Vulnerabilities with EPSS ≥ 0.9: %d\n", epssAbove90)
    fmt.Printf("🔥 Vulnerabilities with Risk ≥ 0.9 : %d\n", riskAbove90)

    // Create SVG histogram
    err = plotSeverityHistogram(severityCount, "./results/severity_histogram.svg")
    if err != nil {
        fmt.Printf("Error generating plot: %v\n", err)
    } else {
        fmt.Println("\n📷 Saved severity histogram to: ./results/severity_histogram.svg")
    }
}

func repeat(char rune, count int) []rune {
    out := make([]rune, count)
    for i := range out {
        out[i] = char
    }
    return out
}

func init() {
	rootCmd.AddCommand(visualizeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// checkCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// checkCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func plotSeverityHistogram(data map[string]int, filename string) error {
    p := plot.New()
    p.Title.Text = "Vulnerability Severity Histogram"
    p.Title.Padding = vg.Points(10)

    p.Y.Label.Text = "Count"
    p.Y.Label.Padding = vg.Points(4)

    p.X.Label.Text = "Severity"
    p.X.Label.Padding = vg.Points(4)

    // Sort severity keys
    keys := make([]string, 0, len(data))
    for k := range data {
        keys = append(keys, k)
    }
    sort.Strings(keys)

    values := make(plotter.Values, len(keys))
    for i, k := range keys {
        values[i] = float64(data[k])
    }

    barWidth := vg.Points(25)
    bar, err := plotter.NewBarChart(values, barWidth)
    if err != nil {
        return err
    }

    bar.LineStyle.Width = vg.Length(0)
    bar.Color = plotutil.Color(0)
    p.Add(bar)

    // Set nominal (category) labels
    p.NominalX(keys...)
    p.X.Tick.Label.Font.Size = vg.Points(10)
    p.Y.Tick.Label.Font.Size = vg.Points(10)

    // Add headroom to avoid bars touching top
    p.Y.Max = maxValue(values) + 1

    // Save plot to file
    return p.Save(6*vg.Inch, 4*vg.Inch, filename)
}

func maxValue(vals plotter.Values) float64 {
    max := 0.0
    for _, v := range vals {
        if v > max {
            max = v
        }
    }
    return max
}

func findLatestReportFile(enablerName string) (string, error) {
    enablersFile, err := os.ReadFile("./config/enablers.json")
    if err != nil {
        return "", fmt.Errorf("unable to read enablers.json: %w", err)
    }

    var enablers EnablerList
    if err := json.Unmarshal(enablersFile, &enablers); err != nil {
        return "", fmt.Errorf("invalid enablers.json: %w", err)
    }

    var matchedEnabler *Enabler
    for _, enabler := range enablers.Enablers {
        if strings.EqualFold(enabler.Name, enablerName) {
            matchedEnabler = &enabler
            break
        }
    }

    if matchedEnabler == nil {
        return "", fmt.Errorf("enabler '%s' not found in enablers.json", enablerName)
    }

    // Normalize identifiers to match report file names
    var identifiers []string
    for _, img := range matchedEnabler.Image {
        // Extract the base name, remove tag and replace dashes with underscores
        imgName := strings.Split(img, "/")
        lastPart := imgName[len(imgName)-1]
        nameOnly := strings.Split(lastPart, ":")[0]
        identifiers = append(identifiers, strings.ToLower(nameOnly))
    }

    var latestFile string
    var latestTime time.Time

    fmt.Printf("🔍 Searching for latest report file for enabler '%s'...\n", enablerName)
    fmt.Printf("  Identifiers: %v\n", identifiers)

    fmt.Println("  Scanning results directory for matching files...")
    fmt.Println(os.Getwd())
    
    // Walk through the results directory to find the latest report file
    err = filepath.WalkDir("./results", func(path string, d fs.DirEntry, err error) error {
        if err != nil || d.IsDir() || !strings.HasSuffix(d.Name(), ".json") {
            return nil
        }

        for _, id := range identifiers {
            if strings.Contains(strings.ToLower(path), strings.ToLower(id)) {
                info, err := os.Stat(path)
                if err != nil {
                    return nil
                }
                if info.ModTime().After(latestTime) {
                    latestTime = info.ModTime()
                    latestFile = path
                }
            }
        }
        return nil
    })

    if err != nil {
        return "", err
    }

    if latestFile == "" {
        return "", errors.New("no report file found for enabler")
    }

    return latestFile, nil
}

type EnablerList struct {
    Enablers []Enabler `json:"enablers"`
}
