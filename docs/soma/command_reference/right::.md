# rights overview

Rights are granted permissions. They can be granted to the following
subjects types:

* user
* admin
* team
* tool

They are either global or scoped to a specific object.
These objects can have the following entitites:

* repository
* bucket
* monitoring
* team

Teams can be both subject and object of the grant, granting the members of
the team the permission in their team's scope.

Rights are runtime definitions, executed via the cli. Rights change what
they grant as the granted permission is remapped.

# SYNOPSIS OVERVIEW

```
soma right list ${category}::${permission}
soma right grant ${category}::${permission} to user|admin|team|tool ${name} [on repository|bucket|monitoring|team ${object}]
soma right revoke ${category}::${permission} from user|admin|team|tool ${name} [on repository|bucket|monitoring|team ${object}]
```

See `soma right help ${command}` for detailed help.
