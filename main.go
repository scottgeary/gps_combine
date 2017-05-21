package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type gpx struct {
	TrackList []trk `xml:"trk"`
}

type trk struct {
	SegmentList []trkseg `xml:"trkseg"`
}

type trkseg struct {
	PointList []trkpt `xml:"trkpt"`
}

type trkpt struct {
	Latitude  float64   `xml:"lat,attr"`
	Longitude float64   `xml:"lon,attr"`
	Timestamp time.Time `xml:"time"`
}

func readGpx(filename string) (gpx, error) {

	fmt.Printf("Opening %s...\n", filename)

	empty := gpx{}

	xmlFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return empty, err
	}
	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)

	var q gpx
	xml.Unmarshal(b, &q)

	// fmt.Println(q)

	return q, nil

}

func isEqual(x trkpt, test trkpt) bool {
	return test.Timestamp.Equal(x.Timestamp)
}

func isBefore(x trkpt, test trkpt) bool {
	return test.Timestamp.Before(x.Timestamp)
}

func isAfter(x trkpt, test trkpt) bool {
	return test.Timestamp.After(x.Timestamp)
}

func isInside(x trkpt, y trkpt, test trkpt) bool {
	if isAfter(x, test) && isBefore(y, test) {
		return true
	} else {
		return false
	}
}

func main() {

	aFile := flag.String("a", "", "A File")
	bFile := flag.String("b", "", "B File")
	outFilename := flag.String("o", "", "Out File")
	flag.Parse()

	aSet, aErr := readGpx(*aFile)
	if aErr != nil {
		fmt.Println(aErr)
		os.Exit(1)
	}
	bSet, bErr := readGpx(*bFile)
	if bErr != nil {
		fmt.Println(bErr)
		os.Exit(2)
	}
	fmt.Println(bSet)

	outSet := aSet

	for _, test_track := range bSet.TrackList {
		for _, test_segment := range test_track.SegmentList {
			for i, test_point := range test_segment.PointList {
				fmt.Printf("Considering bSet index %d\n", i)

				x := trkpt{}
				y := trkpt{}
				changed := false
			OuterLoop:
				for a, known_track := range outSet.TrackList {
					for b, known_segment := range known_track.SegmentList {
						for c, known_point := range known_segment.PointList {
							x = y
							y = known_point
							if isInside(x, y, test_point) {
								fmt.Printf("test_point IS inside a pairing. To insert at %d,%d,%d\n", a, b, c)
								outSet.TrackList[a].SegmentList[b].PointList = append(known_segment.PointList[:c], append([]trkpt{test_point}, known_segment.PointList[c:]...)...)
								changed = true
							}
							if isEqual(x, test_point) {
								fmt.Printf("test_point appears at an equal time. Insert near %d,%d,%d\n", a, b, c)
								outSet.TrackList[a].SegmentList[b].PointList = append(known_segment.PointList[:c], append([]trkpt{test_point}, known_segment.PointList[c:]...)...)
								changed = true
							}
							if changed {
								break OuterLoop
							}
						}
					}
				}
				if changed == false && isAfter(y, test_point) {
					fmt.Printf("test_point IS at the end. Work out the index now.\n")
					a := len(outSet.TrackList) - 1
					b := len(outSet.TrackList[a].SegmentList) - 1
					outSet.TrackList[a].SegmentList[b].PointList = append(outSet.TrackList[a].SegmentList[b].PointList[:], test_point)
				}

			}
		}
	}

	fmt.Println("Marshalling up...")

	aOut, _ := xml.MarshalIndent(aSet, "", "\t")
	// bOut, _ := xml.MarshalIndent(bSet, "", "\t")
	// fmt.Printf("%s", bytesOut)

	fmt.Println("Writing out...")

	out, err := os.Create(*outFilename)
	fmt.Println(err)
	fmt.Fprintf(out, "%s%s", xml.Header, aOut)
	defer out.Close()

}
