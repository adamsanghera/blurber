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

- [X] Ensure that backends are safe to concurrent access from multiple frontends

  - [X] BlurbDB
  - [X] UserDB
  - [X] SubscriptionDB

These will be managed in branches:

1. [X] `protobufs`
1. [X] `split`
    - [X] `split-blurb`
    - [X] `split-user`
    - [X] `split-subscription`
1. `safe`
    - [X] `safe-blurb`
    - [X] `safe-user`
    - [X] `safe-subscription`

## Roadmap for replication

1. [X] Repurpose code from Lab 2 to work in this project
    - [X] Replace labrpc code with gRPC
      - [X] protobuf spec
    - [X] Add service/leader discovery to the protocol, instead of static configurations
    - [X] Write unit tests demonstrating functional replication
    - [X] Add a buffered channel for exporting newly committed commands, to be applied elsewhere
      - [X] Integrate this feature with the recovery calls
1. [X] Graft replication daemon onto subscription database, so that it makes use of replication
    - [X] Background thread for processing the buffered commit channel
    - [X] No state changes outside of background thread
    - [X] Unit tests proving it all works
1. [ ] Node failure
    - [ ] Proving failover security at the daemon level, with unit tests proving efficacy
    - [ ] Proving failover security at the backend-tier level, with unit tests proving efficacy
    - [ ] Integrating leader failover security into the front-end, maintaining a peers list of backend services, to pivot between on lack of contact

## How to run

1. `cd` your way into `./protobufs`, and exec `./compile.sh`.
1. `cd` your way into `./` and exec `./spinBackend.sh`
1. `cd` your way into `./` and exec `.runLocal.sh`, if the backend servers are on your laptop.  Otherwise, exec `./run.sh` if the backend is on the linode server.
1. Point your favorite browser @ `http://localhost:4000`
1. When you're done playing around, exec `./killBackend.sh`

## In the Wild

This app can also be found deployed on heroku, at https://dist-blurber.herokuapp.com/

The Front-End server (everything in the main.go hierarchy) is hosted on heroku, and is set up to continously build on the `/release` branch on our github repo.

The front-end server connects to a backend server hosted by Linode.com (it's cheap).

The three (3!) backend services are hosted on the Linode server, which the heroku app sends gRPC messages to.

## Team

- Adam Sanghera (as8513)
- Hieu Do (hsd258)
