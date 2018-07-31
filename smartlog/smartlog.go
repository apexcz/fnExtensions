package smartlog

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/fnproject/fn/api/models"
	"github.com/fnproject/fn/api/server"
	"github.com/fnproject/fn/fnext"
	"github.com/fsouza/go-dockerclient"
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
	
	cmd := "klar"
	args := []string{"CLAIR_ADDR=localhost", "CLAIR_OUTPUT=High", "CLAIR_THRESHOLD=10", "JSON_OUTPUT=true","postgres:9.5.1", ">", "$HOME/Documents/result.json"}
	
	if err := exec.Command(cmd,args...).Run(); err != nil {
		//log.Fatal(err)
		fmt.Println("Error = ",err)
	}

	return nil
}

func (l *LogListener) AfterCall(ctx context.Context, call *models.Call) error {
	fmt.Println("Triggers after function executes completely")
	return nil
}