package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

// ---------------- STRUCTS ----------------

type Driver struct {
	ID          int      `json:"id"`
	Code        string   `json:"code"`
	Name        string   `json:"name"`
	Nationality string   `json:"nationality"`
	DateOfBirth string   `json:"date_of_birth"`
	Team        string   `json:"team"`
	PastTeams   []string `json:"past_teams"`
	CarNumber   int      `json:"car_number"`
	DebutYear   int      `json:"debut_year"`
	Wins        int      `json:"wins"`
}

// ---------------- TEAM NORMALIZATION ----------------

var teamMap = map[string]string{
	// Red Bull lineage
	"Red Bull Racing": "Red Bull",
	"Red Bull":        "Red Bull",

	// Mercedes lineage
	"Mercedes":              "Mercedes",
	"Mercedes AMG Petronas": "Mercedes",

	// Ferrari lineage
	"Ferrari":          "Ferrari",
	"Scuderia Ferrari": "Ferrari",

	// McLaren
	"McLaren": "McLaren",

	// Aston Martin lineage (was Racing Point, was Force India)
	"Aston Martin":             "Aston Martin",
	"Aston Martin F1 Team":     "Aston Martin",
	"Racing Point":             "Aston Martin",
	"Racing Point F1 Team":     "Aston Martin",
	"Force India":              "Aston Martin",
	"Racing Point Force India": "Aston Martin",

	// Alpine lineage (was Renault)
	"Alpine F1 Team": "Alpine",
	"Alpine":         "Alpine",
	"Renault":        "Alpine",

	// RB lineage (was AlphaTauri, was Toro Rosso)
	"RB F1 Team":          "RB",
	"RB":                  "RB",
	"AlphaTauri":          "RB",
	"Toro Rosso":          "RB",
	"Scuderia Toro Rosso": "RB",

	// Sauber lineage (was Alfa Romeo, now Kick Sauber, becoming Audi)
	"Sauber":      "Sauber",
	"Alfa Romeo":  "Sauber",
	"Kick Sauber": "Sauber",
	"Audi":        "Sauber",

	// Haas
	"Haas F1 Team": "Haas",
	"Haas":         "Haas",

	// Williams
	"Williams":        "Williams",
	"Williams Racing": "Williams",

	// Cadillac (new 2026, formerly Andretti)
	"Cadillac F1 Team": "Cadillac",
	"Andretti":         "Cadillac",
}

func normalizeTeam(name string) string {
	if val, ok := teamMap[name]; ok {
		return val
	}
	return name
}

// ---------------- HELPERS ----------------

// fetchJSON fetches a URL and unmarshals into target.
// Retries once on failure; flat 3s wait on 429.
func fetchJSON(url string, target interface{}) error {
	var lastErr error

	for attempt := 1; attempt <= 2; attempt++ {
		if attempt > 1 {
			fmt.Printf("    ↻ retry after 1s...\n")
			time.Sleep(1 * time.Second)
		}

		resp, err := http.Get(url)
		if err != nil {
			lastErr = fmt.Errorf("http.Get: %w", err)
			continue
		}

		if resp.StatusCode == 429 {
			resp.Body.Close()
			lastErr = fmt.Errorf("rate limited (429)")
			time.Sleep(3 * time.Second)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			lastErr = fmt.Errorf("reading body: %w", err)
			continue
		}

		if err := json.Unmarshal(body, target); err != nil {
			lastErr = fmt.Errorf("json.Unmarshal: %w", err)
			continue
		}

		return nil
	}

	return fmt.Errorf("all retries failed for %s: %w", url, lastErr)
}

func safeMap(v interface{}) (map[string]interface{}, bool) {
	m, ok := v.(map[string]interface{})
	return m, ok && m != nil
}

func safeSlice(v interface{}) ([]interface{}, bool) {
	s, ok := v.([]interface{})
	return s, ok
}

func safeString(v interface{}) string {
	s, _ := v.(string)
	return s
}

func getMRData(res map[string]interface{}) (map[string]interface{}, bool) {
	return safeMap(res["MRData"])
}

// pause sleeps between API calls to stay within rate limits.
func pause() {
	time.Sleep(600 * time.Millisecond)
}

// ---------------- DISCOVER DRIVER IDs ----------------

func getDriverIDsForYear(year int) ([]string, error) {
	url := fmt.Sprintf("https://api.jolpi.ca/ergast/f1/%d/drivers.json?limit=100", year)

	var res map[string]interface{}
	if err := fetchJSON(url, &res); err != nil {
		return nil, err
	}

	mrdata, ok := getMRData(res)
	if !ok {
		return nil, fmt.Errorf("missing MRData for year %d", year)
	}

	driverTable, ok := safeMap(mrdata["DriverTable"])
	if !ok {
		return nil, fmt.Errorf("missing DriverTable for year %d", year)
	}

	drivers, ok := safeSlice(driverTable["Drivers"])
	if !ok {
		return nil, fmt.Errorf("missing Drivers for year %d", year)
	}

	var ids []string
	for _, d := range drivers {
		dm, ok := safeMap(d)
		if !ok {
			continue
		}
		if id := safeString(dm["driverId"]); id != "" {
			ids = append(ids, id)
		}
	}

	return ids, nil
}

func collectAllDriverIDs(fromYear, toYear int) []string {
	seen := map[string]bool{}
	var all []string

	for year := fromYear; year <= toYear; year++ {
		fmt.Printf("  Scanning %d...\n", year)

		ids, err := getDriverIDsForYear(year)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ year %d: %v\n", year, err)
			pause()
			continue
		}

		added := 0
		for _, id := range ids {
			if !seen[id] {
				seen[id] = true
				all = append(all, id)
				added++
			}
		}

		fmt.Printf("    → %d drivers this year, %d new (%d unique total)\n", len(ids), added, len(all))
		pause()
	}

	return all
}

// ---------------- DRIVER BASE ----------------

func getDriverBase(driverID string) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://api.jolpi.ca/ergast/f1/drivers/%s.json", driverID)

	var res map[string]interface{}
	if err := fetchJSON(url, &res); err != nil {
		return nil, err
	}

	mrdata, ok := getMRData(res)
	if !ok {
		return nil, fmt.Errorf("missing MRData")
	}

	driverTable, ok := safeMap(mrdata["DriverTable"])
	if !ok {
		return nil, fmt.Errorf("missing DriverTable")
	}

	drivers, ok := safeSlice(driverTable["Drivers"])
	if !ok || len(drivers) == 0 {
		return nil, fmt.Errorf("no drivers found")
	}

	base, ok := safeMap(drivers[0])
	if !ok {
		return nil, fmt.Errorf("invalid driver entry")
	}

	return base, nil
}

// ---------------- CURRENT TEAM ----------------

func getCurrentTeam(driverID string) (string, error) {
	url := fmt.Sprintf("https://api.jolpi.ca/ergast/f1/current/drivers/%s/constructors.json", driverID)

	var res map[string]interface{}
	if err := fetchJSON(url, &res); err != nil {
		return "", err
	}

	mrdata, ok := getMRData(res)
	if !ok {
		return "", nil
	}

	constructorTable, ok := safeMap(mrdata["ConstructorTable"])
	if !ok {
		return "", nil
	}

	constructors, ok := safeSlice(constructorTable["Constructors"])
	if !ok || len(constructors) == 0 {
		return "", nil
	}

	first, ok := safeMap(constructors[0])
	if !ok {
		return "", nil
	}

	return normalizeTeam(safeString(first["name"])), nil
}

// ---------------- PAST TEAMS ----------------

func getPastTeams(driverID string, current string) ([]string, error) {
	url := fmt.Sprintf("https://api.jolpi.ca/ergast/f1/drivers/%s/constructors.json?limit=100", driverID)

	var res map[string]interface{}
	if err := fetchJSON(url, &res); err != nil {
		return nil, err
	}

	mrdata, ok := getMRData(res)
	if !ok {
		return nil, nil
	}

	constructorTable, ok := safeMap(mrdata["ConstructorTable"])
	if !ok {
		return nil, nil
	}

	constructors, ok := safeSlice(constructorTable["Constructors"])
	if !ok {
		return nil, nil
	}

	unique := map[string]bool{}
	var past []string

	for _, c := range constructors {
		cm, ok := safeMap(c)
		if !ok {
			continue
		}
		name := normalizeTeam(safeString(cm["name"]))
		if name != current && !unique[name] {
			unique[name] = true
			past = append(past, name)
		}
	}

	return past, nil
}

// ---------------- RESULTS ----------------

func getResultsData(driverID string) (wins int, debutYear int, err error) {
	offset := 0
	limit := 100
	debutYear = 9999

	for {
		url := fmt.Sprintf(
			"https://api.jolpi.ca/ergast/f1/drivers/%s/results.json?limit=%d&offset=%d",
			driverID, limit, offset,
		)

		var res map[string]interface{}
		if err = fetchJSON(url, &res); err != nil {
			return 0, 0, err
		}

		mrdata, ok := getMRData(res)
		if !ok {
			break
		}

		raceTable, ok := safeMap(mrdata["RaceTable"])
		if !ok {
			break
		}

		races, ok := safeSlice(raceTable["Races"])
		if !ok || len(races) == 0 {
			break
		}

		for _, race := range races {
			raceMap, ok := safeMap(race)
			if !ok {
				continue
			}

			season, _ := strconv.Atoi(safeString(raceMap["season"]))
			if season > 0 && season < debutYear {
				debutYear = season
			}

			results, ok := safeSlice(raceMap["Results"])
			if !ok {
				continue
			}

			for _, r := range results {
				rm, ok := safeMap(r)
				if !ok {
					continue
				}
				if safeString(rm["position"]) == "1" {
					wins++
				}
			}
		}

		offset += limit
		pause()
	}

	if debutYear == 9999 {
		debutYear = 0
	}

	return wins, debutYear, nil
}

// ---------------- MAIN ----------------

func main() {
	const fromYear = 2018
	const toYear = 2026

	fmt.Printf("=== Discovering drivers from %d to %d ===\n", fromYear, toYear)
	driverIDs := collectAllDriverIDs(fromYear, toYear)
	fmt.Printf("\n→ %d unique drivers discovered\n\n", len(driverIDs))

	if err := os.MkdirAll("data", 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create data directory: %v\n", err)
		os.Exit(1)
	}

	var drivers []Driver
	seq := 1

	for i, driverID := range driverIDs {
		fmt.Printf("[%d/%d] Fetching: %s\n", i+1, len(driverIDs), driverID)

		// --- base info ---
		base, err := getDriverBase(driverID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ getDriverBase: %v — skipping\n", err)
			pause()
			continue
		}
		pause()

		// --- current team: skip if not on a current team (retired/inactive) ---
		current, err := getCurrentTeam(driverID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ getCurrentTeam: %v — skipping\n", err)
			pause()
			continue
		}
		if current == "" {
			fmt.Printf("  ~ no current team (inactive/retired) — skipping\n")
			pause()
			continue
		}
		pause()

		// --- past teams ---
		past, err := getPastTeams(driverID, current)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ getPastTeams: %v — skipping\n", err)
			pause()
			continue
		}
		pause()

		// --- results: skip if we couldn't determine debut year ---
		wins, debut, err := getResultsData(driverID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ getResultsData: %v — skipping\n", err)
			pause()
			continue
		}
		if debut == 0 {
			fmt.Printf("  ~ no debut year found — skipping\n")
			pause()
			continue
		}

		carNumber := 0
		if val := safeString(base["permanentNumber"]); val != "" {
			carNumber, _ = strconv.Atoi(val)
		}

		driver := Driver{
			ID:          seq,
			Code:        safeString(base["code"]),
			Name:        safeString(base["givenName"]) + " " + safeString(base["familyName"]),
			Nationality: safeString(base["nationality"]),
			DateOfBirth: safeString(base["dateOfBirth"]),
			Team:        current,
			PastTeams:   past,
			CarNumber:   carNumber,
			DebutYear:   debut,
			Wins:        wins,
		}

		drivers = append(drivers, driver)
		fmt.Printf("  ✓ %s | Team: %s | Wins: %d | Debut: %d\n",
			driver.Name, driver.Team, driver.Wins, driver.DebutYear)

		seq++

		// Generous pause between drivers to stay well clear of rate limits
		time.Sleep(1200 * time.Millisecond)
	}

	file, err := os.Create("data/drivers.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create output file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(drivers)

	fmt.Printf("\n✅ drivers.json generated with %d drivers\n", len(drivers))
}
