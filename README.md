# GO Scripts

Collection of scripts written in GO for learning or experimental purposes.

## Current Scripts

The following scripts are currently present in the repository.

#### `latency_check.go`

A script to understand the meaningfulness of different single value estimations of latency for monitoring purposes.

With help of T-Digest data structure, we calculate different percentiles for the received latency. We also calculate the average and the maximum to contrast them with the percentile values.

```
avg:  0.13959714512499996
max:  2.32450015
50th percentile:  0.1239733715
75th percentile:  0.1275092825
90th percentile:  0.2502230525
99th percentile:  0.25125328300000005
```

As we see from the results above, maximum and average are misleading since either they are presenting the outlier or taking outliers into consideration. Percentile is a much better figure telling that 99% of the times, the latency was below 0.25 seconds.