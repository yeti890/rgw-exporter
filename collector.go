package main

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"sync"
	"time"

	rgw "github.com/ceph/go-ceph/rgw/admin"
)

var (
	buckets   []rgw.Bucket
	bucketsMu sync.Mutex
)

var (
	usageMap map[UsageKey]*UsageStats
	usageMu  sync.Mutex
)

var (
	users   []UserInfo
	usersMu sync.Mutex
)

var (
	collectUsageDuration   time.Duration
	collectUsageDurationMu sync.Mutex

	collectBucketsDuration   time.Duration
	collectBucketsDurationMu sync.Mutex

	collectUsersDuration   time.Duration
	collectUsersDurationMu sync.Mutex
)

type UsageKey struct {
	User     string
	Bucket   string
	Owner    string
	Category string
}

type UsageStats struct {
	BytesSent     uint64
	BytesReceived uint64
	Ops           uint64
	SuccessfulOps uint64
}

// Информация о пользователе + квоты
type UserInfo struct {
	UserId      string
	DisplayName string
	Suspended   int

	// user_quota
	UserQuotaEnabled      float64
	UserQuotaMaxSizeBytes float64
	UserQuotaMaxObjects   float64

	// bucket_quota (user bucket quota)
	UserBucketQuotaEnabled      float64
	UserBucketQuotaMaxSizeBytes float64
	UserBucketQuotaMaxObjects   float64
}

func startRGWStatCollector(config *Config) {
	conn := getRGWConnection(config)

	tickerUsage := time.NewTicker(time.Duration(config.UsageCollectorInterval) * time.Second)
	tickerBuckets := time.NewTicker(time.Duration(config.BucketsCollectorInterval) * time.Second)
	tickerUsers := time.NewTicker(time.Duration(config.UsersCollectorInterval) * time.Second)

	// usage
	go func() {
		for ; ; <-tickerUsage.C {
			collectUsage(conn, config.SkipWithoutBucket)
		}
	}()

	// buckets
	go func() {
		for ; ; <-tickerBuckets.C {
			collectBuckets(conn)
		}
	}()

	// users
	go func() {
		for ; ; <-tickerUsers.C {
			if config.UsersCollectorEnable {
				collectUsers(conn)
			} else {
				usersMu.Lock()
				users = nil
				usersMu.Unlock()
			}
		}
	}()
}

func getRGWConnection(config *Config) *rgw.API {
	var tr *http.Transport
	if config.Insecure {
		tr = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	} else {
		tr = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false}}
	}

	conn, err := rgw.New(
		config.Endpoint,
		config.AccessKey,
		config.SecretKey,
		&http.Client{
			Timeout:   time.Duration(config.RGWConnectionTimeout) * time.Second,
			Transport: tr,
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	return conn
}

func collectUsage(conn *rgw.API, skipWithoutBucket bool) {
	start := time.Now()

	today := time.Now().UTC().Format(time.DateOnly)
	curUsage, err := conn.GetUsage(context.Background(), rgw.Usage{
		ShowSummary: func() *bool { b := false; return &b }(),
		Start:       today,
	})
	if err != nil {
		log.Println("Unable to get usage stat:", err)
		return
	}

	usageMu.Lock()
	usageMap = sumUsage(curUsage, skipWithoutBucket)
	usageMu.Unlock()

	collectUsageDurationMu.Lock()
	collectUsageDuration = time.Since(start)
	collectUsageDurationMu.Unlock()
}

func collectBuckets(conn *rgw.API) {
	start := time.Now()

	curBuckets, err := conn.ListBucketsWithStat(context.Background())
	if err != nil {
		log.Println("Unable to get bucket stat:", err)
		return
	}

	bucketsMu.Lock()
	buckets = curBuckets
	bucketsMu.Unlock()

	collectBucketsDurationMu.Lock()
	collectBucketsDuration = time.Since(start)
	collectBucketsDurationMu.Unlock()
}

func collectUsers(conn *rgw.API) {
	start := time.Now()
	var curUsers []UserInfo

	curUsersList, err := conn.GetUsers(context.Background())
	if err != nil {
		log.Println("Unable to get users list:", err)
		return
	}

	for _, uid := range *curUsersList {
		curUser, err := conn.GetUser(context.Background(), rgw.User{ID: uid})
		if err != nil {
			log.Println("Unable to get user info for", uid, ":", err)
			continue
		}

		// suspended
		suspended := 0
		if curUser.Suspended != nil {
			suspended = *curUser.Suspended
		}

		// user_quota
		var userQuotaEnabled float64
		var userQuotaMaxSizeBytes float64
		var userQuotaMaxObjects float64

		// curUser.UserQuota — это QuotaSpec (значение), поля внутри — указатели
		if curUser.UserQuota.Enabled != nil && *curUser.UserQuota.Enabled {
			userQuotaEnabled = 1.0
		}

		if curUser.UserQuota.MaxSize != nil {
			userQuotaMaxSizeBytes = float64(*curUser.UserQuota.MaxSize)
		} else if curUser.UserQuota.MaxSizeKb != nil {
			userQuotaMaxSizeBytes = float64(*curUser.UserQuota.MaxSizeKb) * 1024.0
		}

		if curUser.UserQuota.MaxObjects != nil {
			userQuotaMaxObjects = float64(*curUser.UserQuota.MaxObjects)
		}

		// bucket_quota (user bucket quota)
		var userBucketQuotaEnabled float64
		var userBucketQuotaMaxSizeBytes float64
		var userBucketQuotaMaxObjects float64

		// curUser.BucketQuota — тоже QuotaSpec (значение)
		if curUser.BucketQuota.Enabled != nil && *curUser.BucketQuota.Enabled {
			userBucketQuotaEnabled = 1.0
		}

		if curUser.BucketQuota.MaxSize != nil {
			userBucketQuotaMaxSizeBytes = float64(*curUser.BucketQuota.MaxSize)
		} else if curUser.BucketQuota.MaxSizeKb != nil {
			userBucketQuotaMaxSizeBytes = float64(*curUser.BucketQuota.MaxSizeKb) * 1024.0
		}

		if curUser.BucketQuota.MaxObjects != nil {
			userBucketQuotaMaxObjects = float64(*curUser.BucketQuota.MaxObjects)
		}

		user := UserInfo{
			UserId:      curUser.ID,
			DisplayName: curUser.DisplayName,
			Suspended:   suspended,

			UserQuotaEnabled:      userQuotaEnabled,
			UserQuotaMaxSizeBytes: userQuotaMaxSizeBytes,
			UserQuotaMaxObjects:   userQuotaMaxObjects,

			UserBucketQuotaEnabled:      userBucketQuotaEnabled,
			UserBucketQuotaMaxSizeBytes: userBucketQuotaMaxSizeBytes,
			UserBucketQuotaMaxObjects:   userBucketQuotaMaxObjects,
		}

		curUsers = append(curUsers, user)
	}

	usersMu.Lock()
	users = curUsers
	usersMu.Unlock()

	collectUsersDurationMu.Lock()
	collectUsersDuration = time.Since(start)
	collectUsersDurationMu.Unlock()
}

func sumUsage(usage rgw.Usage, skipWithoutBucket bool) map[UsageKey]*UsageStats {
	usageStatsMap := make(map[UsageKey]*UsageStats)

	for _, userUsage := range usage.Entries {
		for _, bucket := range userUsage.Buckets {
			if skipWithoutBucket && (bucket.Bucket == "" || bucket.Bucket == "-") {
				continue
			}
			for _, category := range bucket.Categories {
				key := UsageKey{
					User:     userUsage.User,
					Bucket:   bucket.Bucket,
					Owner:    bucket.Owner,
					Category: category.Category,
				}

				if stats, exists := usageStatsMap[key]; !exists {
					usageStatsMap[key] = &UsageStats{
						BytesSent:     category.BytesSent,
						BytesReceived: category.BytesReceived,
						Ops:           category.Ops,
						SuccessfulOps: category.SuccessfulOps,
					}
				} else {
					stats.BytesSent += category.BytesSent
					stats.BytesReceived += category.BytesReceived
					stats.Ops += category.Ops
					stats.SuccessfulOps += category.SuccessfulOps
				}
			}
		}
	}
	return usageStatsMap
}
