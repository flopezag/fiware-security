package cmd

import (
    "encoding/json"
    "fmt"
    "os"
    "sort"
	"github.com/spf13/cobra"
    "gonum.org/v1/plot"
    "gonum.org/v1/plot/plotter"
    "gonum.org/v1/plot/vg"
    "gonum.org/v1/plot/plotutil" // Add this import
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
		data_analysis()
	},
}


func data_analysis() {
    file, err := os.ReadFile("./results/Keyrock_idm_20250702_1422_grype.json")
    if err != nil {
        panic(err)
    }

    var report Report
    if err := json.Unmarshal(file, &report); err != nil {
        panic(err)
    }

    total := 0
    severityCount := make(map[string]int)
    var epssSum, riskSum float64
    var epssAbove90, riskAbove90 int

    for _, match := range report.Matches {
        vuln := match.Vulnerability
        total++
        severityCount[vuln.Severity]++

        // EPSS — use the first if available
        if len(vuln.EPSS) > 0 {
            epss := vuln.EPSS[0].EPSS
            epssSum += epss
            if epss > 0.9 {
                epssAbove90++
            }
        }

        // Risk
        riskSum += vuln.Risk
        if vuln.Risk > 90 {
            riskAbove90++
        }
    }

    // Output
    fmt.Printf("📌 Total vulnerabilities: %d\n", total)
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
        fmt.Printf("📈 Average EPSS Score: %.5f\n", epssSum/float64(total))
        fmt.Printf("📈 Average Risk Score: %.5f\n", riskSum/float64(total))
    }
    fmt.Printf("🔥 Vulnerabilities with EPSS > 0.9: %d\n", epssAbove90)
    fmt.Printf("⚠️  Vulnerabilities with Risk > 90 : %d\n", riskAbove90)

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
    bar.Color = plotutil.Color(0) // Use plotutil.Color instead
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
