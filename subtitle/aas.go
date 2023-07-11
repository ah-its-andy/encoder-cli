package subtitle

import (
	"bufio"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type AAS struct {
	Subtitle []Subtitle `xml:"Subtitle"`
}

type Subtitle struct {
	StartTime string `xml:"InTC"`
	EndTime   string `xml:"OutTC"`
	Data      string `xml:"Data"`
}

func ConvertAAS2SRT(aasFilePath string, srtFilePath string) error {
	// Read the AAS file
	data, err := ioutil.ReadFile(aasFilePath)
	if err != nil {
		return err
	}

	// Parse the XML data
	var aas AAS
	err = xml.Unmarshal(data, &aas)
	if err != nil {
		return err
	}

	// Create the SRT file
	srtFile, err := os.Create(srtFilePath)
	if err != nil {
		return err
	}
	defer srtFile.Close()

	// Write the SRT data
	writer := bufio.NewWriter(srtFile)
	for i, subtitle := range aas.Subtitle {
		startTime, err := parseTimecode(subtitle.StartTime)
		if err != nil {
			return err
		}

		endTime, err := parseTimecode(subtitle.EndTime)
		if err != nil {
			return err
		}

		subtitleData, err := base64.StdEncoding.DecodeString(subtitle.Data)
		if err != nil {
			return err
		}

		subtitleText := string(subtitleData)

		srtData := fmt.Sprintf("%d\n%s --> %s\n%s\n\n", i+1, startTime, endTime, subtitleText)
		_, err = writer.WriteString(srtData)
		if err != nil {
			return err
		}
	}
	writer.Flush()

	return nil
}

func parseTimecode(timecode string) (string, error) {
	parts := strings.Split(timecode, ":")
	if len(parts) != 4 {
		return "", fmt.Errorf("invalid timecode format: %s", timecode)
	}

	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		return "", err
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}

	seconds, err := strconv.Atoi(parts[2])
	if err != nil {
		return "", err
	}

	frames, err := strconv.Atoi(parts[3])
	if err != nil {
		return "", err
	}

	milliseconds := frames * 40

	timeString := fmt.Sprintf("%02d:%02d:%02d,%03d", hours, minutes, seconds, milliseconds)

	return timeString, nil
}
