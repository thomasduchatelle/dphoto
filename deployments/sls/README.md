DPhoto - VIEWER
=======================================

Online web app, originally designed to view backed up photos and video from the cloud. 

Getting Started
---------------------------------------

App use _serverless.com_ deploys on AWS:

* `api/lambdas`: serverless backend to serve data for the web (golang)
* `web`: create-react-app deployed as static website on S3 (typescript)

They are all deployed as a monolith using _Serverless Framework_. To deploy on dev:

    # test, build, and deploy
    make deploy-app

    # only re-deploy
    sls deploy

    # destroy all software bits (not the data stored on 'infra-data')
    sls remove

Backend code is following the hexagonal architecture: core logic is developed in `domain` and imported here.

Authentication - Design Decisions
---------------------------------------

Authentication requirement is only to use Google oauth. AWS Cognito and Auth0 has been investigated to bring a larger authentication-as-a-service, but final decision is to use directly Google APIs.

AaaS solutions bring sign-in, email/phone number confirmation, password recovery, MFA, ... that are all non-wanted and should be disabled. Ease to use and integrate is the main expectation.

_AWS Cognito_ is complex to provision and does not provide React library (including Amplify) or documentation making the integration simple. The Access Token is not customisable: subject is the user UUID from Cognito (not recognised by DPhoto) and no
claims can support multi-tenancy of DPhoto. A second level of authentication, or verification against database, would be required to use it.

_Auth0_ is much easier to integrate with their JS library, and has customisation claims feature on the access token. Provisioning it, especially from CloudFormation (Serverless Framework), is very complex.

_Google Identity_, retained solution, requires a lot of manual development: on the UI an opensource react component is used to redirect to Google and get the identity token from the user. Then this token is used to authenticate on DPhoto and get an
access token.
