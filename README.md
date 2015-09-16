datadog has poor support for regular statsd ~ graphite style metrics
that have dimensionality using data in metric keys.
You can't wildcard, pattern match, drill down, etc.

so while NSQ has a statsd output, it's a pain to make use of it in datadog,
especially with multiple hosts, topics, channels, clients.

so this project sends the metrics to datadog using their tag based scheme.
