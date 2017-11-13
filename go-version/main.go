package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"

	"go-hep.org/x/hep/hbook"
	"go-hep.org/x/hep/hplot"
)

// nominal model parameters
const (
	mV  = 1000.0
	mDM = 1.0   // mDM is the DarkMatter mediator mass in GeV/c^2
	mt  = 172.0 // mt is the top mass in GeV/c^2
)

func main() {
	totwidth()
	smcouplings()
	brcouplings()
	gammabrcouplings()
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

func BR(gSM, gDM, mV, mDM float64) float64 {
	return WidthInv(gDM, mV, mDM) / WidthTot(gSM, gDM, mV, mDM)
}

func gSMFromBRWidth(br, width, mV, mDM float64) float64 {
	gSM2 := width / PhiVis(mt, mV) * (1 - br)
	return math.Sqrt(gSM2)
}

func gDMFromBRWidth(br, width, mV, mDM float64) float64 {
	gDM2 := width / PhiInv(mt, mV) * br
	return math.Sqrt(gDM2)
}

func totwidth() {
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

	tp, err := hplot.NewTiledPlot(draw.Tiles{Cols: 2, Rows: 1})
	if err != nil {
		log.Fatal(err)
	}

	{
		p := tp.Plot(0, 0)
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
			line.Color = colors(i)
			line.Dashes = dashes(i)
			p.Add(line)
			p.Legend.Add(label, line)
		}
	}

	{
		p := tp.Plot(0, 1)
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
			line.Color = colors(i)
			line.Dashes = dashes(i)
			p.Add(line)
			p.Legend.Add(label, line)
		}
	}

	for _, ext := range []string{".png", ".pdf"} {
		fname := "go-total-width" + ext
		err = tp.Save(-1, 20*vg.Centimeter, fname)
		if err != nil {
			log.Fatalf("error saving %q: %v", fname, err)
		}
	}

}

func smcouplings() {
	gDM := []float64{0.5, 1, 1.5}
	g := floats.Span(make([]float64, 150), 0, 1.5)

	tp, err := hplot.NewTiledPlot(draw.Tiles{Rows: 1, Cols: 2})
	if err != nil {
		log.Fatal(err)
	}

	{
		p := tp.Plot(0, 0)
		p.Title.Text = fmt.Sprintf("m_V = %.0f TeV ; m_DM = %.0f GeV", mV/1000.0, mDM)
		p.X.Label.Text = "g_SM"
		p.Y.Label.Text = "Gamma_V [GeV]"
		p.Legend.Top = true
		p.Legend.Left = true
		for j, gdm := range gDM {
			xys := make(plotter.XYs, len(g))
			for i := range xys {
				xys[i].X = g[i]
				xys[i].Y = WidthTot(g[i], gdm, mV, mDM)
			}
			line, err := hplot.NewLine(xys)
			if err != nil {
				log.Fatal(err)
			}
			line.Color = colors(j)
			p.Add(line)
			p.Legend.Add(fmt.Sprintf("g_DM = %1.1f", gdm), line)
		}
		p.Add(hplot.NewGrid())
	}

	{
		p := tp.Plot(0, 1)
		p.Title.Text = fmt.Sprintf("m_V = %.0f TeV ; m_DM = %.0f GeV", mV/1000.0, mDM)
		p.X.Label.Text = "g_SM"
		p.Y.Label.Text = "BR_chi"
		p.Legend.Top = true
		for j, gdm := range gDM {
			xys := make(plotter.XYs, len(g))
			for i := range xys {
				xys[i].X = g[i]
				xys[i].Y = BR(g[i], gdm, mV, mDM)
			}
			line, err := hplot.NewLine(xys)
			if err != nil {
				log.Fatal(err)
			}
			line.Color = colors(j)
			p.Add(line)
			p.Legend.Add(fmt.Sprintf("g_DM = %1.1f", gdm), line)
		}
		p.Add(hplot.NewGrid())
	}

	for _, ext := range []string{".png", ".pdf"} {
		fname := "go-sm-couplings-dep" + ext
		err = tp.Save(-1, 20*vg.Centimeter, fname)
		if err != nil {
			log.Fatalf("error saving %q: %v", fname, err)
		}
	}

}

func brcouplings() {
	gSM := floats.Span(make([]float64, 150), 0, 1.5)
	gDM := floats.Span(make([]float64, 150), 0, 1.5)
	br := hbook.NewH2D(150, 0, 1.5, 150, 0, 1.5)
	gam := hbook.NewH2D(150, 0, 1.5, 150, 0, 1.5)

	for _, x := range gSM {
		for _, y := range gDM {
			br.Fill(x, y, BR(x, y, mV, mDM))
			gam.Fill(x, y, WidthTot(x, y, mV, mDM))
		}
	}

	tp, err := hplot.NewTiledPlot(draw.Tiles{Rows: 1, Cols: 2})
	if err != nil {
		log.Fatal(err)
	}

	title := fmt.Sprintf("mV = %0.f TeV ; m_DM = %0.f GeV", mV/1000.0, mDM)
	{
		p := tp.Plot(0, 0)
		p.Title.Text = "Gamma_V [GeV] ; " + title
		p.X.Label.Text = "g_SM"
		p.Y.Label.Text = "g_DM"
		p.Add(hplot.NewH2D(gam, nil))
	}
	{
		p := tp.Plot(0, 1)
		p.Title.Text = "BR_chi ; " + title
		p.X.Label.Text = "g_SM"
		p.Y.Label.Text = "g_DM"
		p.Add(hplot.NewH2D(br, nil))
	}

	for _, ext := range []string{".png", ".pdf"} {
		fname := "go-br-gamma-couplings" + ext
		err = tp.Save(-1, 20*vg.Centimeter, fname)
		if err != nil {
			log.Fatalf("error saving %q: %v", fname, err)
		}
	}
}

func gammabrcouplings() {
	gv := floats.Span(make([]float64, 150), 0, 500)
	br := floats.Span(make([]float64, 150), 0, 1)
	gsm := hbook.NewH2D(150, 0, 500, 150, 0, 1)
	gdm := hbook.NewH2D(150, 0, 500, 150, 0, 1)

	for _, x := range gv {
		for _, y := range br {
			gsm.Fill(x, y, gSMFromBRWidth(y, x, mV, mDM))
			gdm.Fill(x, y, gDMFromBRWidth(y, x, mV, mDM))
		}
	}

	tp, err := hplot.NewTiledPlot(draw.Tiles{Rows: 1, Cols: 2})
	if err != nil {
		log.Fatal(err)
	}

	title := fmt.Sprintf("mV = %0.f TeV ; m_DM = %0.f GeV", mV/1000.0, mDM)
	{
		p := tp.Plot(0, 0)
		p.Title.Text = "g_SM ; " + title
		p.X.Label.Text = "Gamma_V [GeV]"
		p.Y.Label.Text = "BR_chi"
		p.Add(hplot.NewH2D(gsm, nil))
	}
	{
		p := tp.Plot(0, 1)
		p.Title.Text = "g_DM ; " + title
		p.X.Label.Text = "Gamma_V [GeV]"
		p.Y.Label.Text = "BR_chi"
		p.Add(hplot.NewH2D(gdm, nil))
	}

	for _, ext := range []string{".png", ".pdf"} {
		fname := "go-couplings-vs-gammabr" + ext
		err = tp.Save(-1, 20*vg.Centimeter, fname)
		if err != nil {
			log.Fatalf("error saving %q: %v", fname, err)
		}
	}
}

func colors(i int) color.Color {
	return plotutil.Color(i)
}

func dashes(i int) []vg.Length {
	if i >= len(plotutil.SoftColors) {
		return plotutil.Dashes(i)
	}
	return plotutil.Dashes(0)
}
