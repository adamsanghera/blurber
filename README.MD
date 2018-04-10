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
1. - [X] BlurbDB
1. - [X] UserDB
1. - [X] SubscriptionDB
1. [ ] Split off backends to be behind a GRPC call:
1. - [ ] BlurbDB
1. - [ ] UserDB
1. - [ ] SubscriptionDB
1. [ ] Ensure that backends are safe to concurrent access from multiple frontends
1. - [ ] BlurbDB
1. - [ ] UserDB
1. - [ ] SubscriptionDB

These will be managed in branches:

1. `protobufs`
1. `split`
1. - `split-blurb`
1. - `split-user`
1. - `split-subscription`
1. `safe`
1. - `safe-blurb`
1. - `safe-user`
1. - `safe-subscription`

## How to run

1. `PORT=4000 go run main.go` or `PORT=4000 DEBUG=1 go run main.go`
1. `http://localhost:4000`

## In the Wild

This app can also be found deployed on heroku, at https://dist-blurber.herokuapp.com/
