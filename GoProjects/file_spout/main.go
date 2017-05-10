package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type FileVol struct {
	filePath string
}

//startContainer is a way which can call docker host to run a container
//para:
//     imageRef:it is a ref for images repository
//     imagesName:it is the name of image,if not set tag,default tag is latest
//     cmdString:the command which can run in docker// addr string, fileName string
func startContainer(imageRef string, imageName string, addr string, fileName string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	_, err = cli.ImagePull(ctx, imageRef, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		//		WorkingDir: "/testfilespoutaa",
		//		Entrypoint: []string{"/bin/sh"},
		Volumes: map[string]struct{}{
			"/home/deepglint/data_0": {},
		},
		Cmd: []string{"./filespoutaa", "-addr", addr, "-file", fileName},
		//		Shell:        []string{"pwd"},
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Tty:          true,
	}, &container.HostConfig{
		Binds: []string{"/home/deepglint/data_0/111.264:/testfilespoutaa/data/111.264", "/home/deepglint/data_0/test.h264:/testfilespoutaa/data/test.h264"},
	}, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, out)

}

func main() {
	addr := flag.String("addr", "", "mserverip:mserverudpport")
	filepath := flag.String("file", "", "h264 file path")
	flag.Parse()
	fmt.Printf("input args:addr = %s, 264filepath = %s \n", *addr, *filepath) //*addr, *filepath
	imageRepository := "docker.io/cjydocker/filespoutaa:1.0"
	imageName := "cjydocker/filespoutaa:1.0"
	fmt.Println(imageRepository)
	startContainer(imageRepository, imageName, *addr, *filepath)
}
