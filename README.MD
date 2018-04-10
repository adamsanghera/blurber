# Blurber 

## Main features

- [x] Login/ register
- [X] Delete account
- [X] Write blurb
- [X] Subscribe to blurbers
- [x] Unsubscribe on feed
- [x] Simple profile page
- [x] Delete personal single blurbs
- [ ] Feed: show all blurbs ('Discover') vs. just subscription

## Roadmap for detaching backend

1. [X] Define protobufs for each backend service:

       - [X] BlurbDB
       - [X] UserDB
       - [X] SubscriptionDB

1. [ ] Split off backends to be behind a GRPC call:

       - [ ] BlurbDB
       - [ ] UserDB
       - [ ] SubscriptionDB

1. [ ] Ensure that backends are safe to concurrent access from multiple frontends
       - [ ] BlurbDB
       - [ ] UserDB
       - [ ] SubscriptionDB

These will be managed in branches:

1. `protobufs`

2. `split`

    - `split-blurb`
    - `split-user`
    - `split-subscription`

3. `safe`

    - `safe-blurb`
    - `safe-user`
    - `safe-subscription`

## How to run

1. `PORT=4000 go run main.go` or `PORT=4000 DEBUG=1 go run main.go`
1. `http://localhost:4000`

## In the Wild

This app can also be found deployed on heroku, at https://dist-blurber.herokuapp.com/
