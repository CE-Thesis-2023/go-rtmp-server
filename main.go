package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"path"
	"runtime"
	"time"

	"github.com/gwuhaolin/livego/configure"
	"github.com/gwuhaolin/livego/protocol/api"
	"github.com/gwuhaolin/livego/protocol/hls"
	"github.com/gwuhaolin/livego/protocol/httpflv"
	"github.com/gwuhaolin/livego/protocol/rtmp"

	log "github.com/sirupsen/logrus"
)

var VERSION = "master"

// func startHls() *hls.Server {
// 	hlsAddr := configure.Config.GetString("hls_addr")
// 	hlsListen, err := net.Listen("tcp", hlsAddr)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	hlsServer := hls.NewServer()
// 	go func() {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				log.Error("HLS server panic: ", r)
// 			}
// 		}()
// 		log.Info("HLS listen On ", hlsAddr)
// 		hlsServer.Serve(hlsListen)
// 	}()
// 	return hlsServer
// }

func startFFMPEG() {

}

func startRtmp(stream *rtmp.RtmpStream, hlsServer *hls.Server) {
	rtmpAddr := configure.Config.GetString("rtmp_addr")
	isRtmps := configure.Config.GetBool("enable_rtmps")

	var rtmpListener net.Listener
	//check for RTMPs
	if isRtmps {
		certPath := configure.Config.GetString("rtmps_cert")
		keyPath := configure.Config.GetString("rtmps_key")
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			log.Fatal(err)
		}

		//sever listen to this addr
		rtmpListener, err = tls.Listen("tcp", rtmpAddr, &tls.Config{
			Certificates: []tls.Certificate{cert},
		})
		if err != nil {
			log.Fatal(err)
		}
		log.Info("RTMPS Listen On ", rtmpAddr)
	} else {
		var err error
		rtmpListener, err = net.Listen("tcp", rtmpAddr)
		if err != nil {
			log.Fatal(err)
		}
		log.Info("RTMP Listen On ", rtmpAddr)

	}

	var rtmpServer *rtmp.Server

	//check for hlsServer
	if hlsServer == nil {
		rtmpServer = rtmp.NewRtmpServer(stream, nil)
		log.Info("HLS server disable....")
	} else {
		rtmpServer = rtmp.NewRtmpServer(stream, hlsServer)
		log.Info("HLS server enable....")
	}

	defer func() {
		if r := recover(); r != nil {
			log.Error("RTMP server panic: ", r)
		}
	}()

	//serve the listener
	rtmpServer.Serve(rtmpListener)
}

func startHTTPFlv(stream *rtmp.RtmpStream) {
	httpflvAddr := configure.Config.GetString("httpflv_addr")

	flvListen, err := net.Listen("tcp", httpflvAddr)
	if err != nil {
		log.Fatal(err)
	}

	hdlServer := httpflv.NewServer(stream)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("HTTP-FLV server panic: ", r)
			}
		}()
		log.Info("HTTP-FLV listen On ", httpflvAddr)
		hdlServer.Serve(flvListen)
	}()
}

func startAPI(stream *rtmp.RtmpStream) {
	apiAddr := configure.Config.GetString("api_addr")   //8090
	rtmpAddr := configure.Config.GetString("rtmp_addr") //1935

	if apiAddr != "" {
		opListen, err := net.Listen("tcp", apiAddr)
		if err != nil {
			log.Fatal(err)
		}
		opServer := api.NewServer(stream, rtmpAddr)
		go func() {
			defer func() { // if this thread failed somehow -> this defer block will be called
				if r := recover(); r != nil {
					log.Error("HTTP-API server panic: ", r)
				}
			}()
			log.Info("HTTP-API listen On ", apiAddr)
			opServer.Serve(opListen)
		}()
	}
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf(" %s:%d", filename, f.Line)
		},
	})
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			log.Error("livego panic: ", r)
			time.Sleep(1 * time.Second)
		}
	}()

	log.Infof(`Start Server: %s`, VERSION)

	apps := configure.Applications{}
	configure.Config.UnmarshalKey("server", &apps)
	for _, app := range apps {
		stream := rtmp.NewRtmpStream()
		// var hlsServer *hls.Server
		// if app.Hls {
		// 	hlsServer = startHls()
		// }
		// if app.Flv {
		// 	startHTTPFlv(stream)
		// }
		if app.Api {
			startAPI(stream)
		}

		startRtmp(stream, nil)
	}
}
