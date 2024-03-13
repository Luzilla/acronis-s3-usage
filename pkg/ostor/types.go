package ostor

// "nr_items": 7,
// "truncated": false,
// "items": [
// 	"s3-usage-8000000000000065-2017-02-01T16:31:54.000Z-1800",
// 	"s3-usage-8000000000000067-2017-02-01T16:30:51.000Z-1800",
// 	"s3-usage-8000000000000068-2017-02-01T16:27:25.000Z-1800",
// 	"s3-usage-8000000000000069-2017-02-01T16:27:24.000Z-1800",
// 	"s3-usage-8000000000000069-2017-02-01T16:31:07.000Z-1800",
// 	"s3-usage-800000000000006a-2017-02-01T16:27:24.000Z-1800",
// 	"s3-usage-800000000000006a-2017-02-01T16:31:08.000Z-1800"
// ]

type OStorResponse struct {
	Count     int      `json:"nr_items"`
	Truncated bool     `json:"truncated"`
	Items     []string `json:"items"`
}

// {
//     "fmt_version": 1,
//     "service_id": 8000000000000065,
//     "start_ts": 1485966714,
//     "period": 1390,
//     "nr_items": 1,
//     "items": [
//         {
//             "key": {
//                 "bucket": "client",
//                 "epoch": 98309,
//                 "user_id": "b81d6c5f895a8c86",
//                 "tag": ""
//             },
//             "counters": {
//                 "ops": {
//                     "put": 1,
//                     "get": 3,
//                     "list": 0,
//                     "other": 0
//                 },
//                 "net_io": {
//                     "uploaded": 41258,
//                     "downloaded": 45511311
//                 }
//             }
//         }
//     ]
// }

type OStorObjectUsageResponse struct {
	Version   int    `json:"fmt_version"`
	ServiceID string `json:"service_id"`
	StartTS   int    `json:"start_ts"`
	Period    int    `json:"period"`
	Count     int    `json:"nr_items"`
	Items     []struct {
		Key      ItemKey `json:"key"`
		Counters struct {
			Operations ItemCountersOps `json:"ops"`
			Net        ItemCountersNet `json:"net_io"`
		} `json:"counters"`
	} `json:"items"`
}

type ItemKey struct {
	Bucket string `json:"bucket"`
	Epoch  int    `json:"epoch"`
	UserID string `json:"user_id"`
	Tag    string `json:"tag"`
}

type ItemCountersOps struct {
	Put   int `json:"put"`
	Get   int `json:"get"`
	List  int `json:"list"`
	Other int `json:"other"`
}

type ItemCountersNet struct {
	Uploaded   int `json:"uploaded"`
	Downloaded int `json:"downloaded"`
}

// { "Users":[
// {
// "UserEmail": "msbnamfus@customer.planetary-networks.de",
// "UserId": "806e7d49f2dd9763",
// "State": "enabled",
// "OwnerId":"0000000000000000",
// "Flags": []
// },

type OstorUsersListResponse struct {
	Users []OstorUser `json:"users"`
}

type OstorUser struct {
	Email      string   `json:"UserEmail"`
	ID         string   `json:"UserId"`
	State      string   `json:"State"`
	Owner      string   `json:OwnerId"`
	Flags      []string `json:"Flags"`
	AccessKeys []struct {
		KeyID       string `json:"AWSAccessKeyId"`
		SecretKeyID string `json:"AWSSecretAccessKey"`
	} `json:"AWSAccessKeys,omitempty"`
	AccountCount string        `json:"AccountCount,omitempty"`
	Accounts     []interface{} `json:"Accounts",omitempty`
}

// {
// 	"UserEmail": "timo@suehl.com",
// 	"UserId": "a69657b97bc522ae",
// 	"State": "enabled",
// 	"OwnerId": "0000000000000000",
// 	"Flags": [],
// 	"AWSAccessKeys": [
// 	{
// 	"AWSAccessKeyId": "a69657b97bc522aeOLQJ",
// 	"AWSSecretAccessKey": "Dxt0wLpukuUrjkwuAnNRPoyNaQ62vWOXtJEYbYxh"
// 	}],
// 	"AccountCount": "0",
// 	"Accounts": [
// 	]
// 	}

// {"Buckets":[
// {    "name": "2023", "epoch": 0, "creation_date": "2023-02-05T09:56:47.000Z", "owner_id": "ca9b037812cbc5d5",      "size": {
//           "current" : 139658698, "hmax": 139658698, "h_integral": 1345907736733, "last_ts": 465442
//         }  },

type OstorBucketListResponse struct {
	Buckets []struct {
		Name      string `json:"name"`
		Epoch     int    `json:"epoc"`
		CreatedAt string `json:"creation_date"`
		OwnerID   string `json:"owner_id"`
		Size      struct {
			Current   int `json:"current"`
			HMax      int `json:"hmax"`
			HIntegral int `json:"h_integral"`
			LastTS    int `json:"last_ts"`
		} `json:"size"`
	}
}
