# monitoring metric definitions

Metric definitions provide the pool of possible measurements
that monitoring systems can declare themselves capable of
monitoring.

# SYNOPSIS OVERVIEW

```
soma metric add ${metric} unit ${unit} description ${text} \
                          [package ${provider}::${package}, ...]
soma metric remove ${metric}
soma metric list
soma metric show ${metric}
```
