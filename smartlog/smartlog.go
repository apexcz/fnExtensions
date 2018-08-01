package smartlog

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
	"io/ioutil"

	"github.com/fnproject/fn/api/models"
	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/fn/fnext"
	"github.com/fsouza/go-dockerclient"
	"github.com/apang1992/jq"

	//_ "github.com/apexcz/fnExtensions/models"
)


func init() {
	server.RegisterExtension(&myLogExt{})
}

type myLogExt struct {

}

func (e *myLogExt) Name() string {
	return "github.com/apexcz/fnExtensions/smartlog"
}

func (e *myLogExt) Setup(s fnext.ExtServer) error {
	s.AddCallListener(&LogListener{})
	return nil
}

type LogListener struct {

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
	
	// Get the path to the asset directory.
	
	assetDir, err := os.Getwd()
	if err != nil {
		return err
	}
	
	fmt.Println("Asset path is ",assetDir)

	go launchClair(call.Image,assetDir)
	time.Sleep(1*time.Second)
	
	/**
	staticDir := filepath.Join(assetDir, "static")
	templateDir := filepath.Join(assetDir, "templates")
	// Make sure all the paths exist.
	tmplPaths := []string{
		staticDir,
		filepath.Join(templateDir, "vulns.html"),
		filepath.Join(templateDir, "repositories.html"),
		filepath.Join(templateDir, "tags.html"),
	}
	for _, path := range tmplPaths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("template %s not found", path)
		}
	}
	*/
	return nil
}

func (l *LogListener) AfterCall(ctx context.Context, call *models.Call) error {
	fmt.Println("Triggers after function executes completely")
	return nil
}

func launchClair(image string,asset string) {
	instruction := fmt.Sprintf("CLAIR_ADDR=localhost CLAIR_OUTPUT=Low CLAIR_THRESHOLD=3 JSON_OUTPUT=true klar postgres:9.5.1 > %s/staticreport.json", asset)
	
	cmd := exec.Command("sh","-c",instruction)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(err.Error())		
	}

	jsonFile, err := os.Open(fmt.Sprintf("%s/staticreport.json", asset))
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	jsonData, _ := ioutil.ReadAll(jsonFile)

	layerCount, err := jq.JsonQuery([]byte(jsonData), ".LayerCount")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		fmt.Println("the layerCount is:", string(layerCount))
	}
}