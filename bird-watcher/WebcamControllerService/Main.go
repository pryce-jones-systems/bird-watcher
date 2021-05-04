package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/blackjack/webcam"
	"image"
	"image/jpeg"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"sort"
	"strconv"
	"time"
)

const VERSION = "bird-watcher-v1.1"

const (
	V4L2_PIX_FMT_PJPG = 0x47504A50
	V4L2_PIX_FMT_YUYV = 0x56595559
)

type FrameSizes []webcam.FrameSize

var supportedFormats = map[webcam.PixelFormat]bool{
	V4L2_PIX_FMT_PJPG: true,
	V4L2_PIX_FMT_YUYV: true,
}
var config Config

func main() {

	// Parse config file
	log.Println("Attempting to read config from /etc/" + VERSION + "/webcam-controller-service-config.json")
	config, err := parseConfigFile("/etc/" + VERSION + "/webcam-controller-service-config.json")
	if err != nil {
		log.Println("Failed to read config from /etc/" + VERSION + "/webcam-controller-service-config.json")
		log.Println("Attempting to read config from config.json")
		config, err = parseConfigFile("config.json")
		if err != nil {
			log.Fatal("Failed to read config from config.json ")
		}
	}
	dev := flag.String("d", config.VideoDevice, "video device to use")
	fmtstr := flag.String("f", "", "video format to use, default first supported")
	single := flag.Bool("m", config.SingleFrameMode, "single image http mode, default mjpeg video")
	addr := flag.String("l", fmt.Sprintf(":%d", config.Port), "addr to listen")
	fps := flag.Bool("p", true, "print fps info")
	flag.Parse()
	log.Println("Successfully parsed config file")

	// Open camera
	cam, err := webcam.Open(*dev)
	if err != nil {
		panic(err.Error())
	}
	defer cam.Close()
	log.Printf("Opened video device %s", config.VideoDevice)

	// Select pixel format
	format_desc := cam.GetSupportedFormats()
	log.Println("Available formats:")
	for _, s := range format_desc {
		log.Printf("\t%s\n", s)
	}

	var format webcam.PixelFormat
FMT:
	for f, s := range format_desc {
		if *fmtstr == "" {
			if supportedFormats[f] {
				format = f
				break FMT
			}

		} else if *fmtstr == s {
			if !supportedFormats[f] {
				log.Fatal(format_desc[f], "format is not supported, exiting")
				return
			}
			format = f
			break
		}
	}
	if format == 0 {
		log.Fatal("No format found, exiting")
		return
	}

	// Select frame size
	frames := FrameSizes(cam.GetSupportedFrameSizes(format))
	sort.Sort(frames)
	log.Println("Supported frame sizes for format:", format_desc[format])
	for _, f := range frames {
		log.Printf("\t%s\n", f.GetString())
	}
	var size *webcam.FrameSize
	if fmt.Sprintf("%dx%d", config.Width, config.Height) == "" {
		size = &frames[len(frames)-1]
	} else {
		for _, f := range frames {
			if fmt.Sprintf("%dx%d", config.Width, config.Height) == f.GetString() {
				size = &f
			}
		}
	}
	if size == nil {
		log.Fatal("No matching frame size, exiting")
		return
	}
	log.Println("Attempting", format_desc[format], fmt.Sprintf("%dx%d", config.Width, config.Height))
	f, w, h, err := cam.SetImageFormat(format, uint32(config.Width), uint32(config.Height))
	if err != nil {
		log.Fatal("Unable to use",fmt.Sprintf("%dx%d", config.Width, config.Height), ". Err msg:", err)
		return

	}
	log.Printf("Success! We're using: %s %dx%d\n", format_desc[f], w, h)

	// Start streaming images over HTTP
	err = cam.StartStreaming()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Begun streaming on port", config.Port)
	var (
		li   chan *bytes.Buffer = make(chan *bytes.Buffer)
		fi   chan []byte        = make(chan []byte)
		back chan struct{}      = make(chan struct{})
	)
	go encodeToImage(cam, back, fi, li, w, h, f)
	if *single {
		go httpImage(*addr, li)
	} else {
		go httpVideo(*addr, li)
	}

	// Loop forever
	timeout := uint32(5) //5 seconds
	start := time.Now()
	var fr time.Duration
	for {

		// Wait for webcam to become ready to deliver a frame
		err = cam.WaitForFrame(timeout)
		if err != nil {
			log.Println(err)
			return
		}
		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			log.Println(err)
			continue
		default:
			log.Println(err)
			return
		}

		// Get frame from webcam
		frame, err := cam.ReadFrame()
		if err != nil {
			log.Println(err)
			return
		}

		// Check that frame is not empty
		if len(frame) != 0 {

			// Log framerate info every 10 seconds
			fr++
			if *fps {
				if d := time.Since(start); d > time.Second*10 {
					log.Println("Streaming at", math.Round(float64(fr)/(float64(d)/float64(time.Second))), "fps")
					start = time.Now()
					fr = 0
				}
			}

			select {
			case fi <- frame:
				<-back
			default:
			}
		}
	}
}

func (slice FrameSizes) Len() int {
	return len(slice)
}

/*
 * Used for sorting
 * @return true if i < j
 */
func (slice FrameSizes) Less(i, j int) bool {
	ls := slice[i].MaxWidth * slice[i].MaxHeight
	rs := slice[j].MaxWidth * slice[j].MaxHeight
	return ls < rs
}

/*
 * Swaps two elements in a slice
 */
func (slice FrameSizes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

/*
 * Asynchronously converts []byte from webcam to image.Image
 */
func encodeToImage(wc *webcam.Webcam, back chan struct{}, fi chan []byte, li chan *bytes.Buffer, w, h uint32, format webcam.PixelFormat) {

	var (
		frame []byte
		img   image.Image
	)

	// Loop forever
	for {

		// Recieve frame from channel
		bframe := <-fi

		// Copy frame
		if len(frame) < len(bframe) {
			frame = make([]byte, len(bframe))
		}
		copy(frame, bframe)
		back <- struct{}{}

		// Perform pixel conversion
		switch format {
		case V4L2_PIX_FMT_YUYV:
			yuyv := image.NewYCbCr(image.Rect(0, 0, int(w), int(h)), image.YCbCrSubsampleRatio422)
			for i := range yuyv.Cb {
				ii := i * 4
				yuyv.Y[i*2] = frame[ii]
				yuyv.Y[i*2+1] = frame[ii+2]
				yuyv.Cb[i] = frame[ii+1]
				yuyv.Cr[i] = frame[ii+3]

			}
			img = yuyv
		default:
			log.Fatal("invalid format ?")
		}

		// Convert to JPEG
		buf := &bytes.Buffer{}
		if err := jpeg.Encode(buf, img, nil); err != nil {
			log.Fatal(err)
			return
		}

		//const N = 50
		// broadcast image up to N ready clients
		nn := 0
	FOR:
		for ; nn < config.MaxSockets; nn++ {
			select {
			case li <- buf:
			default:
				break FOR
			}
		}
		if nn == 0 {
			li <- buf
		}

	}
}

/*
 * Serves a single JPEG frame over HTTP
 * The client must make the request again to get the next frame
 */
func httpImage(addr string, li chan *bytes.Buffer) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Connection from", r.RemoteAddr, r.URL)
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		//remove stale image
		<-li

		img := <-li

		w.Header().Set("Content-Type", "image/jpeg")

		if _, err := w.Write(img.Bytes()); err != nil {
			log.Println(err)
			return
		}

	})

	log.Fatal(http.ListenAndServe(addr, nil))
}

/*
 * Serves an MJPEG stream over HTTP
 * Keeps sending frames until the connection is closed
 */
func httpVideo(addr string, li chan *bytes.Buffer) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Connection from", r.RemoteAddr, r.URL)
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		//remove stale image
		<-li
		const boundary = `frame`
		w.Header().Set("Content-Type", `multipart/x-mixed-replace;boundary=`+boundary)
		multipartWriter := multipart.NewWriter(w)
		multipartWriter.SetBoundary(boundary)
		for {
			img := <-li
			image := img.Bytes()
			iw, err := multipartWriter.CreatePart(textproto.MIMEHeader{
				"Content-type":   []string{"image/jpeg"},
				"Content-length": []string{strconv.Itoa(len(image))},
			})
			if err != nil {
				log.Println(err)
				return
			}
			_, err = iw.Write(image)
			if err != nil {
				log.Println(err)
				return
			}
		}
	})

	log.Fatal(http.ListenAndServe(addr, nil))
}
