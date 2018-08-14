package smartlog

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
	"io/ioutil"
	"html/template"
	"net/http"
	"path/filepath"
	"log"
	"encoding/json"

	"github.com/fnproject/fn/api/models"
	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/fn/fnext"
	"github.com/fsouza/go-dockerclient"
	//"github.com/apang1992/jq"
	"github.com/gorilla/mux"

	apexFn "github.com/apexcz/fnExtensions/models"
)

var tpl *template.Template
var assetDir string

func init() {
	server.RegisterExtension(&staticExt{})

	// Get the path to the asset directory.
	tempAssetDir, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		//return
	}
	assetDir = tempAssetDir

	//tpl = template.Must(template.ParseGlob("/Users/chineduoty/Documents/go/src/github.com/apexcz/fnExtensions/smartlog/templates/*.html"))
	tpl = template.Must(template.ParseGlob(fmt.Sprintf("%s/smartlog/templates/*.html", assetDir)))
}

type staticExt struct {

}

func (e *staticExt) Name() string {
	return "github.com/apexcz/fnExtensions/smartlog"
}

func (e *staticExt) Setup(s fnext.ExtServer) error {
	s.AddCallListener(&LogListener{})
	return nil
}

type LogListener struct {

}

type StaticReportModel struct {
	AssetDir string
	GenReport apexFn.GeneratedReport
	ImageName string
	TotalVuln int
}


func (l *LogListener) BeforeCall(ctx context.Context, call *models.Call) error {
	fmt.Println("Interception before function call occurs")
	fmt.Println("The calling image is ",call.Image)


	//launch docker client to fetch the image
	endpoint := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}

	imgs,err := client.ListImages(docker.ListImagesOptions{All:false})

	if err != nil {
		panic(err)
	}

	img := imgs[0]
	fmt.Println("Last image is ", img.RepoTags)

	//CLAIR_ADDR=localhost CLAIR_OUTPUT=High CLAIR_THRESHOLD=10  klar postgres:9.5.1 > result.json

	fmt.Println("Asset path is ",assetDir)

	go launchClair(call.Image)
	time.Sleep(1*time.Second)

	return nil
}

func (l *LogListener) AfterCall(ctx context.Context, call *models.Call) error {
	fmt.Println("Triggers after function executes completely")
	return nil
}

func (srm *StaticReportModel) StaticHandler(response http.ResponseWriter, request *http.Request) {
	log.Println("Inside handler = ",srm.AssetDir)
	err := tpl.ExecuteTemplate(response,"index.html",srm)
	if err != nil {
		fmt.Println(err)
	}
}

func launchClair(image string) {
	instruction := fmt.Sprintf("CLAIR_ADDR=localhost CLAIR_OUTPUT=Low CLAIR_THRESHOLD=3 JSON_OUTPUT=true klar postgres:9.5.1 > %s/smaticreport.json", assetDir)

	cmd := exec.Command("sh","-c",instruction)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	jsonFile, err := os.Open(fmt.Sprintf("%s/smaticreport.json", assetDir))
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	jsonData, _ := ioutil.ReadAll(jsonFile)

	var report apexFn.GeneratedReport
	json.Unmarshal(jsonData, &report)

	/**
	layerCount, err := jq.JsonQuery([]byte(jsonData), ".LayerCount")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println("the layerCount is:", string(layerCount))
	}
	*/


	staticDir := filepath.Join(assetDir, "smartlog/static")
	templateDir := filepath.Join(assetDir, "smartlog/templates")
	// Make sure all the paths exist.
	tmplPaths := []string{
		staticDir,
		filepath.Join(templateDir, "index.html"),
	}

	fmt.Println("The html paths are => ",tmplPaths)

	for _, path := range tmplPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("template %s not found", path)
			//Errorf
		}
	}

	// Create mux server.
	mux := mux.NewRouter()
	mux.UseEncodedPath()

	vulnCount := len(report.Vulnerabilities.High) + len(report.Vulnerabilities.Low) + len(report.Vulnerabilities.Medium)
	vm := &StaticReportModel{AssetDir: filepath.Join(templateDir, "index.html"), GenReport: report, ImageName: image, TotalVuln:vulnCount}
	mux.HandleFunc("/static-report", vm.StaticHandler)

	// Serve the static assets.
	staticHandler := http.FileServer(http.Dir(staticDir))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))
	mux.Handle("/", staticHandler)

	http.Handle("/",mux)
	http.ListenAndServe(":8580", nil)
}
