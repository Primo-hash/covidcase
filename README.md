# Assignment 1: Context-sensitive Exchange Service

## Application Overview

### Progress
Deployment URL: https://exchangeserve.herokuapp.com/
Repository: https://git.gvk.idi.ntnu.no/Primo/exchangeserve (Internal)

NOTE:
* Try the Deployment URL first for an overview of the API with examples

Missing:
* Keys in JSON for exchangehistory api are not ordered right

### Setup
* Open the /exchangeserve folder in a terminal and execute the following commands
```
cd ./cmd
go build ./exchange_server.go
PORT=8080 ./exchange_server
```
* The service should now be running on http://localhost:8080/

## Assignment Overview

In this assignment, you are going to develop a REST web application in Golang that provides the client to retrieve information about currency exchange rates. For this purpose, you will interrogate an existing web service and return the result in a given output format.

The REST web services you will be using for this purposes are:
* https://exchangeratesapi.io/
* https://restcountries.eu/

The first API focuses on the provision of exchange rates - as you will have explored before, whereas the second one provides country information (including currency information) that you will need in order to complete the assignment.

The API documentation is provided under the corresponding links, and both services vary vastly with respect to feature set and quality of documentation. Use [Postman](https://www.postman.com/) to explore the APIs, but be mindful of rate-limiting.

*A general note: when you develop your services that interrogate existing services, try to find the most efficient way of retrieving the necessary information. This generally means reducing the number of requests to these services to a minimum by using the most suitable endpoint that those APIs provide.*

The final web service should be deployed on [Heroku](https://www.heroku.com/). The initial development should occur on your local machine. For the submission, you will need to provide both a URL to the deployed Heroku service as well as your code repository.

In the following, you will find the specification for the REST API exposed to the user for interrogation/testing.

# Specification

Note: Please post an issue if the specification is unclear - so we can clarify and refine it if needed.

## Endpoints

Your web service will have three resource root paths:

```
/exchange/v1/exchangehistory/
/exchange/v1/exchangeborder/
/exchange/v1/diag/
```

Assuming your web service should run on localhost, port 8080, your resource root paths would look something like this:

```
http://localhost:8080/exchange/v1/exchangehistory/
http://localhost:8080/exchange/v1/exchangeborder/
http://localhost:8080/exchange/v1/diag/
````

The supported request/response pairs are specified in the following.

For the specifications, the following syntax applies:
* ```{:value}``` indicates mandatory input parameters specified by the user (i.e., the one using *your* service).
* ```{value}``` indicates optional input specified by the user (i.e., the one using *your* service), where `value' can itself contain further optional input. The same notation applies for HTTP parameter specifications (e.g., ```{?param}```).

## Exchange Rate History for Given Currency

The initial endpoint focuses on return the history of exchange rates (against EUR as a fixed base currency) for the currency of a given country, where start and end date are provided. The currency is to be determined based on the country name, and where a country has multiple currencies, only the first one is considered.

### Request

```
Method: GET
Path: exchangehistory/{:country_name}/{:begin_date-end_date}
```

```{:country_name}``` refers to the English name for the country as supported by https://restcountries.eu/.

```{:begin_date-end_date}``` indicates the begin date (i.e., the earliest date to be reported) of the exchange rate and the end date (i.e., the latest date of the range) of the period over which exchange rates are reported.

Example request: ```exchangehistory/norway/2020-12-01-2021-01-31```

### Response

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise. Ensure to deal with errors gracefully.

Body (Example):
```
{
    "rates": {
        "2020-12-01": {
            "NOK": 1.1969
        },
        "2020-12-02": {
            "NOK": 1.1633
        },
        ...
        "2021-01-31": {
            "NOK": 1.1754
        }
    },
    "start_at": "2020-12-01",
    "base": "EUR",
    "end_at": "2021-01-31"
}
```

## Current Exchange Rate Bordering Countries

The second endpoint provides an overview of the *current exchange rates* of a given country (which is then the base currency) with all bordering countries.

### Request

```
Method: GET
Path: exchangeborder/{:country_name}{?limit={:number}}
```

```{:country_name}``` refers to the English name for the country as supported by https://restcountries.eu/.

```{?limit={:number}}``` is an optional parameter that limits the number of bordering countries (```number```) for which currencies are reported. **Updated for clarity (if your interpretation deviated based on the previous text), please clarify in your submission)**

Where countries have multiple currencies, only report the first one provided. Where no currency is reported, ignore the country.

Example request: ```exchangeborder/norway?limit=5```

### Response

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise. Ensure to deal with errors gracefully.

Body (Example):
```
{
    "rates": {
        "Sweden": {
            "currency": "SEK", 
            "rate": 1.1703
        },
        "Russia": {
            "currency": "RUB",
            "rate": 72.05
        }, 
        ...
    },
    "base": "NOK"
}
```

## Diagnostics interface

The diagnostics interface indicates the availability of individual services this service depends on. The reporting occurs based on status codes returned by the dependent services, and it further provides information about the uptime of the service.

### Request

```
Method: GET
Path: diag/
```

### Response

* Content type: `application/json`
* Status code: 200 if everything is OK, appropriate error code otherwise.

Body:
```
{
   "exchangeratesapi": "<http status code for exchangeratesapi API>",
   "restcountries": "<http status code for restcountries API>",
   "version": "v1",
   "uptime": <time in seconds from the last service restart>
}
```

Note: ```<some value>``` indicates placeholders for values to be populated by the service.

# Deployment

The service is to be deployed on Heroku. You will need to provide the URL to the deployed service as part of the submission.

# General Aspects

As indicated during the initial sessions, ensure you work with professionalism in mind (see Course Rules). In addition to professionalism, you are at liberty to introduce further features into your service, as long it does not break the specification given above.

Please work in the provided workspace environment (see [here](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/wikis/Workspace-Conventions) - lodge an issue if you have trouble accessing it) for your user and create a project `assignment-1` in this workspace.

As mentioned above, be sensitive to rate limits of external services. If needed, consider mocking the remote services during development.

# Submission

The assignment is an individual assignment. The submission deadline is provided on the course main page. No extensions will be given for late submissions (unless the deadline is collectively extended, i.e., if we agree in class).

As part of the submission you will need to provide:
* a link to your code repository (ensure it is public at that stage)
* a link to the deployed Heroku service

In addition, we will provide you with an option to clarify aspects of your submission (e.g., aspects that don't quite work, or additional features).

The submission occurs via our submission system that not only facilitates the submission, but also the peer review of the assignment. Instructions for the submission system (submission, review) can be found [here](https://git.gvk.idi.ntnu.no/course/prog2005/prog2005-2021/-/wikis/Submission-System), and will also be introduced in class. Early submission is explicitly encouraged.

# Peer Review

After the submission deadline, there will be a second deadline during which you will review other students' submissions. To do this the system provides you with a checklist of aspects to assess. You will need to review *at least two submissions* to meet the mandatory requirements of peer review, but you can review as many submissions as you like, which counts towards your participation mark for the course. The peer-review deadline will be indicated closer to submission time.

