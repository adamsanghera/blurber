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

- [X] Define protobufs for each backend service:

  - [X] BlurbDB
  - [X] UserDB
  - [X] SubscriptionDB

- [ ] Split off backends to be behind a GRPC call:

  - [ ] BlurbDB
  - [ ] UserDB
  - [ ] SubscriptionDB

- [ ] Ensure that backends are safe to concurrent access from multiple frontends

  - [ ] BlurbDB
  - [ ] UserDB
  - [ ] SubscriptionDB

These will be managed in branches:

1. `protobufs` [X]

2. `split` []

    - `split-blurb` []
    - `split-user` []
    - `split-subscription` []

3. `safe` []

    - `safe-blurb` []
    - `safe-user` []
    - `safe-subscription` []

## How to run

1. `cd` your way into `./protobufs`, and run `./compile.sh`.
1. `cd` your way into `./blurb-driver` and run `./run.sh`.
1. `cd` your way into `./` and run `.run.sh` 
1. Point your favorite browser @ `http://localhost:4000`

## In the Wild

This app can also be found deployed on heroku, at https://dist-blurber.herokuapp.com/

## Team

- Adam Sanghera (as8513)
- Hieu Do (hsd258)
