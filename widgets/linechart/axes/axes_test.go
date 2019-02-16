// Copyright 2019 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package axes

import (
	"image"
	"testing"

	"github.com/kylelemons/godebug/pretty"
)

type updateY struct {
	minVal float64
	maxVal float64
}

func TestY(t *testing.T) {
	tests := []struct {
		desc      string
		yp        *YProperties
		cvsAr     image.Rectangle
		wantWidth int
		want      *YDetails
		wantErr   bool
	}{
		{
			desc: "fails on canvas too small",
			yp: &YProperties{
				Min:        0,
				Max:        3,
				ReqXHeight: 2,
			},
			cvsAr:     image.Rect(0, 0, 3, 2),
			wantWidth: 2,
			wantErr:   true,
		},
		{
			desc: "fails on cvsWidth less than required width",
			yp: &YProperties{
				Min:        0,
				Max:        3,
				ReqXHeight: 2,
			},
			cvsAr:     image.Rect(0, 0, 2, 4),
			wantWidth: 2,
			wantErr:   true,
		},
		{
			desc: "fails when max is less than min",
			yp: &YProperties{
				Min:        0,
				Max:        -1,
				ReqXHeight: 2,
			},
			cvsAr:     image.Rect(0, 0, 4, 4),
			wantWidth: 3,
			wantErr:   true,
		},
		{
			desc: "cvsWidth equals required width",
			yp: &YProperties{
				Min:        0,
				Max:        3,
				ReqXHeight: 2,
			},
			cvsAr:     image.Rect(0, 0, 3, 4),
			wantWidth: 2,
			want: &YDetails{
				Width: 2,
				Start: image.Point{1, 0},
				End:   image.Point{1, 2},
				Scale: mustNewYScale(0, 3, 2, nonZeroDecimals, YScaleModeAnchored),
				Labels: []*Label{
					{NewValue(0, nonZeroDecimals), image.Point{0, 1}},
					{NewValue(1.72, nonZeroDecimals), image.Point{0, 0}},
				},
			},
		},
		{
			desc: "success for anchored scale",
			yp: &YProperties{
				Min:        1,
				Max:        3,
				ReqXHeight: 2,
				ScaleMode:  YScaleModeAnchored,
			},
			cvsAr:     image.Rect(0, 0, 3, 4),
			wantWidth: 2,
			want: &YDetails{
				Width: 2,
				Start: image.Point{1, 0},
				End:   image.Point{1, 2},
				Scale: mustNewYScale(0, 3, 2, nonZeroDecimals, YScaleModeAnchored),
				Labels: []*Label{
					{NewValue(0, nonZeroDecimals), image.Point{0, 1}},
					{NewValue(1.72, nonZeroDecimals), image.Point{0, 0}},
				},
			},
		},
		{
			desc: "accommodates X scale that needs more height",
			yp: &YProperties{
				Min:        1,
				Max:        3,
				ReqXHeight: 4,
				ScaleMode:  YScaleModeAnchored,
			},
			cvsAr:     image.Rect(0, 0, 3, 6),
			wantWidth: 2,
			want: &YDetails{
				Width: 2,
				Start: image.Point{1, 0},
				End:   image.Point{1, 2},
				Scale: mustNewYScale(0, 3, 2, nonZeroDecimals, YScaleModeAnchored),
				Labels: []*Label{
					{NewValue(0, nonZeroDecimals), image.Point{0, 1}},
					{NewValue(1.72, nonZeroDecimals), image.Point{0, 0}},
				},
			},
		},
		{
			desc: "success for adaptive scale",
			yp: &YProperties{
				Min:        1,
				Max:        6,
				ReqXHeight: 2,
				ScaleMode:  YScaleModeAdaptive,
			},
			cvsAr:     image.Rect(0, 0, 3, 4),
			wantWidth: 2,
			want: &YDetails{
				Width: 2,
				Start: image.Point{1, 0},
				End:   image.Point{1, 2},
				Scale: mustNewYScale(1, 6, 2, nonZeroDecimals, YScaleModeAdaptive),
				Labels: []*Label{
					{NewValue(1, nonZeroDecimals), image.Point{0, 1}},
					{NewValue(3.88, nonZeroDecimals), image.Point{0, 0}},
				},
			},
		},
		{
			desc: "cvsWidth just accommodates the longest label",
			yp: &YProperties{
				Min:        0,
				Max:        3,
				ReqXHeight: 2,
			},
			cvsAr:     image.Rect(0, 0, 6, 4),
			wantWidth: 2,
			want: &YDetails{
				Width: 5,
				Start: image.Point{4, 0},
				End:   image.Point{4, 2},
				Scale: mustNewYScale(0, 3, 2, nonZeroDecimals, YScaleModeAnchored),
				Labels: []*Label{
					{NewValue(0, nonZeroDecimals), image.Point{3, 1}},
					{NewValue(1.72, nonZeroDecimals), image.Point{0, 0}},
				},
			},
		},
		{
			desc: "cvsWidth is more than we need",
			yp: &YProperties{
				Min:        0,
				Max:        3,
				ReqXHeight: 2,
			},
			cvsAr:     image.Rect(0, 0, 7, 4),
			wantWidth: 2,
			want: &YDetails{
				Width: 5,
				Start: image.Point{4, 0},
				End:   image.Point{4, 2},
				Scale: mustNewYScale(0, 3, 2, nonZeroDecimals, YScaleModeAnchored),
				Labels: []*Label{
					{NewValue(0, nonZeroDecimals), image.Point{3, 1}},
					{NewValue(1.72, nonZeroDecimals), image.Point{0, 0}},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			gotWidth := RequiredWidth(tc.yp.Min, tc.yp.Max)
			if gotWidth != tc.wantWidth {
				t.Errorf("RequiredWidth => got %v, want %v", gotWidth, tc.wantWidth)
			}

			got, err := NewYDetails(tc.cvsAr, tc.yp)
			if (err != nil) != tc.wantErr {
				t.Errorf("Details => unexpected error: %v, wantErr: %v", err, tc.wantErr)
			}
			if err != nil {
				return
			}
			if diff := pretty.Compare(tc.want, got); diff != "" {
				t.Errorf("Details => unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestNewXDetails(t *testing.T) {
	tests := []struct {
		desc             string
		numPoints        int
		yStart           image.Point
		cvsWidth         int
		cvsAr            image.Rectangle
		customLabels     map[int]string
		labelOrientation LabelOrientation
		want             *XDetails
		wantErr          bool
	}{
		{
			desc:      "fails when numPoints is negative",
			numPoints: -1,
			yStart:    image.Point{0, 0},
			cvsAr:     image.Rect(0, 0, 2, 3),
			wantErr:   true,
		},
		{
			desc:      "fails when cvsAr isn't wide enough",
			numPoints: 1,
			yStart:    image.Point{0, 0},
			cvsAr:     image.Rect(0, 0, 1, 3),
			wantErr:   true,
		},
		{
			desc:      "fails when cvsAr isn't tall enough",
			numPoints: 1,
			yStart:    image.Point{0, 0},
			cvsAr:     image.Rect(0, 0, 3, 2),
			wantErr:   true,
		},
		{
			desc:      "works with no data points",
			numPoints: 0,
			yStart:    image.Point{0, 0},
			cvsAr:     image.Rect(0, 0, 2, 3),
			want: &XDetails{
				Start: image.Point{0, 1},
				End:   image.Point{1, 1},
				Scale: mustNewXScale(0, 1, nonZeroDecimals),
				Labels: []*Label{
					{
						Value: NewValue(0, nonZeroDecimals),
						Pos:   image.Point{1, 2},
					},
				},
			},
		},
		{
			desc:             "works with no data points, vertical",
			numPoints:        0,
			yStart:           image.Point{0, 0},
			cvsAr:            image.Rect(0, 0, 2, 3),
			labelOrientation: LabelOrientationVertical,
			want: &XDetails{
				Start: image.Point{0, 1},
				End:   image.Point{1, 1},
				Scale: mustNewXScale(0, 1, nonZeroDecimals),
				Labels: []*Label{
					{
						Value: NewValue(0, nonZeroDecimals),
						Pos:   image.Point{1, 2},
					},
				},
			},
		},
		{
			desc:      "accounts for non-zero yStart",
			numPoints: 0,
			yStart:    image.Point{2, 0},
			cvsAr:     image.Rect(0, 0, 4, 5),
			want: &XDetails{
				Start: image.Point{2, 3},
				End:   image.Point{3, 3},
				Scale: mustNewXScale(0, 1, nonZeroDecimals),
				Labels: []*Label{
					{
						Value: NewValue(0, nonZeroDecimals),
						Pos:   image.Point{3, 4},
					},
				},
			},
		},
		{
			desc:             "accounts for longer vertical labels, the tallest didn't fit",
			numPoints:        1000,
			yStart:           image.Point{2, 0},
			cvsAr:            image.Rect(0, 0, 10, 10),
			labelOrientation: LabelOrientationVertical,
			want: &XDetails{
				Start: image.Point{2, 5},
				End:   image.Point{9, 5},
				Scale: mustNewXScale(1000, 7, nonZeroDecimals),
				Labels: []*Label{
					{
						Value: NewValue(0, nonZeroDecimals),
						Pos:   image.Point{3, 6},
					},
					{
						Value: NewValue(615, nonZeroDecimals),
						Pos:   image.Point{7, 6},
					},
				},
			},
		},
		{
			desc:             "accounts for longer vertical labels, the tallest label fits",
			numPoints:        999,
			yStart:           image.Point{2, 0},
			cvsAr:            image.Rect(0, 0, 10, 10),
			labelOrientation: LabelOrientationVertical,
			want: &XDetails{
				Start: image.Point{2, 6},
				End:   image.Point{9, 6},
				Scale: mustNewXScale(999, 7, nonZeroDecimals),
				Labels: []*Label{
					{
						Value: NewValue(0, nonZeroDecimals),
						Pos:   image.Point{3, 7},
					},
					{
						Value: NewValue(614, nonZeroDecimals),
						Pos:   image.Point{7, 7},
					},
				},
			},
		},
		{
			desc:      "accounts for longer custom labels, vertical",
			numPoints: 2,
			yStart:    image.Point{5, 0},
			cvsAr:     image.Rect(0, 0, 20, 10),
			customLabels: map[int]string{
				0: "start",
				1: "end",
			},
			labelOrientation: LabelOrientationVertical,
			want: &XDetails{
				Start: image.Point{5, 4},
				End:   image.Point{19, 4},
				Scale: mustNewXScale(2, 14, nonZeroDecimals),
				Labels: []*Label{
					{
						Value: NewTextValue("start"),
						Pos:   image.Point{6, 5},
					},
					{
						Value: NewTextValue("end"),
						Pos:   image.Point{19, 5},
					},
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got, err := NewXDetails(tc.numPoints, tc.yStart, tc.cvsAr, tc.customLabels, tc.labelOrientation)
			if (err != nil) != tc.wantErr {
				t.Errorf("NewXDetails => unexpected error: %v, wantErr: %v", err, tc.wantErr)
			}
			if err != nil {
				return
			}

			if diff := pretty.Compare(tc.want, got); diff != "" {
				t.Errorf("NewXDetails => unexpected diff (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestRequiredHeight(t *testing.T) {
	tests := []struct {
		desc             string
		numPoints        int
		customLabels     map[int]string
		labelOrientation LabelOrientation
		want             int
	}{
		{
			desc: "horizontal orientation",
			want: 2,
		},
		{
			desc:             "vertical orientation, no custom labels, need single row for max label",
			numPoints:        9,
			labelOrientation: LabelOrientationVertical,
			want:             2,
		},
		{
			desc:             "vertical orientation, no custom labels, need multiple row for max label",
			numPoints:        100,
			labelOrientation: LabelOrientationVertical,
			want:             4,
		},
		{
			desc:             "vertical orientation, custom labels but all shorter than max label",
			numPoints:        100,
			customLabels:     map[int]string{1: "a", 2: "b"},
			labelOrientation: LabelOrientationVertical,
			want:             4,
		},
		{
			desc:             "vertical orientation, custom labels and some longer than max label",
			numPoints:        100,
			customLabels:     map[int]string{1: "a", 2: "bbbbb"},
			labelOrientation: LabelOrientationVertical,
			want:             6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			got := RequiredHeight(tc.numPoints, tc.customLabels, tc.labelOrientation)
			if got != tc.want {
				t.Errorf("RequiredHeight => %d, want %d", got, tc.want)
			}
		})
	}
}
