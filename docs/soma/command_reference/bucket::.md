# bucket management

Bucket management contains the actions for managing configuration
buckets within repositories.

# SYNOPSIS OVERVIEW

```
soma bucket create ${bucket} in ${repository} environment ${env}
soma bucket destroy ${bucket} [in ${repository}]

soma bucket list in ${repository}
soma bucket show ${bucket} [in ${repository}]
soma bucket dumptree ${bucket} [in ${repository}]
soma bucket search [id] [name] [repository] [env]
soma bucket rename ${bucket} to ${newName} [in ${repository}]

soma bucket property create system  ${system}  on ${repository} view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma bucket property create custom  ${custom}  on ${repository} view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma bucket property create service ${service} on ${repository} view ${view} [inheritance ${inherit}] [childrenonly ${child}]
soma bucket property create oncall  ${oncall}  on ${repository} view ${view} [inheritance ${inherit}] [childrenonly ${child}]

soma bucket property destroy system  ${system}  on ${repository} view ${view}
soma bucket property destroy custom  ${custom}  on ${repository} view ${view}
soma bucket property destroy service ${service} on ${repository} view ${view}
soma bucket property destroy oncall  ${oncall}  on ${repository} view ${view}
```

See `soma bucket help ${command} for detailed help.
