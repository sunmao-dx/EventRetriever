# EventRetriever

#### Intro
EventRetriever is a gitee event spider for datacache.

#### Architect
![](http://assets.processon.com/chart_image/6163e23e0791290cc7819291.png)

#### Domain Model
![](http://assets.processon.com/chart_image/616428d163768921fa176b05.png)

#### How to contribute
If you’re interested in contributing code, the best starting point is to have a look at our Gitee issues to see which tasks are the most urgent. 

Sunmao accepts PR's (pull requests) from all developers.

Issues can be submitted by anyone - either seasoned developers or newbies.

#### Installation

- **Step 1** Setting up the k8s environment, Google GKE or minikube or microk8s are ok for deployment.

- **Step 2** Setting up webhook url in gitee projects, in order to receive issue event requests from you project.

- **Step 3** Setting up `api_url` and `gitee_token` environment variables.

- **Step 4** Using Dockerfile to build docker image and then upload it to DockerHub.

- **Step 5** Deploy the project by yaml on k8s.







