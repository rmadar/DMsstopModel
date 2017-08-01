package main

import (
	"fmt"
	"log"
	"math"

	"go-hep.org/x/hep/hplot"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/plot/palette/brewer"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/vg"
	"github.com/gonum/plot/vg/draw"
)

// nominal model parameters
const (
	mV  = 1000.0
	mDM = 1.0   // mDM is the DarkMatter mediator mass in GeV/c^2
	mt  = 172.0 // mt is the top mass in GeV/c^2
)

func main() {
	gDM := floats.Span(make([]float64, 3), 0.1, 1.6)
	gSM := floats.Span(make([]float64, 3), 0.1, 1.6)
	gs := make([][]float64, 0, 6)
	for _, x := range gDM {
		for _, y := range gSM {
			gs = append(gs, []float64{x, y})
		}
	}

	// mediator masses
	m1 := floats.Span(make([]float64, 1000), mt/5, mt*10)
	m2 := floats.Span(make([]float64, 1000), 500, 5000)

	palette, err := brewer.GetPalette(brewer.TypeAny, "RdYlBu", 11)
	if err != nil {
		log.Fatal(err)
	}

	tp, err := hplot.NewTiledPlot(draw.Tiles{Cols: 2, Rows: 1})
	if err != nil {
		log.Fatal(err)
	}

	p := &tp.Plot(0, 0).Plot
	p.Title.Text = fmt.Sprintf("Relative mediator width; mDM = %0.f GeV", mDM)
	p.X.Scale = &plot.LogScale{}
	p.X.Tick.Marker = &plot.LogTicks{}
	p.Y.Scale = &plot.LogScale{}
	p.Y.Tick.Marker = &plot.LogTicks{}
	p.X.Label.Text = "mV [GeV]"
	p.Y.Label.Text = "GammaV/mV"
	p.Legend.Top = true
	p.Add(hplot.NewGrid())

	for i, g := range gs {
		gSM := g[0]
		gDM := g[1]
		label := fmt.Sprintf("(gSM,gDM) = (%1.1f,%1.1f)", gSM, gDM)
		data := make(plotter.XYs, len(m1))
		for i, x := range m1 {
			data[i].X = x
			data[i].Y = WidthTot(gSM, gDM, x, mDM) / x
		}
		line, err := hplot.NewLine(data)
		if err != nil {
			log.Fatal(err)
		}
		line.Color = palette.Colors()[i]
		p.Add(line)
		p.Legend.Add(label, line)
	}

	p = &tp.Plot(0, 1).Plot
	p.Title.Text = fmt.Sprintf("Relative mediator width; mDM = %0.f GeV", mDM)
	p.X.Scale = &plot.LogScale{}
	p.X.Tick.Marker = &plot.LogTicks{}
	p.X.Min = 400
	p.X.Max = 3e4
	p.Y.Scale = &plot.LogScale{}
	p.Y.Tick.Marker = &plot.LogTicks{}
	p.X.Label.Text = "mV [GeV]"
	p.Y.Label.Text = "GammaV/mV"
	p.Legend.Top = true
	p.Add(hplot.NewGrid())

	for i, g := range gs {
		gSM := g[0]
		gDM := g[1]
		label := fmt.Sprintf("(gSM,gDM) = (%1.1f,%1.1f)", gSM, gDM)
		data := make(plotter.XYs, len(m2))
		for i, x := range m2 {
			data[i].X = x
			data[i].Y = WidthTot(gSM, gDM, x, mDM) / x
		}
		line, err := hplot.NewLine(data)
		if err != nil {
			log.Fatal(err)
		}
		line.Color = palette.Colors()[i]
		p.Add(line)
		p.Legend.Add(label, line)
	}

	for _, fname := range []string{"out.png", "out.pdf"} {
		err = tp.Save(-1, 20*vg.Centimeter, fname)
		if err != nil {
			log.Fatalf("error saving %q: %v", fname, err)
		}
	}
}

// PhiInv returns the phase space for the invisible decay
func PhiInv(mDM, mV float64) float64 {
	r := mDM / mV
	r2 := r * r
	return mV / (12 * math.Pi) * math.Sqrt(1-4*r2) * (1 + 2*r2)
}

// PhiVis returns the phase space for the visible decay
func PhiVis(mt, mV float64) float64 {
	r := mt / mV
	r2 := r * r
	r4 := r2 * r2
	return mV / math.Pi * (1 - r2) * (1 - 0.5*r2 - 0.5*r4)
}

func WidthVis(gSM, mV float64) float64 {
	return gSM * gSM * PhiVis(mt, mV)
}

func WidthInv(gDM, mV, mDM float64) float64 {
	return gDM * gDM * PhiInv(mDM, mV)
}

func WidthTot(gSM, gDM, mV, mDM float64) float64 {
	return WidthVis(gSM, mV) + WidthInv(gDM, mV, mDM)
}
