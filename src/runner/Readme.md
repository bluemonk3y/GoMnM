# Runner is the services runner

* Needs to be run with docker environment initialised: @FOR /f "tokens=*" %i IN ('docker-machine env default') DO @%i

* Builder the Docker Image with layers including 1) The Runner, 2) Your process
* Parameterise the following
  * 1) Service Label & Version(s) 2) Output Service Label & Version(s); Version(s) could be a wildcard
  * a) Admin Topic
  * b) KPI/Semantic Topic
  * Application Label (i.e. Monkey-2016-100)

## Building the microservice instance for /etc.SimpleMs
  * Create image: docker build -t blu3monk3y/simple-ms:v1 .
  * Publish: docker commit --change "ENV DEBUG true" c3f279d17e0a  blu3monk3y/simple-ms:v1
     OR
     [BUILD from Dockerfile and commit]



## Runner process interacts with docker commands

  * [DEPLOY & START in detatched mode - return ID] - >docker run --name ATTACH-ID -dit blu3monk3y/simple-ms:v1
  Have to run interactive and terminal ??
      * [START]- docker [start|stop] <any-name> (start a local container)
  * [Attach to send DATA] -> docker attach --detach-keys=ctrl-c ATTACH-ID
     AND

     [STOP] - docker stop

     [UNDEPLOY] - docker stop XXX ; docker rm <any-container-name> ; docker rmi <any-image-name>

     How do I pipe to container stdin after docker run? - attach

