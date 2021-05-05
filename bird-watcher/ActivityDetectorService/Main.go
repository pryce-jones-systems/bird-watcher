package main

import (
	"errors"
	"fmt"
	"github.com/pryce-jones-systems/go-image-tools/ImageTools"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const VERSION = "bird-watcher-v1.1"

var RAM_DISK_DIR = getRAMDiskDir()
var INPUT_FRAMERATE = 10
var OUTPUT_FRAMERATE = 10 * INPUT_FRAMERATE
var DELAY_BETWEEN_FRAMES = time.Duration(0)//time.Duration(float64(1 /INPUT_FRAMERATE) * 100)
var config Config

type FrameBuffer struct {
	frame [][][]float32
}

func main() {

	// Parse config file
	log.Println("Attempting to read config from /etc/" + VERSION + "/activity-detector-service-config.json")
	config, err := parseConfigFile("/etc/" + VERSION + "/activity-detector-service-config.json")
	if err != nil {
		log.Println("Failed to read config from /etc/" + VERSION + "/activity-detector-service-config.json")
		log.Println("Attempting to read config from config.json")
		config, err = parseConfigFile("config.json")
		if err != nil {
			log.Fatal("Failed to read config from config.json ")
		}
	}
	log.Println("Successfully parsed config file")

	// Create wait group for saving frames to the RAM disk
	var frameSaveWaitGroup sync.WaitGroup

	// Create frame buffer
	frameBuffer := FrameBuffer{}

	// Initially fill frame buffer
	log.Println("Filling the frame buffer")
	counter := 0
	for fr := 0 ; fr < config.FrameBufferSize; fr++ {

		// Get frame from stream
		grayFrame, colFrame, _ := getFrame(config.WebcamURL)

		// Put grayscale frame into buffer
		frameBuffer.frame = append(frameBuffer.frame, grayFrame)

		// Asynchronously save colour frame
		go saveFrameToRAMDisk(frameSaveWaitGroup, colFrame, fmt.Sprintf("%s/fr%08d.jpg", RAM_DISK_DIR, counter))
	}

	// Create channels to pass information about frame activity
	activityChannel, inactivityChanel := make(chan int), make(chan int)

	// Create array to store frame activity
	//  1 means activity in this frame
	//  0 means no activity in this frame
	// -1 means activity in this frame hasn't been checked yet
	activeFrames := make([]int, config.FrameBufferSize)
	for i := 0; i < config.FrameBufferSize; i++ {
		activeFrames[i] = -1
	}

	// Wait for all colour frames to be saved to the RAM disk
	frameSaveWaitGroup.Wait()
	log.Println("Done filling the frame buffer")

	// Loop forever
	log.Println("Object detection started")
	for {

		// Get frame from stream
		grayFrame, colFrame, _ := getFrame(config.WebcamURL)

		// Put grayscale frame into buffer
		frameBuffer = addFrameToBuffer(grayFrame, frameBuffer)

		// Mark current frame as not having been checked for activity (assign a value of -1 in the array)
		activeFrames[counter] = -1

		// Asynchronously save colour frame
		go saveFrameToRAMDisk(frameSaveWaitGroup, colFrame, fmt.Sprintf("%s/fr%08d.jpg", RAM_DISK_DIR, counter))

		// Asynchronously detect activity in frame
		go detectActivity(activityChannel, inactivityChanel, counter, config.ActivityThreshold, grayFrame, frameBuffer)
		select {
		case frameNumber := <- activityChannel:
			activeFrames[frameNumber] = 1
			break
		case frameNumber := <- inactivityChanel:
			activeFrames[frameNumber] = 0
			break
		default:
			break
		}

		// Only frame if buffer is full
		if counter == config.FrameBufferSize - 1 {
			log.Println("Checking if video capture is required")

			// Check if enough consecutive frames are active to trigger recording
			maxConsecutiveFrames := 0
			consecutiveFrames := 0
			if activeFrames[0] == 1 {
				consecutiveFrames = 1
			}
			for fr := 1; fr < config.FrameBufferSize; fr++ {

				// Frame is active and consecutive with the previous active frame
				if activeFrames[fr] == 1 {
					consecutiveFrames++
				}

				if maxConsecutiveFrames < consecutiveFrames {
					maxConsecutiveFrames = consecutiveFrames
				}

				// Run of consecutive active frames ends
				if activeFrames[fr] != 1 {
					consecutiveFrames = 0
				}
			}
			log.Println("\tMaximum consecutive active frames in buffer:", maxConsecutiveFrames)

			// Wait for all colour frames to finish being saved to the RAM disk
			frameSaveWaitGroup.Wait()

			// Save video if required
			// This must be synchronous or we run the risk of race conditions with the routines writing frames
			if maxConsecutiveFrames >= config.ConsecutiveActiveFrames {
				path := fmt.Sprintf("%s/%d-%d-%d-%d-%d-%d.mp4", config.OutputDir, time.Now().Year(), time.Now().YearDay(), time.Now().Hour(), time.Now().Minute(), time.Now().Second(), time.Now().Nanosecond())
				log.Print("\tSaving video at ", path)
				if err := encodeVideo(path); err != nil {
					log.Println(err)
				}
			} else {
				log.Println("\tNo video capture required")
			}

			/*
			// Purge buffer with new frames
			// This stops the same activity from triggering more than one video recording
			counter = 0
			for fr := 0 ; fr < MAX_RAMDISK_FRAMES; fr++ {

				// Get frame from stream
				grayFrame, colFrame, _ := getFrame(MJPEG_STREAM_URL)

				// Put grayscale frame into buffer
				frameBuffer.frame = append(frameBuffer.frame, grayFrame)

				// Asynchronously save colour frame
				go saveFrameToRAMDisk(frameSaveWaitGroup, colFrame, fmt.Sprintf("%s/fr%08d.jpg", RAM_DISK_DIR, counter))

				// Wait so that framerate is maintained
				time.Sleep(DELAY_BETWEEN_FRAMES * time.Millisecond)
			}

			// Wait for all colour frames to be saved to the RAM disk
			frameSaveWaitGroup.Wait()
			 */
		}

		counter++
		counter %= config.FrameBufferSize
	}
}

func detectActivity(activityChan chan int, inactivityChan chan int, frameNumber int, activityThreshold float32, frame [][]float32, frameBuffer FrameBuffer) {

	// Average all frames in the buffer
	background := frameBuffer.frame[0]
	for fr := 1 ; fr < config.FrameBufferSize; fr++ {
		frame := frameBuffer.frame[fr]
		background, _ = ImageTools.Add(background, frame)
	}

	// Subtract background from frame
	subtracted, _ := ImageTools.Subtract(background, frame)

	// Apply threshold to create mask
	mean, std := ImageTools.MeanStd(subtracted)
	threshold := mean + (0.5 * std)
	foregroundMask := ImageTools.SingleThreshold(subtracted, threshold)
	mean, std = ImageTools.MeanStd(foregroundMask)

	//ImageTools.SaveImage("/home/jake/Desktop/subtracted.jpg", subtracted)
	//fmt.Println(mean, std)
	//ImageTools.SaveImage("/home/jake/Desktop/fg.jpg", foregroundMask)

	if std > activityThreshold {
		inactivityChan <- frameNumber
		return
	}
	activityChan <- frameNumber
	return
}

func saveFrameToRAMDisk(waitGroup sync.WaitGroup, frame image.Image, path string) {

	// Add this routine to the wait group
	waitGroup.Add(1)

	// Write timestamp on frame
	bounds := frame.Bounds()
	newFrame := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	draw.Draw(newFrame, newFrame.Bounds(), frame, bounds.Min, draw.Src)
	addLabel(newFrame, bounds.Max.X - 150, bounds.Max.Y - 20, "Pryce-Jones Systems")
	addLabel(newFrame, bounds.Max.X - 150, bounds.Max.Y - 10, strings.Split(time.Now().String(), ".")[0])

	// Open file
	file, err := os.Create(path)
	if err != nil {
		return
	}
	defer file.Close()

	// Encode image
	options := jpeg.Options{Quality: 100}
	err = jpeg.Encode(file, newFrame, &options)

	// Signal that saving is done
	waitGroup.Done()
}

func addLabel(img *image.RGBA, x, y int, label string) {
	col := color.RGBA{0, 0, 255, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}

func encodeVideo(outPath string) error {

	// Use FFMPEG to encode images in ramdisk to a video
	cmd, err := exec.LookPath("ffmpeg")
	if err != nil {
		return err
	}
	args := []string{
		"-y",
		"-framerate", fmt.Sprintf("%d", INPUT_FRAMERATE),
		"-pattern_type", "glob",
		"-i", fmt.Sprintf("%s/*.jpg", RAM_DISK_DIR),
		"-c:v", "libx264",
		"-r", fmt.Sprintf("%d", OUTPUT_FRAMERATE),
		outPath,
	}
	ffmpeg := exec.Command(cmd, args...)
	err = ffmpeg.Start()
	if err != nil {
		return err
	}
	ffmpeg.Wait()
	return nil
}

/*
 * Returns the path of the OS' RAM disk
 */
func getRAMDiskDir() string {

	err := ioutil.WriteFile("/dev/shm/t", []byte{255}, 0600)

	// Use the RAM-disk, if the OS has one
	if err == nil {
		os.Remove("/dev/shm/t")
		return "/dev/shm"
	} else {
		return os.TempDir()
	}
}

func addFrameToBuffer(frame [][]float32, frameBuffer FrameBuffer) FrameBuffer {

	// Iterate over buffer and shift each flame along 1 place
	for fr := len(frameBuffer.frame) - 2; fr > 0; fr-- {
		frameBuffer.frame[fr] = frameBuffer.frame[fr + 1]
	}

	// Add frame to buffer
	frameBuffer.frame[len(frameBuffer.frame) - 1] = frame

	return frameBuffer
}

func getFrame(url string) ([][]float32, image.Image, error) {

	response, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, nil, errors.New(http.StatusText(response.StatusCode))
	}

	img, err := jpeg.Decode(response.Body)
	if err != nil {
		return nil, nil, err
	}

	return ImageTools.Image2Slice(img), img, nil
}

