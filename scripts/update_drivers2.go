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


var teamMap = map[string]string{     // team normalization as api returns diff word-formatting for the same team sometimes, like merc/ mercedes amg petronas/ mercedes etc. 
	//red bull lineage
	"Red Bull Racing": "Red Bull",
	"Red Bull":        "Red Bull",

	//mercedes lineage( later add brawngp and stuff upon expansion)
	"Mercedes":              "Mercedes",
	"Mercedes AMG Petronas": "Mercedes",

	//ferrari lineage
	"Ferrari":          "Ferrari",
	"Scuderia Ferrari": "Ferrari",

	//mclaren lineage
	"McLaren": "McLaren",

	//aston martin(was Racing Point and Force India so accounting for that as well + need to mention this in a rules - popup in the frontend)
	"Aston Martin":             "Aston Martin",
	"Aston Martin F1 Team":     "Aston Martin",
	"Racing Point":             "Aston Martin",
	"Racing Point F1 Team":     "Aston Martin",
	"Force India":              "Aston Martin",
	"Racing Point Force India": "Aston Martin",

	//alpine lineage( was renault and benetton as well but for now its only constructors from 2018 onwards)
	"Alpine F1 Team": "Alpine",
	"Alpine":         "Alpine",
	"Renault":        "Alpine",

	//RB lineage (was AlphaTauri, tororosso, and vcarb now)
	"RB F1 Team":          "RB",
	"RB":                  "RB",
	"AlphaTauri":          "RB",
	"Toro Rosso":          "RB",
	"Scuderia Toro Rosso": "RB",

	//Audi lineage (was Alfa Romeo,kick sauber, now Audi)
	"Sauber":      "Audi",
	"Alfa Romeo":  "Audi",
	"Kick Sauber": "Audi",
	"Audi":        "Audi",

	//haas
	"Haas F1 Team": "Haas",
	"Haas":         "Haas",

	//Williams
	"Williams":        "Williams",
	"Williams Racing": "Williams",

	//Cadillac (new 2026, maybe andretti but it is trivial but jus adding it for now as it doesnt make much of a diff)
	"Cadillac F1 Team": "Cadillac",
	"Andretti":         "Cadillac",
}

func normalizeTeam(name string) string {
	if val, ok := teamMap[name]; ok {
		return val
	}
	return name
}

//helper functions!!!!!!!
//fetchJSON fetches a URL and unmarshals into target
//retries once on failure and a flat 3s wait on 429
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

//pause sleeps between API calls to stay within rate limits as we use a free api
func pause() {
	time.Sleep(600 * time.Millisecond)
}

//we need driver id's and an exact match

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

//Driver database

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

//Past teams of current/past drivers as well

//getTeamHistory returns (latestTeam, pastTeams) by looking at actual race 
//results ordered by season, so we get the true most-recent team 
//but it also returns in alphabetical order from the constructors endpoint (NEED TO FIX THIS IN A LATER PHASE OF DEV!!)
func getTeamHistory(driverID string) (latest string, past []string, err error) {
	//we find the most recent season this driver raced in by fetching their season list
	seasonsURL := fmt.Sprintf("https://api.jolpi.ca/ergast/f1/drivers/%s/seasons.json?limit=100", driverID)

	var seasonsRes map[string]interface{}
	if err = fetchJSON(seasonsURL, &seasonsRes); err != nil {
		return "", nil, err
	}

	mrdata, ok := getMRData(seasonsRes)
	if !ok {
		return "", nil, nil
	}

	seasonTable, ok := safeMap(mrdata["SeasonTable"])
	if !ok {
		return "", nil, nil
	}

	seasons, ok := safeSlice(seasonTable["Seasons"])
	if !ok || len(seasons) == 0 {
		return "", nil, nil
	}

	//seasons are returned in ascending order -> last = most recent
	lastSeasonEntry, ok := safeMap(seasons[len(seasons)-1])
	if !ok {
		return "", nil, nil
	}
	lastSeason := safeString(lastSeasonEntry["season"])

	//get the constructor they drove for in that final season
	pause()
	latestURL := fmt.Sprintf("https://api.jolpi.ca/ergast/f1/%s/drivers/%s/constructors.json", lastSeason, driverID)

	var latestRes map[string]interface{}
	if err = fetchJSON(latestURL, &latestRes); err != nil {
		return "", nil, err
	}

	mrdata2, ok := getMRData(latestRes)
	if !ok {
		return "", nil, nil
	}

	ctTable, ok := safeMap(mrdata2["ConstructorTable"])
	if !ok {
		return "", nil, nil
	}

	ctors, ok := safeSlice(ctTable["Constructors"])
	if !ok || len(ctors) == 0 {
		return "", nil, nil
	}

	latestEntry, ok := safeMap(ctors[len(ctors)-1])
	if !ok {
		return "", nil, nil
	}
	latest = normalizeTeam(safeString(latestEntry["name"]))

	//get all constructors across career for past teams
	pause()
	allURL := fmt.Sprintf("https://api.jolpi.ca/ergast/f1/drivers/%s/constructors.json?limit=100", driverID)

	var allRes map[string]interface{}
	if err = fetchJSON(allURL, &allRes); err != nil {
		return latest, nil, nil //have latest at least for now, past teams are a kinda bonus but still needed
	}

	mrdata3, ok := getMRData(allRes)
	if !ok {
		return latest, nil, nil
	}

	allCtTable, ok := safeMap(mrdata3["ConstructorTable"])
	if !ok {
		return latest, nil, nil
	}

	allCtors, ok := safeSlice(allCtTable["Constructors"])
	if !ok {
		return latest, nil, nil
	}

	unique := map[string]bool{latest: true}
	for _, c := range allCtors {
		cm, ok := safeMap(c)
		if !ok {
			continue
		}
		name := normalizeTeam(safeString(cm["name"]))
		if !unique[name] {
			unique[name] = true
			past = append(past, name)
		}
	}

	return latest, past, nil
}

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

func main() { // main func
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
		base, err := getDriverBase(driverID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ getDriverBase: %v — skipping\n", err)
			pause()
			continue
		}
		pause()
		//team history: latest team + all past teams in one call
		latest, past, err := getTeamHistory(driverID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ getTeamHistory: %v — skipping\n", err)
			pause()
			continue
		}
		if latest == "" {
			fmt.Fprintf(os.Stderr, "  ✗ no team history found — skipping\n")
			pause()
			continue
		}
		pause()

		//skip if debut date not found  (the API gives debutn ot found for drivers who just took 1 session like fp1 fp2 etc. like a jak crawford etc. so theyu get skipped)
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
			Team:        latest,
			PastTeams:   past,
			CarNumber:   carNumber,
			DebutYear:   debut,
			Wins:        wins,
		}

		drivers = append(drivers, driver)
		fmt.Printf("  ✓ %s | Team: %s | Wins: %d | Debut: %d\n",
			driver.Name, driver.Team, driver.Wins, driver.DebutYear)

		seq++

		//pause between drivers to clear rate limits
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

	fmt.Printf("\ndrivers.json generated with %d drivers\n", len(drivers))
}
