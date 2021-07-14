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

// TextMarker
type TextMarker struct {
	sm.MapObject
	Position   s2.LatLng
	Text       string
	TextWidth  float64
	TextHeight float64
	TipSize    float64
	LineWidth  float64
}

func NewTextMarker(pos s2.LatLng, text string) *TextMarker {
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

func (s *TextMarker) ExtraMarginPixels() (float64, float64, float64, float64) {
	w := math.Max(4.0+s.TextWidth, 2*s.TipSize)
	h := s.TipSize + s.TextHeight + 4.0
	return w * 0.5, h, w * 0.5, 0.0
}

func (s *TextMarker) Bounds() s2.Rect {
	r := s2.EmptyRect()
	r = r.AddPoint(s.Position)
	return r
}

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

func buildMap(arr []*Displayable, outfilepth *string, useMarker *bool) {
	ctx := sm.NewContext()
	ctx.SetSize(400, 300)
	const ZOOM_OUT_FACTOR = 1000

	ctx.SetSize(2.5*ZOOM_OUT_FACTOR, ZOOM_OUT_FACTOR)
	Lat_Ip := make(map[float64]string, len(arr))
	Lon_Ip := make(map[float64]string, len(arr))  // maps to string uuid
	Uuid_Map := make(map[string]string, len(arr)) // uuid maps to stringified list of ip sep by comma
	Uuid_Org := make(map[string]string, len(arr)) // uuid maps to stringified list of Org sep by comma
	Uuid_ASN := make(map[string]string, len(arr)) // uuid maps to stringified list of ASN sep by comma

	for _, el := range arr {

		Cityrecord := el.City
		if (""+Cityrecord.Country.Names["en"] == "") && (""+Cityrecord.City.Names["en"] == "") {
			fmt.Printf("Warning: no information found for %v\n", el.IP_address)
			continue
		}

		ASNrecord := el.Asn
		ASN_Number := strconv.FormatUint(uint64(ASNrecord.AutonomousSystemNumber), 10)
		Org := ASNrecord.AutonomousSystemOrganization

		cor_uuid_lat, containedLat := Lat_Ip[Cityrecord.Location.Latitude]
		_, containedLon := Lon_Ip[Cityrecord.Location.Longitude]
		useip := el.IP_address
		useorg := Org
		useasn := ASN_Number

		if !containedLat && !containedLon {
			uuidS := uuid.New().String()
			Lat_Ip[Cityrecord.Location.Latitude] = uuidS
			Lon_Ip[Cityrecord.Location.Longitude] = uuidS
			Uuid_Map[uuidS] = el.IP_address
			Uuid_ASN[uuidS] = ASN_Number
			Uuid_Org[uuidS] = Org
			if el.Connected_to != nil {
				for _, d := range *el.Connected_to {
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
			if !strings.Contains(Uuid_Map[cor_uuid_lat], el.IP_address) {
				Uuid_Map[cor_uuid_lat] += ", " + el.IP_address
			}
			if !strings.Contains(Uuid_ASN[cor_uuid_lat], ASN_Number) {
				Uuid_ASN[cor_uuid_lat] += ", " + ASN_Number
			}
			if !strings.Contains(Uuid_Org[cor_uuid_lat], Org) {
				Uuid_Org[cor_uuid_lat] += ", " + Org
			}
			useip = Uuid_Map[cor_uuid_lat]
			useasn = Uuid_ASN[cor_uuid_lat]
			useorg = Uuid_Org[cor_uuid_lat]

		}

		if *useMarker {
			rec := sm.NewMarker(
				s2.LatLngFromDegrees(Cityrecord.Location.Latitude, Cityrecord.Location.Longitude),
				color.RGBA{0xff, 0, 0, 0xff},
				16.0*(ZOOM_OUT_FACTOR/300),
			)
			ctx.AddObject(rec)

		} else {
			Write_text := "IP: " + useip + "\nCountry: " + Cityrecord.Country.Names["en"]
			if len(Cityrecord.Subdivisions) > 0 {
				Write_text += "\nRegion: " + Cityrecord.Subdivisions[0].Names["en"]
			}
			Write_text += "\nCity: " + Cityrecord.City.Names["en"] + "\nOrg: " + useorg + "\nASN: " + useasn
			if el.Exit_node {
				Write_text += "\n " + el.Service_type + " exit node"
			}

			rec := NewTextMarker(s2.LatLngFromDegrees(Cityrecord.Location.Latitude, Cityrecord.Location.Longitude), Write_text)
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
