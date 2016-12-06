# Runner is the services runner

* Needs to be run with docker environment initialised: @FOR /f "tokens=*" %i IN ('docker-machine env default') DO @%i

* Builder the Docker Image with layers including 1) The Runner, 2) Your process
* Parameterise the following
  * 1) Service Label & Version(s) 2) Output Service Label & Version(s); Version(s) could be a wildcard
  * a) Admin Topic
  * b) KPI/Semantic Topic
  * Application Label (i.e. Monkey-2016-100)

## Building the microservice instance for /etc.SimpleMs
  * 1. Create image: docker build -t blu3monk3y/simple-ms:v1 .
  * Publish: docker commit --change "ENV DEBUG true" c3f279d17e0a  blu3monk3y/simple-ms:v1
     OR
     [BUILD from Dockerfile and commit]



## Runner process interacts with docker commandsk

  *  2.  [DEPLOY & START] - >docker run --name simple-ms -it blu3monk3y/simple-ms:v1
     ** Look at docker exec instead - then docker attach

     OR

     [START]- docker [start|stop] <any-name> (start a local container)

     Then

     [DATA] - pipe to stdin: docker attach --detach-keys=ctrl-a c4ca4f19d4cd
     AND

     [STOP] - docker stop

     [UNDEPLOY] - docker stop XXX ; docker rm <any-container-name> ; docker rmi <any-image-name>

     How do I pipe to container stdin after docker run? - attach

