package segment

type Logo struct {
	Default string `json:"default"`
	Mark    string `json:"mark"`
	Alt     string `json:"alt"`
}

type Pagination struct {
	Current      string  `json:"current"`
	Next         *string `json:"next,omitempty"`
	TotalEntries int     `json:"totalEntries"`
}

type IntegrationOption struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Required     bool   `json:"required"`
	Description  string `json:"description"`
	DefaultValue string `json:"defaultValue"`
	Label        string `json:"label"`
}

var (
	SourceSlugs = []string{
		"active-campaign",
		"adwords",
		"aircall",
		"airship",
		"amazon-s3",
		"amp",
		"amplitude-cohorts",
		"android",
		"autopilothq",
		"beamer",
		"blueshift",
		"braze",
		"candu",
		"chatlio",
		"clojure",
		"customerio",
		"delighted",
		"drip",
		"facebook-ads",
		"facebook-lead-ads",
		"factual-engine",
		"foursquare-pilgrim",
		"friendbuy",
		"go",
		"herow",
		"http-api",
		"hubspot",
		"intercom",
		"ios",
		"iterable",
		"java",
		"javascript",
		"klaviyo",
		"klenty",
		"kotlin",
		"kotlin-android",
		"launchdarkly",
		"leanplum",
		"looker",
		"mailchimp",
		"mailjet",
		"mandrill",
		"marketo",
		"mixpanel-cohorts-source",
		"moesif-api-analytics",
		"net",
		"node.js",
		"nudgespot",
		"pendo",
		"php",
		"pixel-tracking-api",
		"project",
		"provesource",
		"python",
		"radar",
		"react-native",
		"refiner",
		"regal-voice",
		"roku",
		"ruby",
		"salesforce",
		"salesforce-marketing-cloud",
		"selligent-marketing-cloud",
		"sendgrid",
		"shopify-littledata",
		"snowflake",
		"stripe",
		"swift",
		"twilio",
		"twilio-event-streams",
		"vero",
		"wootric",
		"xamarin",
		"youbora",
		"zendesk",
	}
	WarehouseSlugs = []string{
		"azuresqldb",
		"azuresqldw",
		"bigquery",
		"db2",
		"postgres",
		"redshift",
		"snowflake",
	}
)
