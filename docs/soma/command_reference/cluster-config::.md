# cluster management

```
soma cluster create ${cluster} in ${bucket}
soma cluster destroy ${cluster} in ${bucket}
soma cluster list in ${bucket}
soma cluster show ${cluster} in ${bucket}
soma cluster dumptree ${cluster} in ${bucket}
soma cluster member list of ${cluster} in ${bucket}
soma cluster member assign ${node} to ${cluster} [in ${bucket}]
soma cluster member unassign ${node} from ${cluster} [in ${bucket}]
soma cluster property create system  ${system}  on ${cluster} in ${bucket} view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma cluster property create custom  ${custom}  on ${cluster} in ${bucket} view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma cluster property create service ${service} on ${cluster} in ${bucket} view ${view} [inheritance ${inherit}] [childrenonly ${child}]
soma cluster property create oncall  ${oncall}  on ${cluster} in ${bucket} view ${view} [inheritance ${inherit}] [childrenonly ${child}]
soma cluster property update system  ${system}  on ${cluster} in ${bucket} view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma cluster property update custom  ${custom}  on ${cluster} in ${bucket} view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma cluster property destroy system  ${system}  on ${cluster} in ${bucket} view ${view}
soma cluster property destroy custom  ${custom}  on ${cluster} in ${bucket} view ${view}
soma cluster property destroy service ${service} on ${cluster} in ${bucket} view ${view}
soma cluster property destroy oncall  ${oncall}  on ${cluster} in ${bucket} view ${view}
```

See `soma cluster help ${command}` for detailed help.
