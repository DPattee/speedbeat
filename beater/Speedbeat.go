package beater

import (
	"errors"
	"fmt"
	"time"
	"os/exec"
	"encoding/json"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"

	"github.com/dpattee/speedbeat/config"
)

// speedbeat configuration.
type speedbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}



// New creates an instance of speedbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &speedbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts speedbeat.
func (bt *speedbeat) Run(b *beat.Beat) error {
	logp.Info("speedbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

    var testResult Speedresult
    testResult, err = bt.DoTest()

    if err != nil {
			logp.Warn("Event not sent")
    } else {
			event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"isp": testResult.ISP,
				"resulturl": testResult.Result.URL,
				"test_timestamp": testResult.Timestamp,
				"ping.latency": testResult.Ping.Latency,
				"ping.jitter": testResult.Ping.Jitter,
				"packetloss": testResult.PacketLoss,
				"download.bandwidth": testResult.Download.Bandwidth,
				"download.bandwidth_mbps": testResult.Download.Bandwidth*8,
				"download.elapsed": testResult.Download.Elapsed,
				"upload.bandwidth": testResult.Upload.Bandwidth,
				"upload.bandwidth_mbps": testResult.Upload.Bandwidth*8,
				"upload.elapsed": testResult.Upload.Elapsed,
				"testserver.name": testResult.Server.Name,
				"testserver.location": testResult.Server.Location,
				"testserver.ip": testResult.Server.IP,
			},
		}
		  bt.client.Publish(event)
		  logp.Info("Event sent")
		}
	}
}

// Stop stops speedbeat.
func (bt *speedbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

// Run the CLI and parse results
func (bt *speedbeat) DoTest() (Speedresult, error) {
	logp.Info("Speedtest starting")
	var tempspeedresult Speedresult

  //set up the command line
	cmd := exec.Command("speedtest", "--format=json")
logp.Info(cmd.String())
	//run command
	 output, err := cmd.Output()
	 if err != nil {
    logp.Warn("Speedtest Error")
		fmt.Println( "Error:", err)
		fmt.Println( "Output:", output)
		return tempspeedresult, errors.New("Speedtest CLI Error")
	}

	logp.Info("Speedtest Success")
	fmt.Printf( "Output: %s\n", output)

	//parse the results

	tempResultJSON := output
err = json.Unmarshal([]byte(tempResultJSON), &tempspeedresult)

  if err != nil {
		logp.Warn("Test Result Parsing Error")
		fmt.Println( "Error:", err)
		return tempspeedresult, errors.New("Speedtest Result Parsing Error")
	}

  return tempspeedresult, nil
}
