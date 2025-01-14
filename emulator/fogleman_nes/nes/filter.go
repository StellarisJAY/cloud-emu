package nes

//Copyright (C) 2015 Michael Fogleman https://github.com/fogleman/nes
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.
import "math"

type Filter interface {
	Step(x float32) float32
}

// First order filters are defined by the following parameters.
// y[n] = B0*x[n] + B1*x[n-1] - A1*y[n-1]
type FirstOrderFilter struct {
	B0    float32
	B1    float32
	A1    float32
	prevX float32
	prevY float32
}

func (f *FirstOrderFilter) Step(x float32) float32 {
	y := f.B0*x + f.B1*f.prevX - f.A1*f.prevY
	f.prevY = y
	f.prevX = x
	return y
}

// sampleRate: samples per second
// cutoffFreq: oscillations per second
func LowPassFilter(sampleRate float32, cutoffFreq float32) Filter {
	c := sampleRate / math.Pi / cutoffFreq
	a0i := 1 / (1 + c)
	return &FirstOrderFilter{
		B0: a0i,
		B1: a0i,
		A1: (1 - c) * a0i,
	}
}

func HighPassFilter(sampleRate float32, cutoffFreq float32) Filter {
	c := sampleRate / math.Pi / cutoffFreq
	a0i := 1 / (1 + c)
	return &FirstOrderFilter{
		B0: c * a0i,
		B1: -c * a0i,
		A1: (1 - c) * a0i,
	}
}

type FilterChain []Filter

func (fc FilterChain) Step(x float32) float32 {
	if fc != nil {
		for i := range fc {
			x = fc[i].Step(x)
		}
	}
	return x
}
