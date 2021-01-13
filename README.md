# Url-Collector

## What's missing
- Graceful shutdown
- time.Time cannot currenlty be mocked, so we could introduce some kind of Clock interface; This prevens us from testing future dates provided in the request's URL

> What if we were to change the NASA API to some other images provider?

You simply need to implement Fetcher interface for another provider
>What if, apart from using NASA API, we would want to have another microservice fetching urls from
    European Space Agency. How much code could be reused?

Fetcher and client are very flexible interfaces and they could be reused  but we would need a different form of executing GetImages task, right now server depends on Fetcher and it directly calls Fetcher's GetImages, I'd like to avoid any logic in server's handler, so we would some kind of wrapper to wrap our fetchers and schedule our execution

>What if we wanted to add some more query params to narrow down lists of urls - for example,
selecting only images taken by certain person. (field copyright in the API response)

It is already possible with (filters ...Filter) that are passed to Fetcher's GetImages method