// This is an example on how to create a custom text marker.

package main

import (
	"fmt"
	sm "github.com/flopp/go-staticmaps"
	"github.com/fogleman/gg"
	"github.com/golang/geo/s2"
	"github.com/google/uuid"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"image/color"
	"math"
	"strconv"
	"strings"
)

// TextMarker is struct containing info used to draw on the map
type TextMarker struct {
	sm.MapObject
	Position   s2.LatLng
	Text       string
	TextWidth  float64
	TextHeight float64
	TipSize    float64
	LineWidth  float64
}

// InfoTextMarker is a TextMarker with area text written above it, fitted to the width of the max line length, and to the height of the number of lines, plus a margin
func InfoTextMarker(pos s2.LatLng, text string) *TextMarker {
	s := new(TextMarker)
	s.Position = pos
	s.Text = text
	s.TipSize = 8.0
	s.LineWidth = 10.0

	d := &font.Drawer{
		// Face: basicfont.Face7x13,
		Face: inconsolata.Bold8x16,
	}
	sp := strings.Split(s.Text, "\n")
	var maxSize float32 = 0
	for _, elem := range sp {
		if float32(len(elem)) > maxSize {
			maxSize = float32(len(elem))
		}
	}
	var useStr string = s.Text
	if maxSize > 0 {
		useStr = s.Text[0:int(maxSize)]
	}

	marginHeight := 1.0
	if len(sp) > 1 {
		marginHeight = 2.0
	}
	s.TextWidth = float64(d.MeasureString(useStr+"margin") >> 6)
	s.TextHeight = 13.0 * (float64(strings.Count(text, "\n")) + marginHeight)
	return s
}

// ExtraMarginPixels is a method used to add a margin to the text area
func (s *TextMarker) ExtraMarginPixels() (float64, float64, float64, float64) {
	w := math.Max(4.0+s.TextWidth, 2*s.TipSize)
	h := s.TipSize + s.TextHeight + 4.0
	return w * 0.5, h, w * 0.5, 0.0
}

// Bounds is a method used to create empty bounds for the text area
func (s *TextMarker) Bounds() s2.Rect {
	r := s2.EmptyRect()
	r = r.AddPoint(s.Position)
	return r
}

// Draw is an implemented functionality so it can be rendered
func (s *TextMarker) Draw(gc *gg.Context, trans *sm.Transformer) {
	if !sm.CanDisplay(s.Position) {
		return
	}

	w := math.Max(4.0+s.TextWidth, 2*s.TipSize)
	h := s.TextHeight + 4.0
	x, y := trans.LatLngToXY(s.Position)
	gc.ClearPath()
	gc.SetLineWidth(1)
	gc.SetLineCap(gg.LineCapRound)
	gc.SetLineJoin(gg.LineJoinRound)
	gc.LineTo(x, y)
	gc.LineTo(x-s.TipSize, y-s.TipSize)
	gc.LineTo(x-w*0.5, y-s.TipSize)
	gc.LineTo(x-w*0.5, y-s.TipSize-h)
	gc.LineTo(x+w*0.5, y-s.TipSize-h)
	gc.LineTo(x+w*0.5, y-s.TipSize)
	gc.LineTo(x+s.TipSize, y-s.TipSize)
	gc.LineTo(x, y)
	gc.SetColor(color.RGBA{0xff, 0xff, 0xff, 0xff})
	gc.FillPreserve()
	gc.SetColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.Stroke()

	gc.SetRGBA(0.0, 0.0, 0.0, 1.0)
	sp := strings.Split(s.Text, "\n")
	for i, j := 0, len(sp)-1; i < j; i, j = i+1, j-1 {
		sp[i], sp[j] = sp[j], sp[i]
	}
	for index, element := range sp {
		gc.DrawString(strings.TrimSpace(element), x-s.TextWidth*0.5, (y - ((float64(index)*1.3 + 1.0) * (s.TipSize + 4.0))))
	}
	// gc.DrawString(s.Text, x-s.TextWidth*0.5, y-s.TipSize-4.0)
}

// buildMap is a method used to iterate over the input IPs, query their location in the DB, and draw them on the street map
func buildMap(arr []*Displayable, outfilepth *string, useMarker *bool) {
	ctx := sm.NewContext()
	ctx.SetSize(400, 300)
	const ZOOM_OUT_FACTOR = 1000

	ctx.SetSize(2.5*ZOOM_OUT_FACTOR, ZOOM_OUT_FACTOR)
	LatIP := make(map[float64]string, len(arr))
	LonIP := make(map[float64]string, len(arr))  // maps to string uuid
	UUIDMap := make(map[string]string, len(arr)) // uuid maps to stringified list of ip sep by comma
	UUIDOrg := make(map[string]string, len(arr)) // uuid maps to stringified list of Org sep by comma
	UUIDASN := make(map[string]string, len(arr)) // uuid maps to stringified list of ASN sep by comma

	for _, el := range arr {

		Cityrecord := el.City
		if (""+Cityrecord.Country.Names["en"] == "") && (""+Cityrecord.City.Names["en"] == "") {
			fmt.Printf("Warning: no information found for %v\n", el.IPAddress)
			continue
		}

		ASNrecord := el.Asn
		ASNNumber := strconv.FormatUint(uint64(ASNrecord.AutonomousSystemNumber), 10)
		Org := ASNrecord.AutonomousSystemOrganization

		corUUIDLat, containedLat := LatIP[Cityrecord.Location.Latitude]
		_, containedLon := LonIP[Cityrecord.Location.Longitude]
		useip := el.IPAddress
		useorg := Org
		useasn := ASNNumber

		if !containedLat && !containedLon {
			uuidS := uuid.New().String()
			LatIP[Cityrecord.Location.Latitude] = uuidS
			LonIP[Cityrecord.Location.Longitude] = uuidS
			UUIDMap[uuidS] = el.IPAddress
			UUIDASN[uuidS] = ASNNumber
			UUIDOrg[uuidS] = Org
			if el.ConnectedTo != nil {
				for _, d := range *el.ConnectedTo {
					if (""+d.City.Country.Names["en"] == "") && (""+d.City.City.Names["en"] == "") {
						continue
					}
					platlon := []s2.LatLng{s2.LatLngFromDegrees(Cityrecord.Location.Latitude, Cityrecord.Location.Longitude), s2.LatLngFromDegrees(d.City.Location.Latitude, d.City.Location.Longitude)}
					col := color.RGBA{0xff, 0, 0, 0xff}
					pth := sm.NewPath(platlon, col, 2.0)
					ctx.AddObject(pth)
				}
			}

		} else {
			if !strings.Contains(UUIDMap[corUUIDLat], el.IPAddress) {
				UUIDMap[corUUIDLat] += ", " + el.IPAddress
			}
			if !strings.Contains(UUIDASN[corUUIDLat], ASNNumber) {
				UUIDASN[corUUIDLat] += ", " + ASNNumber
			}
			if !strings.Contains(UUIDOrg[corUUIDLat], Org) {
				UUIDOrg[corUUIDLat] += ", " + Org
			}
			useip = UUIDMap[corUUIDLat]
			useasn = UUIDASN[corUUIDLat]
			useorg = UUIDOrg[corUUIDLat]

		}

		if *useMarker {
			rec := sm.NewMarker(
				s2.LatLngFromDegrees(Cityrecord.Location.Latitude, Cityrecord.Location.Longitude),
				color.RGBA{0xff, 0, 0, 0xff},
				16.0*(ZOOM_OUT_FACTOR/300),
			)
			ctx.AddObject(rec)

		} else {
			WriteText := "IP: " + useip + "\nCountry: " + Cityrecord.Country.Names["en"]
			if len(Cityrecord.Subdivisions) > 0 {
				WriteText += "\nRegion: " + Cityrecord.Subdivisions[0].Names["en"]
			}
			WriteText += "\nCity: " + Cityrecord.City.Names["en"] + "\nOrg: " + useorg + "\nASN: " + useasn
			if el.ExitNode {
				WriteText += "\n " + el.ServiceType + " exit node"
			}

			rec := InfoTextMarker(s2.LatLngFromDegrees(Cityrecord.Location.Latitude, Cityrecord.Location.Longitude), WriteText)
			ctx.AddObject(rec)
		}
	}

	img, err := ctx.Render()
	if err != nil {
		panic(err)
	}

	if err := gg.SavePNG(*outfilepth, img); err != nil {
		panic(err)
	}
}
