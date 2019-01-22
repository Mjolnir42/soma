# group management

```
soma group create ${group} in ${bucket}
soma group destroy ${group} in ${bucket}
soma group list in ${bucket}
soma group show ${group} in ${bucket}
soma group dumptree ${group} in ${bucket}
soma group member assign group ${child-group} to ${group} in ${bucket}
soma group member unassign group ${child-group} from ${group} in ${bucket}
soma group property create system  ${system}  on ${group} in ${bucket} view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma group property create custom  ${custom}  on ${group} in ${bucket} view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma group property create service ${service} on ${group} in ${bucket} view ${view} [inheritance ${inherit}] [childrenonly ${child}]
soma group property create oncall  ${oncall}  on ${group} in ${bucket} view ${view} [inheritance ${inherit}] [childrenonly ${child}]
soma group property destroy system  ${system}  on ${group} in ${bucket} view ${view}
soma group property destroy custom  ${custom}  on ${group} in ${bucket} view ${view}
soma group property destroy service ${service} on ${group} in ${bucket} view ${view}
soma group property destroy oncall  ${oncall}  on ${group} in ${bucket} view ${view}
```

See `soma group help ${command}` for detailed help.
