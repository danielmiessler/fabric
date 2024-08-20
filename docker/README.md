## Building and Running the Fabric Application with Docker

This guide provides step-by-step instructions on how to build a Fabric image and run inside a docker container. The image is built using a provided Dockerfile.

### Prerequisites:

Ensure you have Docker installed on your system.
Navigate to the Fabric project directory containing the `Dockerfile` and other required files such as `.env`.

Make sure to copy the `.env_example` file to `.env` and add your API keys in the `.env` file or run `fabric --setup and add your API keys once the container is up and running. This will automatically generate the .env file for you.

### Building the Image:

Execute the following command in your terminal (Linux or Mac or wsl on Windows):

``` $ docker build -t fabric:latest .```

Replace < latest tag > with any version identifier, like v1.0 from fabric project.

The . at the end of the command indicates that you're building the image using `DockerFile` available in the current directory.

### Running the Container:

After building the image, run the container using the following command:

```
$ docker run --rm -i -t -d --name fabric -v ./fabric/:/home/appuser/.config/fabric --network traefik fabic:latest
```
Replace < --network traefik > with a desired name for your docker `network`. I have `ollama`, `open-webui`, `whisper` etc. running in the same network called `traefik` so I am using that. You can delete it if you don't want `fabric` to be in a specific network.

 Docker should create a network for you automatically. The container will be detached and run in the background (with the -d flag). Container will be named as `fabric`. The volume mount is used to share the config folder with the container. This will be useful when you want to add your own `.env` file or save container files/folders such as `pattern`. The `-i` flag keeps the container open for interactive commands. The `-t` flag allocates a `pseudo-tty` to the container. `--rm` flag removes the container when it exits.

 #### Docker Compose:

 I have provided a `docker-compose.yml` file as well. You can use it to run the container using `docker compose up`.

#### Tips:

To check if the container is running, use the command: 

`$ docker ps.`

To stop the container, use the command: `$ docker stop fabric`.

To remove the container and its associated resources, use the command: 

`$ docker rm fabric.`
