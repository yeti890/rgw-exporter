package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type RGWExporter struct {
	config Config

	// usage
	ops_total            *prometheus.Desc
	successful_ops_total *prometheus.Desc
	sent_bytes_total     *prometheus.Desc
	received_bytes_total *prometheus.Desc

	// bucket
	bucket_quota_enabled            *prometheus.Desc
	bucket_quota_size               *prometheus.Desc
	bucket_quota_objects            *prometheus.Desc
	bucket_size                     *prometheus.Desc
	bucket_actual_size              *prometheus.Desc
	bucket_objects                  *prometheus.Desc
	bucket_num_shards               *prometheus.Desc
	bucket_objects_per_shard        *prometheus.Desc
	buckets_total                   *prometheus.Desc
	buckets_size_total_bytes        *prometheus.Desc
	buckets_actual_size_total_bytes *prometheus.Desc
	bucket_quotas_size_total_bytes  *prometheus.Desc
	objects_total                   *prometheus.Desc

	// user
	user_suspended *prometheus.Desc

	user_quota_enabled     *prometheus.Desc
	user_quota_size_bytes  *prometheus.Desc
	user_quota_max_objects *prometheus.Desc

	user_bucket_quota_enabled     *prometheus.Desc
	user_bucket_quota_size_bytes  *prometheus.Desc
	user_bucket_quota_max_objects *prometheus.Desc

	users_total                  *prometheus.Desc
	user_buckets_total           *prometheus.Desc
	user_quotas_size_total_bytes *prometheus.Desc
	user_used_size_bytes         *prometheus.Desc

	// percent of usage quota
	bucket_quota_usage_percent *prometheus.Desc
	user_quota_usage_percent   *prometheus.Desc

	// service metrics
	collector_buckets_duration_seconds *prometheus.Desc
	collector_usage_duration_seconds   *prometheus.Desc
	collector_users_duration_seconds   *prometheus.Desc
}

func NewRGWExporter(config *Config) *RGWExporter {
	return &RGWExporter{
		config: *config,

		// usage — add uid
		ops_total: prometheus.NewDesc(
			"radosgw_usage_ops_total",
			"Number of requests",
			[]string{"region", "cluster", "endpoint", "uid", "bucket", "category"},
			nil,
		),
		successful_ops_total: prometheus.NewDesc(
			"radosgw_usage_successful_ops_total",
			"Number of successful requests",
			[]string{"region", "cluster", "endpoint", "uid", "bucket", "category"},
			nil,
		),
		sent_bytes_total: prometheus.NewDesc(
			"radosgw_usage_sent_bytes_total",
			"Bytes sent by the RGW",
			[]string{"region", "cluster", "endpoint", "uid", "bucket", "category"},
			nil,
		),
		received_bytes_total: prometheus.NewDesc(
			"radosgw_usage_received_bytes_total",
			"Bytes received by the RGW",
			[]string{"region", "cluster", "endpoint", "uid", "bucket", "category"},
			nil,
		),

		// bucket-level — add uid
		bucket_quota_enabled: prometheus.NewDesc(
			"radosgw_usage_bucket_quota_enabled",
			"Quota enabled for bucket",
			[]string{"region", "cluster", "endpoint", "bucket", "uid"},
			nil,
		),
		bucket_quota_size: prometheus.NewDesc(
			"radosgw_usage_bucket_quota_size",
			"Max allowed bucket size bytes (bucket quota)",
			[]string{"region", "cluster", "endpoint", "bucket", "uid"},
			nil,
		),
		bucket_quota_objects: prometheus.NewDesc(
			"radosgw_usage_bucket_quota_objects",
			"Max allowed objects in bucket",
			[]string{"region", "cluster", "endpoint", "bucket", "uid"},
			nil,
		),
		bucket_size: prometheus.NewDesc(
			"radosgw_usage_bucket_size",
			"Bucket size bytes (logical)",
			[]string{"region", "cluster", "endpoint", "bucket", "uid"},
			nil,
		),
		bucket_actual_size: prometheus.NewDesc(
			"radosgw_usage_bucket_actual_size",
			"Bucket actual size bytes (on disk)",
			[]string{"region", "cluster", "endpoint", "bucket", "uid"},
			nil,
		),
		bucket_objects: prometheus.NewDesc(
			"radosgw_usage_bucket_objects",
			"Bucket objects count",
			[]string{"region", "cluster", "endpoint", "bucket", "uid"},
			nil,
		),
		bucket_num_shards: prometheus.NewDesc(
			"radosgw_usage_bucket_num_shards",
			"Number of bucket index shards",
			[]string{"region", "cluster", "endpoint", "bucket", "uid"},
			nil,
		),
		bucket_objects_per_shard: prometheus.NewDesc(
			"radosgw_usage_bucket_objects_per_shard",
			"Number of objects per shard (objects / num_shards)",
			[]string{"region", "cluster", "endpoint", "bucket", "uid"},
			nil,
		),

		// aggregate for buckets
		buckets_total: prometheus.NewDesc(
			"radosgw_usage_buckets_total",
			"Total number of buckets",
			[]string{"region", "cluster", "endpoint"},
			nil,
		),
		buckets_size_total_bytes: prometheus.NewDesc(
			"radosgw_usage_buckets_size_total_bytes",
			"Total logical size of all buckets in bytes",
			[]string{"region", "cluster", "endpoint"},
			nil,
		),
		buckets_actual_size_total_bytes: prometheus.NewDesc(
			"radosgw_usage_buckets_actual_size_total_bytes",
			"Total actual size of all buckets in bytes",
			[]string{"region", "cluster", "endpoint"},
			nil,
		),
		bucket_quotas_size_total_bytes: prometheus.NewDesc(
			"radosgw_usage_bucket_quotas_size_total_bytes",
			"Total configured bucket quotas size in bytes (enabled and >0)",
			[]string{"region", "cluster", "endpoint"},
			nil,
		),
		objects_total: prometheus.NewDesc(
			"radosgw_usage_objects_total",
			"Total number of objects across all buckets",
			[]string{"region", "cluster", "endpoint"},
			nil,
		),

		// user-level
		user_suspended: prometheus.NewDesc(
			"radosgw_usage_user_suspended",
			"1 - suspended, 0 - active",
			[]string{"region", "cluster", "endpoint", "uid", "display_name"},
			nil,
		),

		user_quota_enabled: prometheus.NewDesc(
			"radosgw_usage_user_quota_enabled",
			"User quota enabled: 1 - enabled, 0 - disabled",
			[]string{"region", "cluster", "endpoint", "uid"},
			nil,
		),
		user_quota_size_bytes: prometheus.NewDesc(
			"radosgw_usage_user_quota_size_bytes",
			"User quota max size in bytes",
			[]string{"region", "cluster", "endpoint", "uid"},
			nil,
		),
		user_quota_max_objects: prometheus.NewDesc(
			"radosgw_usage_user_quota_objects",
			"User quota max objects",
			[]string{"region", "cluster", "endpoint", "uid"},
			nil,
		),

		user_bucket_quota_enabled: prometheus.NewDesc(
			"radosgw_usage_user_bucket_quota_enabled",
			"User bucket quota enabled: 1 - enabled, 0 - disabled",
			[]string{"region", "cluster", "endpoint", "uid"},
			nil,
		),
		user_bucket_quota_size_bytes: prometheus.NewDesc(
			"radosgw_usage_user_bucket_quota_size_bytes",
			"User bucket quota max size in bytes",
			[]string{"region", "cluster", "endpoint", "uid"},
			nil,
		),
		user_bucket_quota_max_objects: prometheus.NewDesc(
			"radosgw_usage_user_bucket_quota_objects",
			"User bucket quota max objects",
			[]string{"region", "cluster", "endpoint", "uid"},
			nil,
		),

		users_total: prometheus.NewDesc(
			"radosgw_usage_users_total",
			"Total number of users",
			[]string{"region", "cluster", "endpoint"},
			nil,
		),
		user_buckets_total: prometheus.NewDesc(
			"radosgw_usage_user_buckets_total",
			"Total number of buckets owned by user",
			[]string{"region", "cluster", "endpoint", "uid"},
			nil,
		),
		user_quotas_size_total_bytes: prometheus.NewDesc(
			"radosgw_usage_user_quotas_size_total_bytes",
			"Total configured user quotas size in bytes (enabled and >0)",
			[]string{"region", "cluster", "endpoint"},
			nil,
		),
		user_used_size_bytes: prometheus.NewDesc(
			"radosgw_usage_user_used_size_bytes",
			"Total logical used size by user (sum of bucket sizes), in bytes",
			[]string{"region", "cluster", "endpoint", "uid"},
			nil,
		),

		bucket_quota_usage_percent: prometheus.NewDesc(
			"radosgw_usage_bucket_quota_usage_percent",
			"Bucket quota usage in percent (0-100), size-based",
			[]string{"region", "cluster", "endpoint", "bucket", "uid"},
			nil,
		),
		user_quota_usage_percent: prometheus.NewDesc(
			"radosgw_usage_user_quota_usage_percent",
			"User quota usage in percent (0-100), size-based",
			[]string{"region", "cluster", "endpoint", "uid"},
			nil,
		),

		collector_buckets_duration_seconds: prometheus.NewDesc(
			"radosgw_usage_collector_buckets_duration_seconds",
			"Buckets collector duration time",
			[]string{"region", "cluster", "endpoint"},
			nil,
		),
		collector_usage_duration_seconds: prometheus.NewDesc(
			"radosgw_usage_collector_usage_duration_seconds",
			"Usage collector duration time",
			[]string{"region", "cluster", "endpoint"},
			nil,
		),
		collector_users_duration_seconds: prometheus.NewDesc(
			"radosgw_usage_collector_users_duration_seconds",
			"Users collector duration time",
			[]string{"region", "cluster", "endpoint"},
			nil,
		),
	}
}

func (collector *RGWExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.ops_total
	ch <- collector.successful_ops_total
	ch <- collector.sent_bytes_total
	ch <- collector.received_bytes_total

	ch <- collector.bucket_quota_enabled
	ch <- collector.bucket_quota_size
	ch <- collector.bucket_quota_objects
	ch <- collector.bucket_size
	ch <- collector.bucket_actual_size
	ch <- collector.bucket_objects
	ch <- collector.bucket_num_shards
	ch <- collector.bucket_objects_per_shard

	ch <- collector.buckets_total
	ch <- collector.buckets_size_total_bytes
	ch <- collector.buckets_actual_size_total_bytes
	ch <- collector.bucket_quotas_size_total_bytes
	ch <- collector.objects_total

	ch <- collector.user_suspended
	ch <- collector.user_quota_enabled
	ch <- collector.user_quota_size_bytes
	ch <- collector.user_quota_max_objects
	ch <- collector.user_bucket_quota_enabled
	ch <- collector.user_bucket_quota_size_bytes
	ch <- collector.user_bucket_quota_max_objects

	ch <- collector.users_total
	ch <- collector.user_buckets_total
	ch <- collector.user_quotas_size_total_bytes
	ch <- collector.user_used_size_bytes

	ch <- collector.bucket_quota_usage_percent
	ch <- collector.user_quota_usage_percent

	ch <- collector.collector_buckets_duration_seconds
	ch <- collector.collector_usage_duration_seconds
	ch <- collector.collector_users_duration_seconds
}

func (collector *RGWExporter) Collect(ch chan<- prometheus.Metric) {
	region := collector.config.Region
	cluster := collector.config.ClusterName
	endpoint := collector.config.PubEndpoint

	// ---------- buckets: per-bucket & aggregate ----------

	bucketsMu.Lock()

	bucketsTotal := 0
	totalBucketSize := 0.0
	totalBucketActualSize := 0.0
	totalBucketQuotasSize := 0.0
	totalObjects := 0.0

	userBucketCount := make(map[string]float64)
	userUsedSize := make(map[string]float64)

	for _, bucket := range buckets {
		bucketsTotal++

		quotaEnabled := 0.0
		quotaSize := 0.0
		quotaObjects := 0.0

		if bucket.BucketQuota.Enabled != nil && *bucket.BucketQuota.Enabled {
			quotaEnabled = 1.0
		}

		if bucket.BucketQuota.MaxSize != nil {
			quotaSize = float64(*bucket.BucketQuota.MaxSize)
		} else if bucket.BucketQuota.MaxSizeKb != nil {
			quotaSize = float64(*bucket.BucketQuota.MaxSizeKb) * 1024.0
		}

		if bucket.BucketQuota.MaxObjects != nil {
			quotaObjects = float64(*bucket.BucketQuota.MaxObjects)
		}

		bucketSize := 0.0
		if bucket.Usage.RgwMain.Size != nil {
			bucketSize = float64(*bucket.Usage.RgwMain.Size)
		}

		bucketActualSize := 0.0
		if bucket.Usage.RgwMain.SizeActual != nil {
			bucketActualSize = float64(*bucket.Usage.RgwMain.SizeActual)
		}

		bucketObjects := 0.0
		if bucket.Usage.RgwMain.NumObjects != nil {
			bucketObjects = float64(*bucket.Usage.RgwMain.NumObjects)
		}

		// num_shards
		numShards := -1.0
		if bucket.NumShards != nil {
			numShards = float64(*bucket.NumShards)
		}

		objectsPerShard := 0.0
		if numShards > 0 {
			objectsPerShard = bucketObjects / numShards
		}

		// aggragates
		totalBucketSize += bucketSize
		totalBucketActualSize += bucketActualSize
		totalObjects += bucketObjects

		if quotaEnabled == 1.0 && quotaSize > 0 {
			totalBucketQuotasSize += quotaSize
		}

		uid := bucket.Owner

		if uid != "" {
			userBucketCount[uid]++
			userUsedSize[uid] += bucketSize
		}

		// per-bucket metrics (add uid)
		ch <- prometheus.MustNewConstMetric(
			collector.bucket_quota_enabled,
			prometheus.GaugeValue,
			quotaEnabled,
			region, cluster, endpoint, bucket.Bucket, uid,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.bucket_quota_size,
			prometheus.GaugeValue,
			quotaSize,
			region, cluster, endpoint, bucket.Bucket, uid,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.bucket_quota_objects,
			prometheus.GaugeValue,
			quotaObjects,
			region, cluster, endpoint, bucket.Bucket, uid,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.bucket_size,
			prometheus.GaugeValue,
			bucketSize,
			region, cluster, endpoint, bucket.Bucket, uid,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.bucket_actual_size,
			prometheus.GaugeValue,
			bucketActualSize,
			region, cluster, endpoint, bucket.Bucket, uid,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.bucket_objects,
			prometheus.GaugeValue,
			bucketObjects,
			region, cluster, endpoint, bucket.Bucket, uid,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.bucket_num_shards,
			prometheus.GaugeValue,
			numShards,
			region, cluster, endpoint, bucket.Bucket, uid,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.bucket_objects_per_shard,
			prometheus.GaugeValue,
			objectsPerShard,
			region, cluster, endpoint, bucket.Bucket, uid,
		)

		quotaUsagePercent := 0.0
		if quotaEnabled == 1.0 && quotaSize > 0 {
			quotaUsagePercent = (bucketSize / quotaSize) * 100.0
		}

		ch <- prometheus.MustNewConstMetric(
			collector.bucket_quota_usage_percent,
			prometheus.GaugeValue,
			quotaUsagePercent,
			region, cluster, endpoint, bucket.Bucket, uid,
		)
	}

	bucketsMu.Unlock()

	// aggregate from buckets & objects
	ch <- prometheus.MustNewConstMetric(
		collector.buckets_total,
		prometheus.GaugeValue,
		float64(bucketsTotal),
		region, cluster, endpoint,
	)

	ch <- prometheus.MustNewConstMetric(
		collector.buckets_size_total_bytes,
		prometheus.GaugeValue,
		totalBucketSize,
		region, cluster, endpoint,
	)

	ch <- prometheus.MustNewConstMetric(
		collector.buckets_actual_size_total_bytes,
		prometheus.GaugeValue,
		totalBucketActualSize,
		region, cluster, endpoint,
	)

	ch <- prometheus.MustNewConstMetric(
		collector.bucket_quotas_size_total_bytes,
		prometheus.GaugeValue,
		totalBucketQuotasSize,
		region, cluster, endpoint,
	)

	ch <- prometheus.MustNewConstMetric(
		collector.objects_total,
		prometheus.GaugeValue,
		totalObjects,
		region, cluster, endpoint,
	)

	// ---------- usage ----------

	usageMu.Lock()
	for key, stats := range usageMap {
		ch <- prometheus.MustNewConstMetric(
			collector.sent_bytes_total,
			prometheus.CounterValue,
			float64(stats.BytesSent),
			region, cluster, endpoint, key.User, key.Bucket, key.Category,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.received_bytes_total,
			prometheus.CounterValue,
			float64(stats.BytesReceived),
			region, cluster, endpoint, key.User, key.Bucket, key.Category,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.ops_total,
			prometheus.CounterValue,
			float64(stats.Ops),
			region, cluster, endpoint, key.User, key.Bucket, key.Category,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.successful_ops_total,
			prometheus.CounterValue,
			float64(stats.SuccessfulOps),
			region, cluster, endpoint, key.User, key.Bucket, key.Category,
		)
	}
	usageMu.Unlock()

	// ---------- users ----------

	usersMu.Lock()
	usersTotal := len(users)
	totalUserQuotasSize := 0.0

	for _, user := range users {
		ch <- prometheus.MustNewConstMetric(
			collector.user_suspended,
			prometheus.GaugeValue,
			float64(user.Suspended),
			region, cluster, endpoint, user.UserId, user.DisplayName,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.user_quota_enabled,
			prometheus.GaugeValue,
			user.UserQuotaEnabled,
			region, cluster, endpoint, user.UserId,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.user_quota_size_bytes,
			prometheus.GaugeValue,
			user.UserQuotaMaxSizeBytes,
			region, cluster, endpoint, user.UserId,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.user_quota_max_objects,
			prometheus.GaugeValue,
			user.UserQuotaMaxObjects,
			region, cluster, endpoint, user.UserId,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.user_bucket_quota_enabled,
			prometheus.GaugeValue,
			user.UserBucketQuotaEnabled,
			region, cluster, endpoint, user.UserId,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.user_bucket_quota_size_bytes,
			prometheus.GaugeValue,
			user.UserBucketQuotaMaxSizeBytes,
			region, cluster, endpoint, user.UserId,
		)

		ch <- prometheus.MustNewConstMetric(
			collector.user_bucket_quota_max_objects,
			prometheus.GaugeValue,
			user.UserBucketQuotaMaxObjects,
			region, cluster, endpoint, user.UserId,
		)

		// total buckets from uid
		if cnt, ok := userBucketCount[user.UserId]; ok {
			ch <- prometheus.MustNewConstMetric(
				collector.user_buckets_total,
				prometheus.GaugeValue,
				cnt,
				region, cluster, endpoint, user.UserId,
			)
		} else {
			ch <- prometheus.MustNewConstMetric(
				collector.user_buckets_total,
				prometheus.GaugeValue,
				0,
				region, cluster, endpoint, user.UserId,
			)
		}

		// used size by uid
		used := userUsedSize[user.UserId]

		ch <- prometheus.MustNewConstMetric(
			collector.user_used_size_bytes,
			prometheus.GaugeValue,
			used,
			region, cluster, endpoint, user.UserId,
		)

		// percent usage user quota by uid
		quotaUsagePercent := 0.0
		if user.UserQuotaEnabled == 1.0 && user.UserQuotaMaxSizeBytes > 0 {
			quotaUsagePercent = (used / user.UserQuotaMaxSizeBytes) * 100.0
			totalUserQuotasSize += user.UserQuotaMaxSizeBytes
		}

		ch <- prometheus.MustNewConstMetric(
			collector.user_quota_usage_percent,
			prometheus.GaugeValue,
			quotaUsagePercent,
			region, cluster, endpoint, user.UserId,
		)
	}
	usersMu.Unlock()

	ch <- prometheus.MustNewConstMetric(
		collector.users_total,
		prometheus.GaugeValue,
		float64(usersTotal),
		region, cluster, endpoint,
	)

	ch <- prometheus.MustNewConstMetric(
		collector.user_quotas_size_total_bytes,
		prometheus.GaugeValue,
		totalUserQuotasSize,
		region, cluster, endpoint,
	)

	// ---------- service metrics ----------

	ch <- prometheus.MustNewConstMetric(
		collector.collector_buckets_duration_seconds,
		prometheus.GaugeValue,
		collectBucketsDuration.Seconds(),
		region, cluster, endpoint,
	)

	ch <- prometheus.MustNewConstMetric(
		collector.collector_usage_duration_seconds,
		prometheus.GaugeValue,
		collectUsageDuration.Seconds(),
		region, cluster, endpoint,
	)

	ch <- prometheus.MustNewConstMetric(
		collector.collector_users_duration_seconds,
		prometheus.GaugeValue,
		collectUsersDuration.Seconds(),
		region, cluster, endpoint,
	)
}
