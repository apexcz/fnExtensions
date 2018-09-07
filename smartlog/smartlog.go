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
	"encoding/json"

	"github.com/fnproject/fn/api/models"
	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/fn/fnext"
	"github.com/fsouza/go-dockerclient"
	"github.com/gorilla/mux"

	apexFn "github.com/apexcz/fnExtensions/models"
)

var tpl *template.Template
var assetDir string
var vulnThreshold int
var vulnLevel string
var staticMode int

var containerStartedAt time.Time
var client *docker.Client

func init() {
	// Registers the extension
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

	vulnThreshold = 1
	vulnLevel = "Low"
	staticMode = 0 // 0 - learning , 1 - enforce
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
	fmt.Println("Interception before function call occurs. \n The calling image is ",call.Image)

	//time function call started
	containerStartedAt = time.Now()

	//launch docker client to fetch the image
	endpoint := "unix:///var/run/docker.sock"
	newClient, err := docker.NewClient(endpoint)
	if err != nil {
		panic(err)
	}
	client = newClient

	imgId,err := client.InspectImage(call.Image)
	if err != nil {
		panic(err)
	}

	imageId := imgId.ID[7:19]
	fmt.Println("Image id is ", imageId)
	//CLAIR_ADDR=localhost CLAIR_OUTPUT=High CLAIR_THRESHOLD=10  klar postgres:9.5.1 > result.json

	if staticMode == 0 {
		go launchClair(call.Image,imageId)
		time.Sleep(1*time.Second)
	}else{
		risksCount := launchClair(call.Image,imageId)
		if risksCount > vulnThreshold {
			err := client.RemoveImage(call.Image)
			if err != nil {
				panic(err)
			}
		}
	}

	return nil
}

func (l *LogListener) AfterCall(ctx context.Context, call *models.Call) error {

	callDuration := time.Now().Sub(containerStartedAt)
	fmt.Println("Function call executed for ",callDuration)
	fmt.Printf("Call Model is %v \n",call)

	/**
	container,err := client.InspectContainer(call.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println("\n ========================== \n")
	fmt.Printf("Container details is %v \n",container)
	*/
	return nil
}

func (srm *StaticReportModel) StaticHandler(response http.ResponseWriter, request *http.Request) {
	err := tpl.ExecuteTemplate(response,"index.html",srm)
	if err != nil {
		fmt.Println(err)
	}
}

func launchClair(image string,imageid string) (vulnCount int) {
	instruction := fmt.Sprintf("CLAIR_ADDR=localhost CLAIR_OUTPUT=%s CLAIR_THRESHOLD=%d JSON_OUTPUT=true klar %s > %s/%s.json",vulnLevel,vulnThreshold, image,assetDir,imageid)

	startScan := time.Now()

	cmd := exec.Command("sh","-c",instruction)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())
	}

	scanDuration := time.Now().Sub(startScan)

	fmt.Println("Duration of static scan is ",scanDuration)

	jsonFile, err := os.Open(fmt.Sprintf("%s/%s.json", assetDir,imageid))
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	jsonData, _ := ioutil.ReadAll(jsonFile)

	var report apexFn.GeneratedReport
	json.Unmarshal(jsonData, &report)

	staticDir := filepath.Join(assetDir, "smartlog/static")
	templateDir := filepath.Join(assetDir, "smartlog/templates")
	// Make sure all the paths exist.
	tmplPaths := []string{
		staticDir,
		filepath.Join(templateDir, "index.html"),
	}

	for _, path := range tmplPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			fmt.Printf("template %s not found", path)
		}
	}

	// Create mux server.
	mux := mux.NewRouter()
	mux.UseEncodedPath()

	vulnCount = len(report.Vulnerabilities.High) + len(report.Vulnerabilities.Low) + len(report.Vulnerabilities.Medium)
	vm := &StaticReportModel{AssetDir: filepath.Join(templateDir, "index.html"), GenReport: report, ImageName: image, TotalVuln:vulnCount}
	mux.HandleFunc(fmt.Sprintf("/static-report/%s",imageid), vm.StaticHandler)

	// Serve the static assets.
	staticHandler := http.FileServer(http.Dir(staticDir))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", staticHandler))
	mux.Handle("/", staticHandler)
	http.Handle("/",mux)
	http.ListenAndServe(":8580", nil)
	return
}
