/**
* @license
* Copyright 2020 Dynatrace LLC
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
 */

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dtcookie/dynatrace/api/config/dashboards"
	"github.com/dtcookie/dynatrace/api/config/maintenancewindows"
	"github.com/dtcookie/dynatrace/api/config/managementzones"
	"github.com/dtcookie/dynatrace/api/config/requestattributes"
	"github.com/dtcookie/dynatrace/terraform"
)

// func rnd(enums interface{}) string {
// 	rv := reflect.ValueOf(enums)
// 	return fmt.Sprintf("%v", rv.Field(rand.Intn(rv.NumField())).Interface())
// }

// func rndAggregationRate() *dashboards.AggregationRate {
// 	r := dashboards.AggregationRate(rnd(dashboards.AggregationRates))
// 	return &r
// }

// func rndLeftAxisCustomUnit() *dashboards.LeftAxisCustomUnit {
// 	r := dashboards.LeftAxisCustomUnit(rnd(dashboards.LeftAxisCustomUnits))
// 	return &r
// }

// func rndRightAxisCustomUnit() *dashboards.RightAxisCustomUnit {
// 	r := dashboards.RightAxisCustomUnit(rnd(dashboards.RightAxisCustomUnits))
// 	return &r
// }
func usage(args []string) bool {
	fmt.Printf("Usage: %s convert <inputfile.json> [-o <outputfile>] [-f <tf|json>\n", args[0])
	return true
}

func convert(args []string) bool {
	var err error
	var bytes []byte
	var outFile *os.File
	outFile = nil
	format := "tf"
	resType := ""

	if len(args) == 1 {
		return false
	}
	if strings.TrimSpace(args[1]) != "convert" {
		return false
	}
	if len(args) < 3 {
		return usage(args)
	}
	inFile := args[2]
	expectOutFile := false
	expectFormat := false
	expectType := false

	for idx := 3; idx < len(args); idx++ {
		if expectOutFile {
			if outFile, err = os.Create(strings.TrimSpace(args[idx])); err != nil {
				fmt.Println(err.Error())
				return true
			}
			defer outFile.Close()
			expectOutFile = false
		} else if expectFormat {
			format = strings.TrimSpace(args[idx])
			if (format != "tf") && (format != "json") {
				fmt.Printf("unsupported format '%v' (either 'tf' or 'json' are supported)\n", format)
				return true
			}
			expectFormat = false
		} else if expectType {
			resType = strings.TrimSpace(args[idx])
			if (resType != "dashboard") && (resType != "managementzone") && (resType != "maintenancewindow") && (resType != "requestattribute") {
				fmt.Printf("unsupported resource type '%v' (either 'dashboard' or 'managementzone' or 'maintenancewindow' or 'requestattribute' are supported)\n", resType)
				return true
			}
			expectType = false
		} else if strings.TrimSpace(args[idx]) == "-o" {
			expectOutFile = true
		} else if strings.TrimSpace(args[idx]) == "-f" {
			expectFormat = true
		} else if strings.TrimSpace(args[idx]) == "-t" {
			expectType = true
		}
	}
	if (resType == "") || expectOutFile || expectFormat || expectType {
		return usage(args)
	}

	if bytes, err = ioutil.ReadFile(inFile); err != nil {
		fmt.Println(err.Error())
		return true
	}
	if resType == "dashboard" {
		dashboard := new(dashboards.Dashboard)
		if e2 := json.Unmarshal(bytes, dashboard); e2 != nil {
			panic(e2)
		}

		if format == "tf" {
			if bytes, err = terraform.Marshal(dashboard, "dynatrace_dashboard", dashboard.Metadata.Name); err != nil {
				panic(err)
			}
		} else {
			if bytes, err = terraform.MarshalJSON(dashboard, "dynatrace_dashboard", dashboard.Metadata.Name); err != nil {
				panic(err)
			}
		}
	} else if resType == "managementzone" {
		managementzone := new(managementzones.ManagementZone)
		if e2 := json.Unmarshal(bytes, managementzone); e2 != nil {
			panic(e2)
		}

		if format == "tf" {
			if bytes, err = terraform.Marshal(managementzone, "dynatrace_management_zone", managementzone.Name); err != nil {
				panic(err)
			}
		} else {
			if bytes, err = terraform.MarshalJSON(managementzone, "dynatrace_management_zone", managementzone.Name); err != nil {
				panic(err)
			}
		}
	} else if resType == "maintenancewindow" {
		maintenancewindow := new(maintenancewindows.MaintenanceWindow)
		if e2 := json.Unmarshal(bytes, maintenancewindow); e2 != nil {
			panic(e2)
		}

		if format == "tf" {
			if bytes, err = terraform.Marshal(maintenancewindow, "dynatrace_maintenance_window", maintenancewindow.Name); err != nil {
				panic(err)
			}
		} else {
			if bytes, err = terraform.MarshalJSON(maintenancewindow, "dynatrace_maintenance_window", maintenancewindow.Name); err != nil {
				panic(err)
			}
		}
	} else if resType == "requestattribute" {
		requestattribute := new(requestattributes.RequestAttribute)
		if e2 := json.Unmarshal(bytes, requestattribute); e2 != nil {
			panic(e2)
		}

		if format == "tf" {
			if bytes, err = terraform.Marshal(requestattribute, "dynatrace_request_attribute", requestattribute.Name); err != nil {
				panic(err)
			}
		} else {
			if bytes, err = terraform.MarshalJSON(requestattribute, "dynatrace_request_attribute", requestattribute.Name); err != nil {
				panic(err)
			}
		}
	}

	if outFile == nil {
		fmt.Println(string(bytes))
	} else {
		fmt.Fprintln(outFile, string(bytes))
	}

	return true
}

// func rndStr() *string {
// 	s := uuid.New().String()
// 	return &s
// }

// func rndStr() *string {
// 	s := uuid.New().String()
// 	return &s
// }

// dashboards.AggregationRate(rnd(dashboards.AggregationRates))

// func doTestOutput() {
// 	dashboard := &dashboards.Dashboard{
// 		ConfigurationMetadata: &api.ConfigurationMetadata{
// 			CurrentConfigurationVersions: []string{uuid.New().String()},
// 			ClusterVersion:               uuid.New().String(),
// 			ConfigurationVersions:        []int64{rand.Int63n(100)},
// 		},
// 		ID: rndStr(),
// 		Metadata: &dashboards.DashboardMetadata{
// 			Name:   uuid.New().String(),
// 			Shared: opt.NewBool(true),
// 			Owner:  rndStr(),
// 			SharingDetails: &dashboards.SharingInfo{
// 				LinkShared: opt.NewBool(true),
// 				Published:  opt.NewBool(true),
// 			},
// 			Filter: &dashboards.DashboardFilter{
// 				Timeframe: rndStr(),
// 				ManagementZone: &api.EntityShortRepresentation{
// 					ID:          uuid.New().String(),
// 					Name:        uuid.New().String(),
// 					Description: uuid.New().String(),
// 				},
// 			},
// 			Tags:            []string{uuid.New().String()},
// 			Preset:          opt.NewBool(true),
// 			ValidFilterKeys: []string{uuid.New().String()},
// 		},
// 		Tiles: []dashboards.Tile{
// 			&dashboards.CustomChartingTile{
// 				AbstractTile: dashboards.AbstractTile{
// 					Name:       uuid.New().String(),
// 					TileType:   "CUSTOM_CHARTING",
// 					Configured: opt.NewBool(true),
// 					Bounds: &dashboards.TileBounds{
// 						Top:    rand.Int31n(100),
// 						Left:   rand.Int31n(100),
// 						Width:  rand.Int31n(100),
// 						Height: rand.Int31n(100),
// 					},
// 					Filter: &dashboards.TileFilter{
// 						Timeframe: rndStr(),
// 						ManagementZone: &api.EntityShortRepresentation{
// 							ID:          uuid.New().String(),
// 							Name:        uuid.New().String(),
// 							Description: uuid.New().String(),
// 						},
// 					},
// 				},
// 				FilterConfig: &dashboards.CustomFilterConfig{
// 					Type:        dashboards.CustomFilterConfigType(rnd(dashboards.CustomFilterConfigTypes)),
// 					CustomName:  uuid.New().String(),
// 					DefaultName: uuid.New().String(),
// 					ChartConfig: &dashboards.CustomFilterChartConfig{
// 						LegendShown: opt.NewBool(true),
// 						Type:        dashboards.CustomFilterChartConfigType(rnd(dashboards.CustomFilterChartConfigTypes)),
// 						Series: []dashboards.CustomFilterChartSeriesConfig{
// 							{
// 								Metric:      uuid.New().String(),
// 								Aggregation: dashboards.Aggregation(rnd(dashboards.Aggregations)),
// 								Percentile:  opt.NewInt64(rand.Int63n(100)),
// 								Type:        dashboards.CustomFilterChartSeriesConfigType(rnd(dashboards.CustomFilterChartSeriesConfigTypes)),
// 								EntityType:  uuid.New().String(),
// 								Dimensions: []dashboards.CustomFilterChartSeriesDimensionConfig{
// 									{
// 										ID:              uuid.New().String(),
// 										Name:            rndStr(),
// 										Values:          []string{uuid.New().String()},
// 										EntityDimension: opt.NewBool(true),
// 									},
// 								},
// 								SortAscending:   true,
// 								SortColumn:      true,
// 								AggregationRate: rndAggregationRate(),
// 							},
// 						},
// 						ResultMetadata: map[string]dashboards.CustomChartingItemMetadataConfig{
// 							uuid.New().String(): {
// 								LastModified: opt.NewInt64(rand.Int63n(100)),
// 								CustomColor:  uuid.New().String(),
// 							},
// 						},
// 						AxisLimits: map[string]float64{
// 							uuid.New().String(): rand.Float64(),
// 						},
// 						LeftAxisCustomUnit:  rndLeftAxisCustomUnit(),
// 						RightAxisCustomUnit: rndRightAxisCustomUnit(),
// 					},
// 					FiltersPerEntityType: map[string]map[string][]string{
// 						uuid.New().String(): {
// 							uuid.New().String(): []string{uuid.New().String()},
// 						},
// 					},
// 				},
// 			},
// 			&dashboards.MarkdownTile{
// 				AbstractTile: dashboards.AbstractTile{
// 					Name:       uuid.New().String(),
// 					TileType:   "MARKDOWN",
// 					Configured: opt.NewBool(true),
// 					Bounds: &dashboards.TileBounds{
// 						Top:    rand.Int31n(100),
// 						Left:   rand.Int31n(100),
// 						Width:  rand.Int31n(100),
// 						Height: rand.Int31n(100),
// 					},
// 					Filter: &dashboards.TileFilter{
// 						Timeframe: rndStr(),
// 						ManagementZone: &api.EntityShortRepresentation{
// 							ID:          uuid.New().String(),
// 							Name:        uuid.New().String(),
// 							Description: uuid.New().String(),
// 						},
// 					},
// 				},
// 				Markdown: uuid.New().String(),
// 			},
// 			&dashboards.UserSessionQueryTile{
// 				AbstractTile: dashboards.AbstractTile{
// 					Name:       uuid.New().String(),
// 					TileType:   "DTAQL",
// 					Configured: opt.NewBool(true),
// 					Bounds: &dashboards.TileBounds{
// 						Top:    rand.Int31n(100),
// 						Left:   rand.Int31n(100),
// 						Width:  rand.Int31n(100),
// 						Height: rand.Int31n(100),
// 					},
// 					Filter: &dashboards.TileFilter{
// 						Timeframe: rndStr(),
// 						ManagementZone: &api.EntityShortRepresentation{
// 							ID:          uuid.New().String(),
// 							Name:        uuid.New().String(),
// 							Description: uuid.New().String(),
// 						},
// 					},
// 				},
// 				CustomName:     uuid.New().String(),
// 				Query:          uuid.New().String(),
// 				Type:           dashboards.UserSessionQueryTileType(rnd(dashboards.UserSessionQueryTileTypes)),
// 				TimeFrameShift: rndStr(),
// 				VisualizationConfig: &dashboards.UserSessionQueryTileConfiguration{
// 					HasAxisBucketing: opt.NewBool(true),
// 				},
// 				Limit: opt.NewInt32(rand.Int31n(100)),
// 			},
// 			&dashboards.FilterableEntityTile{
// 				AbstractTile: dashboards.AbstractTile{
// 					Name:       uuid.New().String(),
// 					TileType:   "APPLICATIONS",
// 					Configured: opt.NewBool(true),
// 					Bounds: &dashboards.TileBounds{
// 						Top:    rand.Int31n(100),
// 						Left:   rand.Int31n(100),
// 						Width:  rand.Int31n(100),
// 						Height: rand.Int31n(100),
// 					},
// 					Filter: &dashboards.TileFilter{
// 						Timeframe: rndStr(),
// 						ManagementZone: &api.EntityShortRepresentation{
// 							ID:          uuid.New().String(),
// 							Name:        uuid.New().String(),
// 							Description: uuid.New().String(),
// 						},
// 					},
// 				},
// 				FilterConfig: &dashboards.CustomFilterConfig{
// 					Type:        dashboards.CustomFilterConfigType(rnd(dashboards.CustomFilterConfigTypes)),
// 					CustomName:  uuid.New().String(),
// 					DefaultName: uuid.New().String(),
// 					ChartConfig: &dashboards.CustomFilterChartConfig{
// 						LegendShown: opt.NewBool(true),
// 						Type:        dashboards.CustomFilterChartConfigType(rnd(dashboards.CustomFilterChartConfigTypes)),
// 						Series: []dashboards.CustomFilterChartSeriesConfig{
// 							{
// 								Metric:      uuid.New().String(),
// 								Aggregation: dashboards.Aggregation(rnd(dashboards.Aggregations)),
// 								Percentile:  opt.NewInt64(rand.Int63n(100)),
// 								Type:        dashboards.CustomFilterChartSeriesConfigType(rnd(dashboards.CustomFilterChartSeriesConfigTypes)),
// 								EntityType:  uuid.New().String(),
// 								Dimensions: []dashboards.CustomFilterChartSeriesDimensionConfig{
// 									{
// 										ID:              uuid.New().String(),
// 										Name:            rndStr(),
// 										Values:          []string{uuid.New().String()},
// 										EntityDimension: opt.NewBool(true),
// 									},
// 								},
// 								SortAscending:   true,
// 								SortColumn:      true,
// 								AggregationRate: rndAggregationRate(),
// 							},
// 						},
// 						ResultMetadata: map[string]dashboards.CustomChartingItemMetadataConfig{
// 							uuid.New().String(): {
// 								LastModified: opt.NewInt64(rand.Int63()),
// 								CustomColor:  uuid.New().String(),
// 							},
// 						},
// 						AxisLimits: map[string]float64{
// 							uuid.New().String(): rand.Float64(),
// 						},
// 						LeftAxisCustomUnit:  rndLeftAxisCustomUnit(),
// 						RightAxisCustomUnit: rndRightAxisCustomUnit(),
// 					},
// 					FiltersPerEntityType: map[string]map[string][]string{
// 						uuid.New().String(): {
// 							uuid.New().String(): []string{uuid.New().String()},
// 						},
// 					},
// 				},
// 				ChartVisible: opt.NewBool(true),
// 			},
// 			&dashboards.AssignedEntitiesWithMetricTile{
// 				AbstractTile: dashboards.AbstractTile{
// 					Name:       uuid.New().String(),
// 					TileType:   "APPLICATION_WORLDMAP",
// 					Configured: opt.NewBool(true),
// 					Bounds: &dashboards.TileBounds{
// 						Top:    rand.Int31n(100),
// 						Left:   rand.Int31n(100),
// 						Width:  rand.Int31n(100),
// 						Height: rand.Int31n(100),
// 					},
// 					Filter: &dashboards.TileFilter{
// 						Timeframe: rndStr(),
// 						ManagementZone: &api.EntityShortRepresentation{
// 							ID:          uuid.New().String(),
// 							Name:        uuid.New().String(),
// 							Description: uuid.New().String(),
// 						},
// 					},
// 				},
// 				AssignedEntities: []string{uuid.New().String()},
// 				Metric:           rndStr(),
// 			},
// 			&dashboards.AssignedEntitiesTile{
// 				AbstractTile: dashboards.AbstractTile{
// 					Name:       uuid.New().String(),
// 					TileType:   "USERS",
// 					Configured: opt.NewBool(true),
// 					Bounds: &dashboards.TileBounds{
// 						Top:    rand.Int31n(100),
// 						Left:   rand.Int31n(100),
// 						Width:  rand.Int31n(100),
// 						Height: rand.Int31n(100),
// 					},
// 					Filter: &dashboards.TileFilter{
// 						Timeframe: rndStr(),
// 						ManagementZone: &api.EntityShortRepresentation{
// 							ID:          uuid.New().String(),
// 							Name:        uuid.New().String(),
// 							Description: uuid.New().String(),
// 						},
// 					},
// 				},
// 				AssignedEntities: []string{uuid.New().String()},
// 			},
// 			&dashboards.SyntheticSingleWebCheckTile{
// 				AbstractTile: dashboards.AbstractTile{
// 					Name:       uuid.New().String(),
// 					TileType:   "SYNTHETIC_SINGLE_WEBCHECK",
// 					Configured: opt.NewBool(true),
// 					Bounds: &dashboards.TileBounds{
// 						Top:    rand.Int31n(100),
// 						Left:   rand.Int31n(100),
// 						Width:  rand.Int31n(100),
// 						Height: rand.Int31n(100),
// 					},
// 					Filter: &dashboards.TileFilter{
// 						Timeframe: rndStr(),
// 						ManagementZone: &api.EntityShortRepresentation{
// 							ID:          uuid.New().String(),
// 							Name:        uuid.New().String(),
// 							Description: uuid.New().String(),
// 						},
// 					},
// 				},
// 				AssignedEntities:          []string{uuid.New().String()},
// 				ExcludeMaintenanceWindows: opt.NewBool(true),
// 			},
// 		},
// 	}

// 	if bytes, err := terraform.MarshalJSON(dashboard, "dynatrace_dashboards", dashboard.Metadata.Name); err != nil {
// 		panic(err)
// 	} else {
// 		fmt.Println(string(bytes))
// 	}
// }
