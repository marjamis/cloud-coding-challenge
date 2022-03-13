# Stan Cloud Engineering Challenge

---

**NOTICE:** This is my submission from many years ago and therefore isn't currently up to newer standards and/or approaches. I may over time go about updating this repo, such as using CDK instead of CFN templates, but there is no guarantee and should be reviewed with this in mind.

---

## Overview

The purpose of this challenge is to assess your skills and approach to automated provisioning, deployment and configuration management.

Please complete the basic goals below and deliver the outputs as requested. You don't need to spend more than a couple of hours on the challenge.

Additional challenges are available if you have time and want to show us what you can do.

Completing the basic goals in an automated, reliable, reproducible way is preferable to completing any of the additional challenges.

## Basic Goals

* Automate the deployment of secure, publicly available, load-balanced web servers that return the instance ID of the host that served the request.
* Ensure that the web servers are available in two AWS availability zones and will automatically rebalance themselves if there is no healthy web server instance in either availability zone.
* Redirect any HTTP requests to HTTPS. Self-signed certificates are acceptable.
* Write one or two paragraphs about how you might further improve the automation and management of the infrastructure if you were to take it into production.

*Note:* All of the services required are eligible for the [AWS Free Tier](https://aws.amazon.com/free/).

## Additional Challenges

* Provide basic automated tests to cover included scripts, templates, manifests, recipes, code, etc.
* Return a custom static page for 4XX/5XX errors.
* Add a database to the infrastructure and have your application serve something from the database in addition to the instance ID.
* Write one or two paragraphs about why you solved the basic goals and additional challenges the way you did.

## Output

Please include in your e-mail:

1. a public URL to access your deployment
1. your written answers
1. any scripts, config files, manifests, recipes, or source code you used to achieve the goals, in the form of
    * a public source repository URL,
    * a zipped file,
    * or a zipped file via download link

We'll complete our review of the live environment within 2 business days so that you can tear it down promptly.

---

# My notes

## General Notes

* Missed the Additional Challenge for puppet due to personal time contraints also to ensure I submitted in a reasonable time but my plan was to have puppet run locally to perform all the on-instance configurations. In that way puppet is still controlling the configuration without the need for the normal server/client setup for this assessment.

## Output

1. [http://metis.qc.to/](http://metis.qc.to/) / [https://metis.qc.to/](https://metis.qc.to/)
2.

* AWS_ACCESS_KEY=```<redacted>```
* AWS_SECRET_ACCESS_KEY=```<redacted>```

3. SSH Keys located at: ./ssh/
4. Makefiles have targets to validate/test their respective areas. The application has a simple functions_test.go with some basic unit tests in the application src directory.
5. For the most part I believe the setup is pretty straight forward. There are Makefiles that control the creation of the infrastructure via CloudFormation API calls. If using the primary Makfile(in the root of the zip) with the "all" target it should perform all the operations including building/pushing any required docker images to an ECR Repo. Some general notes about the process:

* Currently requires and already created S3 Bucket and ECR repository
* The Makefiles have some environment variables that need to be updated to ensure the correct resources/locations are used.
* To access the instances requires use of the tempoary bastion I created with the instanceId:

```<redacted>```

6. How would you further automate the management of the infrastructure?

There are quite a few things that I would add and thought a list would be easier with brief explanations as required for each point.

* During environment creation increased signalling to CFN after steps completed in the instance configuration. Currently I just have the one at the end of my UserData but greater control would be good even when using other systems, such as puppet.
* Secure transmission/pulling of credentials from a source to the instances for use. Such as with S3/KMS or ParameterStore.
* Automate NameServers being added. Currently as Im using FreeDNS I had to do this manually when the HostedZone was created.
* Once the IAM certificate is added to the account I currently have to update the parameters for the certificates ARN, this should be automated and apart of the below listed pipeline.
* Proper Bastion Hosts for access via CFN
  * While the instances should have SSH blocked(via SG's) it's a good coverall to be able to add SG's rules to allow SSH access to the instance. This shouldn't be done unless absolutely necessary but it's good for all setups to be there to be able to SSH into a machine after temporarily renabling som SG rules that allow access from the bastion
* Improve the orchestration for the containers. Such as using an ECS Cluster for the environment to allow a higher resource utilisation(assuming this app doesn't take all the resources of a machine). To simplify this assessment I took the UserData approach which really robust at all but does the trick and does implicity have the instances as immutable during updates as they're replaced.
* Additional logging such as with CloudWatch logs and ALB/ELB/CLB Access Logs to centralise all logs/metrics without a need to ever go onto a machine.
* Refined scaling on a metric that makes sense for the application, an example potentially being latency, rather than CPUUtilization.
* Resource Tagging where ever possible for both billing, permissions and general management.
* Modularise the templates more for repeat use and just as a good practice for control. Such as having security related components, IAM roles and the like in seperate Templates that Security control independent of the environment for controls of accesses.
* Cross regional environment for uptime with the use of R53 using alias records and health checks
* RDS DB while on an internal Subnet still has the CNAME leaked via a dig to the public Route53 HostedZone so this really should be broken into Private and Public HostedZones as required.
* Connections from the ELB to the back-end instances also to be HTTPs
* After the lambda function creates the initial table disable its ability to access the database completely
* Use of CFN change-sets to ensure the changes that are made to a Stack are know.
* Complete CI/CD pipelines for both the application and the stack
  * Stack Pipeline
    * Allows monitor of a source, such as a repo.
    * Perform checks against templates and the outputs to ensure compliance with any infrastructure requirements.
    * Build and then test infrastructure before changes are made to other environments
    * Treat the Templates very much like code and allow very similar checks as you would code at each stage.
* End-to-End testing of the application prior to any push to production, Performance testing included as well
* Lock down IAM permissions for all roles a lot more
* Overall refine the process quite a bit more and make it more extensible

7. There are also quite a few things I would do from the application side, some of these being:

* Fix the database connection logic for retries on failures and healthcheck against this and other things. Currently the connection pool is made and if theres an issue there isn't an in application attempt to create a new connection pool
* Improved log/error output for both formats and details
  * Timestamps in calls being from handlers and between other functions, though a lot of this would be better served via profiling in something like xray or opentracing
* HTTPS connections to the RDS DB as well.
* Improve DB queries and the configuration, such as prepared statements, especially as a security measure
* Improve the test coverage from unit tests
