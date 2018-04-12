# Blurber

## Main features

- [x] Login/ register
- [X] Delete account
- [X] Write blurb
- [X] Subscribe to blurbers
- [x] Unsubscribe on feed
- [x] Simple profile page
- [x] Delete personal single blurbs
- [ ] OPTIONAL: Feed option to show all blurbs ('Discover' page) vs. just subscription

## Roadmap for detaching backend

- [X] Define protobufs for each backend service:

  - [X] BlurbDB
  - [X] UserDB
  - [X] SubscriptionDB

- [x] Split off backends to be behind a GRPC call:

  - [X] BlurbDB
  - [x] UserDB
  - [X] SubscriptionDB

- [ ] Ensure that backends are safe to concurrent access from multiple frontends

  - [ ] BlurbDB
  - [ ] UserDB
  - [ ] SubscriptionDB

These will be managed in branches:

1. `protobufs` [X]

2. `split` [x]

    - `split-blurb` [X]
    - `split-user` [x]
    - `split-subscription` [X]

3. `safe` []

    - `safe-blurb` []
    - `safe-user` []
    - `safe-subscription` []

## How to run

1. `cd` your way into `./protobufs`, and exec `./compile.sh`.
1. `cd` your way into `./` and exec `./spinBackend.sh`
1. Note that you'll need to manually kill the processes that are spun off here!!
1. `cd` your way into `./` and exec `.runLocal.sh`, if the backend servers are on your laptop.  Otherwise, exec `./run.sh` if the backend is on the linode server.
1. Point your favorite browser @ `http://localhost:4000`

## In the Wild

This app can also be found deployed on heroku, at https://dist-blurber.herokuapp.com/

## Team

- Adam Sanghera (as8513)
- Hieu Do (hsd258)
