# SOMA Setup Guide

## Initial Bootstrap

1. Install pgSQL database, perform instructions from somadbctl/README.md

    1. connect as `pgsql` or another superuser and execute the following queries:
```
CREATE ROLE soma_dba WITH LOGIN PASSWORD 'somepassword';
CREATE ROLE soma_svc WITH LOGIN PASSWORD 'otherpassword';
CREATE ROLE soma_inv WITH LOGIN PASSWORD 'randomstring';
CREATE DATABASE soma WITH OWNER soma_dba;
GRANT CONNECT ON DATABASE soma TO soma_dba;
GRANT CONNECT ON DATABASE soma TO soma_svc;
GRANT CONNECT ON DATABASE soma TO soma_inv;
\connect soma;
CREATE EXTENSION pgcrypto SCHEMA public;
```
    2. update `pg_hba.conf` so all three users can connect to the
       database via TCP socket

The `soma_dba` account is the account used by this tool, it creates the
schemas, tables, indexes etc. It also issues the correct GRANT
statements for the other, restricted users.

`soma_svc` is the user of the SOMA service daemon.

`soma_inv` is a currently unused third account, to allow for inventory
updates to be split into another daemon.

2. `mkdir -p ~/.soma/dbctl`

3. `vim ~/.soma/dbctl/somadbctl.conf`

```
	# somadbctl configuration file
	environment: production
	timeout: 5
	tlsmode: disable
	database: {
	  host: localhost
	  user: soma_dba
	  dbname: soma
	  port: 5432
	  password: *************
	}
```

4. Initialize database:
```
	% somadbctl init
	...
	The generated boostrap token was: 79f4fce4c2ef288bbea0ddd53067f86da90ea662ab57f20c60d5ad74729f98fd
```

5. Create server directory structure for default instance `huxley`

```
    * mkdir -p /srv/soma/huxley/conf/
    * mkdir -p /srv/soma/huxley/log/job
    * mkdir -p /srv/soma/huxley/log/repository
```

6. vim /srv/soma/huxley/conf/soma.conf --- the '' around ldap.base.dn are important! (suspected UCL bug)

```
	# soma configuration file
	environment: production
	readonly: false

	open.door.policy: false
	database: {
	  host: localhost
	  user: soma_svc
	  database: soma
	  port: 5432
	  password: ********
	  timeout: 5
	  tlsmode: disable
	}
	daemon: {
	  listen: localhost
	  port: 8888
	  tls: true
	  cert.file: /srv/soma/huxley/conf/soma.pem
	  key.file: /srv/soma/huxley/conf/soma.key.pem
	}
	authentication: {
	  kex.expiry: 60
	  token.expiry: 43200
	  credential.expiry: 365
	  activation.mode: ldap
	  # dd if=/dev/random bs=1M count=16 2>/dev/null | sha512 | cut -c 1-64
	  token.seed: 5ae10f15a8a341d67fd2ed3fb18176f8ccdb2d82383304e64fcbed6f7d3f6eb3
	  token.key: 9b77ee4fc8cc433624559b2bbbaa7eb761749f2754a1fd248b602cf14ea80a5f
	}
	ldap: {
	  uid.attribute: uid
	  base.dn: 'o=foobar,c=SNAFU'
	  user.dn: sn=userDN
	  address: ldap.example.org
	  port: 636
	  tls: true
	  cert.file: /srv/soma/huxley/conf/ldap.example.org.chain.pem
	  insecure: false
	}
```

7. Generate self-signed SSL certificate to `localhost`

```
	% openssl genrsa -out /srv/soma/huxley/conf/soma.key.pem 4096
	% openssl req -new -x509 -key /srv/soma/huxley/conf/soma.key.pem -out /srv/soma/huxley/conf/soma.pem -days 365 -config /etc/ssl/openssl.cnf
```

8. Copy certificate chain for `ldap.example.org` to `/srv/soma/huxley/conf/ldap.example.org.chain.pem`

9. Startup SOMA

```
somad -nopoke -config /srv/soma/huxley/conf/soma.conf
```

9. Initialize soma cli files

```
	% soma init
	% cp /srv/soma/huxley/conf/soma.pem ~/.soma/adm/ca.pem
```

10. vim ~/.soma/adm/somaadm.conf --- high boltdb.open.timeout if somaadm is called in parallel a lot

```
	# somaadm configuration file
	timeout: 5
	api: https://localhost:8888/
	cert: ca.pem
	logdir: logs
	activation.mode: ldap
	boltdb: {
	  path: db
	  file: somaadm.db
	  mode: 0600
	  open.timeout: 30
	}
  auth: {
    user: root
  # pass: example_password
  # token: 294f2fc329dbc725ae267bc3318f247fa9a37089d7e1eb04f24a72a2651a01f8
  }
```

11. Activate root account

```
	% soma ops bootstrap

Welcome to SOMA!

This dialogue will guide you to set up the system's root account of
your new instance.

As first step, enter the root password you want to set.
Enter password:
Repeat password:
✔ Entered passwords match
Password score    (0-4): 4
Estimated entropy (bit): 43.509000
Estimated time to crack: 2.0 years
Select this password? (y/n): y

Very good. Now enter the bootstrap token printed by somadbctl at the
end of the schema installation.
Enter password:
Repeat password:
✔ Entered passwords match

Alright. Let's sully that pristine database. Here we go!

Generating keypair: ✔  OK
Initiating key exchange: ✔  OK
Sending bootstrap request: ✔  OK
Validating received token: Writing token to local cache: ✔  OK

All done. Thank you for flying with SOMA.
Suggested next steps:
        - create system_admin permission
        - create your team
        - create your user
        - grant system_admin to your user
        - activate your user
        - switch to using your user instead of root
```

12. Validate again, because you did not see the first validation and distrust the machine:

```
  % soma -u root permission list in omnipotence | jq .
	{
		"statusCode": 200,
		"statusText": "OK",
		"errors": [],
		"permissions": [
			{
				"id": "00000000-0000-0000-0000-000000000000",
				"name": "omnipotence",
				"category": "omnipotence"
			}
		]
	}
```

## Installation of the permission system

1. Create the permission scope categories

```
soma category add global
soma category add identity
soma category add monitoring
soma category add operation
soma category add permission
soma category add repository
soma category add self
soma category add team
```

2. Create the sections for each category

```
soma section add action to permission
soma section add admin-mgmt to identity
soma section add attribute to global
soma section add bucket to repository
soma section add capability to monitoring
soma section add category to permission
soma section add check-config to repository
soma section add cluster to repository
soma section add datacenter to global
soma section add deployment to monitoring
soma section add entity to global
soma section add environment to global
soma section add group to repository
soma section add hostdeployment to global
soma section add instance to repository
soma section add instance-mgmt to global
soma section add job to self
soma section add job-mgmt to global
soma section add job-result-mgmt to global
soma section add job-status-mgmt to global
soma section add job-type-mgmt to global
soma section add level to global
soma section add metric to global
soma section add mode to global
soma section add monitoringsystem to monitoring
soma section add monitoringsystem-mgmt to global
soma section add node to team
soma section add node-config to repository
soma section add node-mgmt to global
soma section add oncall to global
soma section add permission to permission
soma section add predicate to global
soma section add property-custom to repository
soma section add property-mgmt to global
soma section add property-native to global
soma section add property-service to team
soma section add property-system to global
soma section add property-template to global
soma section add provider to global
soma section add repository to team
soma section add repository-config to repository
soma section add repository-mgmt to global
soma section add right to permission
soma section add section to permission
soma section add server to global
soma section add state to global
soma section add status to global
soma section add system to operation
soma section add team to self
soma section add team-mgmt to identity
soma section add unit to global
soma section add user to self
soma section add user-mgmt to identity
soma section add validity to global
soma section add view to global
soma section add workflow to operation
```

3. Create the actions for each section

```
soma action add add to action
soma action add add to admin-mgmt
soma action add add to capability
soma action add add to category
soma action add add to datacenter
soma action add add to entity
soma action add add to environment
soma action add add to job-result-mgmt
soma action add add to job-status-mgmt
soma action add add to job-type-mgmt
soma action add add to metric
soma action add add to mode
soma action add add to monitoringsystem-mgmt
soma action add add to node-mgmt
soma action add add to oncall
soma action add add to permission
soma action add add to predicate
soma action add add to property-custom
soma action add add to property-mgmt
soma action add add to property-native
soma action add add to property-service
soma action add add to property-system
soma action add add to property-template
soma action add add to provider
soma action add add to section
soma action add add to server
soma action add add to state
soma action add add to status
soma action add add to team-mgmt
soma action add add to unit
soma action add add to user-mgmt
soma action add add to validity
soma action add add to view
soma action add all to instance-mgmt
soma action add assemble to hostdeployment
soma action add assign to node
soma action add assign to node-config
soma action add audit to repository
soma action add create to bucket
soma action add create to check-config
soma action add create to cluster
soma action add create to group
soma action add create to repository-mgmt
soma action add destroy to bucket
soma action add destroy to check-config
soma action add destroy to cluster
soma action add destroy to group
soma action add destroy to repository
soma action add failed to deployment
soma action add filter to deployment
soma action add get to hostdeployment
soma action add grant to right
soma action add insert-null to server
soma action add list to action
soma action add list to attribute
soma action add list to bucket
soma action add list to capability
soma action add list to category
soma action add list to check-config
soma action add list to cluster
soma action add list to datacenter
soma action add list to deployment
soma action add list to entity
soma action add list to environment
soma action add list to group
soma action add list to instance
soma action add list to job
soma action add list to job-mgmt
soma action add list to job-result-mgmt
soma action add list to job-status-mgmt
soma action add list to job-type-mgmt
soma action add list to level
soma action add list to metric
soma action add list to mode
soma action add list to monitoringsystem
soma action add list to node
soma action add list to oncall
soma action add list to permission
soma action add list to predicate
soma action add list to property-custom
soma action add list to property-mgmt
soma action add list to property-native
soma action add list to property-service
soma action add list to property-system
soma action add list to property-template
soma action add list to provider
soma action add list to repository
soma action add list to repository-config
soma action add list to right
soma action add list to section
soma action add list to server
soma action add list to state
soma action add list to status
soma action add list to team-mgmt
soma action add list to unit
soma action add list to user-mgmt
soma action add list to validity
soma action add list to view
soma action add list to workflow
soma action add map to permission
soma action add member-assign to cluster
soma action add member-assign to group
soma action add member-assign to oncall
soma action add member-list to cluster
soma action add member-list to group
soma action add member-list to oncall
soma action add member-list to team-mgmt
soma action add member-unassign to cluster
soma action add member-unassign to group
soma action add member-unassign to oncall
soma action add pending to deployment
soma action add property-create to bucket
soma action add property-create to cluster
soma action add property-create to group
soma action add property-create to node-config
soma action add property-create to repository-config
soma action add property-destroy to bucket
soma action add property-destroy to cluster
soma action add property-destroy to group
soma action add property-destroy to node-config
soma action add property-destroy to repository-config
soma action add property-update to bucket
soma action add property-update to cluster
soma action add property-update to group
soma action add property-update to node-config
soma action add property-update to repository-config
soma action add purge to node-mgmt
soma action add purge to server
soma action add purge to team-mgmt
soma action add purge to user-mgmt
soma action add rebuild-repository to system
soma action add remove to action
soma action add remove to admin-mgmt
soma action add remove to attribute
soma action add remove to capability
soma action add remove to category
soma action add remove to datacenter
soma action add remove to entity
soma action add remove to environment
soma action add remove to job-result-mgmt
soma action add remove to job-status-mgmt
soma action add remove to job-type-mgmt
soma action add remove to metric
soma action add remove to mode
soma action add remove to monitoringsystem-mgmt
soma action add remove to node-mgmt
soma action add remove to oncall
soma action add remove to permission
soma action add remove to predicate
soma action add remove to property-custom
soma action add remove to property-mgmt
soma action add remove to property-native
soma action add remove to property-service
soma action add remove to property-system
soma action add remove to property-template
soma action add remove to provider
soma action add remove to section
soma action add remove to server
soma action add remove to show
soma action add remove to state
soma action add remove to status
soma action add remove to team-mgmt
soma action add remove to unit
soma action add remove to user-mgmt
soma action add remove to validity
soma action add remove to view
soma action add rename to cluster
soma action add rename to datacenter
soma action add rename to entity
soma action add rename to environment
soma action add rename to repository
soma action add rename to state
soma action add rename to view
soma action add repossess to repository
soma action add restart-repository to system
soma action add retry to workflow
soma action add revoke to right
soma action add search to action
soma action add search to bucket
soma action add search to capability
soma action add search to check-config
soma action add search to cluster
soma action add search to group
soma action add search to job
soma action add search to job-result-mgmt
soma action add search to job-status-mgmt
soma action add search to job-type-mgmt
soma action add search to monitoringsystem
soma action add search to monitoringsystem-mgmt
soma action add search to node
soma action add search to oncall
soma action add search to permission
soma action add search to property-custom
soma action add search to property-mgmt
soma action add search to property-native
soma action add search to property-service
soma action add search to property-system
soma action add search to property-template
soma action add search to repository
soma action add search to repository-config
soma action add search to right
soma action add search to section
soma action add search to server
soma action add search to show
soma action add search to team
soma action add search to team-mgmt
soma action add search to user
soma action add search to user-mgmt
soma action add search to workflow
soma action add search-all to monitoringsystem-mgmt
soma action add set to workflow
soma action add show to action
soma action add show to admin-mgmt
soma action add show to attribute
soma action add show to bucket
soma action add show to capability
soma action add show to category
soma action add show to check-config
soma action add show to cluster
soma action add show to datacenter
soma action add show to deployment
soma action add show to entity
soma action add show to environment
soma action add show to group
soma action add show to instance
soma action add show to instance-mgmt
soma action add show to job
soma action add show to job-result-mgmt
soma action add show to job-status-mgmt
soma action add show to job-type-mgmt
soma action add show to level
soma action add show to metric
soma action add show to mode
soma action add show to monitoringsystem
soma action add show to node
soma action add show to oncall
soma action add show to permission
soma action add show to predicate
soma action add show to property-custom
soma action add show to property-mgmt
soma action add show to property-native
soma action add show to property-service
soma action add show to property-system
soma action add show to property-template
soma action add show to provider
soma action add show to repository
soma action add show to repository-config
soma action add show to right
soma action add show to section
soma action add show to server
soma action add show to state
soma action add show to status
soma action add show to team
soma action add show to team-mgmt
soma action add show to unit
soma action add show to user
soma action add show to user-mgmt
soma action add show to validity
soma action add show to view
soma action add show-config to node
soma action add shutdown to system
soma action add stop-repository to system
soma action add success to deployment
soma action add summary to workflow
soma action add sync to datacenter
soma action add sync to node-mgmt
soma action add sync to server
soma action add sync to team-mgmt
soma action add sync to user-mgmt
soma action add token to system
soma action add tree to bucket
soma action add tree to cluster
soma action add tree to group
soma action add tree to node-config
soma action add tree to repository-config
soma action add unassign to node
soma action add unassign to node-config
soma action add unmap to permission
soma action add update to bucket
soma action add update to check-config
soma action add update to cluster
soma action add update to group
soma action add update to node-mgmt
soma action add update to oncall
soma action add update to repository-config
soma action add update to server
soma action add update to team-mgmt
soma action add update to user-mgmt
soma action add use to monitoringsystem
soma action add versions to instance
soma action add wait to job
soma action add wait to job-mgmt
```

4. Create permissions within their scope. These are default permissions
   that must be created.

```
soma permission add browse to global
soma permission add information to self
soma permission add viewer to permission
soma permission add designer to permission
soma permission add auditor to permission
soma permission add admin to permission
```

5. Map actions to the required default permissions

```
soma permission map attribute::list to global::browse
soma permission map attribute::show to global::browse
soma permission map datacenter::list to global::browse
soma permission map datacenter::show to global::browse
soma permission map entity::list to global::browse
soma permission map entity::show to global::browse
soma permission map environment::list to global::browse
soma permission map environment::show to global::browse
soma permission map job-result-mgmt::list to global::browse
soma permission map job-result-mgmt::search to global::browse
soma permission map job-result-mgmt::show to global::browse
soma permission map job-status-mgmt::list to global::browse
soma permission map job-status-mgmt::search to global::browse
soma permission map job-status-mgmt::show to global::browse
soma permission map job-type-mgmt::list to global::browse
soma permission map job-type-mgmt::search to global::browse
soma permission map job-type-mgmt::show to global::browse
soma permission map level::list to global::browse
soma permission map level::show to global::browse
soma permission map metric::list to global::browse
soma permission map metric::show to global::browse
soma permission map mode::list to global::browse
soma permission map mode::show to global::browse
soma permission map oncall::list to global::browse
soma permission map oncall::search to global::browse
soma permission map oncall::show to global::browse
soma permission map predicate::list to global::browse
soma permission map predicate::show to global::browse
soma permission map property-mgmt::list to global::browse
soma permission map property-mgmt::search to global::browse
soma permission map property-mgmt::show to global::browse
soma permission map property-native::list to global::browse
soma permission map property-native::search to global::browse
soma permission map property-native::show to global::browse
soma permission map property-system::list to global::browse
soma permission map property-system::search to global::browse
soma permission map property-system::show to global::browse
soma permission map property-template::list to global::browse
soma permission map property-template::search to global::browse
soma permission map property-template::show to global::browse
soma permission map provider::list to global::browse
soma permission map provider::show to global::browse
soma permission map server::list to global::browse
soma permission map server::search to global::browse
soma permission map server::show to global::browse
soma permission map state::list to global::browse
soma permission map state::show to global::browse
soma permission map status::list to global::browse
soma permission map status::show to global::browse
soma permission map unit::list to global::browse
soma permission map unit::show to global::browse
soma permission map validity::list to global::browse
soma permission map validity::show to global::browse
soma permission map view::list to global::browse
soma permission map view::show to global::browse

soma permission map job::list to self::information
soma permission map job::search to self::information
soma permission map job::show to self::information
soma permission map job::wait to self::information
soma permission map team::search to self::information
soma permission map team::show to self::information
soma permission map user::search to self::information
soma permission map user::show to self::information

soma permission map category::list to permission::viewer
soma permission map category::show to permission::viewer
soma permission map section::list to permission::viewer
soma permission map section::show to permission::viewer
soma permission map section::search to permission::viewer
soma permission map action::list to permission::viewer
soma permission map action::show to permission::viewer
soma permission map action::search to permission::viewer
soma permission map permission::list to permission::viewer
soma permission map permission::show to permission::viewer
soma permission map permission::search to permission::viewer

soma permission map permission::add to permission::designer
soma permission map permission::map to permission::designer
soma permission map permission::unmap to permission::designer

soma permission map right::list to permission::auditor
soma permission map right::show to permission::auditor
soma permission map right::search to permission::auditor

soma permission map permission::remove to permission::admin
soma permission map right::grant to permission::admin
soma permission map right::revoke to permission::admin
```

6. Create static data schema. This is the schema that must be created
   for SOMA to work properly.

```
soma state add clustered
soma state add grouped
soma state add standalone
soma state add unassigned

soma entity add repository
soma entity add bucket
soma entity add group
soma entity add cluster
soma entity add node
soma entity add template
soma entity add server
soma entity add monitoring
soma entity add team

soma view add any
soma view add local

soma status add awaiting_computation
soma status add computed
soma status add awaiting_rollout
soma status add rollout_in_progress
soma status add rollout_failed
soma status add active
soma status add awaiting_deprovision
soma status add deprovision_in_progress
soma status add deprovision_failed
soma status add deprovisioned
soma status add awaiting_deletion
soma status add blocked
soma status add none

soma mode add private
soma mode add public

soma attribute add credential_password cardinality once
soma attribute add credential_user cardinality once

soma property-mgmt native add environment
soma property-mgmt native add hardware_node
soma property-mgmt native add state
soma property-mgmt native add entity

soma property-mgmt system add tag
soma validity add tag direct true inherited false on repository
soma validity add tag direct true inherited true on bucket
soma validity add tag direct true inherited true on group
soma validity add tag direct true inherited true on cluster
soma validity add tag direct true inherited true on node

soma property-mgmt system add disable_check_configuration
soma validity add disable_check_configuration direct true inherited false on repository
soma validity add disable_check_configuration direct true inherited true on bucket
soma validity add disable_check_configuration direct true inherited true on group
soma validity add disable_check_configuration direct true inherited true on cluster
soma validity add disable_check_configuration direct true inherited true on node

soma property-mgmt system add disable_all_monitoring
soma validity add disable_all_monitoring direct true inherited false on repository
soma validity add disable_all_monitoring direct true inherited true on bucket
soma validity add disable_all_monitoring direct true inherited true on group
soma validity add disable_all_monitoring direct true inherited true on cluster
soma validity add disable_all_monitoring direct true inherited true on node

soma property-mgmt system add dns_zone
soma validity add dns_zone direct true inherited false on repository
soma validity add dns_zone direct true inherited true on bucket
soma validity add dns_zone direct true inherited true on group
soma validity add dns_zone direct true inherited true on cluster
soma validity add dns_zone direct true inherited true on node

soma property-mgmt system add fqdn
soma validity add fqdn direct true inherited false on repository
soma validity add fqdn direct true inherited true on bucket
soma validity add fqdn direct true inherited true on group
soma validity add fqdn direct true inherited true on cluster
soma validity add fqdn direct true inherited true on node

soma environment add default

soma job result-mgmt add pending
soma job result-mgmt add success
soma job result-mgmt add failed

soma job status-mgmt add queued
soma job status-mgmt add in_progress
soma job status-mgmt add processed

soma job type-mgmt add bucket::create
soma job type-mgmt add bucket::destroy
soma job type-mgmt add bucket::property-create
soma job type-mgmt add bucket::property-destroy
soma job type-mgmt add bucket::property-update
soma job type-mgmt add bucket::rename
soma job type-mgmt add check-config::create
soma job type-mgmt add check-config::destroy
soma job type-mgmt add cluster::create
soma job type-mgmt add cluster::destroy
soma job type-mgmt add cluster::member-assign
soma job type-mgmt add cluster::member-unassign
soma job type-mgmt add cluster::property-create
soma job type-mgmt add cluster::property-destroy
soma job type-mgmt add cluster::property-update
soma job type-mgmt add group::create
soma job type-mgmt add group::destroy
soma job type-mgmt add group::member-assign
soma job type-mgmt add group::member-unassign
soma job type-mgmt add group::property-create
soma job type-mgmt add group::property-destroy
soma job type-mgmt add group::property-update
soma job type-mgmt add node-config::assign
soma job type-mgmt add node-config::property-create
soma job type-mgmt add node-config::property-destroy
soma job type-mgmt add node-config::property-update
soma job type-mgmt add node-config::unassign
soma job type-mgmt add repository-config::property-create
soma job type-mgmt add repository-config::property-destroy
soma job type-mgmt add repository-config::property-update
soma job type-mgmt add repository::destroy
soma job type-mgmt add repository::rename
soma job type-mgmt add repository::repossess
```

7. Create site-specific data schema

```
soma environment add ${env}

soma level add ${lvl} shortname ${short} numeric ${num}

soma predicate add ${pred}

soma view add ${view}

soma datacenter add ${locode}

soma server null datacenter ${locode}

soma unit add ${symbol} name ${name}

soma metric add ${name} unit ${symbol} description "${text}"
```

8. Import inventory information

```
soma team-mgmt add ${team-name} ldap ${ldapID}

soma user-mgmt add ${username} firstname ${fname} lastname ${lname} employeenr ${num} mailaddr ${addr} team ${team-name}

soma oncall add ${oncallduty} phone ${extension}

soma oncall member assign ${username} to ${oncallduty}

soma node add ${node-name} team ${team-name} online ${isOnline} server ${server-name} assetid ${assetID}
```

9. Create monitoring system definitions

```
soma monitoringsystem-mgmt add ${name} mode public contact ${username} team ${teamname}
soma capability declare ${monitoringsystem-name} view ${view} metric ${metric} thresholds 3
```

10. Elevate initial user account to admin

```
soma user-mgmt admin grant ${username}
soma right grant system::global to admin admin_${username}
soma right grant system::repository to admin admin_${username}
soma right grant system::team to admin admin_${username}
soma right grant system::identity to admin admin_${username}
soma right grant system::monitoring to admin admin_${username}
soma right grant system::permission to admin admin_${username}
soma right grant system::operation to admin admin_${username}
soma right grant system::self to admin admin_${username}
```

11. Permissions to grant to regulat users

```
soma --admin right grant global::browse to user ${username}
soma --admin right grant permission::viewer to user ${username}
soma --admin right grant self::information to user ${username}
```
