package static

import (
	"github.com/docker/docker/reference"
	"github.com/coreos/clair/api/v3/clairpb"
)

func (a *apiV3) Push(image *docker.Image) error {
	req := &clairpb.PostAncestryRequest{
		Format:       "Docker",
		//AncestryName: image.Name,
		AncestryName: reference.NamedTagged
	}

	ls := make([]*clairpb.PostAncestryRequest_PostLayer, len(image.FsLayers))
	for i := 0; i < len(image.FsLayers); i++ {
		ls[i] = newLayerV3(image, i)
	}
	req.Layers = ls
	_, err := a.client.PostAncestry(context.Background(), req)
	return err
}