# Manual for reproducing my steps

**TO RUN `image_for_collecting_data` AND `duckietown_diploma_impl` IMAGE YOU NEED THE DUCKIEBOT ITSELF**

## Github
Repository is available on my [GitHub](https://github.com/SPGC/DoubleProjectArtifact)

The dataset described in report can also be found on my [GitHub](https://github.com/SPGC/DuckiebotDataset) 

## Prerequisites

To be able to build images for Duckietown it's important to use Ubuntu 20.04, 
newer versions of Ubuntu aren't support libraries required for build, older ones weren't tested.
Windows, MacOS and other Linux distributions are not supported at all (you can try, 
but nothing will probably work).

Dependencies:
1. Docker
2. Duckietown shell (DTS)
3. GoLang (optional)
4. Python 3

Firstly you have to follow all steps in [Duckietown manual](https://docs.duckietown.com/daffy/opmanual-duckiebot/setup/setup_laptop/index.html) to install DTS.
When done, ensure that you have installed Python 3 and if you want to run server on you computer without Docker install Go.

If you want to run the code without changing the bot setup, here's bot data:
+ WiFi: SSID: duckie_net, password: quack_quack
+ Bot name: spgc0duckie

## Collecting data

To collect data you have to run server on your PC and run image with data sender node on the Duckiebot.

### Running the server
To run the server you should open `DuckieDatasetServer` folder and type:
`go run main.go`
There are some flags that you can use:
+ `-p` or `--port` to specify port to use in format `:NNNN`
+ `-o` or `--output_path` to specify where to save images
+ `-s` or `--start_number` to specify starting number for images name

You can also run the server in docker. Don't forget to specify port, in Docker port `:8080` is used
and data is saved in folder `/Data`

Don't forget, that every time you change the light in the room you have to save images in a new folder,
otherwise the marking up algorithm can have a lot of noise.

### Running the Duckiebot image

To run the Duckiebot image for collecting data you have to build it, run it and specify the IP address of the server:
1. To build image you have to type: `dts devel build -f -H BOT_NAME`
2. To run you have to type: `dts devel run -M -f -H BOT_NAME -L camera_sender`
3. To specify the IP, you have to type following command, while image is running:
```bash
dts start_gui_tools BOT_NAME
rosparam set /camera_sender_node/ip IP:PORT
```

### Collecting images
To collect images it's useful to operate the Duckiebot with joystick to do that type:
`dts duckiebot keyborad_control BOT_NAME`

## Marking up the dataset
To mark up the dataset use `MarkingUpTheDataset.ipynb` notebook. Don't forget to change folders' paths in cells 2 and 11

## Training the NN
To train NN and choose hyperparameters use `DL_CV.ipynb` notebook. Follow all the instructions in the notebook

## Resulting image
To see the result of the project open `duckietown_diploma_impl` folder:
1. To build image you have to type: `dts devel build -f -H BOT_NAME`
2. To run you have to type: `dts devel run -M -f -H BOT_NAME -L jass`

It will take a while to build and run image.
While image is starting (in running step, not in building) you can open joystick and utility for showing images from camera and NN.
1. `dts duckiebot keyborad_control BOT_NAME` for open joystick
2. To open utility to show images from camera and NN type:
```bash
dts start_gui_tools BOT_NAME
rqt_image_view &
``` 

After you see this line in terminal `Node init successful` you can open Joystick and press `a` key to start autonomous driving

In utility with images you can open `/BOT_NAME/camera_node/image/compressed` to see raw images or NN output

**Important**

In `duckietown_diploma_impl` most of the code isn't written by me, but it's necessary for running the image.

## Dependencies

Most of the dependencies are listed in requirements.txt, however a few can be missed, use pip to install them.

To install dependencies from file type:
```bash
pip install -r requirements.txt
```