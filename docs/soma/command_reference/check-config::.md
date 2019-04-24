# check configuration

```
soma check-config create ${check} in ${repository} on ${entityType} ${entityName} with ${capability} threshold predicate '>=' value ${value} level info interval 60 [${constraints}] 
soma check-config destroy ${check} in repository ${repository}
soma check-config list in ${check}
soma check-config show ${check} in ${bucket}
```
${entityType} can be any of repository|bucket|group|cluster|node

See `soma check-config help ${command}` for detailed help.