package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// -------------------------------------------------
// STRUCT AND GLOBAL DECLARATIONS
// -------------------------------------------------
var proteins map[string]ProteinData
var tmScoreMatrix map[string]map[string]float32
var transcripts_table map[string]string = map[string]string{} //uses transcripts as keys and hash_ids as values
var stats DBStats

// struct for protein data
type ProteinData struct {
	ID         string `json:"id"`
	Species    string `json:"species"`
	CommonName string `json:"common_name"`
	Gene       string `json:"gene"`
	Transcript string `json:"transcript"`
	CDNA       string `json:"cDNA"`
	Protein    string `json:"protein"`
	Conf       struct {
		UseTemplates bool   `json:"use_templates"`
		UseAmber     bool   `json:"use_amber"`
		MsaMode      string `json:"msa_mode"`
		ModelType    string `json:"model_type"`
		NumModels    int    `json:"num_models"`
		NumRecycles  int    `json:"num_recycles"`
		RankBy       string `json:"rank_by"`
		PairMode     string `json:"pair_mode"`
	} `json:"conf"`
}

// struct for general database statistics that will appear on the header of the index page
type DBStats struct {
	NumTranscripts int
	NumGenes       int
	NumSpecies     int
	NumStructures  int
}

// -------------------------------------------------
// HANDLER FUNCTIONS PAGES
// -------------------------------------------------
func indexPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[+] Received connection to main page")
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, stats)
}

func proteinPage(w http.ResponseWriter, r *http.Request) {
	tmp := strings.Split(r.URL.Path, "/")
	protein := tmp[2]
	fmt.Println("[+] Received connection to protein page (" + protein + ")")

	if output, ok := proteins[protein]; ok {
		tmpl := template.Must(template.ParseFiles("templates/protein.html"))
		tmpl.Execute(w, output)
	} else if id, ok := transcripts_table[protein]; ok {
		http.Redirect(w, r, "/protein/"+id, 303)
	} else {
		http.Redirect(w, r, "/", 303)
	}
}

func searchPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[+] Received connection to search page")

	type Output struct {
		Query string
		Data  []ProteinData
	}

	tmp, ok := r.URL.Query()["query"]
	output := Output{"", []ProteinData{}}
	if ok == true {
		output.Query = strings.ToLower(tmp[0])
		for _, v := range proteins {
			if strings.Contains(strings.ToLower(v.Gene), output.Query) || strings.Contains(strings.ToLower(v.Transcript), output.Query) || strings.Contains(strings.ToLower(v.Species), output.Query) || strings.Contains(strings.ToLower(v.CommonName), output.Query) {
				output.Data = append(output.Data, v)
			}
		}
	} else {
		for _, v := range proteins {
			output.Data = append(output.Data, v)
		}
	}

	tmpl := template.Must(template.ParseFiles("templates/search.html"))
	tmpl.Execute(w, output)
}

func apiPage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[+] Received connection to API page")
	tmpl := template.Must(template.ParseFiles("templates/api.html"))
	tmpl.Execute(w, stats)
}

// -------------------------------------------------
// HANDLER FUNCTIONS API ENDPOINTS
// -------------------------------------------------
func listProteinsAPI(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(&proteins)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func proteinInfoAPI(w http.ResponseWriter, r *http.Request) {
	tmp := strings.Split(r.URL.Path, "/")
	protein := tmp[3]
	fmt.Println("[+] Received connection to protein API (" + protein + ")")

	if output, ok := proteins[protein]; ok {
		data, _ := json.Marshal(&output)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: Invalid protein supplied")
	}
}

func tmScoreAPI(w http.ResponseWriter, r *http.Request) {
	tmp := strings.Split(r.URL.Path, "/")
	protein1 := tmp[len(tmp)-1]
	protein2 := tmp[len(tmp)-2]
	fmt.Println("[+] Received connection to TM-SCORE API (" + protein1 + ", " + protein2 + ")")

	if output, ok := tmScoreMatrix[protein1][protein2]; ok {
		data, _ := json.Marshal(&output)
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: Invalid proteins supplied")
	}
}

// -------------------------------------------------
// MAIN
// -------------------------------------------------
func main() {
	// load proteins data from json
	data, err := ioutil.ReadFile("./proteins.json")
	if err != nil {
		log.Fatal("Error: Cannot read proteins.json", err)
	}
	json.Unmarshal(data, &proteins)

	// load TM-Score Matrix
	data, err = ioutil.ReadFile("./tm-matrix.json")
	if err != nil {
		log.Fatal("Error: Cannot read tm-matrix.json", err)
	}
	json.Unmarshal(data, &tmScoreMatrix)

	// create transcripts look up table
	for id, protein := range proteins {
		transcripts_table[protein.Transcript] = id
	}

	// get db stats
	var genes = make(map[string]bool)
	var species = make(map[string]bool)
	var transcripts = make(map[string]bool)
	var structures = make(map[string]bool)
	for _, v := range proteins {
		genes[v.Gene] = false
		species[v.Species] = false
		transcripts[v.Transcript] = false
		structures[v.Protein] = false
	}
	stats.NumGenes = len(genes)
	stats.NumSpecies = len(species)
	stats.NumTranscripts = len(transcripts)
	stats.NumStructures = len(structures)

	// -------------------------------------------------
	// WEBSERVER STUFF
	// -------------------------------------------------
	// serve static directories
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css/"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js/"))))
	http.Handle("/.well-known/", http.StripPrefix("/.well-known/", http.FileServer(http.Dir(".well-known/"))))

	// serve templates
	http.HandleFunc("/", indexPage)
	http.HandleFunc("/protein/", proteinPage)
	http.HandleFunc("/search", searchPage)
	http.HandleFunc("/api", apiPage)

	// API endpoints
	http.HandleFunc("/api/protein/", proteinInfoAPI)
	http.HandleFunc("/api/list", listProteinsAPI)
	http.HandleFunc("/api/tmscore/", tmScoreAPI)

	// Listen
	fmt.Println("[*] Starting server live at http://127.0.01:8002")
	http.ListenAndServe(":8002", nil)
}
