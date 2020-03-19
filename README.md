# ubc-queue-alert
Alerts TA slack channels when the queue has hit a threshold

## Requirement
All you need is go 1.11+ with go mod setup (should be default)

## How to use

Setup the following environment variables:
* `TOKEN`: Your own private token retrieved from the queue platfrom
* `URL`: The url for the queue platform you are using (ex: "https://queue.students.cs.ubc.ca")
* `COURSE`: The course you are querying (ex: "CPSC 213")
* `DELAY`: The delay in minutes of which the program will query the queue
* `THRESHOLD`: The threshold that the app will notify your slack channel once the queue length exceeds
* `WEBHOOK`: A slack app [webhook](https://api.slack.com/messaging/webhooks) that points to a specific channel, this is where the app will be posting 

Then simply run:
```$ go run main.go```

## Contribution
This was made to support CPSC 213's move online with the COVID-19 situation. If you feel like this can help support your own course in any way, please feel free to use it!
If the need arises for adding more features that interact with the queue API, I will be adding it here and updating the README.
