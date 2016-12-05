# Runner is the services runner

* Needs to be run with docker environment initialised: @FOR /f "tokens=*" %i IN ('docker-machine env default') DO @%i

* Builder the Docker Image with layers including 1) The Runner, 2) Your process
* Parameterise the following
  * 1) Service Label & Version(s) 2) Output Service Label & Version(s); Version(s) could be a wildcard
  * a) Admin Topic
  * b) KPI/Semantic Topic
  * Application Label (i.e. Monkey-2016-100)
