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

package dql

import (
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/opt"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/dynatrace/rest"
	"github.com/dynatrace-oss/terraform-provider-dynatrace/provider/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ENV_VAR_POLL_SLEEP_DURATION = "DYNATRACE_DQL_POLL_SLEEP_DURATION"
const DEFAULT_POLL_SLEEP_DURATION = 5000
const MIN_POLL_SLEEP_DURATION = 0
const MAX_POLL_SLEEP_DURATION = 60000

var POLL_SLEEP_DURATION = evalPollSleepDuration()

func evalPollSleepDuration() int {
	value := os.Getenv(ENV_VAR_POLL_SLEEP_DURATION)
	if len(value) == 0 {
		return DEFAULT_POLL_SLEEP_DURATION
	}
	iValue, err := strconv.Atoi(value)
	if err != nil {
		return DEFAULT_POLL_SLEEP_DURATION
	}
	if iValue < 0 {
		return DEFAULT_POLL_SLEEP_DURATION
	}
	if iValue > MAX_POLL_SLEEP_DURATION {
		return DEFAULT_POLL_SLEEP_DURATION
	}
	return iValue
}

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceRead,
		Schema: map[string]*schema.Schema{
			"query": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "example: fetch events | filter event.type == \"davis\" AND davis.status != \"CLOSED\" | fields timestamp, davis.title, davis.underMaintenance, davis.status | sort timestamp | limit 10",
			},
			"default_timeframe_start": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The query timeframe 'start' timestamp in ISO-8601 or RFC3339 format. If the timeframe 'end' parameter is missing, the whole timeframe is ignored. Note that if a timeframe is specified within the query string (query) then it has precedence over this query request parameter",
			},
			"default_timeframe_end": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The query timeframe 'end' timestamp in ISO-8601 or RFC3339 format. If the timeframe 'start' parameter is missing, the whole timeframe is ignored. Note that if a timeframe is specified within the query string (query) then it has precedence over this query request parameter",
			},
			"timezone": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "UTC",
				Description: "The query timezone. If none is specified, UTC is used as fallback. The list of valid input values matches that of the IANA Time Zone Database (TZDB). It accepts values in their canonical names like 'Europe/Paris', the abbreviated version like CET or the UTC offset format like '+01:00'",
			},
			"locale": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The query locale. If none specified, then a language/country neutral locale is chosen. The input values take the ISO-639 Language code with an optional ISO-3166 country code appended to it with an underscore. For instance, both values are valid 'en' or 'en_US'",
			},
			"max_result_records": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The maximum number of result records that this query will return",
			},
			"max_result_bytes": {
				Type:        schema.TypeInt,
				Description: "The maximum number of result bytes that this query will return",
				Optional:    true,
			},
			"fetch_timeout_seconds": {
				Type:        schema.TypeInt,
				Description: "The query will stop reading data after reaching the fetch-timeout. The query execution will continue, providing a partial result based on the read data",
				Optional:    true,
			},
			"default_sampling_ratio": {
				Type:        schema.TypeFloat,
				Description: "In case not specified in the DQL string, the sampling ratio defined here is applied. Note that this is only applicable to log queries",
				Optional:    true,
			},
			"default_scan_limit_gbytes": {
				Type:        schema.TypeInt,
				Description: "Limit in gigabytes for the amount data that will be scanned during read",
				Optional:    true,
			},
			"records": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func processHeredocString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

func DataSourceRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	creds, err := config.Credentials(m, config.CredValDefault)
	if err != nil {
		return diag.FromErr(err)
	}

	restClient, err := rest.CreatePlatformClient(ctx, creds.OAuth.EnvironmentURL, creds)
	if err != nil {
		return diag.FromErr(err)
	}

	client := NewClient(restClient)

	dqlRequest := DQLRequest{}
	if v, ok := d.GetOk("query"); ok {
		dqlRequest.Query = processHeredocString(v.(string))
	}
	if v, ok := d.GetOk("default_timeframe_start"); ok && len(v.(string)) > 0 {
		dqlRequest.DefaultTimeframeStart = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("default_timeframe_end"); ok && len(v.(string)) > 0 {
		dqlRequest.DefaultTimeframeEnd = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("timezone"); ok && len(v.(string)) > 0 {
		dqlRequest.Timezone = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("locale"); ok && len(v.(string)) > 0 {
		dqlRequest.Locale = opt.NewString(v.(string))
	}
	if v, ok := d.GetOk("max_result_records"); ok {
		dqlRequest.MaxResultRecords = opt.NewInt(v.(int))
	}
	if v, ok := d.GetOk("max_result_bytes"); ok {
		dqlRequest.MaxResultBytes = opt.NewInt(v.(int))
	}
	if v, ok := d.GetOk("fetch_timeout_seconds"); ok {
		dqlRequest.FetchTimeoutSeconds = opt.NewInt(v.(int))
	}
	if v, ok := d.GetOk("default_sampling_ratio"); ok {
		dqlRequest.DefaultSamplingRatio = opt.NewInt(v.(int))
	}
	if v, ok := d.GetOk("default_scan_limit_gbytes"); ok {
		dqlRequest.DefaultScanLimitGbytes = opt.NewInt(v.(int))
	}

	defer func(dqlRequest DQLRequest) {
		d.SetId(generateHash(dqlRequest))
	}(dqlRequest)

	data, _ := json.Marshal(dqlRequest)

	response, err := client.Fetch(ctx, []byte(data))
	if err != nil {
		return diag.FromErr(err)
	}
	var dqlResponse DQLResponse
	if err := json.Unmarshal(response.Data, &dqlResponse); err != nil {
		return diag.FromErr(err)
	}

	for {
		if dqlResponse.State == "SUCCEEDED" {
			d.Set("records", string(dqlResponse.Result.Records))
			return diag.Diagnostics{}
		}
		switch dqlResponse.State {
		case "NOT_STARTED", "RUNNING":
			if len(dqlResponse.RequestToken) == 0 {
				return diag.FromErr(fmt.Errorf("query is running but no request token for result polling was provided by REST API"))
			}
			time.Sleep(time.Duration(POLL_SLEEP_DURATION) * time.Millisecond)
			response, err := client.Poll(ctx, dqlResponse.RequestToken)
			if err != nil {
				return diag.FromErr(err)
			}
			dqlResponse = DQLResponse{}
			if err := json.Unmarshal(response.Data, &dqlResponse); err != nil {
				return diag.FromErr(err)
			}
		case "CANCELLED":
			return diag.FromErr(fmt.Errorf("query got cancelled unexpectedly"))
		case "FAILED":
			return diag.FromErr(fmt.Errorf("query failed"))
		default:
		}
	}

}

func generateHash(dqlRequest DQLRequest) string {
	data, _ := json.Marshal(dqlRequest)
	h := fnv.New128()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum([]byte{}))
}

type DQLResponse struct {
	State        string `json:"state"`
	RequestToken string `json:"requestToken"`
	Progress     int    `json:"progress"`
	Result       struct {
		Records json.RawMessage `json:"records"`
	}
}

type DQLRequest struct {
	Query                      string  `json:"query"`
	DefaultTimeframeStart      *string `json:"defaultTimeframeStart,omitempty"`
	DefaultTimeframeEnd        *string `json:"defaultTimeframeEnd,omitempty"`
	Timezone                   *string `json:"timezone,omitempty"`
	Locale                     *string `json:"locale,omitempty"`
	MaxResultRecords           *int    `json:"maxResultRecords,omitempty"`
	MaxResultBytes             *int    `json:"maxResultBytes,omitempty"`
	FetchTimeoutSeconds        *int    `json:"fetchTimeoutSeconds,omitempty"`
	RequestTimeoutMilliseconds int     `json:"requestTimeoutMilliseconds,omitempty"`
	EnablePreview              bool    `json:"enablePreview"`
	DefaultSamplingRatio       *int    `json:"defaultSamplingRatio,omitempty"`
	DefaultScanLimitGbytes     *int    `json:"defaultScanLimitGbytes,omitempty"`
	QueryOptions               any     `json:"queryOptions"`
	FilterSegments             any     `json:"filterSegments"`
}
